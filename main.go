// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: Backend for BlueNode Server OS

package main

import (
	"bluenode-helper/handlers"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// Unix socket path
	socketPath = "/var/run/bnhelper.sock"
	// Graceful shutdown timeout
	shutdownTimeout = 30 * time.Second
)

func main() {
	// Parse command-line flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("bluenode-helper version %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	// Remove existing socket if it exists
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove existing socket: %v", err)
	}

	// Create Unix socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to create Unix socket: %v", err)
	}
	defer listener.Close()

	// Set socket permissions
	if err := os.Chmod(socketPath, 0660); err != nil {
		log.Fatalf("Failed to set socket permissions: %v", err)
	}

	log.Printf("Unix socket created at: %s", socketPath)

	// Create HTTP server
	mux := http.NewServeMux()

	// Register Docker API handlers
	dockerHandler := handlers.NewDockerHandler()
	dockerHandler.RegisterRoutes(mux)

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "BlueNode Helper API")
	})

	server := &http.Server{
		Handler: mux,
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Println("Server is starting...")
		serverErrors <- server.Serve(listener)
	}()

	// Channel to listen for interrupt or terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received or server error occurs
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed, forcing shutdown: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Error closing server: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}

	// Cleanup socket file
	if err := os.Remove(socketPath); err != nil {
		log.Printf("Warning: Failed to remove socket file: %v", err)
	}
}
