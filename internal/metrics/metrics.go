package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetrics Metrics

var (
	CertificateEpoch = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certalert_certificate_epoch_seconds",
			Help: "The expiration date of the certificate as a epoch",
		},
		[]string{"instance", "subject", "type", "reason"},
	)

	// New metric to track failed certificate extractions
	CertificateExtractionStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certalert_certificate_extraction_status",
			Help: "Status of certificate extraction (0=success, 1=failure)",
		},
		[]string{"instance", "subject", "type", "reason"},
	)
)

type Metrics struct {
	Registry *prometheus.Registry
}

func init() {
	PromMetrics = *NewMetrics()
}

// NewMetrics registers all prometheus metrics
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.Register(CertificateEpoch)            // Register the global metric
	reg.Register(CertificateExtractionStatus) // Register the new metric

	return &Metrics{
		Registry: reg,
	}
}
