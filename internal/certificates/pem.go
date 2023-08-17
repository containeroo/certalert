package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ExtractPEMCertificatesInfo extracts certificate information from the given PEM data
func ExtractPEMCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo
	var counter int

	// handleError is a helper function to handle failOnError
	handleError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warningf("Failed to extract certificate information: %v", errMsg)
		certInfoList = append(certInfoList, CertificateInfo{
			Name:  name,
			Type:  "pem",
			Error: errMsg,
		})
		return nil
	}

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
				if handleError(fmt.Sprintf("Failed to parse/decrypt private key '%s': %v", name, err)) != nil {
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
			if handleError(fmt.Sprintf("Failed to parse certificate '%s': %v", name, err)) != nil {
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
		if err := handleError(fmt.Sprintf("Failed to decode any certificate in '%s'", name)); err != nil {
			return certInfoList, err
		}
	}

	return certInfoList, nil
}
