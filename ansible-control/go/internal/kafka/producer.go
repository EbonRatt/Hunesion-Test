package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer represents a Kafka producer
type Producer struct {
	writer *kafka.Writer
	config Config
}

// NewProducer creates a new Kafka producer with the given configuration
func NewProducer(config Config) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Broker),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		config: config,
	}
}

// SendMessage sends a message to Kafka
func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	var valueBytes []byte
	var err error

	// If value is already bytes, use it directly
	if bytes, ok := value.([]byte); ok {
		valueBytes = bytes
	} else {
		// Otherwise, marshal to JSON
		valueBytes, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Successfully sent message to Kafka topic %s with key: %s\n", p.config.Topic, key)
	return nil
}

// SendMessageToTopic sends a message to a specific topic (different from the producer's default topic)
func (p *Producer) SendMessageToTopic(ctx context.Context, topic, key string, value interface{}) error {
	var valueBytes []byte
	var err error

	// If value is already bytes, use it directly
	if bytes, ok := value.([]byte); ok {
		valueBytes = bytes
	} else {
		// Otherwise, marshal to JSON
		valueBytes, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
	}

	// Create a temporary writer for this specific topic
	writer := &kafka.Writer{
		Addr:     kafka.TCP(p.config.Broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	msg := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	log.Printf("Successfully sent message to Kafka topic %s with key: %s\n", topic, key)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}
