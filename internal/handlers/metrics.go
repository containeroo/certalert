package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/metrics"
	"certalert/internal/server"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	server.Register("/metrics", "Delivers metrics for Prometheus to scrape", Metrics, "GET", "POST")
}

// setMetricsForCertificateInfo sets metrics for a given certificate info.
//
// It takes a CertificateInfo object and sets metrics in Prometheus for the
// certificate extraction status, epoch, and error reason (if any).
//
// Parameters:
//   - ci: certificates.CertificateInfo
//     The CertificateInfo object for which metrics should be set.
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

// Metrics is an HTTP handler for the /metrics route.
//
// This handler returns the metrics for Prometheus to scrape. It processes the
// configured certificates and sets metrics based on the extraction status, epoch,
// and error reason (if any).
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
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
