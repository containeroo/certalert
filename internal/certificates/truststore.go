package certificates

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func init() {
	registerCertificateType("truststore", ExtractTrustStoreCertificatesInfo, "ts", "truststore")
}

// ExtractP12CertificatesInfo extracts certificate information from a P12 TrustStore file.
func ExtractTrustStoreCertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Decode the P12 data
	certificates, err := pkcs12.DecodeTrustStore(certificateData, cert.Password)
	if err != nil {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "truststore", fmt.Sprintf("Failed to decode P12 file '%s': %v", cert.Name, err), failOnError)
	}

	// Extract certificates
	for _, certificate := range certificates {
		subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

		certificateInfo := CertificateInfo{
			Name:    cert.Name,
			Subject: subject,
			Epoch:   certificate.NotAfter.Unix(),
			Type:    "truststore",
		}
		certificateInfoList = append(certificateInfoList, certificateInfo)

		log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
	}

	return certificateInfoList, nil
}
