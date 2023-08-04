package certificates

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
)

// ExtractJksCertificatesInfo reads the JKS file, extracts certificate information, and returns a list of CertificateInfo
func ExtractJKSCertificatesInfo(name string, certData []byte, password string) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	ks := keystore.New()
	if err := ks.Load(bytes.NewReader(certData), []byte(password)); err != nil {
		return nil, fmt.Errorf("Failed to load JKS file '%s': %w", name, err)
	}

	for _, alias := range ks.Aliases() {
		var certificates []keystore.Certificate
		// Check the entry type and get certificates accordingly
		if ks.IsPrivateKeyEntry(alias) {
			entry, err := ks.GetPrivateKeyEntry(alias, []byte(password))
			if err != nil {
				return nil, fmt.Errorf("Failed to get entries in JKS file '%s': %w", name, err)
			}
			certificates = entry.CertificateChain
		} else if ks.IsTrustedCertificateEntry(alias) {
			entry, err := ks.GetTrustedCertificateEntry(alias)
			if err != nil {
				return nil, fmt.Errorf("Failed to get certificates in JKS file '%s': %w", name, err)
			}
			certificates = []keystore.Certificate{entry.Certificate}
		} else {
			return nil, fmt.Errorf("Unknown entry type for alias '%s' in JKS file '%s'", alias, name)
		}

		// Iterate over all certificates in the entry
		for idx, cert := range certificates {
			certificate, err := x509.ParseCertificate(cert.Content)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse certificate '%s': %w", name, err)
			}

			subject := certificate.Subject.CommonName
			if subject == "" {
				subject = fmt.Sprint(idx)
			}

			certInfo := CertificateInfo{
				Name:    name,
				Subject: subject,
				Epoch:   certificate.NotAfter.Unix(),
				Type:    "jks",
			}

			certInfoList = append(certInfoList, certInfo)
		}
	}

	return certInfoList, nil
}
