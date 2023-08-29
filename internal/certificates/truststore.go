package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo extracts certificate information from a P12 file
func ExtractTrustStoreCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// Decode the P12 data
	certs, err := pkcs12.DecodeTrustStore(certData, password)
	if err != nil {
		return certInfoList, handleError(&certInfoList, name, "truststore", fmt.Sprintf("Failed to decode P12 file '%s': %v", name, err), failOnError)
	}

	// Extract certificates
	for _, cert := range certs {
		subject := cert.Subject.CommonName
		if subject == "" {
			subject = fmt.Sprintf("%d", len(certInfoList)+1)
		}
		certInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   cert.NotAfter.Unix(),
			Type:    "truststore",
		}
		certInfoList = append(certInfoList, certInfo)

		log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())
	}

	return certInfoList, nil
}
