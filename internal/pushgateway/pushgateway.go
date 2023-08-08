package pushgateway

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"fmt"
)

// Send pushes the certificate information to the Pushgateway
func Send(address string, jobName string, auth config.Auth, certs []certificates.Certificate, insecureSkipVerify bool, failOnError bool) error {
	pusher := createPusher(address, jobName, auth, insecureSkipVerify)

	certificatesInfo, err := certificates.Process(certs, failOnError)
	if err != nil {
		// err is only returned if failOnError is true
		return fmt.Errorf("Failed to process certificates: %w", err)
	}

	for _, certificateInfo := range certificatesInfo {
		if err := pushToGateway(pusher, certificateInfo); err != nil {
			return fmt.Errorf("Failed to push certificate info to gateway: %w", err)
		}
	}
	return nil
}
