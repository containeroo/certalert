package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetrics Metrics

const (
	CertalertMetricName = "certalert_certificate_epoch_seconds"
	CertalertMetricHelp = "The epoch of the certificate"
)

// Metrics holds all prometheus metrics and the custom registry
type Metrics struct {
	Registry         *prometheus.Registry
	CertificateEpoch *prometheus.GaugeVec
}

// NewMetrics registers all prometheus metrics
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry() // Create a new registry
	m := &Metrics{
		Registry: reg, // Store the registry in the Metrics struct
		CertificateEpoch: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: CertalertMetricName,
				Help: CertalertMetricHelp,
			},
			[]string{"instance", "subject", "type"},
		),
	}
	reg.Register(m.CertificateEpoch)
	return m
}
