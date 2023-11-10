package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCertificateType("pem", ExtractPEMCertificatesInfo, "pem", "crt")
}

// ExtractPEMCertificatesInfo extracts certificate information from a P7 file.
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

			log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
		default:
			log.Debugf("Skip PEM block of type '%s'", block.Type)
		}

		certificateData = rest // Move to the next PEM block
		continue
	}

	if len(certificateInfoList) == 0 {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", cert.Name), failOnError)
	}

	return certificateInfoList, nil
}
