package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Auth Service...")

	// TODO: Load configuration
	// TODO: Initialize logger
	// TODO: Connect to PostgreSQL
	// TODO: Setup HTTP routes (Gin)
	// TODO: Setup gRPC server
	// TODO: Start servers

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Auth Service shutting down...")
}
