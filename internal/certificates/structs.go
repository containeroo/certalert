package certificates

import (
	"certalert/internal/utils"
	"fmt"
	"sort"
	"strings"
	"time"
)

// FileExtensionsTypes contains a sorted list of unique certificate types extracted from 'FileExtensionsToType'
var FileExtensionsTypes = []string{}

// LenFileExtensionsTypes holds the length of 'FileExtensionsTypes'
var LenFileExtensionsTypes = len(FileExtensionsTypes)

// FileExtensionsTypesString is a formatted string containing the sorted certificate types for user-friendly display
var FileExtensionsTypesString string

// init initializes 'FileExtensionsTypes', 'LenFileExtensionsTypes', and 'FileExtensionsTypesString'
func init() {
	FileExtensionsTypes = utils.ExtractMapKeys(FileExtensionsToType)
	// Sort the list of certificate types
	sort.Strings(FileExtensionsTypes)

	LenFileExtensionsTypes = len(FileExtensionsTypes)

	FileExtensionsTypesString = fmt.Sprintf("'%s' or '%s'", strings.Join(FileExtensionsTypes[:LenFileExtensionsTypes-1], "', '"), FileExtensionsTypes[LenFileExtensionsTypes-1])
}

// extractFunction is a function type representing the signature for extracting certificate information.
// It takes parameters 'name', 'certData', 'password', and 'failOnError', and returns a slice of 'CertificateInfo' and an error.
type extractFunction func(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error)

// TypeToExtractionFunction maps each certificate type to its corresponding extraction function.
// The map allows dynamic selection of the appropriate extraction function based on the certificate type.
var TypeToExtractionFunction = map[string]extractFunction{
	"p12":        ExtractP12CertificatesInfo,
	"pem":        ExtractPEMCertificatesInfo,
	"jks":        ExtractJKSCertificatesInfo,
	"p7":         ExtractP7CertificatesInfo,
	"truststore": ExtractTrustStoreCertificatesInfo,
}

// FileExtensionsToType maps each file extension to its canonical certificate type.
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
