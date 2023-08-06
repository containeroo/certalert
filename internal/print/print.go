package print

import (
	"certalert/internal/certificates"
	"fmt"
)

func EvaluateOutputFormat(outputFormat string, certs []certificates.Certificate) (string, error) {
	certificatesInfo, err := certificates.Process(certs)
	if err != nil {
		return "", err
	}
	switch outputFormat {
	case "yaml":
		return outputAsYaml(certificatesInfo)
	case "json":
		return outputAsJson(certificatesInfo)
	case "text":
		return outputAsText(certificatesInfo)
	default:
		return "", fmt.Errorf("Unsupported output format: %s", outputFormat)
	}
}
