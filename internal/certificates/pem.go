package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ExtractPEMCertificatesInfo extracts certificate information from a P7 file
func ExtractPEMCertificatesInfo(cert Certificate, certData []byte, failOnError bool) ([]CertificateInfo, error) {
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
			c, err := x509.ParseCertificate(certData)
			if err != nil {
				if err := handleFailOnError(&certInfoList, cert.Name, "pem", fmt.Sprintf("Failed to parse certificate '%s': %v", cert.Name, err), failOnError); err != nil {
					return certInfoList, err
				}
			}

			subject := c.Subject.CommonName
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certInfoList)+1)
			}
			certInfo := CertificateInfo{
				Name:    cert.Name,
				Subject: subject,
				Epoch:   c.NotAfter.Unix(),
				Type:    "pem",
			}
			certInfoList = append(certInfoList, certInfo)

			log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())
		default:
			log.Debugf("Skip PEM block of type '%s'", block.Type)
		}

		certData = rest // Move to the next PEM block
	}

	if len(certInfoList) == 0 {
		return certInfoList, handleFailOnError(&certInfoList, cert.Name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", cert.Name), failOnError)
	}

	return certInfoList, nil
}
