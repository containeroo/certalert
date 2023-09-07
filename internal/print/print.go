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

// ConvertCertificatesToFormat converts the provided certificates to the specified output format
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
