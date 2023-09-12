package pushgateway

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/metrics"
	"certalert/internal/utils"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/push"
)

// createPusher creates a configured pusher
func createPusher(address, job string, auth config.Auth, insecureSkipVerify bool) *push.Pusher {
	var httpClient *http.Client
	if insecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	} else {
		httpClient = &http.Client{}
	}

	pusher := push.New(address, job).
		Collector(metrics.CertificateEpoch).
		Client(httpClient)

	if utils.HasStructField(auth, "Bearer.Token") && auth.Bearer.Token != "" {
		pusher = pusher.BasicAuth("Bearer", auth.Bearer.Token)
	} else if utils.HasStructField(auth, "Basic.Username") && auth.Basic.Username != "" {
		pusher = pusher.BasicAuth(auth.Basic.Username, auth.Basic.Password)
	}

	return pusher
}

// pushToGateway pushes the certificate information to the pushgateway
func pushToGateway(pusher *push.Pusher, cert certificates.CertificateInfo) error {
	gauge := metrics.CertificateEpoch.WithLabelValues(cert.Name, cert.Type, cert.Subject)
	gauge.Set(float64(cert.Epoch))

	if err := pusher.Push(); err != nil {
		return fmt.Errorf("Could not push to Pushgateway: %w", err)
	}

	return nil
}
