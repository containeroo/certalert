package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo reads the P12 file, extracts certificate information, and returns a list of CertificateInfo
func ExtractP12CertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// Decode the P12 data
	_, certificate, caCerts, err := pkcs12.DecodeChain(certData, password)
	if err != nil {
		return certInfoList, handleError(&certInfoList, name, "p12", fmt.Sprintf("Failed to decode P12 file '%s': %v", name, err), failOnError)
	}

	// Prepare for extraction
	certs := append(caCerts, certificate)
	var counter int

	// Extract certificates
	for _, cert := range certs {
		counter++
		subject := cert.Subject.CommonName
		if subject == "" {
			subject = fmt.Sprint(counter)
		}

		certInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   cert.NotAfter.Unix(),
			Type:    "p12",
		}
		log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())

		certInfoList = append(certInfoList, certInfo)
	}

	return certInfoList, nil
}
