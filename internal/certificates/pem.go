package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ExtractPEMCertificatesInfo reads the PEM file, extracts certificate information, and returns a list of CertificateInfo
func ExtractPEMCertificatesInfo(name string, certData []byte, password string) ([]CertificateInfo, error) {
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
			if password != "" {
				_, err := x509.DecryptPEMBlock(block, []byte(password))
				if err != nil {
					log.Warningf("Failed to decrypt private key '%s': %v", name, err)
				}
			} else {
				_, err := x509.ParsePKCS8PrivateKey(block.Bytes)
				if err != nil {
					log.Warningf("Failed to parse private key '%s': %v", name, err)
				}
			}
		}

		// skip if is not a certificate
		if block.Type != "CERTIFICATE" {
			certData = rest // Move to the next PEM block
			continue
		}

		counter++

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse certificate '%s': %w", name, err)
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
	if certInfoList == nil {
		return nil, fmt.Errorf("Failed to decode certificate '%s'", name)
	}

	return certInfoList, nil
}
