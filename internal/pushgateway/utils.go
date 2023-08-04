package pushgateway

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/metrics"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// createPusher creates a new pusher with the necessary configuration
func createPusher(address, job string, auth config.Auth) *push.Pusher {
	certificateExpirationEpoch := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.CertalertMetricName,
			Help: metrics.CertalertMetricHelp,
		},
		[]string{},
	)

	pusher := push.New(address, job).
		Collector(certificateExpirationEpoch)

	if auth.Bearer.Token != "" {
		pusher = pusher.BasicAuth("Bearer", auth.Bearer.Token)
	} else if auth.Basic.Username != "" {
		pusher = pusher.BasicAuth(auth.Basic.Username, auth.Basic.Password)
	}

	return pusher
}

// pushToGateway pushes the certificate information to the Pushgateway
func pushToGateway(pusher *push.Pusher, cert certificates.CertificateInfo) error {
	pusher = pusher.
		Grouping("instance", cert.Name).
		Grouping("type", cert.Type).
		Grouping("subject", cert.Subject)

	if err := pusher.Push(); err != nil {
		return fmt.Errorf("Could not push to Pushgateway: %w", err)
	}

	return nil
}
