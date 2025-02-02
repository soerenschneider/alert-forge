package alertmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"

	"github.com/soerenschneider/alert-forge/internal"
	"github.com/soerenschneider/alert-forge/internal/model"
	"go.uber.org/multierr"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AlertmanagerOpt func(client *AlertmanagerClient) error

func WithClient(client HttpClient) AlertmanagerOpt {
	return func(am *AlertmanagerClient) error {
		if client == nil {
			return errors.New("empty client")
		}
		am.client = client
		return nil
	}
}

func WithBlacklist(labelsBlacklist map[string]string) AlertmanagerOpt {
	return func(a *AlertmanagerClient) error {
		a.blacklistedLabels = labelsBlacklist
		return nil
	}
}

func NewAlertmanagerClient(instances []string, opts ...AlertmanagerOpt) (*AlertmanagerClient, error) {
	ret := &AlertmanagerClient{
		client: &http.Client{
			Timeout: 5,
		},
		instances: instances,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(ret); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return ret, errs
}

type AlertmanagerClient struct {
	client            HttpClient
	instances         []string
	blacklistedLabels map[string]string
}

// FetchActiveAlerts queries Alertmanager for active alerts and returns them.
func (c *AlertmanagerClient) fetch(ctx context.Context, url string) ([]*model.Alert, error) {
	url = url + "/api/v2/alerts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alerts: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		return nil, errors.New(errMsg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var alerts []*model.Alert
	if err := json.Unmarshal(body, &alerts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return alerts, nil
}

func (c *AlertmanagerClient) iterateOverInstances(ctx context.Context) ([]*model.Alert, error) {
	for _, instance := range c.instances {
		alerts, err := c.fetch(ctx, instance)
		if err != nil {
			slog.Warn("failed to get alerts from alertmanager", "instance", instance)
		} else {
			return alerts, nil
		}
	}

	return nil, errors.New("no instances available to fetch alerts from")
}

func (c *AlertmanagerClient) GetActiveAlerts(ctx context.Context) ([]*model.Alert, error) {
	alerts, err := c.iterateOverInstances(ctx)
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		alert.Fingerprint = alert.Hash()
	}

	alerts = c.Filter(alerts)

	slices.SortFunc(alerts, func(a, b *model.Alert) int {
		if a.Labels["severity"] != b.Labels["severity"] {
			return int(internal.DefaultSeverityMapping[b.Labels["severity"]] - internal.DefaultSeverityMapping[a.Labels["severity"]])
		}

		if a.StartsAt != b.StartsAt {
			if a.StartsAt.Before(b.StartsAt) {
				return -1
			}
			return 1
		}
		return 0
	})

	return alerts, nil
}

func (c *AlertmanagerClient) Filter(alerts []*model.Alert) []*model.Alert {
	ret := make([]*model.Alert, 0, len(alerts))

	for _, alert := range alerts {
		filter := false
		for label, value := range c.blacklistedLabels {
			if alert.Labels[label] == value {
				filter = true
				continue
			}
		}

		if !filter {
			ret = append(ret, alert)
		}
	}

	return ret
}
