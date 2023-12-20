package handlers

import (
	"bytes"
	"certalert/internal/config"
	"certalert/internal/server"
	"net/http"

	"gopkg.in/yaml.v3"
)

func init() {
	server.Register("/config", "Provides the currently active configuration file. Plaintext passwords are redacted", Config, "GET", "POST")
}

// Config is an HTTP handler for the /config route.
//
// This handler returns the currently active configuration file in YAML format.
// It encodes the configuration using a YAML encoder and writes the result to the
// HTTP response writer. If an error occurs during encoding or writing, it returns
// an HTTP 500 Internal Server Error response.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
func Config(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	defer yamlEncoder.Close()
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(&config.AppCopy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b.Bytes())
}
