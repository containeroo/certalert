package certificates

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
)

func ExtractJKSCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// handleError is a helper function to handle failOnError
	handleError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		certInfoList = append(certInfoList, CertificateInfo{
			Name:  name,
			Type:  "jks",
			Error: errMsg,
		})
		return nil
	}

	ks := keystore.New()
	if err := ks.Load(bytes.NewReader(certData), []byte(password)); err != nil {
		return certInfoList, handleError(fmt.Sprintf("Failed to load JKS file '%s': %v", name, err))
	}

	for _, alias := range ks.Aliases() {
		var certificates []keystore.Certificate

		if ks.IsPrivateKeyEntry(alias) {
			entry, err := ks.GetPrivateKeyEntry(alias, []byte(password))
			if err != nil {
				if handleError(fmt.Sprintf("Failed to get private key in JKS file '%s': %v", name, err)) != nil {
					return certInfoList, nil
				}
				continue
			}
			certificates = entry.CertificateChain
		} else if ks.IsTrustedCertificateEntry(alias) {
			entry, err := ks.GetTrustedCertificateEntry(alias)
			if err != nil {
				if handleError(fmt.Sprintf("Failed to get certificates in JKS file '%s': %v", name, err)) != nil {
					return certInfoList, nil
				}
				continue
			}
			certificates = []keystore.Certificate{entry.Certificate}
		} else {
			if handleError(fmt.Sprintf("Unknown entry type for alias '%s' in JKS file '%s'", alias, name)) != nil {
				return certInfoList, nil
			}
			continue
		}

		for _, cert := range certificates {
			certificate, err := x509.ParseCertificate(cert.Content)
			if err != nil {
				if handleError(fmt.Sprintf("Failed to parse certificate '%s': %v", name, err)) != nil {
					return certInfoList, nil
				}
				continue
			}

			subject := certificate.Subject.CommonName
			if subject == "" {
				subject = fmt.Sprintf("%d", len(certInfoList)+1)
			}

			certInfoList = append(certInfoList, CertificateInfo{
				Name:    name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "jks",
			})
		}
	}
	return certInfoList, nil
}
