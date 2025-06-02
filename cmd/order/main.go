package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Order Service...")

	// TODO: Load configuration
	// TODO: Initialize logger
	// TODO: Connect to PostgreSQL
	// TODO: Connect to NATS
	// TODO: Setup gRPC server
	// TODO: Start server

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Order Service shutting down...")
}
