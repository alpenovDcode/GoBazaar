package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Payment Service...")

	// TODO: Load configuration
	// TODO: Initialize logger
	// TODO: Connect to Stripe
	// TODO: Connect to NATS
	// TODO: Setup HTTP routes (Gin)
	// TODO: Setup gRPC client for Order Service
	// TODO: Start server

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Payment Service shutting down...")
}
