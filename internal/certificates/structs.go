package certificates

import (
	"time"
)

type extractFunction func(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error)

// Map each certificate type to its extraction function
var TypeToExtractionFunction = map[string]extractFunction{
	"p12":        ExtractP12CertificatesInfo,
	"pem":        ExtractPEMCertificatesInfo,
	"jks":        ExtractJKSCertificatesInfo,
	"p7":         ExtractP7CertificatesInfo,
	"truststore": ExtractTrustStoreCertificatesInfo,
}

// Map each file extension to its canonical certificate type
var FileExtensionsToType = map[string]string{
	"p12":        "p12",
	"pkcs12":     "p12",
	"pfx":        "p12",
	"pem":        "pem",
	"crt":        "pem",
	"jks":        "jks",
	"p7":         "p7",
	"p7b":        "p7",
	"p7c":        "p7",
	"truststore": "truststore",
	"ts":         "truststore",
}

// Certificate represents a certificate configuration
type Certificate struct {
	Name     string `mapstructure:"name"`
	Enabled  *bool  `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty"`
	Path     string `mapstructure:"path"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
}

// CertificateInfo represents the extracted certificate information
type CertificateInfo struct {
	Name    string `mapstructure:"name"`
	Subject string `mapstructure:"subject"`
	Epoch   int64  `mapstructure:"epoch"`
	Type    string `mapstructure:"type"`
	Error   string `mapstructure:"error"`
}

// ExpiryAsTime returns the expiry date as a time.Time
func (ci *CertificateInfo) ExpiryAsTime() time.Time {
	return time.Unix(ci.Epoch, 0)
}
