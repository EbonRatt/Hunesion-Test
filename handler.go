package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// PlaybookHandler processes playbook tasks
type PlaybookHandler struct {
}

// NewPlaybookHandler creates a new playbook handler
func NewPlaybookHandler() *PlaybookHandler {
	return &PlaybookHandler{}
}

// HandleMessage processes incoming Kafka messages
func (ph *PlaybookHandler) HandleMessage(message []byte) error {
	// Try to parse as CreateUserPlaybook first
	playbook, err := ParseCreateUserPlaybook(message)
	if err == nil {
		return ph.HandleCreateUserPlaybook(playbook)
	}

	// Try to parse as generic Task
	task, err := ParseTaskFromKafkaMessage(message)
	if err == nil {
		return ph.HandleTask(task)
	}

	return fmt.Errorf("failed to parse message as playbook or task")
}

// HandleCreateUserPlaybook handles the create user playbook
func (ph *PlaybookHandler) HandleCreateUserPlaybook(playbook *CreateUserPlaybook) error {
	log.Printf("Executing create user playbook: IP=%s, Username=%s, UserID=%s",
		playbook.IP, playbook.Username, playbook.UserID)

	// Set defaults if not provided
	sshUser := playbook.SSHUser
	if sshUser == "" {
		sshUser = "root" // Default SSH user
	}

	err := CreateUser(
		playbook.IP,
		sshUser,
		playbook.SSHPass,
		playbook.SSHKey,
		playbook.Username,
		playbook.UserID,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("Successfully executed create user playbook for %s on %s", playbook.Username, playbook.IP)
	return nil
}

// HandleTask handles generic tasks
func (ph *PlaybookHandler) HandleTask(task *Task) error {
	log.Printf("Executing task: IP=%s, Action=%s", task.IP, task.Action)

	switch task.Action {
	case "create_user":
		username, ok := task.Params["username"]
		if !ok {
			return fmt.Errorf("username parameter missing")
		}
		userID, ok := task.Params["user_id"]
		if !ok {
			return fmt.Errorf("user_id parameter missing")
		}

		return CreateUser(
			task.IP,
			task.Username,
			task.Password,
			task.SSHKey,
			username,
			userID,
		)

	default:
		return fmt.Errorf("unknown action: %s", task.Action)
	}
}

// ExampleMessageFormat returns example JSON formats for messages
func ExampleMessageFormat() {
	// Example 1: CreateUserPlaybook format
	example1 := CreateUserPlaybook{
		IP:       "10.1.1.1",
		Username: "user10",
		UserID:   "10",
		SSHUser:  "root",
		SSHPass:  "password",
	}
	json1, _ := json.MarshalIndent(example1, "", "  ")
	fmt.Println("Example CreateUserPlaybook format:")
	fmt.Println(string(json1))

	// Example 2: Task format
	example2 := Task{
		IP:     "10.1.1.1",
		Action: "create_user",
		Params: map[string]string{
			"username": "user10",
			"user_id":  "10",
		},
		Username: "root",
		Password: "password",
	}
	json2, _ := json.MarshalIndent(example2, "", "  ")
	fmt.Println("\nExample Task format:")
	fmt.Println(string(json2))
}
