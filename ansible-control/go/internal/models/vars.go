package models

import (
	"encoding/json"
	"strconv"
	"time"
)

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

// UnmarshalJSON implements custom JSON unmarshaling to handle both camelCase and kebab-case field names,
// and to handle timestamp as both number (Unix timestamp) and string.
func (e *KafkaEvent) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to handle flexible field names
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Helper function to get string value from map, checking multiple keys
	getString := func(keys ...string) string {
		for _, key := range keys {
			if val, ok := raw[key]; ok && val != nil {
				if str, ok := val.(string); ok {
					return str
				}
			}
		}
		return ""
	}

	// Helper function to get string array from map, checking multiple keys
	getStringArray := func(keys ...string) []string {
		for _, key := range keys {
			if val, ok := raw[key]; ok && val != nil {
				if arr, ok := val.([]interface{}); ok {
					result := make([]string, 0, len(arr))
					for _, item := range arr {
						if str, ok := item.(string); ok {
							result = append(result, str)
						}
					}
					return result
				}
			}
		}
		return nil
	}

	// Helper function to get raw JSON for payload
	getRawJSON := func(key string) json.RawMessage {
		if val, ok := raw[key]; ok && val != nil {
			if data, err := json.Marshal(val); err == nil {
				return json.RawMessage(data)
			}
		}
		return nil
	}

	// Helper function to handle timestamp (number or string)
	getTimestamp := func() string {
		if val, ok := raw["timestamp"]; ok && val != nil {
			switch v := val.(type) {
			case string:
				return v
			case float64:
				// Unix timestamp (seconds with decimal precision)
				sec := int64(v)
				nsec := int64((v - float64(sec)) * 1e9)
				t := time.Unix(sec, nsec)
				return t.Format(time.RFC3339Nano)
			case int64:
				return time.Unix(v, 0).Format(time.RFC3339)
			case int:
				return time.Unix(int64(v), 0).Format(time.RFC3339)
			}
		}
		return ""
	}

	// Extract fields, checking both camelCase and kebab-case variants
	e.EventType = getString("event-type", "eventType", "EventType")
	e.TargetServer = getStringArray("target-server", "targetServer", "TargetServer")
	e.OSType = getString("OS-type", "osType", "OsType", "ostype")
	e.User = getString("user", "User")
	e.Password = getString("password", "Password")
	e.Playbook = getString("playbook", "Playbook")
	e.Timestamp = getTimestamp()

	// Handle payload - try to get it as raw JSON
	if payload := getRawJSON("payload"); payload != nil {
		e.Payload = payload
	}

	return nil
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

// UnmarshalJSON implements custom JSON unmarshaling to handle ansible_port as both string and number.
// This is necessary because Java Map<String, String> sends all values as strings.
func (s *ServerInfo) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to handle flexible field types
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Helper function to get string value
	getString := func(key string) string {
		if val, ok := raw[key]; ok && val != nil {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	// Helper function to get int value (handles both string and number)
	// This is critical for Java Map<String, String> which sends ports as strings
	getInt := func(key string, defaultValue int) int {
		if val, ok := raw[key]; ok && val != nil {
			switch v := val.(type) {
			case int:
				return v
			case int64:
				return int(v)
			case float64:
				// JSON numbers are unmarshaled as float64
				return int(v)
			case string:
				// Handle string format (from Java Map<String, String>)
				if parsed, err := strconv.Atoi(v); err == nil {
					return parsed
				}
			}
		}
		return defaultValue
	}

	// Extract all fields
	s.Alias = getString("alias")
	s.AnsibleHost = getString("ansible_host")
	s.AnsibleUser = getString("ansible_user")
	s.AnsiblePort = getInt("ansible_port", 22) // Handles both "22" (string) and 22 (number)
	s.AnsiblePassword = getString("ansible_password")
	s.AnsibleBecomeMethod = getString("ansible_become_method")
	s.AnsibleBecomePassword = getString("ansible_become_password")
	s.Group = getString("group")

	return nil
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
