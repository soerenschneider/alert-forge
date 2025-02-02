package report

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/soerenschneider/alert-forge/internal/metrics"
	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/model/digest"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/multierr"
)

type ReportManagerOpt func(a *ReportManager) error

func WithBlacklist(labelsBlacklist map[string]string) ReportManagerOpt {
	return func(a *ReportManager) error {
		a.blacklistedLabels = labelsBlacklist
		return nil
	}
}

func WithTimezone(timezone *time.Location) ReportManagerOpt {
	return func(a *ReportManager) error {
		if timezone == nil {
			return errors.New("nil timezone provided")
		}
		a.timeZone = timezone
		return nil
	}
}

type Receivers interface {
	SendReport(ctx context.Context, subject, body string) error
}

type Alertmanager interface {
	GetActiveAlerts(ctx context.Context) ([]*model.Alert, error)
}

type ReportDatabase interface {
	GetLatestReport(ctx context.Context) (digest.StatusReport, error)
	SaveReport(ctx context.Context, r digest.StatusReport) error
	GetTodaysAlerts(ctx context.Context) ([]model.Alert, error)
}

type RenderEngine interface {
	Render(data digest.DigestData) ([]byte, error)
}

type ReportManager struct {
	receivers []Receivers

	client            Alertmanager
	db                ReportDatabase
	blacklistedLabels map[string]string
	renderEngine      RenderEngine
	cronExpression    string
	timeZone          *time.Location
}

type ReportManagerParams struct {
	Receivers []Receivers

	Client         Alertmanager
	Db             ReportDatabase
	RenderEngine   RenderEngine
	CronExpression string
}

func NewReportManager(params ReportManagerParams, opts ...ReportManagerOpt) (*ReportManager, error) {
	ret := &ReportManager{
		client:         params.Client,
		db:             params.Db,
		receivers:      params.Receivers,
		cronExpression: params.CronExpression,
		renderEngine:   params.RenderEngine,
		timeZone:       time.UTC,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(ret); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return ret, errs
}

func (a *ReportManager) Start(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	opts := []gocron.SchedulerOption{
		gocron.WithLocation(a.timeZone),
	}
	scheduler, err := gocron.NewScheduler(opts...)
	if err != nil {
		return err
	}

	_, err = scheduler.NewJob(
		gocron.CronJob(a.cronExpression, false),
		gocron.NewTask(
			func(ctx context.Context) {
				if err := a.run(ctx); err != nil {
					metrics.ReportsErrors.Inc()
					slog.Error("could not build and dispatch report", "err", err)
				}
			},
			ctx,
		),
	)
	if err != nil {
		return err
	}

	scheduler.Start()
	<-ctx.Done()

	return scheduler.Shutdown()
}

func (a *ReportManager) run(ctx context.Context) error {
	activeAlerts, err := a.client.GetActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("error fetching alerts: %w", err)
	}

	filteredAlerts := a.Filter(activeAlerts)

	reportActiveAlerts := NewStatusReport(filteredAlerts)
	reportPreviousAlerts, err := a.db.GetLatestReport(ctx)
	if err != nil {
		metrics.ReportsErrors.Inc()
		slog.Error("could not get most recent record from db", "err", err)
		reportPreviousAlerts = reportActiveAlerts
	} else {
		slog.Info("Loaded status report", "date", reportPreviousAlerts.Date)
	}

	if err := a.db.SaveReport(ctx, reportActiveAlerts); err != nil {
		metrics.ReportsErrors.Inc()
		slog.Error("could not save report", "err", err)
	}

	resolvedToday, err := a.db.GetTodaysAlerts(ctx)
	if err != nil {
		metrics.ReportsErrors.Inc()
		slog.Error("could not get today's resolved alerts", "err", err)
	}

	if reportActiveAlerts.AlertCount() == 0 && reportPreviousAlerts.AlertCount() == 0 && len(resolvedToday) == 0 {
		slog.Debug("no current or previous active alerts")
		return nil
	}

	if err := a.buildAndDispatchReports(ctx, filteredAlerts, resolvedToday, reportActiveAlerts, reportPreviousAlerts); err != nil {
		return fmt.Errorf("could not send report: %w", err)
	}

	return nil
}

func (a *ReportManager) buildAndDispatchReports(ctx context.Context, activeAlerts []model.Alert, resolvedToday []model.Alert, reportActiveAlerts, reportPreviousAlerts digest.StatusReport) error {
	categorizeReports := CompareStatusReports(reportActiveAlerts, reportPreviousAlerts)

	for severity, counts := range categorizeReports {
		total := 0
		for val := range maps.Values(counts) {
			total += val
		}
		slog.Info("Comparing alerts", "severity", severity, "total", total, "new", counts["new"], "gone", counts["gone"], "seen", counts["seen"])
	}

	groupedAlerts := digest.GroupAlertsBySeverity(activeAlerts)

	severityInfo := make([]digest.SeverityInfo, 0, len(groupedAlerts))
	for severity := range groupedAlerts {
		info := digest.SeverityInfo{
			Severity: severity,
		}
		digest.Convert(&info, categorizeReports)
		severityInfo = append(severityInfo, info)
	}

	digestData := digest.DigestData{
		Severities:     severityInfo,
		FiringAlerts:   activeAlerts,
		ResolvedAlerts: resolvedToday,
	}
	data, err := a.renderEngine.Render(digestData)
	if err != nil {
		return err
	}

	subject := FormatSubjectLine(activeAlerts)
	p := pool.New().WithErrors().WithMaxGoroutines(2).WithContext(ctx)
	for _, sink := range a.receivers {
		p.Go(func(ctx context.Context) error {
			return sink.SendReport(ctx, subject, string(data))
		})
	}

	return p.Wait()
}

func (a *ReportManager) Filter(alerts []*model.Alert) []model.Alert {
	ret := make([]model.Alert, 0, len(alerts))

	for _, alert := range alerts {
		filter := false
		for label, value := range a.blacklistedLabels {
			if alert.Labels[label] == value {
				filter = true
				continue
			}
		}

		if !filter {
			ret = append(ret, *alert)
		}
	}

	return ret
}

func NewStatusReport(alerts []model.Alert) digest.StatusReport {
	report := digest.StatusReport{
		Date:          time.Now(),
		SeverityCount: make(map[string]map[string]struct{}),
	}

	// Loop through all the alerts and count severities, also track fingerprints
	for _, alert := range alerts {
		severity := alert.Labels["severity"]
		if severity != "" {
			_, found := report.SeverityCount[severity]
			if !found {
				report.SeverityCount[severity] = make(map[string]struct{})
			}
			report.SeverityCount[severity][alert.Fingerprint] = struct{}{}
		}
	}

	return report
}
