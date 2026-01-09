package main

import (
	"encoding/json"
	"fmt"
)

// Task represents a task to be executed
type Task struct {
	IP       string            `json:"ip"`
	Action   string            `json:"action"`
	Params   map[string]string `json:"params"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	SSHKey   string            `json:"ssh_key,omitempty"`
}

// Playbook defines a set of tasks to execute
type Playbook struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tasks       []string `json:"tasks"`
}

// CreateUserPlaybook is a predefined playbook for creating users
type CreateUserPlaybook struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	SSHUser  string `json:"ssh_user,omitempty"`
	SSHPass  string `json:"ssh_pass,omitempty"`
	SSHKey   string `json:"ssh_key,omitempty"`
}

// ParseTaskFromKafkaMessage parses a Kafka message into a Task
func ParseTaskFromKafkaMessage(message []byte) (*Task, error) {
	var task Task
	err := json.Unmarshal(message, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// ParseCreateUserPlaybook parses a create user playbook from JSON
func ParseCreateUserPlaybook(message []byte) (*CreateUserPlaybook, error) {
	var playbook CreateUserPlaybook
	err := json.Unmarshal(message, &playbook)
	if err != nil {
		return nil, fmt.Errorf("failed to parse create user playbook: %w", err)
	}
	return &playbook, nil
}
