package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenschneider/alert-forge/internal"
	"github.com/soerenschneider/alert-forge/internal/db/sqlite/generated"
	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/model/digest"
	"github.com/soerenschneider/alert-forge/pkg"
	"go.uber.org/multierr"
)

func NewSqliteStorage(conf *internal.Config) (*SqliteStore, error) {
	if conf == nil {
		return nil, errors.New("nil conf provided")
	}
	db, err := sql.Open("sqlite3", conf.SqliteDatabaseFile)
	if err != nil {
		return nil, err
	}

	gen := generated.New(db)
	ret := &SqliteStore{
		db:         db,
		generated:  gen,
		timeSource: &pkg.DefaultTime{},
		timeZone:   *time.UTC,
		conf:       conf,
	}
	return ret, ret.Migrate(context.Background())
}

type TimeSource interface {
	Now() time.Time
}

type SqliteStore struct {
	db         *sql.DB
	generated  *generated.Queries
	timeSource TimeSource
	timeZone   time.Location
	conf       *internal.Config
}

func (s *SqliteStore) AnalyzeLongestDuration(ctx context.Context) ([]generated.AnalyzeLongestDurationRow, error) {
	return s.generated.AnalyzeLongestDuration(ctx)

}
func (s *SqliteStore) SaveAlert(ctx context.Context, alert model.Alert) error {
	// update fingerprint
	alert.Fingerprint = alert.Hash()

	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	params := generated.SaveAlertParams{
		Name:    alert.Labels["alertname"],
		Started: alert.StartsAt.Unix(),
		Ended: sql.NullInt64{
			Valid: false,
		},
		Instance: pkg.RemovePort(alert.Labels["instance"]),
		Data:     data,
	}

	severityMapped := s.conf.SeverityMapping[strings.ToLower(alert.Labels["severity"])]
	params.Severity = severityMapped

	if !alert.EndsAt.IsZero() {
		params.Ended = sql.NullInt64{
			Int64: alert.EndsAt.Unix(),
			Valid: true,
		}
	}

	return s.generated.SaveAlert(ctx, params)
}

