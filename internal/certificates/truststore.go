package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ExtractP12CertificatesInfo extracts certificate information from a P12 file
func ExtractTrustStoreCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Decode the P12 data
	certificates, err := pkcs12.DecodeTrustStore(certData, password)
	if err != nil {
		return certificateInfoList, handleFailOnError(&certificateInfoList, name, "truststore", fmt.Sprintf("Failed to decode P12 file '%s': %v", name, err), failOnError)
	}

	// Extract certificates
	for _, certificate := range certificates {
		subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

		certificateInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   certificate.NotAfter.Unix(),
			Type:    "truststore",
		}
		certificateInfoList = append(certificateInfoList, certificateInfo)

		log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
	}

	return certificateInfoList, nil
}
