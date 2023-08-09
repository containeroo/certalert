package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"net/http"
)

// HealthzHandler returns the status of the application
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := certificates.Process(config.App.Certs, true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
