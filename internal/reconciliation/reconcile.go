package reconciliation

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/soerenschneider/alert-forge/internal/metrics"
	"github.com/soerenschneider/alert-forge/internal/model"
	"go.uber.org/multierr"
)

const (
	defaultReconciliationInterval = 15 * time.Minute
)

type db interface {
	GetActiveAlerts(ctx context.Context) ([]model.Alert, error)
	SaveAlert(ctx context.Context, alert model.Alert) error
}

type alertmanagerClient interface {
	GetActiveAlerts(ctx context.Context) ([]*model.Alert, error)
}

type AlertReconciler struct {
	db       db
	client   alertmanagerClient
	interval time.Duration
}

func NewReconciler(db db, client alertmanagerClient) (*AlertReconciler, error) {
	if db == nil {
		return nil, errors.New("nil db supplied")
	}

	if client == nil {
		return nil, errors.New("nil client supplied")
	}

	ret := &AlertReconciler{
		db:       db,
		client:   client,
		interval: defaultReconciliationInterval,
	}

	return ret, nil
}

func (m *AlertReconciler) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(m.interval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := m.reconcile(ctx); err != nil {
				slog.Error("error encountered during reconciliation", "err", err)
				metrics.ReconciliationsErrorsTotal.Inc()
			}
		}
	}
}

func (m *AlertReconciler) reconcile(ctx context.Context) error {
	metrics.ReconciliationsTotal.Inc()
	alertsFromDb, err := m.db.GetActiveAlerts(ctx)
	if err != nil {
		return err
	}

	alertsFromAm, err := m.client.GetActiveAlerts(ctx)
	if err != nil {
		return err
	}

	resolvedAlerts := findUnresolvedAlerts(alertsFromAm, alertsFromDb)
	metrics.ReconciledAlertsTotal.Set(float64(len(resolvedAlerts)))
	if len(resolvedAlerts) > 0 {
		slog.Warn("discovered alerts that are marked as active but are resolved", "amount", len(resolvedAlerts))
	}

	var errs error
	for _, alert := range resolvedAlerts {
		// resolve alert by setting end time
		alert.EndsAt = time.Now()
		if err := m.db.SaveAlert(ctx, alert); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func findUnresolvedAlerts(am []*model.Alert, db []model.Alert) []model.Alert {
	amIds := make(map[string]struct{})
	for _, alert := range am {
		amIds[alert.Fingerprint] = struct{}{}
	}

	var result []model.Alert
	for _, alert := range db {
		if _, found := amIds[alert.Fingerprint]; !found {
			result = append(result, alert)
		}
	}

	return result
}
