package certificates

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
	log "github.com/sirupsen/logrus"
)

// ExtractJKSCertificatesInfo extracts certificate information from a JKS file
func ExtractJKSCertificatesInfo(cert Certificate, certData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	ks := keystore.New()
	if err := ks.Load(bytes.NewReader(certData), []byte(cert.Password)); err != nil {
		return certInfoList, handleFailOnError(&certInfoList, cert.Name, "jks", fmt.Sprintf("Failed to load JKS file '%s': %v", cert.Name, err), failOnError)
	}

	for _, alias := range ks.Aliases() {
		var certificates []keystore.Certificate

		if ks.IsPrivateKeyEntry(alias) {
			entry, err := ks.GetPrivateKeyEntry(alias, []byte(cert.Password))
			if err != nil {
				if err := handleFailOnError(&certInfoList, cert.Name, "jks", fmt.Sprintf("Failed to get private key in JKS file '%s': %v", cert.Name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}
			certificates = entry.CertificateChain
		} else if ks.IsTrustedCertificateEntry(alias) {
			entry, err := ks.GetTrustedCertificateEntry(alias)
			if err != nil {
				if err := handleFailOnError(&certInfoList, cert.Name, "jks", fmt.Sprintf("Failed to get certificates in JKS file '%s': %v", cert.Name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}
			certificates = []keystore.Certificate{entry.Certificate}
		} else {
			if err := handleFailOnError(&certInfoList, cert.Name, "jks", fmt.Sprintf("Unknown entry type for alias '%s' in JKS file '%s'", alias, cert.Name), failOnError); err != nil {
				return certInfoList, err
			}
			continue
		}

		for _, c := range certificates {
			certificate, err := x509.ParseCertificate(c.Content)
			if err != nil {
				if err := handleFailOnError(&certInfoList, cert.Name, "jks", fmt.Sprintf("Failed to parse certificate '%s': %v", cert.Name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}

			subject := certificate.Subject.CommonName
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certInfoList)+1)
			}
			certInfo := CertificateInfo{
				Name:    cert.Name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "jks",
			}
			certInfoList = append(certInfoList, certInfo)

			log.Debugf("Certificate '%s' expires on %s", subject, certInfo.ExpiryAsTime())
		}
	}

	return certInfoList, nil
}
