package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystemReconciler = "reconciler"

var (
	ReconciliationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemReconciler,
		Name:      "reconciliations_total",
		Help:      "The amount of times the reconciler ran",
	})

	ReconciliationsErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemReconciler,
		Name:      "reconciliations_errors_total",
		Help:      "The amount of times the reconciler encountered errors",
	})

	ReconciledAlertsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystemReconciler,
		Name:      "reconciled_items_total",
		Help:      "Total amount of reconciled items",
	})
)
