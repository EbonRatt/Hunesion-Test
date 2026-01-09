package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// KafkaListener handles Kafka message consumption
type KafkaListener struct {
	consumer        sarama.ConsumerGroup
	topic           string
	playbookHandler *PlaybookHandler
}

// NewKafkaListener creates a new Kafka listener
func NewKafkaListener(brokers, groupID, topic string) (*KafkaListener, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumerGroup([]string{brokers}, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaListener{
		consumer:        consumer,
		topic:           topic,
		playbookHandler: NewPlaybookHandler(),
	}, nil
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	playbookHandler *PlaybookHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}
			log.Printf("Received message from topic %s, partition %d, offset %d: %s",
				message.Topic, message.Partition, message.Offset, string(message.Value))

			// Process the message
			if err := h.playbookHandler.HandleMessage(message.Value); err != nil {
				log.Printf("Error handling message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

// Start begins listening for messages
func (kl *KafkaListener) Start(ctx context.Context) error {
	log.Printf("Starting Kafka listener on topic: %s", kl.topic)

	handler := &consumerGroupHandler{
		playbookHandler: kl.playbookHandler,
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka listener...")
			return nil
		default:
			err := kl.consumer.Consume(ctx, []string{kl.topic}, handler)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
				return err
			}
		}
	}
}

// Close closes the Kafka consumer
func (kl *KafkaListener) Close() error {
	return kl.consumer.Close()
}
