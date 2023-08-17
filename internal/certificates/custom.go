package certificates

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type ContentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"optional,explicit,tag:0"`
}

func extractCertificateFromPKCS7(data []byte) (*x509.Certificate, error) {
	var contentInfo ContentInfo

	// Parse the ContentInfo
	_, err := asn1.Unmarshal(data, &contentInfo)
	if err != nil {
		return nil, err
	}

	// Check if it's PKCS#7 data
	if !contentInfo.ContentType.Equal(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}) {
		return nil, errors.New("not a PKCS#7 signed data content info")
	}

	// Parse the SignedData
	var signedDataContent asn1.RawValue
	_, err = asn1.Unmarshal(contentInfo.Content.Bytes, &signedDataContent)
	if err != nil {
		return nil, err
	}

	// Extract the actual certificate
	return x509.ParseCertificate(signedDataContent.Bytes)
}

func ExtractCustomCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// handleError is a helper function to handle failOnError.
	handleError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warningf("Failed to extract certificate information: %v", errMsg)
		certInfoList = append(certInfoList, CertificateInfo{
			Name:  name,
			Type:  "p12",
			Error: errMsg,
		})
		return nil
	}

	cert, err := extractCertificateFromPKCS7(certData)
	if err != nil {
		fmt.Println("Failed to extract the certificate:", err)
		return certInfoList, handleError(fmt.Sprintf("Failed to extract the certificate: %v", err))
	}

	// Print the certificate details
	fmt.Println("Certificate Subject:", cert.Subject)

	return certInfoList, nil
}
