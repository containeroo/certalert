package handlers

import (
	"bytes"
	"certalert/internal/config"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Config returns the config as yaml
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
