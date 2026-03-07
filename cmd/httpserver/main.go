package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JStephen72/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	srv, err := server.New(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	srv.Close()
	log.Println("Server gracefully stopped")
}
