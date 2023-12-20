package certificates

import (
	"fmt"
	"sort"
	"strings"
)

// FileExtensionsTypes contains a sorted list of unique certificate types.
type FileExtensionsTypes []string

// String returns a string representation of the FileExtensionsTypes.
func (f FileExtensionsTypes) String() string {
	sort.Strings(f)

	lenSortedSlice := len(f)
	return fmt.Sprintf("'%s' or '%s'", strings.Join(f[:lenSortedSlice-1], "', '"), f[lenSortedSlice-1])
}

var FileExtensionsTypesSorted FileExtensionsTypes

// extractFunction 	is a function type representing the signature for extracting certificate information.
type extractFunction func(cert Certificate, certData []byte, failOnError bool) ([]CertificateInfo, error)

// TypeToExtractionFunction maps each certificate type to its corresponding extraction function.
// The map allows dynamic selection of the appropriate extraction function based on the certificate type.
var TypeToExtractionFunction = map[string]extractFunction{}

// FileExtensionsToType maps each file extension to its canonical certificate type.
// The canonical type is used to select the appropriate extraction function from 'TypeToExtractionFunction'.
var FileExtensionsToType = map[string]string{}

// registerCertificateType registers a certificate type along with its extraction function and associated file extensions.
//
// This function adds an entry to the TypeToExtractionFunction map, mapping the provided certType to the
// corresponding extractFunction. It also updates the FileExtensionsToType map to associate the specified
// file extensions with the given certificate type. The sorted list of file extensions, FileExtensionsTypesSorted,
// is updated accordingly.
//
// If the certType is already registered or if an extension is already mapped to a different certificate type,
// the function panics with a corresponding error message.
//
// Parameters:
//   - certType: string
//     The certificate type to register.
//   - e: extractFunction
//     The extraction function associated with the certificate type.
//   - extensions: ...string
//     Variable number of file extensions associated with the certificate type.
//
// Panics:
//   - If certType is already registered.
//   - If an extension is already mapped to a different certificate type.
func registerCertificateType(certType string, e extractFunction, extensions ...string) {
	if _, exists := TypeToExtractionFunction[certType]; exists {
		panic(fmt.Sprintf("Certificate type '%s' is already registered", certType))
	}

	TypeToExtractionFunction[certType] = e // Register the extraction function

	// Add the file extensions to the map
	for _, ext := range extensions {
		// Check if the extension is already mapped to a different certificate type
		if existingCertType, exists := FileExtensionsToType[ext]; exists && existingCertType != certType {
			panic(fmt.Sprintf("Extension '%s' is already mapped to certificate type '%s'", ext, existingCertType))
		}
		FileExtensionsToType[ext] = certType
	}

	// Append the extensions to the sorted list
	FileExtensionsTypesSorted = append(FileExtensionsTypesSorted, extensions...)
}
