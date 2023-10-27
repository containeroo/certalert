package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ExtractPEMCertificatesInfo extracts certificate information from a P7 file
func ExtractPEMCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Parse all PEM blocks and filter by type
	for {
		block, rest := pem.Decode(certData)
		if block == nil {
			break
		}

		switch block.Type {
		case "CERTIFICATE":
			certificate, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, name, "pem", fmt.Sprintf("Failed to parse certificate '%s': %v", name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
			}

			subject := certificate.Subject.ToRDNSequence().String()
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certificateInfoList)+1)
			}
			certificateInfo := CertificateInfo{
				Name:    name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "pem",
			}
			certificateInfoList = append(certificateInfoList, certificateInfo)

			log.Debugf("Certificate '%s' expires on %s", certificateInfo.Subject, certificateInfo.ExpiryAsTime())

			certData = rest // Move to the next PEM block
		default:
			log.Debugf("Skip PEM block of type '%s'", block.Type)
		}
		certData = rest
		continue
	}

	if len(certificateInfoList) == 0 {
		return certificateInfoList, handleFailOnError(&certificateInfoList, name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", name), failOnError)
	}

	return certificateInfoList, nil
}
