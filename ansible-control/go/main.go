package main

import (
	"context"
	"log"

	"test-kafka/internal/handler"
	"test-kafka/internal/kafka"
)

func main() {
	// Create Kafka producer for sending responses
	producerConfig := kafka.Config{
		Broker: "localhost:29092",
		Topic:  "producer-topic", // Response topic
	}
	producer := kafka.NewProducer(producerConfig)
	defer producer.Close()

	// Create Kafka consumer with custom handler
	consumerConfig := kafka.Config{
		Broker: "localhost:29092",
		Topic:  "consumer-topic",
	}

	// Create handler with producer closure
	kafkaHandler := handler.AnsibleEventHandler(producer, "producer-topic")
	consumer := kafka.NewConsumer(consumerConfig, kafkaHandler)

	// Start consuming (blocks until stopped)
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
