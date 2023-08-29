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

// Run starts the HTTP server
func Run(hostname string, port int) {
	// Use a buffered channel for runChan to prevent signal drops
	runChan := make(chan os.Signal, 1)
	signal.Notify(runChan, os.Interrupt, syscall.SIGTERM)

	// Create a cancelable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the router and configure routes
	router := newRouter()

	// Create the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", hostname, port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// Run the server in a new goroutine
	go func() {
		log.Infof("Server is starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for signals
	select {
	case sig := <-runChan:
		log.Printf("Received signal: %v. Shutting down gracefully...", sig)
		cancel() // Cancel context to trigger graceful shutdown
		// Force shutdown after 10 seconds
		time.AfterFunc(10*time.Second, func() {
			log.Fatal("Timed out waiting for server to shut down")
		})
	}

	// Shutdown the server and wait for it to finish
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully shut down")
}

// NewRouter generates the router used in the HTTP Server
func newRouter() *mux.Router {
	router := mux.NewRouter()

	// Add the handlers to the router
	for _, h := range handlers.Handlers {
		router.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}

	// Custom 404 handler
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router
}

// notFoundHandler handles 404 responses
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	body := `
		<html>
			<head>
				<title>404 Not Found</title>
			</head>
			<body>
				<h1>404 Not Found</h1>
				<p>The page you requested could not be found.</p>
			</body>
		</html>
	`
	w.Write([]byte(body))
}
