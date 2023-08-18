package certificates

import (
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mozilla.org/pkcs7"
)

// ExtractP7CertificatesInfo extracts certificate information from a P7 file
func ExtractP7CertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	handleError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warningf("Failed to extract certificate information: %v", errMsg)
		return nil
	}

	// Parse all PEM blocks and filter by type
	for {
		block, remainingData := pem.Decode(certData)
		if block == nil {
			break
		}

		switch block.Type {
		case "PKCS7":
			certData = block.Bytes
			// Parse the P7B data
			p7, err := pkcs7.Parse(certData)
			if err != nil {
				return certInfoList, handleError(fmt.Sprintf("Failed to parse P7B file '%s': %v", name, err))
			}

			// Extract certificates
			for _, cert := range p7.Certificates {
				subject := cert.Subject.CommonName
				certInfo := CertificateInfo{
					Name:    name,
					Subject: subject,
					Epoch:   cert.NotAfter.Unix(),
					Type:    "p7b",
				}
				log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())
				certInfoList = append(certInfoList, certInfo)
			}
		default:
			log.Warningf("Unknown PEM block type '%s' in P7B file", block.Type)
		}
		certData = remainingData
	}

	if len(certInfoList) == 0 {
		if err := handleError(fmt.Sprintf("Failed to decode any certificate in '%s'", name)); err != nil {
			return certInfoList, err
		}
	}

	return certInfoList, nil
}
