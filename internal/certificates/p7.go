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
				return certInfoList, handleError(&certInfoList, name, "p7", fmt.Sprintf("Failed to parse P7B file '%s': %v", name, err), failOnError)
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
		return certInfoList, handleError(&certInfoList, name, "p7", fmt.Sprintf("Failed to decode any certificate in '%s'", name), failOnError)
	}

	return certInfoList, nil
}
