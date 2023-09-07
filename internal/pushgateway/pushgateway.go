package pushgateway

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Send pushes the certificate information to the Pushgateway
func Send(address string, jobName string, auth config.Auth, certs []certificates.Certificate, insecureSkipVerify bool, failOnError bool) error {
	if address == "" {
		return fmt.Errorf("Pushgateway address is empty")
	}

	if !utils.IsValidURL(address) {
		return fmt.Errorf("Invalid pushgateway address '%s'", address)
	}

	pusher := createPusher(address, jobName, auth, insecureSkipVerify)

	certificatesInfo, err := certificates.Process(certs, failOnError)
	if err != nil {
		return fmt.Errorf("Failed to process certificates: %w", err)
	}

	for _, certificateInfo := range certificatesInfo {
		if err := pushToGateway(pusher, certificateInfo); err != nil {
			return fmt.Errorf("Failed to push certificate info to gateway: %w", err)
		}
		log.Debugf("Pushed certificate '%s' (%s, %s, %s) expiration epoch '%d'", certificateInfo.Name, certificateInfo.Name, certificateInfo.Type, certificateInfo.Subject, certificateInfo.Epoch)
	}

	return nil
}
