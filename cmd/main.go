package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/soerenschneider/alert-forge/internal"
	"github.com/soerenschneider/alert-forge/internal/alertmanager"
	"github.com/soerenschneider/alert-forge/internal/db/sqlite"
	"github.com/soerenschneider/alert-forge/internal/metrics"
	"github.com/soerenschneider/alert-forge/internal/reconciliation"
	"github.com/soerenschneider/alert-forge/internal/report"
	"github.com/soerenschneider/alert-forge/internal/report/receivers"
	"github.com/soerenschneider/alert-forge/internal/webhooks"
)

var (
	configFile string
	httpClient *http.Client = retryablehttp.NewClient().HTTPClient
)

func parseFlags() {
	flag.StringVar(&configFile, "config-file", "/etc/alertmanager-reports.yaml", "Config file")
	flag.Parse()
}

func main() {
	parseFlags()

	cfg, err := internal.ReadConfig(configFile)
	if err != nil {
		log.Fatal("could not read config file", err)
	}

	db, err := sqlite.NewSqliteStorage(cfg)
	if err != nil {
		log.Fatal("could not build sqlite storage", err)
	}

	clientOpts := []alertmanager.AlertmanagerOpt{
		alertmanager.WithClient(httpClient),
		alertmanager.WithBlacklist(cfg.IgnoredAlertLabelPairs),
	}
	client, err := alertmanager.NewAlertmanagerClient(cfg.AlertmanagerInstances, clientOpts...)
	if err != nil {
		log.Fatal("could not build Alertmanager client", err)
	}

	server, err := buildServer(cfg, db)
	if err != nil {
		log.Fatal("could not startup http server", err)
	}

	reportManager, err := buildReportManager(cfg, client, db)
	if err != nil {
		log.Fatal("could not build app", "err", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	go func() {
		if err := reportManager.Start(ctx, wg); err != nil {
			slog.Error("report manager error", "err", err)
		}
	}()

	reconciler, err := reconciliation.NewReconciler(db, client)
	if err != nil {
		log.Fatal("could not build reconciliation controller", err)
	}

	go func() {
		reconciler.Start(ctx, wg)
	}()

	go func() {
		if err := server.Listen(ctx, wg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("could not listen", err)
		}
	}()

	go func() {
		if cfg.MetricsAddr == "" {
			return
		}

		metricsServer, err := metrics.New(cfg.MetricsAddr)
		if err != nil {
			log.Fatal("could not create metrics server", err)
		}

		if err := metricsServer.StartServer(ctx, wg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("could not start metrics server", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigs
	cancel()

	// wait on all members of the waitgroup but end forcefully after the timeout has passed
	gracefulExitDone := make(chan struct{})

	go func() {
		slog.Info("Waiting for components to shut down gracefully")
		wg.Wait()
		close(gracefulExitDone)
	}()

	select {
	case <-gracefulExitDone:
		slog.Info("All components shut down gracefully within the timeout")
	case <-time.After(30 * time.Second):
		slog.Error("Components could not be shutdown within timeout, killing process forcefully")
	}
}

func buildServer(cfg *internal.Config, db webhooks.Db) (*webhooks.HttpServer, error) {
	renderer, err := webhooks.NewAlertRenderer()
	if err != nil {
		return nil, fmt.Errorf("could not build alert renderer: %w", err)
	}

	alertWebhook, err := webhooks.NewReceiver(db, renderer)
	if err != nil {
		return nil, fmt.Errorf("could not build webhook receiver: %w", err)
	}

	return webhooks.New(cfg.HttpServer.Address, alertWebhook)
}

func buildReportManager(cfg *internal.Config, client report.Alertmanager, db report.ReportDatabase) (*report.ReportManager, error) {
	var reportManagerOpts []report.ReportManagerOpt
	if len(cfg.IgnoredAlertLabelPairs) > 0 {
		reportManagerOpts = append(reportManagerOpts, report.WithBlacklist(cfg.IgnoredAlertLabelPairs))
	}

	var timezone *time.Location
	if cfg.Reports.Timezone != "" {
		var err error
		timezone, err = time.LoadLocation(cfg.Reports.Timezone)
		if err != nil {
			return nil, fmt.Errorf("could not parse provided timezone: %w", err)
		}
		reportManagerOpts = append(reportManagerOpts, report.WithTimezone(timezone))
	}

	digestCreator, err := report.NewDigestCreator()
	if err != nil {
		return nil, fmt.Errorf("could not build report schedule: %w", err)
	}

	reportReceivers, err := buildReportReceivers(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not build report receivers: %w", err)
	}

	reportManagerArgs := report.ReportManagerParams{
		Receivers:      reportReceivers,
		CronExpression: cfg.Reports.CronExpression,
		Client:         client,
		Db:             db,
		RenderEngine:   digestCreator,
	}

	return report.NewReportManager(reportManagerArgs, reportManagerOpts...)
}

func buildReportReceivers(cfg *internal.Config) ([]report.Receivers, error) {
	var sinks []report.Receivers

	if len(cfg.Reports.AwtrixInstances) > 0 {
		awtrix, err := receivers.NewAwtrix(cfg.Reports.AwtrixInstances[0], httpClient)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, awtrix)
	}

	if cfg.Reports.Email != nil {
		var opts []receivers.EmailOpt

		if cfg.Reports.Email.Port != 0 {
			opts = append(opts, receivers.WithPort(cfg.Reports.Email.Port))
		}

		from, err := cfg.Reports.Email.GetFrom()
		if err != nil {
			return nil, err
		}

		to, err := cfg.Reports.Email.GetTo()
		if err != nil {
			return nil, err
		}

		user, err := cfg.Reports.Email.GetUsername()
		if err != nil {
			return nil, err
		}

		password, err := cfg.Reports.Email.GetPassword()
		if err != nil {
			return nil, err
		}

		emailSink, err := receivers.NewEmail(from, to, cfg.Reports.Email.Server, user, password, opts...)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, emailSink)
	}
	return sinks, nil
}
