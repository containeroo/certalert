package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo extracts certificate information from a P12 file
func ExtractTrustStoreCertificatesInfo(cert Certificate, certData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// Decode the P12 data
	certs, err := pkcs12.DecodeTrustStore(certData, cert.Password)
	if err != nil {
		return certInfoList, handleFailOnError(&certInfoList, cert.Name, "truststore", fmt.Sprintf("Failed to decode P12 file '%s': %v", cert.Name, err), failOnError)
	}

	// Extract certificates
	for _, c := range certs {
		subject := c.Subject.CommonName
		if subject == "" {
			subject = fmt.Sprintf("%d", len(certInfoList)+1)
		}
		certInfo := CertificateInfo{
			Name:    cert.Name,
			Subject: subject,
			Epoch:   c.NotAfter.Unix(),
			Type:    "truststore",
		}
		certInfoList = append(certInfoList, certInfo)

		log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())
	}

	return certInfoList, nil
}
