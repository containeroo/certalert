package certificates

import (
	"time"
)

// Certificate represents a certificate configuration
type Certificate struct {
	Name     string `mapstructure:"name"`
	Enabled  *bool  `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty"`
	Path     string `mapstructure:"path"`
	Password string `mapstructure:"password,omitempty" yaml:"password,omitempty"`
	Type     string `mapstructure:"type" yaml:"type,omitempty"`
}

// CertificateInfo represents the extracted certificate information
type CertificateInfo struct {
	Name    string `mapstructure:"name"`
	Subject string `mapstructure:"subject"`
	Epoch   int64  `mapstructure:"epoch"`
	Type    string `mapstructure:"type,omitempty"`
	Error   string `mapstructure:"error"`
}

// ExpiryAsTime returns the expiry date as a time.Time
func (ci *CertificateInfo) ExpiryAsTime() time.Time {
	return time.Unix(ci.Epoch, 0)
}
