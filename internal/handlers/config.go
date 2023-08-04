package handlers

import (
	"bytes"
	"certalert/internal/config"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// redactVariable redacts sensitive data from the config if it is not prefixed with env: or file:
func redactVariable(s string) string {
	if strings.HasPrefix(s, "env:") || strings.HasPrefix(s, "file:") {
		return s
	}
	return "<REDACTED>"
}

// ConfigHandler returns the config as yaml
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	defer yamlEncoder.Close()
	yamlEncoder.SetIndent(2)

	// copy config and remove sensitive data
	configCopy := config.App

	configCopy.Pushgateway.Address = redactVariable(configCopy.Pushgateway.Address)
	configCopy.Pushgateway.Auth.Basic.Username = redactVariable(configCopy.Pushgateway.Auth.Basic.Username)
	configCopy.Pushgateway.Auth.Basic.Password = redactVariable(configCopy.Pushgateway.Auth.Basic.Password)
	configCopy.Pushgateway.Auth.Bearer.Token = redactVariable(configCopy.Pushgateway.Auth.Bearer.Token)

	for idx, cert := range configCopy.Certs {
		configCopy.Certs[idx].Password = redactVariable(cert.Password)
	}

	if err := yamlEncoder.Encode(&config.App); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b.Bytes())
}
