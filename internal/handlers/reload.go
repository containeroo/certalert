package handlers

import (
	"certalert/internal/config"
	"certalert/internal/utils"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Reload is the handler for the /reload route
// It reloads the configuration file
func Reload(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Force reloading configuration")

	if err := config.App.Read(viper.ConfigFileUsed()); err != nil {
		log.Fatalf("Unable to read config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.App.Parse(); err != nil {
		log.Fatalf("Unable to parse config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.DeepCopy(config.App, &config.AppCopy); err != nil {
		log.Fatalf("Unable to copy config: %s", err)
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
