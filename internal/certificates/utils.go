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
	for _, cert := range certificates {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debugf("Skip certificate '%s' as it is disabled", cert.Name)
			continue
		}
		if cert.Valid != nil && !*cert.Valid {
			log.Debugf("Skip certificate '%s' as it is not valid", cert.Name)
			continue
		}

		log.Debugf("Processing certificate '%s'", cert.Name)

		var certificateInfo []CertificateInfo
		certData, err := os.ReadFile(cert.Path)
		if err != nil {
			return nil, fmt.Errorf("Failed to read certificate file '%s': %w", cert.Path, err)
		}

		extractFunc, found := TypeToExtractionFunction[cert.Type]
		if !found {
			return nil, fmt.Errorf("Unknown certificate type '%s'", cert.Type)
		}

		certificateInfo, err = extractFunc(cert.Name, certData, cert.Password, failOnError)
		if err != nil {
			// err is only returned if failOnError is true
			return nil, fmt.Errorf("Error extracting certificate information: %v", err)
		}

		certificatesInfo = append(certificatesInfo, certificateInfo...)
	}

	return certificatesInfo, nil
}
