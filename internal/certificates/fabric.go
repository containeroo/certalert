package certificates

import (
	"certalert/internal/utils"
	"fmt"
	"sort"
	"strings"
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
type extractFunction func(cert Certificate, certData []byte, failOnError bool) ([]CertificateInfo, error)

// ExtractionFunctionFabric maps each certificate type to its corresponding extraction function.
// The map allows dynamic selection of the appropriate extraction function based on the certificate type.
var ExtractionFunctionFabric = map[string]extractFunction{
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
