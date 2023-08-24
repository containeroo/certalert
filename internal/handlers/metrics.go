package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics is the handler for the /metrics route
// It returns the metrics for Prometheus to scrape
func Metrics(w http.ResponseWriter, r *http.Request) {
	certificateInfos, err := certificates.Process(config.App.Certs, config.App.FailOnError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ci := range certificateInfos {
		if ci.Error != "" {
			continue
		}
		metrics.CertificateEpoch.With(
			prometheus.Labels{
				"instance": ci.Name,
				"subject":  ci.Subject,
				"type":     ci.Type,
			},
		).Set(float64(ci.Epoch))
	}

	// Serve metrics
	promhttp.HandlerFor(metrics.PromMetrics.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
