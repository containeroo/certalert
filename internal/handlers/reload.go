package handlers

import (
	"certalert/internal/config"
	"log"
	"net/http"
)

// ReloadHandler is a handler function that reloads the application configuration
func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	if err := config.ParseConfig(&config.App); err != nil {
		log.Fatalf("Unable to parse config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration reloaded successfully"))
}
