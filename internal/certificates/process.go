package certificates

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// Process processes a list of certificates and extracts certificate information.
//
// This function takes a slice of Certificate structs, indicating the certificates to process,
// and a flag indicating whether to fail on error. It returns a slice of CertificateInfo containing
// information about each certificate.
//
// The function iterates through each certificate, checking for disabled status and logging
// processing details. It reads the raw certificate data from the specified file, infers the type
// if not explicitly specified, and calls the corresponding extraction function. The extracted
// certificate information is then added to the result list.
//
// Parameters:
//   - certificates: []Certificate
//     A slice of Certificate structs representing the certificates to process.
//   - failOnError: bool
//     A flag indicating whether to fail immediately on encountering an error.
//
// Returns:
//   - []CertificateInfo
//     A slice of CertificateInfo structs containing information about each processed certificate.
//   - error
//     An error, if any, encountered during the processing. If failOnError is false, the function may
//     return a non-nil error along with the partial list of CertificateInfo.
func Process(certificates []Certificate, failOnError bool) (certificatesInfo []CertificateInfo, err error) {
	var certInfoList []CertificateInfo

	for _, cert := range certificates {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debug().Msgf("Skip certificate '%s' as it is disabled", cert.Name)
			continue
		}

		log.Debug().Msgf("Processing certificate '%s'", cert.Name)

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

		extractFunc, found := TypeToExtractionFunction[inferredType]
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
