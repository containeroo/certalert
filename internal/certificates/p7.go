package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mozilla.org/pkcs7"
)

// ExtractP7CertificatesInfo extracts certificate information from a P7 file
func ExtractP7CertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Parse all PEM blocks and filter by type
	for {
		block, rest := pem.Decode(certificateData)
		if block == nil {
			break
		}

		certificateData = block.Bytes
		switch block.Type {
		case "PKCS7":
			p7, err := pkcs7.Parse(certificateData)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "p7", fmt.Sprintf("Failed to parse P7B file '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
			}

			for _, certificate := range p7.Certificates {
				subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

				certificateInfo := CertificateInfo{
					Name:    cert.Name,
					Subject: subject,
					Epoch:   certificate.NotAfter.Unix(),
					Type:    "p7",
				}
				certificateInfoList = append(certificateInfoList, certificateInfo)

				log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
			}
		case "CERTIFICATE":
			certificate, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "p7", fmt.Sprintf("Failed to parse certificate '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
			}

			subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

			certificateInfo := CertificateInfo{
				Name:    cert.Name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "p7",
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
