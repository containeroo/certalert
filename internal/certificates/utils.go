package certificates

import (
	"fmt"

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

// certExistsInSlice checks if a certificate exists in a slice by comparing the name, subject and type.
// The certificate epoch is not compared!
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
