package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/rs/zerolog/log"
)

func init() {
	registerCertificateType("pem", ExtractPEMCertificatesInfo, "pem", "crt")
}

// ExtractPEMCertificatesInfo extracts certificate information from a PEM-encoded file.
//
// This function takes a Certificate struct, the raw certificate data as a byte slice, and a
// flag indicating whether to fail on error. It returns a slice of CertificateInfo containing
// information about each certificate found in the PEM file.
//
// The function parses all PEM blocks from the input certificateData, filters by type ("CERTIFICATE"),
// and extracts certificate information. It logs information about each certificate, including its
// subject, expiration time, and type.
//
// Parameters:
//   - cert: Certificate
//     A Certificate struct representing the PEM file, including its name and other details.
//   - certificateData: []byte
//     The raw binary data of the PEM file.
//   - failOnError: bool
//     A flag indicating whether to fail immediately on encountering an error.
//
// Returns:
//   - []CertificateInfo
//     A slice of CertificateInfo structs containing information about each certificate in the PEM file.
//   - error
//     An error, if any, encountered during the extraction process. If failOnError is false, the
//     function may return a non-nil error along with the partial list of CertificateInfo.
func ExtractPEMCertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Parse all PEM blocks and filter by type
	for {
		block, rest := pem.Decode(certificateData)
		if block == nil {
			break
		}

		switch block.Type {
		case "CERTIFICATE":
			certificate, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "pem", fmt.Sprintf("Failed to parse certificate '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
			}

			subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

			certificateInfo := CertificateInfo{
				Name:    cert.Name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "pem",
			}
			certificateInfoList = append(certificateInfoList, certificateInfo)

			log.Debug().Msgf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
		default:
			log.Debug().Msgf("Skip PEM block of type '%s'", block.Type)
		}

		certificateData = rest // Move to the next PEM block
		continue
	}

	if len(certificateInfoList) == 0 {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", cert.Name), failOnError)
	}

	return certificateInfoList, nil
}
