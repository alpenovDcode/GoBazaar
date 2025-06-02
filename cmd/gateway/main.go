package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting API Gateway...")

	// TODO: Load configuration
	// TODO: Initialize logger
	// TODO: Setup routing to microservices
	// TODO: Setup JWT middleware
	// TODO: Setup rate limiting
	// TODO: Start server

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("API Gateway shutting down...")
}
