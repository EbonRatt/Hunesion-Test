package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"test-kafka/internal/ansible"
	"test-kafka/internal/kafka"
	"test-kafka/internal/models"

	kafkago "github.com/segmentio/kafka-go"
)

// AnsibleEventHandler processes Kafka messages and routes them to Ansible playbooks
// based on the event-type field in the JSON
func AnsibleEventHandler(producer *kafka.Producer, responseTopic string) func(ctx context.Context, msg kafkago.Message) error {
	return func(ctx context.Context, msg kafkago.Message) error {
		// Step 1: Check if message is valid JSON
		msgValue := string(msg.Value)

		var event models.KafkaEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			// Not valid JSON - log the error and the message
			log.Printf("Failed to parse JSON: %v\n", err)
			log.Printf("Message length: %d bytes\n", len(msg.Value))

			// Show first 200 chars safely
			previewLen := 200
			if len(msgValue) < previewLen {
				previewLen = len(msgValue)
			}
			log.Printf("First %d chars: %q\n", previewLen, msgValue[:previewLen])

			// Try to find and report problematic characters (smart quotes, etc.)
			problematicChars := []struct {
				pos  int
				char rune
			}{}
			for i, r := range msgValue {
				// Check for smart quotes and other problematic Unicode characters
				if r == 0x201C || r == 0x201D || r == 0x2018 || r == 0x2019 ||
					(r > 127 && r < 256 && (r < 0x80 || r > 0x9F)) {
					problematicChars = append(problematicChars, struct {
						pos  int
						char rune
					}{i, r})
				}
			}

			if len(problematicChars) > 0 {
				log.Printf("Found %d potentially problematic characters:\n", len(problematicChars))
				for _, pc := range problematicChars {
					log.Printf("  Position %d: %q (U+%04X) - Replace with straight quote \"\n", pc.pos, pc.char, pc.char)
				}
			}

			log.Printf("Received plain text message (JSON parse failed): %s\n", msgValue)
			return nil // Don't treat as error, just skip it
		}

		// Step 2: Validate that it's actually a KafkaEvent (has required fields)
		if event.EventType == "" && event.Playbook == "" {
			// Looks like JSON but not our event format
			log.Printf("Received JSON but not a valid event: %s\n", string(msg.Value))
			return nil
		}

		// Step 3: Log the received event
		log.Printf("Received event: type=%s, playbook=%s, servers=%v, timestamp=%s\n",
			event.EventType, event.Playbook, event.TargetServer, event.Timestamp)
		log.Printf("Payload: %s\n", string(event.Payload))

		// Step 4: Route based on event-type
		var err error
		switch strings.ToUpper(event.EventType) {
		case "CREATE":
			err = HandleCreateEvent(ctx, &event, producer, responseTopic)
		case "UPDATE":
			err = HandleUpdateEvent(ctx, &event, producer, responseTopic)
		case "DELETE":
			err = HandleDeleteEvent(ctx, &event, producer, responseTopic)
		default:
			log.Printf("Unknown event-type: %s\n", event.EventType)
			err = fmt.Errorf("unknown event-type: %s", event.EventType)
		}

		// Send error response if processing failed
		if err != nil {
			SendResponse(ctx, producer, responseTopic, &event, "error", fmt.Sprintf("Failed to process %s event", event.EventType), err.Error())
		}

		return err
	}
}

// HandleCreateEvent processes "Create" events
func HandleCreateEvent(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing CREATE event for playbook: %s\n", event.Playbook)

	// Handle adding servers to hosts file
	if event.Playbook == "add_server" {
		return HandleAddServer(ctx, event, producer, responseTopic)
	}

	// Handle creating users
	if event.Playbook == "create_user" || event.Playbook == "create_user_file" {
		return HandleCreateUser(ctx, event, producer, responseTopic)
	}

	// Add more playbook handlers here...
	log.Printf("Unhandled playbook for CREATE: %s\n", event.Playbook)
	return nil
}

