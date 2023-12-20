package handlers

import (
	"certalert/internal/config"
	"certalert/internal/server"
	"certalert/internal/utils"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	server.Register("/-/reload", "Reloads the configuration", Reload, "GET", "POST")
}

// Reload is an HTTP handler for the /reload route.
//
// This handler reloads the configuration file. It reads, parses, and redacts the
// configuration. It also updates the copy of the configuration used for exposing
// the current configuration via the /config route.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
func Reload(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("Force reloading configuration")

	if err := config.App.Read(viper.ConfigFileUsed()); err != nil {
		log.Fatal().Msgf("Unable to read config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.App.Parse(); err != nil {
		log.Fatal().Msgf("Unable to parse config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.DeepCopy(config.App, &config.AppCopy); err != nil {
		log.Fatal().Msgf("Unable to copy config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.RedactConfig(&config.AppCopy); err != nil {
		log.Fatal().Msgf("Unable to redact config: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration reloaded successfully"))
}
