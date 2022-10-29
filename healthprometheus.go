package healthprometheus

import (
	"context"

	"github.com/alexliesenfeld/health"
	"github.com/prometheus/client_golang/prometheus"
)

var healthStatus = []health.AvailabilityStatus{
	health.StatusDown,
	health.StatusUnknown,
	health.StatusUp,
}

// MetricNameOpts
// options struct for customizing metric names
type MetricNameOpts struct {
	// Name for availability status gauge
	AvailabilityStatusGauge string
}

// DefaultMetricNameOpts
// MetricNameOpts with default values
var DefaultMetricNameOpts = MetricNameOpts{
	AvailabilityStatusGauge: "health",
}

type HealthPrometheus struct {
	availabilityStatusGauge *prometheus.GaugeVec
}

// New create new instance of HealthPrometheus
func New(namespace, subsystem string, nameOpts MetricNameOpts) *HealthPrometheus {
	return &HealthPrometheus{
		availabilityStatusGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      nameOpts.AvailabilityStatusGauge,
				Help:      "Availability status of service identified by name label",
			},
			[]string{"name", "status"},
		),
	}
}

// New create new instance of HealthPrometheus with default metric names
func NewDefault(namespace, subsystem string) *HealthPrometheus {
	return New(namespace, subsystem, DefaultMetricNameOpts)
}

// Interceptor implements health.Interceptor
//
// Will update underlying metrics as availability status changes
func (m *HealthPrometheus) Interceptor(next health.InterceptorFunc) health.InterceptorFunc {
	return func(ctx context.Context, name string, state health.CheckState) health.CheckState {
		for _, s := range healthStatus {
			mm := m.availabilityStatusGauge.With(prometheus.Labels{
				"name":   name,
				"status": string(s),
			})

			if s == state.Status {
				mm.Set(1)
			} else {
				mm.Set(0)
			}
		}

		return next(ctx, name, state)
	}
}

// Collectors
// return all Prometheus collectors
func (m *HealthPrometheus) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.availabilityStatusGauge,
	}
}

func (m *HealthPrometheus) AvailabilityStatusCollector() prometheus.Collector {
	return m.availabilityStatusGauge
}
