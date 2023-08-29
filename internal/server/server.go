package server

import (
	"certalert/internal/handlers"
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

// RunServer starts the HTTP server
func RunServer(hostname string, port int) {
	// Set up a channel to listen for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context for graceful server shutdown
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

	// Log the server start
	log.Printf("Server is starting on %s", server.Addr)

	// Run the server in a new goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start due to error: %v", err)
		}
	}()

	// Wait for a signal
	interrupt := <-runChan

	// Log and then gracefully terminate the server
	log.Printf("Server is shutting down due to: %v", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shut down: %v", err)
	}
}

// NewRouter generates the router used in the HTTP Server
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// Add the handlers to the router
	for _, h := range handlers.Handlers {
		router.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}

	return router
}
