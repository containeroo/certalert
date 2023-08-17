package handlers

import (
	"certalert/internal/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ReloadHandler is a handler function that reloads the application configuration
func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Force reloading configuration")

	if err := config.ParseConfig(&config.App, config.FailOnError); err != nil {
		log.Fatalf("Unable to parse config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.RedactConfig(&config.AppCopy); err != nil {
		log.Fatalf("Unable to redact config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration reloaded successfully"))
}
