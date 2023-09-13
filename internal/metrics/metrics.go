package metrics

import (
	"log"

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

// Metrics represents the prometheus metrics
type Metrics struct {
	Registry *prometheus.Registry
}

// Init initializes the prometheus metrics
func init() {
	PromMetrics = *NewMetrics()
}

// NewMetrics registers all prometheus metrics
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()

	if err := reg.Register(CertificateEpoch); err != nil {
		log.Fatalf("Unable to register CertificateEpoch metric: %s", err)
	}
	if err := reg.Register(CertificateExtractionStatus); err != nil {
		log.Fatalf("Unable to register CertificateExtractionStatus metric: %s", err)
	}

	return &Metrics{
		Registry: reg,
	}
}
