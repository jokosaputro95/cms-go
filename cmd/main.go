package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jokosaputro95/cms-go/cmd/app"
)

func main() {
	container, err := app.StartServer(false, ".env")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
        if err := container.Start(); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()

	// Wait for interrupt signal
	<-quit
	log.Println("Received interrupt signal, shutting down...")
	

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	if err := container.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}