package models

import "encoding/json"

// Vars represents variables for the create_user playbook.
// Variables are passed to the playbook via the --extra-vars flag.
type Vars struct {
	Username      string   `json:"username"`
	UserShell     string   `json:"user_shell,omitempty"`
	UserGroups    []string `json:"user_groups,omitempty"`
	PlainPassword string   `json:"plain_password,omitempty"`
}

// PasswordVars represents variables for the change_password playbook.
type PasswordVars struct {
	Username      string `json:"username"`
	PlainPassword string `json:"plain_password"`
}

// SudoBlacklistVars represents variables for the sudo_blacklist playbook.
type SudoBlacklistVars struct {
	Username          string   `json:"username"`
	BlacklistCommands []string `json:"blacklist_commands"`
}

// KafkaEvent represents the JSON structure received from Kafka
type KafkaEvent struct {
	EventType    string          `json:"event-type"`    // "Create", "Update", or "Delete"
	TargetServer []string        `json:"target-server"` // Array of server names
	OSType       string          `json:"OS-type"`       // "Linux" or "Windows"
	User         string          `json:"user"`          // e.g., "admin"
	Password     string          `json:"password"`      // e.g., "123"
	Playbook     string          `json:"playbook"`      // e.g., "template-playbook"
	Payload      json.RawMessage `json:"payload"`       // Dynamic value (can be any JSON)
	Timestamp    string          `json:"timestamp"`     // Timestamp string
}

// ServerInfo represents server information for adding to Ansible hosts file
type ServerInfo struct {
	Alias                 string `json:"alias"`                             // Server alias (e.g., "target1")
	AnsibleHost           string `json:"ansible_host"`                      // IP address or hostname
	AnsibleUser           string `json:"ansible_user"`                      // SSH user
	AnsiblePort           int    `json:"ansible_port"`                      // SSH port (default: 22)
	AnsiblePassword       string `json:"ansible_password,omitempty"`        // SSH password (optional)
	AnsibleBecomeMethod   string `json:"ansible_become_method,omitempty"`   // Become method: "su" or "sudo" (optional)
	AnsibleBecomePassword string `json:"ansible_become_password,omitempty"` // Sudo/su password (optional)
	Group                 string `json:"group,omitempty"`                   // Ansible group name (default: "targets")
}

// KafkaResponseEvent represents the response sent back to Kafka after processing an event
type KafkaResponseEvent struct {
	Status       string   `json:"status"`          // "success" or "error"
	EventType    string   `json:"event-type"`      // Original event type (e.g., "Create", "Delete")
	Playbook     string   `json:"playbook"`        // Original playbook name
	TargetServer []string `json:"target-server"`   // Servers that were processed
	Message      string   `json:"message"`         // Human-readable message
	Timestamp    string   `json:"timestamp"`       // Response timestamp
	Error        string   `json:"error,omitempty"` // Error message if status is "error"
}
