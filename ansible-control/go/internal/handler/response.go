package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"test-kafka/internal/kafka"
	"test-kafka/internal/models"
)

// SendResponse sends a response message back to Kafka
func SendResponse(ctx context.Context, producer *kafka.Producer, responseTopic string, event *models.KafkaEvent, status string, message string, errMsg string) {
	if producer == nil {
		log.Printf("Producer not available, skipping response message\n")
		return
	}

	response := models.KafkaResponseEvent{
		Status:       status,
		EventType:    event.EventType,
		Playbook:     event.Playbook,
		TargetServer: event.TargetServer,
		Message:      message,
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	if errMsg != "" {
		response.Error = errMsg
	}

	// Use response topic if provided, otherwise use the same topic as input
	topic := responseTopic
	if topic == "" {
		topic = "test-response" // Default response topic
	}

	// Use event timestamp or playbook as key for message ordering
	key := fmt.Sprintf("%s-%s", event.EventType, event.Playbook)
	if len(event.TargetServer) > 0 {
		key = fmt.Sprintf("%s-%s-%s", event.EventType, event.Playbook, event.TargetServer[0])
	}

	if err := producer.SendMessageToTopic(ctx, topic, key, response); err != nil {
		log.Printf("Failed to send response message: %v\n", err)
	} else {
		log.Printf("Successfully sent response message to topic %s\n", topic)
	}
}
