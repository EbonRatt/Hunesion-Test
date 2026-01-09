package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command line flags
	kafkaBrokers := flag.String("brokers", "localhost:9092", "Kafka broker addresses")
	kafkaGroupID := flag.String("group", "ansible-go-group", "Kafka consumer group ID")
	kafkaTopic := flag.String("topic", "ansible-tasks", "Kafka topic to consume from")
	flag.Parse()

	log.Println("Starting Ansible-like system with Kafka listener...")
	log.Printf("Kafka Brokers: %s", *kafkaBrokers)
	log.Printf("Kafka Group ID: %s", *kafkaGroupID)
	log.Printf("Kafka Topic: %s", *kafkaTopic)

	// Create Kafka listener
	listener, err := NewKafkaListener(*kafkaBrokers, *kafkaGroupID, *kafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka listener: %v", err)
	}
	defer listener.Close()

	// Print example message formats
	log.Println("\n=== Example Message Formats ===")
	ExampleMessageFormat()
	log.Println("\n=== Listening for messages... ===")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start listener in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- listener.Start(ctx)
	}()

	// Wait for interrupt or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Listener error: %v", err)
		}
	}

	log.Println("Shutdown complete")
}
