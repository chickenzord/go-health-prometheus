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

type HealthPrometheus struct {
	availabilityGauge *prometheus.GaugeVec
}

// New
// create new instance of HealthPrometheus with metricsName
func New(metricName string) *HealthPrometheus {
	return &HealthPrometheus{
		availabilityGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metricName,
				Help: "Availability status of service identified by name label",
			},
			[]string{"name", "status"},
		),
	}
}

// Interceptor
// implements health.Interceptor that will update underlying metrics as availability status changes
func (m *HealthPrometheus) Interceptor(next health.InterceptorFunc) health.InterceptorFunc {
	return func(ctx context.Context, name string, state health.CheckState) health.CheckState {
		for _, s := range healthStatus {
			mm := m.availabilityGauge.With(prometheus.Labels{
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
// return Prometheus collectors containing health metrics
func (m *HealthPrometheus) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.availabilityGauge,
	}
}
