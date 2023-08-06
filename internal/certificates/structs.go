package certificates

import "time"

var ValidTypes = []string{"p12", "pkcs12", "pfx", "pem", "crt", "jks"}

// Certificate represents a certificate configuration
type Certificate struct {
	Name     string `yaml:"name"`
	Enabled  *bool  `yaml:"enabled,omitempty"`
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
	Type     string `yaml:"type"`
}

// CertificateInfo represents the extracted certificate information
type CertificateInfo struct {
	Name    string `yaml:"name"`
	Subject string `yaml:"subject"`
	Epoch   int64  `yaml:"epoch"`
	Type    string `yaml:"type"`
}

// ExpiryAsTime returns the expiry date as a time.Time
func (ci *CertificateInfo) ExpiryAsTime() time.Time {
	return time.Unix(ci.Epoch, 0)
}
