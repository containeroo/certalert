package certificates

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
	log "github.com/sirupsen/logrus"
)

// ExtractJKSCertificatesInfo extracts certificate information from a JKS file
func ExtractJKSCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	ks := keystore.New()
	if err := ks.Load(bytes.NewReader(certData), []byte(password)); err != nil {
		return certInfoList, handleFailOnError(&certInfoList, name, "jks", fmt.Sprintf("Failed to load JKS file '%s': %v", name, err), failOnError)
	}

	for _, alias := range ks.Aliases() {
		var certificates []keystore.Certificate

		if ks.IsPrivateKeyEntry(alias) {
			entry, err := ks.GetPrivateKeyEntry(alias, []byte(password))
			if err != nil {
				if err := handleFailOnError(&certInfoList, name, "jks", fmt.Sprintf("Failed to get private key in JKS file '%s': %v", name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}
			certificates = entry.CertificateChain
		} else if ks.IsTrustedCertificateEntry(alias) {
			entry, err := ks.GetTrustedCertificateEntry(alias)
			if err != nil {
				if err := handleFailOnError(&certInfoList, name, "jks", fmt.Sprintf("Failed to get certificates in JKS file '%s': %v", name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}
			certificates = []keystore.Certificate{entry.Certificate}
		} else {
			if err := handleFailOnError(&certInfoList, name, "jks", fmt.Sprintf("Unknown entry type for alias '%s' in JKS file '%s'", alias, name), failOnError); err != nil {
				return certInfoList, err
			}
			continue
		}

		for _, cert := range certificates {
			certificate, err := x509.ParseCertificate(cert.Content)
			if err != nil {
				if err := handleFailOnError(&certInfoList, name, "jks", fmt.Sprintf("Failed to parse certificate '%s': %v", name, err), failOnError); err != nil {
					return certInfoList, err
				}
				continue
			}

			subject := certificate.Subject.ToRDNSequence().String()
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certInfoList)+1)
			}
			certInfo := CertificateInfo{
				Name:    name,
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
