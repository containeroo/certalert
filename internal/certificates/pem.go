package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ExtractPEMCertificatesInfo extracts certificate information from a P7 file
func ExtractPEMCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// Parse all PEM blocks and filter by type
	for {
		block, rest := pem.Decode(certData)
		if block == nil {
			break
		}

		certData = block.Bytes
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				if err := handleError(&certInfoList, name, "pem", fmt.Sprintf("Failed to parse certificate '%s': %v", name, err), failOnError); err != nil {
					return certInfoList, err
				}
			}

			subject := cert.Subject.CommonName
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certInfoList)+1)
			}
			certInfo := CertificateInfo{
				Name:    name,
				Subject: subject,
				Epoch:   cert.NotAfter.Unix(),
				Type:    "pem",
			}
			certInfoList = append(certInfoList, certInfo)

			log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())

			certData = rest // Move to the next PEM block
		default:
			log.Debug("Skip PEM block of type '%s'", block.Type)
		}
		certData = rest
		continue
	}

	if len(certInfoList) == 0 {
		return certInfoList, handleError(&certInfoList, name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", name), failOnError)
	}

	return certInfoList, nil
}
