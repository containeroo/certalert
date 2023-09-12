package certificates

import (
	"certalert/internal/utils"
	"fmt"
	"sort"
	"strings"
	"time"
)

// FileExtensionsTypesSorted contains a sorted list of unique certificate types extracted from 'FileExtensionsToType'
var FileExtensionsTypesSorted []string

// FileExtensionsTypesSortedString is a formatted string containing the sorted certificate types for user-friendly display
var FileExtensionsTypesSortedString string

// init initializes 'FileExtensionsTypesSorted' and 'FileExtensionsTypesSortedString'
func init() {
	FileExtensionsTypesSorted = utils.ExtractMapKeys(FileExtensionsToType)
	sort.Strings(FileExtensionsTypesSorted)

	lenFileExtensionsTypesSorted := len(FileExtensionsTypesSorted)

	FileExtensionsTypesSortedString = fmt.Sprintf("'%s' or '%s'", strings.Join(FileExtensionsTypesSorted[:lenFileExtensionsTypesSorted-1], "', '"), FileExtensionsTypesSorted[lenFileExtensionsTypesSorted-1])
}

// extractFunction is a function type representing the signature for extracting certificate information.
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
// The canonical type is used to select the appropriate extraction function from 'TypeToExtractionFunction'.
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
