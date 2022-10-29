# go-health-prometheus
Go library for integrating [alexliesenfeld/health](https://github.com/alexliesenfeld/health) with Prometheus. 
Implements both Health Interceptor and Prometheus Collector.

[![Go Reference](https://pkg.go.dev/badge/github.com/chickenzord/go-health-prometheus.svg)](https://pkg.go.dev/github.com/chickenzord/go-health-prometheus)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/go-health-prometheus)](https://goreportcard.com/report/github.com/chickenzord/go-health-prometheus)

## Example usage

```go
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
		health.WithInterceptors(healthProm.Interceptor), // Use the interceptor to record health metrics
		// ... checks omitted for brevity
	)

	// Setup Prometheus
	registry := prometheus.NewRegistry()
	registry.MustRegister(healthProm.Collectors()...) // Register the health metric collectors
	// ... you can register another collectors here (e.g. Go process collector) 

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle("/health", health.NewHandler(healthChecker))
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	fmt.Println("Listening on :9000")
	if err := http.ListenAndServe(":9000", mux); err != nil {
		panic(err)
	}
}
```

See `example` folder for more info on how to use this library.

## Example metrics

```
myapp_health{name="database" status="up"} 1
myapp_health{name="database" status="down"} 0
myapp_health{name="database" status="unknown"} 0

myapp_health{name="redis" status="up"} 0
myapp_health{name="redis" status="down"} 1
myapp_health{name="redis" status="unknown"} 0
```
