package print

import (
	"certalert/internal/certificates"
	"fmt"
)

const (
	YAMLFormat = "yaml"
	JSONFormat = "json"
	TextFormat = "text"
)

// ConvertCertificatesToFormat converts the provided certificates to the specified output format
func ConvertCertificatesToFormat(outputFormat string, certs []certificates.Certificate, failOnError bool) (string, error) {
	certificatesInfo, err := certificates.Process(certs, failOnError)
	if err != nil {
		return "", err
	}

	formatHandlers := map[string]func(interface{}) (string, error){
		YAMLFormat: convertToYaml,
		JSONFormat: convertToJson,
		TextFormat: convertToTable,
	}

	if handler, exists := formatHandlers[outputFormat]; exists {
		return handler(certificatesInfo)
	}
	return "", fmt.Errorf("Unsupported output format: %s", outputFormat)
}
