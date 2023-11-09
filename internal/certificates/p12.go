package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo extracts certificate information from a P12 file
func ExtractP12CertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Decode the P12 data
	_, certificate, caCerts, err := pkcs12.DecodeChain(certificateData, cert.Password)
	if err != nil {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "p12", fmt.Sprintf("Failed to decode P12 file '%s': %v", cert.Name, err), failOnError)
	}

	// Prepare for extraction
	certificates := append(caCerts, certificate)

	// Extract certificates
	for _, certificate := range certificates {
		subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

		certificateInfo := CertificateInfo{
			Name:    cert.Name,
			Subject: subject,
			Epoch:   certificate.NotAfter.Unix(),
			Type:    "p12",
		}
		certificateInfoList = append(certificateInfoList, certificateInfo)

		log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
	}

	return certificateInfoList, nil
}
