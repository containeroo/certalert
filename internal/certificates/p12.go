package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo extracts certificate information from a P12 file
func ExtractP12CertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// Decode the P12 data
	_, certificate, caCerts, err := pkcs12.DecodeChain(certData, password)
	if err != nil {
		return certInfoList, handleFailOnError(&certInfoList, name, "p12", fmt.Sprintf("Failed to decode P12 file '%s': %v", name, err), failOnError)
	}

	// Prepare for extraction
	certs := append(caCerts, certificate)

	// Extract certificates
	for _, cert := range certs {
		subject := cert.Subject.ToRDNSequence().String()
		if subject == "" {
			subject = fmt.Sprintf("%d", len(certInfoList)+1)
		}
		certInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   cert.NotAfter.Unix(),
			Type:    "p12",
		}
		certInfoList = append(certInfoList, certInfo)

		log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())
	}

	return certInfoList, nil
}
