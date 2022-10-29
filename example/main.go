package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexliesenfeld/health"
	healthprometheus "github.com/chickenzord/go-health-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	healthProm := healthprometheus.New("myapp_health")

	// Setup health checker
	healthChecker := health.NewChecker(
		health.WithCheck(health.Check{
			Name: "database",
			Check: func(ctx context.Context) error {
				// always up
				return nil
			},
		}),
		health.WithCheck(health.Check{
			Name: "redis",
			Check: func(ctx context.Context) error {
				// always down
				return fmt.Errorf("connection error")
			},
		}),
		health.WithInterceptors(healthProm.Interceptor), // Use the interceptor to record health metrics
	)

	// Setup Prometheus
	registry := prometheus.NewRegistry()
	registry.MustRegister(healthProm.Collectors()...) // Register the health metric collectors

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle("/health", health.NewHandler(healthChecker))
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	fmt.Println("Listening on :9000")
	if err := http.ListenAndServe(":9000", mux); err != nil {
		panic(err)
	}
}
