package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// ExampleProducer sends example messages to Kafka
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run example_producer.go <broker> <topic>")
		fmt.Println("Example: go run example_producer.go localhost:9092 ansible-tasks")
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Example 1: Create user playbook format
	playbook1 := CreateUserPlaybook{
		IP:       "10.1.1.1",
		Username: "user10",
		UserID:   "10",
		SSHUser:  "root",
		SSHPass:  "your-password-here",
	}

	message1, _ := json.Marshal(playbook1)
	msg1 := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message1),
	}

	partition, offset, err := producer.SendMessage(msg1)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent successfully! Topic: %s, Partition: %d, Offset: %d\n", topic, partition, offset)
	fmt.Printf("Message content: %s\n", string(message1))
}