func (s *SqliteStore) GetAlertsTodayResolved(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlertsTodayResolved(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlerts(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetActiveAlerts(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetActiveAlerts(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetTodaysAlerts(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlertsToday(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlertsYesterday(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlertsYesterday(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlertsToday(ctx context.Context) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlertsToday(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlertsBySeverity(ctx context.Context, severity string) ([]model.Alert, error) {
	mapped, found := s.conf.SeverityMapping[severity]
	if !found {
		return nil, model.ErrNotFound
	}

	alerts, err := s.generated.GetAlertsBySeverity(ctx, mapped)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlertsByInstance(ctx context.Context, instance string) ([]model.Alert, error) {
	alerts, err := s.generated.GetAlertsByInstance(ctx, instance)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Alert, len(alerts))
	for idx, alert := range alerts {
		converted := model.Alert{}
		err := json.Unmarshal(alert.([]byte), &converted)
		if err != nil {
			return nil, err
		}
		ret[idx] = converted
	}
	return ret, nil
}

func (s *SqliteStore) GetAlert(ctx context.Context, fingerprint string) (model.Alert, error) {
	data, err := s.generated.GetAlert(ctx, fingerprint)
	if err != nil {
		return model.Alert{}, err
	}

	converted := model.Alert{}

	return converted, json.Unmarshal(data.([]byte), &converted)
}

func (s *SqliteStore) SaveReport(ctx context.Context, report digest.StatusReport) error {
	severity, err := json.Marshal(report.SeverityCount)
	if err != nil {
		return err
	}

	params := generated.SaveReportParams{
		Date: sql.NullString{
			String: report.Date.Format(time.RFC3339),
			Valid:  true,
		},
		SeverityCount: sql.NullString{
			String: string(severity),
			Valid:  true,
		},
	}

	return s.generated.SaveReport(ctx, params)
}

func (s *SqliteStore) GetLatestReport(ctx context.Context) (digest.StatusReport, error) {
	report, err := s.generated.GetLatestReport(ctx)
	if err != nil {
		return digest.StatusReport{}, err
	}

	return convertReport(report)
}

func convertReport(report generated.Report) (digest.StatusReport, error) {
	date, err := time.Parse(time.RFC3339, report.Date.String)
	if err != nil {
		return digest.StatusReport{}, fmt.Errorf("failed to parse date: %w", err)
	}

	ret := digest.StatusReport{
		Id:   report.ID,
		Date: date,
	}

	err = json.Unmarshal([]byte(report.SeverityCount.String), &ret.SeverityCount)
	if err != nil {
		return digest.StatusReport{}, fmt.Errorf("failed to unmarshal SeverityCount JSON: %w", err)
	}

	return ret, nil
}

func (s *SqliteStore) GetReport(ctx context.Context, id int64) (digest.StatusReport, error) {
	report, err := s.generated.GetReport(ctx, id)
	if err != nil {
		return digest.StatusReport{}, err
	}

	return convertReport(report)
}

func (s *SqliteStore) GetAlertsPerDayStats(ctx context.Context) ([]model.AlertsByDay, error) {
	stats, err := s.generated.AnalyzeAlertsPerDay(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AlertsByDay, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AlertsByDay{
			AlertDate:  stat.AlertDate,
			AlertCount: stat.AlertCount,
		}
	}

	return ret, nil
}

func (s *SqliteStore) GetAlertsPerInstanceStats(ctx context.Context) ([]model.AlertsByInstance, error) {
	stats, err := s.generated.AnalyzeInstances(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AlertsByInstance, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AlertsByInstance{
			Instance:      stat.Instance,
			InstanceCount: stat.InstanceCount,
		}
	}

	return ret, nil
}

func (s *SqliteStore) GetAverageAlertDuration(ctx context.Context) ([]model.AlertDuration, error) {
	stats, err := s.generated.AnalyzeAverageAlertDuration(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AlertDuration, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AlertDuration{
			Name:            stat.Name,
			AverageDuration: stat.AverageDuration,
		}
	}

	return ret, nil
}

func (s *SqliteStore) GetAlertsWithLongestDuration(ctx context.Context) ([]model.AlertDuration, error) {
	stats, err := s.generated.AnalyzeLongestDuration(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AlertDuration, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AlertDuration{
			Name:            stat.Name,
			AverageDuration: stat.Duration,
		}
	}

	return ret, nil
}

func (s *SqliteStore) GetAverageAlertDurationBySeverity(ctx context.Context) ([]model.AverageAlertDurationBySeverity, error) {
	stats, err := s.generated.AnalyzeAverageAlertDurationBySeverity(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AverageAlertDurationBySeverity, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AverageAlertDurationBySeverity{
			Severity:        stat.Severity,
			AverageDuration: stat.AverageDuration,
		}
	}

	return ret, nil
}

func (s *SqliteStore) GetStats(ctx context.Context) (model.AlertStats, error) {
	ret := model.AlertStats{}
	var errs, err error

	ret.AlertsByInstance, err = s.GetAlertsPerInstanceStats(ctx)
	errs = multierr.Append(errs, err)

	ret.AlertDuration, err = s.GetAverageAlertDuration(ctx)
	errs = multierr.Append(errs, err)

	ret.AlertsBySeverity, err = s.GetAlertsBySeverityStats(ctx)
	errs = multierr.Append(errs, err)

	ret.AlertsByDay, err = s.GetAlertsPerDayStats(ctx)
	errs = multierr.Append(errs, err)

	ret.AverageAlertDurationBySeverity, err = s.GetAverageAlertDurationBySeverity(ctx)
	errs = multierr.Append(errs, err)

	return ret, errs
}

func (s *SqliteStore) GetAlertsBySeverityStats(ctx context.Context) ([]model.AlertsBySeverity, error) {
	stats, err := s.generated.AnalyzeSeverities(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]model.AlertsBySeverity, len(stats))
	for idx, stat := range stats {
		ret[idx] = model.AlertsBySeverity{
			Severity:   stat.Severity,
			AlertCount: stat.AlertCount,
		}
	}

	return ret, nil
}

func (s *SqliteStore) Migrate(ctx context.Context) error {
	if schemaVersionReadError != nil {
		return schemaVersionReadError
	}

	var currentVersion int
	_ = s.db.QueryRowContext(ctx, `SELECT version FROM schema_version`).Scan(&currentVersion)

	slog.Info("Checking db schema", "current", currentVersion, "latest", schemaVersion)
	if currentVersion >= schemaVersion {
		return nil
	}

	migrations, err := GetMigrations()
	if err != nil {
		return err
	}

	for version := currentVersion; version < schemaVersion; version++ {
		newVersion := version + 1

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("can not start transaction %w", err)
		}

		sql := migrations[version]
		_, err = tx.ExecContext(ctx, string(sql))
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM schema_version`); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES ($1)`, newVersion); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}
		slog.Info("Successfully migrated DB", "version", newVersion)
	}

	return nil
}
