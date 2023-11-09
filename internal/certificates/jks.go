package certificates

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
	log "github.com/sirupsen/logrus"
)

// ExtractJKSCertificatesInfo extracts certificate information from a JKS file
func ExtractJKSCertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	ks := keystore.New()
	if err := ks.Load(bytes.NewReader(certificateData), []byte(cert.Password)); err != nil {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "jks", fmt.Sprintf("Failed to load JKS file '%s': %v", cert.Name, err), failOnError)
	}

	for _, alias := range ks.Aliases() {
		var certificates []keystore.Certificate

		if ks.IsPrivateKeyEntry(alias) {
			entry, err := ks.GetPrivateKeyEntry(alias, []byte(cert.Password))
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "jks", fmt.Sprintf("Failed to get private key in JKS file '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
				continue
			}
			certificates = entry.CertificateChain
		} else if ks.IsTrustedCertificateEntry(alias) {
			entry, err := ks.GetTrustedCertificateEntry(alias)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "jks", fmt.Sprintf("Failed to get certificates in JKS file '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
				continue
			}
			certificates = []keystore.Certificate{entry.Certificate}
		} else {
			if err := handleFailOnError(&certificateInfoList, cert.Name, "jks", fmt.Sprintf("Unknown entry type for alias '%s' in JKS file '%s'", alias, cert.Name), failOnError); err != nil {
				return certificateInfoList, err
			}
			continue
		}

		for _, certificate := range certificates {
			x509Cert, err := x509.ParseCertificate(certificate.Content)
			if err != nil {
				if err := handleFailOnError(&certificateInfoList, cert.Name, "jks", fmt.Sprintf("Failed to parse certificate '%s': %v", cert.Name, err), failOnError); err != nil {
					return certificateInfoList, err
				}
				continue
			}

			subject := generateCertificateSubject(x509Cert.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)
			certificateInfo := CertificateInfo{
				Name:    cert.Name,
				Subject: subject,
				Epoch:   x509Cert.NotAfter.Unix(),
				Type:    "jks",
			}
			certificateInfoList = append(certificateInfoList, certificateInfo)

			log.Debugf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
		}
	}

	// no check of len of certificateInfoList needed here, because if the JKS file is empty,
	// ks.Load will throw an error and we will never get here

	return certificateInfoList, nil
}
