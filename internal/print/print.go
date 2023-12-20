package print

import (
	"certalert/internal/certificates"
	"fmt"
)

// FormatHandlers maps each output format to its corresponding conversion function.
var FormatHandlers = map[string]func(interface{}) (string, error){
	"yaml": convertToYaml,
	"json": convertToJson,
	"text": convertToTable,
}

// ConvertCertificatesToFormat converts the provided certificates to the specified output format.
// It processes the certificates, and based on the specified output format, it invokes the corresponding
// conversion function from FormatHandlers to generate the formatted output.
//
// Parameters:
//   - outputFormat: string
//     The desired output format ("yaml", "json", or "text").
//   - certs: []certificates.Certificate
//     The list of certificates to convert.
//   - failOnError: bool
//     A flag indicating whether to fail on errors during certificate processing.
//
// Returns:
//   - string
//     The formatted output as a string.
//   - error
//     An error if certificate processing or conversion fails.
func ConvertCertificatesToFormat(outputFormat string, certs []certificates.Certificate, failOnError bool) (string, error) {
	certificatesInfo, err := certificates.Process(certs, failOnError)
	if err != nil {
		return "", err
	}

	if handler, exists := FormatHandlers[outputFormat]; exists {
		return handler(certificatesInfo)
	}
	return "", fmt.Errorf("Unsupported output format: %s", outputFormat)
}
