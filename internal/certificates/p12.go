package certificates

import (
	"fmt"

	"github.com/rs/zerolog/log"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func init() {
	registerCertificateType("p12", ExtractP12CertificatesInfo, "p12", "pkcs12", "pfx")
}

// ExtractP12CertificatesInfo extracts certificate information from a PKCS#12 (P12) file.
//
// This function takes a Certificate struct, the raw certificate data as a byte slice, and a
// flag indicating whether to fail on error. It returns a slice of CertificateInfo containing
// information about each certificate found in the P12 file.
//
// The function decodes the P12 data, extracts the main certificate and any associated CA certificates,
// and prepares for extraction. It then iterates through the certificates, logging information about
// each certificate, including its subject, expiration time, and type.
//
// Parameters:
//   - cert: Certificate
//     A Certificate struct representing the P12 file, including its name and other details.
//   - certificateData: []byte
//     The raw binary data of the P12 file.
//   - failOnError: bool
//     A flag indicating whether to fail immediately on encountering an error.
//
// Returns:
//   - []CertificateInfo
//     A slice of CertificateInfo structs containing information about each certificate in the P12 file.
//   - error
//     An error, if any, encountered during the extraction process. If failOnError is false, the
//     function may return a non-nil error along with the partial list of CertificateInfo.
func ExtractP12CertificatesInfo(cert Certificate, certificateData []byte, failOnError bool) ([]CertificateInfo, error) {
	var certificateInfoList []CertificateInfo

	// Decode the P12 data
	_, certificate, caCerts, err := pkcs12.DecodeChain(certificateData, cert.Password)
	if err != nil {
		return certificateInfoList, handleFailOnError(&certificateInfoList, cert.Name, "p12", fmt.Sprintf("Failed to decode P12 file '%s': %v", cert.Name, err), failOnError)
	}

	// Prepare for extraction
	certificates := append(caCerts, certificate)

	// Extract certificates
	for _, certificate := range certificates {
		subject := generateCertificateSubject(certificate.Subject.ToRDNSequence().String(), len(certificateInfoList)+1)

		certificateInfo := CertificateInfo{
			Name:    cert.Name,
			Subject: subject,
			Epoch:   certificate.NotAfter.Unix(),
			Type:    "p12",
		}
		certificateInfoList = append(certificateInfoList, certificateInfo)

		log.Debug().Msgf("Certificate '%s' expires on %s", subject, certificateInfo.ExpiryAsTime())
	}

	return certificateInfoList, nil
}
