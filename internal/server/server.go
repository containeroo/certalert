package server

import (
	"certalert/internal/handlers"
	"certalert/internal/metrics"
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// NewRouter generates the router used in the HTTP Server
func NewRouter() *mux.Router {
	// Create router and define routes and return that router
	router := mux.NewRouter()

	metrics.PromMetrics = *metrics.NewMetrics()

	//register handlers
	router.HandleFunc("/", handlers.Home).Methods("GET", "POST")
	router.HandleFunc("/-/reload", handlers.Reload).Methods("GET", "POST")
	router.HandleFunc("/config", handlers.Config).Methods("GET", "POST")
	router.HandleFunc("/certificates", handlers.Certificates).Methods("GET", "POST")
	router.HandleFunc("/healthz", handlers.Healthz).Methods("GET", "POST")
	router.Handle("/metrics", http.HandlerFunc(handlers.Metrics)).Methods("GET", "POST")

	return router
}

// Run will run the HTTP Server
func RunServer(hostname string, port int) {

	// Set up a channel to listen to for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	// Define server options
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", hostname, port),
		Handler:      NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Alert the user that the server is starting
	log.Infof("Server is starting on %s", server.Addr)

	// Run the server on a new goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
		}
	}()

	// Block on this channel listeninf for those previously defined syscalls assign
	// to variable so we can let the user know why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	log.Infof("Server is shutting down due to %+v", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
