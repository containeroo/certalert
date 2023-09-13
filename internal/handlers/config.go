package handlers

import (
	"bytes"
	"certalert/internal/config"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func init() {
	Register("/config", Config, "GET", "POST")
}

// Config is the handler for the /config route
// It returns the currently active configuration file
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
	if _, err := w.Write(b.Bytes()); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
