package healthmetrics

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

type HealthMetrics struct {
	metrics *prometheus.GaugeVec
}

// New
// create new instance of HealthMetrics with metricsName
func New(metricsName string) *HealthMetrics {
	return &HealthMetrics{
		metrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metricsName,
				Help: "Availability status of service identified in name label",
			},
			[]string{"name", "status"},
		),
	}
}

// HealthInterceptor
// implements health.HealthInterceptor that will update underlying gauge metrics as availability status changes
func (m *HealthMetrics) HealthInterceptor(next health.InterceptorFunc) health.InterceptorFunc {
	return func(ctx context.Context, name string, state health.CheckState) health.CheckState {
		for _, s := range healthStatus {
			mm := m.metrics.With(prometheus.Labels{
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

// PrometheusCollector
// return Prometheus collector containing gauge health metrics
func (m *HealthMetrics) PrometheusCollector() prometheus.Collector {
	return m.metrics
}
