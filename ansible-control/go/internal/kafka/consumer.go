package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

// Config holds Kafka consumer configuration
type Config struct {
	Broker string // Kafka broker address (e.g., "localhost:29092")
	Topic  string // Kafka topic name to consume from
}

// MessageHandler is a function type that processes Kafka messages
// This allows you to customize what happens when a message is received
type MessageHandler func(ctx context.Context, msg kafka.Message) error

// Consumer represents a Kafka consumer
type Consumer struct {
	reader  *kafka.Reader
	config  Config
	handler MessageHandler
}

// NewConsumer creates a new Kafka consumer with the given configuration
func NewConsumer(config Config, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.Broker},
		Topic:       config.Topic,
		StartOffset: kafka.LastOffset, // Start from the latest offset (only new messages)
		GroupID:     "agent-consumer", // Empty = no consumer group, reads directly from topic
		MinBytes:    1,                // Read even smallest messages
		MaxBytes:    10e6,             // 10MB max
	})

	return &Consumer{
		reader:  reader,
		config:  config,
		handler: handler,
	}
}

// Start begins consuming messages from Kafka
// It runs until the context is cancelled or an error occurs
func (c *Consumer) Start(ctx context.Context) error {
	// Always close the reader when we're done
	defer c.reader.Close()

	// Set up graceful shutdown on Ctrl+C
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle shutdown signals in a goroutine
	go func() {
		<-sigChan
		fmt.Println("\nStopping consumer...")
		cancel()
	}()

	// Log startup information
	fmt.Println("Starting Kafka Consumer...")
	fmt.Printf("Broker: %s\n", c.config.Broker)
	fmt.Printf("Topic: %s\n", c.config.Topic)
	fmt.Println("Waiting for NEW messages only (old messages will be skipped)...")
	fmt.Println()

	// Main consumption loop
	for {
		// ReadMessage blocks until a message is available or context is cancelled
		msg, err := c.reader.ReadMessage(ctx)

		// Check if context was cancelled (user pressed Ctrl+C)
		if err == context.Canceled || ctx.Err() == context.Canceled {
			fmt.Println("Consumer stopped.")
			return nil
		}

		// Handle other errors
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			if ctx.Err() != nil {
				return fmt.Errorf("consumer stopped due to context cancellation: %w", err)
			}
			continue // Try next message
		}

		// Process the message using the handler
		if c.handler != nil {
			if err := c.handler(ctx, msg); err != nil {
				log.Printf("Error processing message: %v\n", err)
				// Continue processing other messages even if one fails
			}
		}
	}
}

// DefaultHandler is a simple message handler that prints message details
// You can use this as a starting point or create your own handler
func DefaultHandler(ctx context.Context, msg kafka.Message) error {
	fmt.Println("--- Message Received ---")
	fmt.Printf("Topic:     %s\n", msg.Topic)
	fmt.Printf("Partition: %d\n", msg.Partition)
	fmt.Printf("Offset:    %d\n", msg.Offset)
	fmt.Printf("Key:       %s\n", string(msg.Key))
	fmt.Printf("Value:     %s\n", string(msg.Value))
	fmt.Println("------------------------")
	fmt.Println()
	return nil
}

func RunConsumer(broker, topic string) error {
	cfg := Config{
		Broker: broker,
		Topic:  topic,
	}
	c := NewConsumer(cfg, DefaultHandler)
	return c.Start(context.Background())
}
