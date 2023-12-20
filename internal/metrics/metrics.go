package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetrics Metrics

var (
	// New metric to track certificate expiration date as epoch
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

// NewMetrics creates a new instance of the Metrics struct, initializing a Prometheus registry,
// and registering global metrics like CertificateEpoch and CertificateExtractionStatus.
//
// Returns:
//   - *Metrics
//     A pointer to the newly created Metrics instance.
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.Register(CertificateEpoch)            // Register the global metric
	reg.Register(CertificateExtractionStatus) // Register the new metric

	return &Metrics{
		Registry: reg,
	}
}
