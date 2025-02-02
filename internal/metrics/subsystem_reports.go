package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystemVaultRenewal = "reports"

var (
	ReportsAttempted = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemVaultRenewal,
		Name:      "requests_total",
		Help:      "The amount of times a report was requested",
	})

	ReportsErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemVaultRenewal,
		Name:      "errors_total",
		Help:      "Total errors while trying to craft a report",
	})
)
