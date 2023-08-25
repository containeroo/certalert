package certificates

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// GetByName returns the certificate with the given name
func GetCertificateByName(name string, certificates []Certificate) (*Certificate, error) {
	for _, cert := range certificates {
		if cert.Name == name {
			return &cert, nil
		}
	}
	return nil, fmt.Errorf("Certificate '%s' not found", name)
}

// Process extracts certificate information from the certificates and updates the Prometheus metrics
func Process(certificates []Certificate, failOnError bool) (certificatesInfo []CertificateInfo, err error) {
	var certInfoList []CertificateInfo

	for _, cert := range certificates {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debugf("Skip certificate '%s' as it is disabled", cert.Name)
			continue
		}

		log.Debugf("Processing certificate '%s'", cert.Name)

		certData, err := os.ReadFile(cert.Path)
		if err != nil {
			// Accessibility of the file is checked in the config validation, if reached
			// here, the file exists but can't be read for some reason.
			if err := handleError(&certInfoList, cert.Name, cert.Type, fmt.Sprintf("Failed to read certificate file '%s'. %v", cert.Path, err), failOnError); err != nil {
				return nil, err
			}
			continue
		}

		extractFunc, found := TypeToExtractionFunction[cert.Type]
		if !found {
			// This should never happen as the config validation ensures that the type is valid
			if err := handleError(&certInfoList, cert.Name, cert.Type, fmt.Sprintf("Unknown certificate type '%s'", cert.Type), failOnError); err != nil {
				return nil, err
			}
			continue
		}

		certs, err := extractFunc(cert.Name, certData, cert.Password, failOnError)
		if err != nil {
			// err is only returned if failOnError is true
			return nil, fmt.Errorf("Error extracting certificate information: %v", err)
		}

		certInfoList = append(certInfoList, certs...)
	}

	return certInfoList, nil
}

// certExistsInSlice checks if a certificate exists in a slice
func certExistsInSlice(cert CertificateInfo, slice []CertificateInfo) bool {
	for _, c := range slice {
		if cert.Name == c.Name && cert.Subject == c.Subject && cert.Type == c.Type {
			return true
		}
	}
	return false
}

// handleError is a helper function to handle failOnError
func handleError(certInfoList *[]CertificateInfo, certName, certType, errMsg string, failOnError bool) error {
	if failOnError {
		return fmt.Errorf(errMsg)
	}
	log.Warningf("Failed to extract certificate information: %v", errMsg)
	*certInfoList = append(*certInfoList, CertificateInfo{
		Name:  certName,
		Type:  certType,
		Error: errMsg,
	})
	return nil
}
