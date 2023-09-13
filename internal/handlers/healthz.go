package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	Register("/healthz", "Returns the health of the application", Healthz, []string{"GET", "POST"})
}

// Healthz returns the status of the application
// It returns a 200 if the application is healthy
func Healthz(w http.ResponseWriter, r *http.Request) {
	if _, err := certificates.Process(config.App.Certs, config.App.FailOnError); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
