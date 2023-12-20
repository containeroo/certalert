package pushgateway

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/utils"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Send pushes certificate information to the specified Pushgateway.
//
// Parameters:
//   - address: string
//     The address of the Pushgateway.
//   - jobName: string
//     The job label to associate with the pushed metrics.
//   - auth: config.Auth
//     The authentication configuration for communicating with the Pushgateway.
//   - certs: []certificates.Certificate
//     The list of certificates to process and push to the Pushgateway.
//   - insecureSkipVerify: bool
//     Whether to skip TLS certificate verification when communicating with the Pushgateway.
//   - failOnError: bool
//     Whether to fail on processing errors for individual certificates.
//
// Returns:
//   - error
//     An error if the push to the Pushgateway fails or if there are errors processing individual certificates.
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
		log.Debug().Msgf("Pushed certificate '%s' (%s, %s, %s) expiration epoch '%d'", certificateInfo.Name, certificateInfo.Name, certificateInfo.Type, certificateInfo.Subject, certificateInfo.Epoch)
	}

	return nil
}
