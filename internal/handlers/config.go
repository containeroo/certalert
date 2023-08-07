package handlers

import (
	"bytes"
	"certalert/internal/config"
	"certalert/internal/utils"
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

	if utils.HasKey(configCopy.Pushgateway, "Address") {
		configCopy.Pushgateway.Address = redactVariable(configCopy.Pushgateway.Address)
	}

	if utils.HasKey(configCopy.Pushgateway, "Basic.Username") {
		configCopy.Pushgateway.Auth.Basic.Username = redactVariable(configCopy.Pushgateway.Auth.Basic.Username)
	}

	if utils.HasKey(configCopy.Pushgateway, "Basic.Password") {
		configCopy.Pushgateway.Auth.Basic.Password = redactVariable(configCopy.Pushgateway.Auth.Basic.Password)
	}

	if utils.HasKey(configCopy.Pushgateway, "Bearer.Token") {
		configCopy.Pushgateway.Auth.Bearer.Token = redactVariable(configCopy.Pushgateway.Auth.Bearer.Token)
	}

	for idx, cert := range configCopy.Certs {
		configCopy.Certs[idx].Password = redactVariable(cert.Password)
	}

	if err := yamlEncoder.Encode(&configCopy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b.Bytes())
}
