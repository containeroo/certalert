package certificates

import "time"

var ValidTypes = []string{"p12", "pkcs12", "pfx", "pem", "crt", "jks"}

// Certificate represents a certificate configuration
type Certificate struct {
	Name     string `json:"name"`
	Enabled  *bool  `json:"enabled,omitempty"`
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
}

// ExpiryAsTime returns the expiry date as a time.Time
func (ci *CertificateInfo) ExpiryAsTime() time.Time {
	return time.Unix(ci.Epoch, 0)
}