// HandleUpdateEvent processes "Update" events
func HandleUpdateEvent(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing UPDATE event for playbook: %s\n", event.Playbook)

	// Handle enabling users (re-enable after soft delete)
	if event.Playbook == "enable_user" || event.Playbook == "enable_user_file" {
		return HandleEnableUser(ctx, event, producer, responseTopic)
	}

	// Example: If playbook is "change_password", update password
	if event.Playbook == "change_password" {
		var passwordVars models.PasswordVars
		if err := json.Unmarshal(event.Payload, &passwordVars); err != nil {
			return fmt.Errorf("failed to parse change_password payload: %w", err)
		}

		// Save vars file
		varsDir := filepath.Join("internal", "ansible", "vars")
		if err := os.MkdirAll(varsDir, 0755); err != nil {
			return err
		}

		fileName := "job-change-password-" + passwordVars.Username + ".json"
		hostVarsPath := filepath.Join(varsDir, fileName)

		b, err := json.MarshalIndent(passwordVars, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(hostVarsPath, b, 0644); err != nil {
			return err
		}

		// Run playbook on target servers
		containerPath := "/work/internal/ansible/vars/" + fileName
		limit := strings.Join(event.TargetServer, ",")
		// Use relative path since working_dir in container is /work
		out, errOut, err := ansible.RunPlaybookChangePassword(ctx, "internal/ansible/playbooks/change_password.yml", containerPath, limit)

		log.Printf("UPDATE result - STDOUT: %s\n", out)
		if errOut != "" {
			log.Printf("UPDATE result - STDERR: %s\n", errOut)
		}

		// Check for "no hosts matched" or host not found errors (even if playbook didn't fail)
		combinedOutput := strings.ToLower(out + errOut)
		if strings.Contains(combinedOutput, "no hosts matched") ||
			strings.Contains(combinedOutput, "could not match supplied host pattern") ||
			strings.Contains(combinedOutput, "provided hosts list is empty") {
			errorMsg := fmt.Sprintf("Target server(s) not found in inventory. Please register the server(s) first: %v", event.TargetServer)
			log.Printf("CHANGE PASSWORD ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		// Only check for errors if the playbook actually failed
		if err != nil {
			// Check if user not found or other errors by examining the output
			usernameLower := strings.ToLower(passwordVars.Username)

			// Check for "user not found" or "does not exist" patterns (only in error context)
			userNotFoundPatterns := []string{
				"user '" + usernameLower + "' does not exist on",
				"failed: user '" + usernameLower + "' does not exist",
				"id: '" + usernameLower + "': no such user",
			}

			for _, pattern := range userNotFoundPatterns {
				if strings.Contains(combinedOutput, pattern) {
					errorMsg := fmt.Sprintf("User '%s' does not exist on target server(s)", passwordVars.Username)
					log.Printf("CHANGE PASSWORD ERROR: %s\n", errorMsg)
					// Return error - main handler will send response
					return fmt.Errorf(errorMsg)
				}
			}

			// Check for other specific failure patterns (avoid matching "failed=0" which means success)
			if strings.Contains(combinedOutput, "fatal:") ||
				strings.Contains(combinedOutput, "error!") ||
				strings.Contains(combinedOutput, "unreachable=") {
				errorMsg := fmt.Sprintf("Failed to change password for user '%s': %v", passwordVars.Username, err)
				log.Printf("CHANGE PASSWORD ERROR: %s\n", errorMsg)
				return fmt.Errorf(errorMsg)
			}

			log.Printf("UPDATE result - ERROR: %v\n", err)
			return err
		}

		// Success - clear password from vars file for security
		clearPasswordFromVarsFile(hostVarsPath, "plain_password")

		// Success - send success response
		successMsg := fmt.Sprintf("Successfully changed password for user '%s' on target server(s)", passwordVars.Username)
		log.Printf("CHANGE PASSWORD SUCCESS: %s\n", successMsg)
		SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
		return nil
	}

	log.Printf("Unhandled playbook for UPDATE: %s\n", event.Playbook)
	return nil
}

// HandleDeleteEvent processes "Delete" events
func HandleDeleteEvent(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing DELETE event for playbook: %s\n", event.Playbook)

	// Handle removing servers from hosts file (also handle typo "rmeove_server")
	if event.Playbook == "remove_server" || event.Playbook == "delete_server" || event.Playbook == "rmeove_server" {
		return HandleRemoveServer(ctx, event, producer, responseTopic)
	}

	// Handle deleting users
	if event.Playbook == "delete_user" || event.Playbook == "remove_user" {
		return HandleDeleteUser(ctx, event, producer, responseTopic)
	}

	// Add more delete handlers here...
	log.Printf("Unhandled playbook for DELETE: %s\n", event.Playbook)
	return nil
}
