package main

import (
	"context"
	"log"

	"test-kafka/config"
	"test-kafka/internal/handler"
	"test-kafka/internal/kafka"
)

func main() {
	// Load configuration from config.json
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}

	log.Printf("Loaded configuration:")
	log.Printf("  Kafka Broker: %s", cfg.Kafka.Broker)
	log.Printf("  Consumer Topic: %s", cfg.Kafka.ConsumerTopic)
	log.Printf("  Producer Topic: %s", cfg.Kafka.ProducerTopic)
	log.Println()

	// Create Kafka producer for sending responses
	producerConfig := kafka.Config{
		Broker: cfg.Kafka.Broker,
		Topic:  cfg.Kafka.ProducerTopic,
	}
	producer := kafka.NewProducer(producerConfig)
	defer producer.Close()

	// Create Kafka consumer with custom handler
	consumerConfig := kafka.Config{
		Broker: cfg.Kafka.Broker,
		Topic:  cfg.Kafka.ConsumerTopic,
	}

	// Create handler with producer closure
	kafkaHandler := handler.AnsibleEventHandler(producer, cfg.Kafka.ProducerTopic)
	consumer := kafka.NewConsumer(consumerConfig, kafkaHandler)

	// Start consuming (blocks until stopped)
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
