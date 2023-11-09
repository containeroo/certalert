package certificates

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

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
			if err := handleFailOnError(&certInfoList, cert.Name, cert.Type, fmt.Sprintf("Failed to read certificate file '%s'. %v", cert.Path, err), failOnError); err != nil {
				return nil, err
			}
			continue
		}

		// If user specify the type, we need to convert it to the canonical type
		inferredType, found := FileExtensionsToType[cert.Type]
		if !found {
			// This should never happen as the config validation ensures that the type is valid
			if err := handleFailOnError(&certInfoList, cert.Name, cert.Type, fmt.Sprintf("Unknown certificate type '%s'", cert.Type), failOnError); err != nil {
				return nil, err
			}
			continue
		}

		extractFunc, found := ExtractionFunctionFabric[inferredType]
		if !found {
			// This should never happen as the config validation ensures that the type is valid
			if err := handleFailOnError(&certInfoList, cert.Name, cert.Type, fmt.Sprintf("Unknown certificate type '%s'", cert.Type), failOnError); err != nil {
				return nil, err
			}
			continue
		}

		certs, err := extractFunc(cert, certData, failOnError)
		if err != nil {
			// err is only returned if failOnError is true
			return nil, fmt.Errorf("Error extracting certificate information: %v", err)
		}

		certInfoList = append(certInfoList, certs...)
	}

	return certInfoList, nil
}
