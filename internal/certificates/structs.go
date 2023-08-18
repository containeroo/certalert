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
	Name     string `json:"name"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Valid    *bool  `json:"-" yaml:"-"`
	Path     string `json:"path"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

// CertificateInfo represents the extracted certificate information
type CertificateInfo struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Epoch   int64  `json:"epoch"`
	Type    string `json:"type"`
	Error   string `json:"error"`
}

// ExpiryAsTime returns the expiry date as a time.Time
func (ci *CertificateInfo) ExpiryAsTime() time.Time {
	return time.Unix(ci.Epoch, 0)
}
