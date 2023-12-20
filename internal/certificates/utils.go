package certificates

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// GetCertificateByName retrieves a certificate by its name from a list of certificates.
//
// This function searches the provided list of certificates for a certificate with the specified name.
// If found, a pointer to the matching Certificate is returned; otherwise, an error is returned indicating
// that the certificate with the given name was not found.
//
// Parameters:
//   - name: string
//     The name of the certificate to retrieve.
//   - certificates: []Certificate
//     The list of certificates to search.
//
// Returns:
//   - *Certificate
//     A pointer to the Certificate with the specified name if found; otherwise, nil.
//   - error
//     An error indicating the reason if the certificate with the given name is not found.
func GetCertificateByName(name string, certificates []Certificate) (*Certificate, error) {
	for _, cert := range certificates {
		if cert.Name == name {
			return &cert, nil
		}
	}
	return nil, fmt.Errorf("Certificate '%s' not found", name)
}

// certExistsInSlice checks if a given certificate (CertificateInfo) exists in a slice of certificates.
//
// This function iterates through the provided slice of CertificateInfo and compares the name, subject, and type
// of the given certificate with each element in the slice. If a matching certificate is found, the function returns true;
// otherwise, it returns false, indicating that the certificate is not present in the slice.
//
// Parameters:
//   - cert: CertificateInfo
//     The certificate to check for existence in the slice.
//   - slice: []CertificateInfo
//     The slice of certificates to search for the specified certificate.
//
// Returns:
//   - bool
//     True if the given certificate exists in the slice; otherwise, false.
func certExistsInSlice(cert CertificateInfo, slice []CertificateInfo) bool {
	for _, c := range slice {
		if cert.Name == c.Name && cert.Subject == c.Subject && cert.Type == c.Type {
			return true
		}
	}
	return false
}

// handleFailOnError handles errors encountered during the extraction of certificate information.
//
// This function is used to manage errors that may occur during the extraction of certificate information.
// If failOnError is true, it returns an error with the provided error message (errMsg). If failOnError is false,
// it logs a warning message with the error and appends a CertificateInfo entry to the specified slice of certificates
// (certInfoList) containing details about the failed extraction.
//
// Parameters:
//   - certInfoList: *[]CertificateInfo
//     A pointer to the slice of CertificateInfo, where information about the failed extraction is appended.
//   - certName: string
//     The name of the certificate for which the extraction failed.
//   - certType: string
//     The type of the certificate for which the extraction failed.
//   - errMsg: string
//     The error message describing the reason for the extraction failure.
//   - failOnError: bool
//     A boolean indicating whether the extraction failure should result in an error or a warning log.
//
// Returns:
//   - error
//     If failOnError is true, returns an error with the provided error message (errMsg); otherwise, returns nil.
func handleFailOnError(certInfoList *[]CertificateInfo, certName, certType, errMsg string, failOnError bool) error {
	if failOnError {
		return fmt.Errorf(errMsg)
	}
	log.Warn().Msgf("Failed to extract certificate information: %v", errMsg)
	*certInfoList = append(*certInfoList, CertificateInfo{
		Name:  certName,
		Type:  certType,
		Error: errMsg,
	})
	return nil
}

// generateCertificateSubject generates a certificate subject string based on the given default subject
// and an index. If the default subject is empty, it constructs a default subject using the index.
//
// Parameters:
//   - defaultSubject: string
//     The default subject string to use, or an empty string.
//   - index: int
//     The index to use when constructing the default subject.
//
// Returns:
//   - string
//     The generated certificate subject string.
func generateCertificateSubject(defaultSubject string, index int) string {
	if defaultSubject == "" {
		defaultSubject = fmt.Sprintf("Certificate %d", index)
	}
	return defaultSubject
}
