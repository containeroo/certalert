package certificates

import (
	"fmt"

	"github.com/rs/zerolog/log"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func init() {
	registerCertificateType("truststore", ExtractTrustStoreCertificatesInfo, "ts", "truststore")
}

// ExtractTrustStoreCertificatesInfo extracts certificate information from a TrustStore file (typically in P12 format).
//
// This function decodes the provided certificateData using the given password and extracts certificate information
// such as subject, expiry time, and type. The extracted information is returned as a list of CertificateInfo structs.
//
// Parameters:
//   - cert: Certificate
//     The configuration for the certificate, including its name, path, and optional password.
//   - certificateData: []byte
//     The binary data of the TrustStore file.
//   - failOnError: bool
//     A flag indicating whether the function should fail immediately on encountering an error.
//
// Returns:
//   - []CertificateInfo
//     A list of CertificateInfo structs containing information about the extracted certificates.
//   - error
//     An error if the extraction process encounters issues. If failOnError is true, this error will be non-nil.

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

		log.Debug().Msgf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
	}

	return certificateInfoList, nil
}
