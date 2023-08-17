package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetrics Metrics

var (
	CertificateEpoch = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certalert_certificate_epoch_seconds",
			Help: "The epoch of the certificate",
		},
		[]string{"instance", "subject", "type"},
	)
)

type Metrics struct {
	Registry *prometheus.Registry
}

// NewMetrics registers all prometheus metrics
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.Register(CertificateEpoch) // Register the global metric

	return &Metrics{
		Registry: reg,
	}
}
