package webhooks

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/model/webhook"
	"go.uber.org/multierr"
)

type Db interface {
	SaveAlert(ctx context.Context, alert model.Alert) error
	GetAlert(ctx context.Context, fingerprint string) (model.Alert, error)
	GetAlerts(ctx context.Context) ([]model.Alert, error)
	GetAlertsByInstance(ctx context.Context, instance string) ([]model.Alert, error)
	GetAlertsBySeverity(ctx context.Context, severity string) ([]model.Alert, error)
	GetAlertsToday(ctx context.Context) ([]model.Alert, error)
	GetAlertsYesterday(ctx context.Context) ([]model.Alert, error)
}

type Stats interface {
	GetStats(ctx context.Context) (model.AlertStats, error)
}

type Templating interface {
	Render(alerts []model.Alert) ([]byte, error)
}

type Receiver struct {
	db          Db
	templating  Templating
	statsSource Stats
	stats       *StatsRenderer
}

func (r *Receiver) Statistics(w http.ResponseWriter, req *http.Request, params StatisticsParams) {
	stats, err := r.statsSource.GetStats(req.Context())
	if err != nil {
		slog.Error("could not get stats", "err", err)
		w.WriteHeader(500)
		return
	}

	if isBrowser(req) {
		data, err := r.stats.Render(stats)
		if err != nil {
			slog.Error("could render statistics", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(stats)
		if err != nil {
			slog.Error("could not marshal statistics", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func (r *Receiver) GetAlertsByInstance(w http.ResponseWriter, req *http.Request, instance string, params GetAlertsByInstanceParams) {
	alerts, err := r.db.GetAlertsByInstance(req.Context(), instance)
	if err != nil {
		slog.Error("could not get alerts", "err", err)
		if errors.Is(err, model.ErrNotFound) {
			w.WriteHeader(404)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	if isBrowser(req) {
		data, err := r.templating.Render(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func (r *Receiver) GetAlertsBySeverity(w http.ResponseWriter, req *http.Request, severity string, params GetAlertsBySeverityParams) {
	alerts, err := r.db.GetAlertsBySeverity(req.Context(), severity)
	if err != nil {
		slog.Error("could not get alerts", "err", err)
		w.WriteHeader(500)
		return
	}

	if isBrowser(req) {
		data, err := r.templating.Render(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func (r *Receiver) GetAlertsToday(w http.ResponseWriter, req *http.Request, params GetAlertsTodayParams) {
	alerts, err := r.db.GetAlertsToday(req.Context())
	if err != nil {
		slog.Error("could not get alerts", "err", err)
		w.WriteHeader(500)
		return
	}

	if isBrowser(req) {
		data, err := r.templating.Render(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func (r *Receiver) GetAlertsYesterday(w http.ResponseWriter, req *http.Request, params GetAlertsYesterdayParams) {
	alerts, err := r.db.GetAlertsYesterday(req.Context())
	if err != nil {
		slog.Error("could not get alerts", "err", err)
		w.WriteHeader(500)
		return
	}

	if isBrowser(req) {
		data, err := r.templating.Render(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func (r *Receiver) CreateAlert(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		slog.Error("could not read body", "err", err)
		w.WriteHeader(500)
		return
	}
	defer func() {
		_ = req.Body.Close()
	}()

	var root webhook.Root
	if err := json.Unmarshal(body, &root); err != nil {
		slog.Error("could not unmarshal alert", "err", err)
		w.WriteHeader(500)
		return
	}

	var errs error
	for _, alert := range root.Alerts {
		converted := model.Alert{
			Annotations: alert.Annotations,
			EndsAt:      alert.EndsAt,
			Fingerprint: alert.Fingerprint,
			StartsAt:    alert.StartsAt,
			Status: model.Status{
				InhibitedBy: nil,
				SilencedBy:  nil,
				State:       alert.Status,
			},
			GeneratorURL: alert.GeneratorURL,
			Labels:       alert.Labels,
		}

		if err := r.db.SaveAlert(req.Context(), converted); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if errs != nil {
		slog.Error("could not save alert", "err", errs)
		w.WriteHeader(500)
		return
	}
}

func (r *Receiver) GetAllAlerts(w http.ResponseWriter, req *http.Request, params GetAllAlertsParams) {
	alerts, err := r.db.GetAlerts(req.Context())
	if err != nil {
		slog.Error("could not get alerts", "err", err)
		w.WriteHeader(500)
		return
	}

	if isBrowser(req) {
		data, err := r.templating.Render(alerts)
		if err != nil {
			slog.Error("could render alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	} else {
		data, err := json.Marshal(alerts)
		if err != nil {
			slog.Error("could not marshal alerts", "err", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func NewReceiver(db Db, templating Templating) (*Receiver, error) {
	if db == nil {
		return nil, errors.New("nil db passed")
	}

	if templating == nil {
		return nil, errors.New("nil templating passed")
	}

	// TODO
	stats, _ := NewStatsRenderer()
	return &Receiver{
		db:          db,
		statsSource: db.(Stats),
		templating:  templating,
		stats:       stats,
	}, nil
}

func isBrowser(req *http.Request) bool {
	acceptHeader := req.Header.Get("Accept")
	return strings.Contains(acceptHeader, "text/html") || strings.Contains(acceptHeader, "application/xhtml+xml")
}
