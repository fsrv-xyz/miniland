package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func serveMetricsHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func ServeMetrics(addr string) {
	err := http.ListenAndServe(addr, serveMetricsHandler())
	if err != nil {
		return
	}
}
