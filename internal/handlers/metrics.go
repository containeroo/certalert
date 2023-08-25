package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetMetricsForCertificateInfo sets metrics for a given certificate info.
func setMetricsForCertificateInfo(ci certificates.CertificateInfo) {
	labels := prometheus.Labels{
		"instance": ci.Name,
		"subject":  ci.Subject,
		"type":     ci.Type,
		"reason":   "none", // default value
	}

	if ci.Error != "" {
		// Add the reason only for CertificateExtractionStatus when there's an error
		labels["reason"] = ci.Error
		metrics.CertificateExtractionStatus.With(labels).Set(1)
	} else {
		// Set without reason label
		metrics.CertificateExtractionStatus.With(labels).Set(0)
		metrics.CertificateEpoch.With(labels).Set(float64(ci.Epoch))
	}
}

// Metrics is the handler for the /metrics route
// It returns the metrics for Prometheus to scrape
func Metrics(w http.ResponseWriter, r *http.Request) {
	certificateInfos, err := certificates.Process(config.App.Certs, config.App.FailOnError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ci := range certificateInfos {
		setMetricsForCertificateInfo(ci)
	}

	// Serve metrics
	promhttp.HandlerFor(metrics.PromMetrics.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
