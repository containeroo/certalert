package print

import (
	"certalert/internal/certificates"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCertificatesToFormat(t *testing.T) {
	certs := []certificates.Certificate{
		{
			Name: "TestCert",
			Path: "../../tests/certs/p12/without_password.p12",
			Type: "p12",
		},
	}

	// Test valid formats
	for _, format := range []string{YAMLFormat, JSONFormat, TextFormat} {
		_, err := ConvertCertificatesToFormat(format, certs, true)
		assert.Nil(t, err)
	}

	// Test unsupported format
	_, err := ConvertCertificatesToFormat("unsupported", certs, true)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported output format: unsupported", err.Error())

}
