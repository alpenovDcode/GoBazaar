package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Cart Service...")

	// TODO: Load configuration
	// TODO: Initialize logger
	// TODO: Connect to Redis
	// TODO: Connect to NATS
	// TODO: Setup HTTP routes (Gin)
	// TODO: Start server

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Cart Service shutting down...")
}
