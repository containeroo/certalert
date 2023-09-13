package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	Register("/", "Shows this page (the endpoints)", Home, []string{"GET", "POST"})
}

// Home is the handler for the / route
// It displays all the endpoints
func Home(w http.ResponseWriter, r *http.Request) {
	tplData := TemplateData{
		Endpoints: Handlers,
		CSS:       CSS,
	}
	tpl, err := renderTemplate(tplBase, tplEndpoints, tplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(tpl)); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
