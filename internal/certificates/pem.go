package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// ExtractPEMCertificatesInfo extracts certificate information from the given PEM data
func ExtractPEMCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo
	var counter int

	// Decode PEM and extract certificates
	for {
		block, rest := pem.Decode(certData)
		if block == nil {
			break
		}

		// if block is a private key, try to parse it
		if block.Type == "PRIVATE KEY" {
			var err error
			if password != "" {
				_, err = x509.DecryptPEMBlock(block, []byte(password))
			} else {
				_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			}
			if err != nil {
				if err := handleError(&certInfoList, name, "pem", fmt.Sprintf("Failed to parse/decrypt private key '%s': %v", name, err), failOnError); err != nil {
					return certInfoList, err
				}
				certData = rest
				continue
			}
		}

		// skip if is not a certificate
		if block.Type != "CERTIFICATE" {
			certData = rest
			continue
		}

		counter++

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			if err := handleError(&certInfoList, name, "pem", fmt.Sprintf("Failed to parse certificate '%s': %v", name, err), failOnError); err != nil {
				return certInfoList, err
			}
			certData = rest
			continue
		}

		subject := cert.Subject.CommonName
		if subject == "" {
			subject = fmt.Sprint(counter)
		}

		certInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   cert.NotAfter.Unix(),
			Type:    "pem",
		}

		certInfoList = append(certInfoList, certInfo)

		certData = rest // Move to the next PEM block
	}
	if len(certInfoList) == 0 {
		return certInfoList, handleError(&certInfoList, name, "pem", fmt.Sprintf("Failed to decode any certificate in '%s'", name), failOnError)
	}

	return certInfoList, nil
}
