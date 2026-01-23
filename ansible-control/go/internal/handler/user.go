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
)

// clearPasswordFromVarsFile clears password fields from vars JSON file for security
func clearPasswordFromVarsFile(filePath, passwordField string) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Warning: Failed to read vars file for password clearing: %v\n", err)
		return
	}

	// Parse JSON
	var varsData map[string]interface{}
	if err := json.Unmarshal(data, &varsData); err != nil {
		log.Printf("Warning: Failed to parse vars file for password clearing: %v\n", err)
		return
	}

	// Clear password fields
	if passwordField != "" {
		if _, exists := varsData[passwordField]; exists {
			varsData[passwordField] = ""
		}
	} else {
		// Clear all password-related fields
		passwordFields := []string{"plain_password", "password", "ansible_password", "ansible_become_password"}
		for _, field := range passwordFields {
			if _, exists := varsData[field]; exists {
				varsData[field] = ""
			}
		}
	}

	// Write back to file
	clearedData, err := json.MarshalIndent(varsData, "", "  ")
	if err != nil {
		log.Printf("Warning: Failed to marshal cleared vars data: %v\n", err)
		return
	}

	if err := os.WriteFile(filePath, clearedData, 0644); err != nil {
		log.Printf("Warning: Failed to write cleared vars file: %v\n", err)
		return
	}

	log.Printf("Cleared password fields from vars file: %s\n", filePath)
}

// HandleCreateUser processes create_user playbook events
func HandleCreateUser(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	var userVars models.Vars
	if err := json.Unmarshal(event.Payload, &userVars); err != nil {
		return fmt.Errorf("failed to parse create_user payload: %w", err)
	}

	// Save vars file using fixed template name
	varsDir := filepath.Join("internal", "ansible", "vars")
	if err := os.MkdirAll(varsDir, 0755); err != nil {
		return err
	}

	// Use fixed template file name: job-create-users.json
	fileName := "job-create-users.json"
	hostVarsPath := filepath.Join(varsDir, fileName)

	b, err := json.MarshalIndent(userVars, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(hostVarsPath, b, 0644); err != nil {
		return err
	}

	log.Printf("Created/Updated vars file: %s\n", hostVarsPath)

	// Check if we should skip playbook execution (file-only mode)
	// Skip if target-server is empty or playbook name contains "_file"
	skipExecution := len(event.TargetServer) == 0 || strings.Contains(event.Playbook, "_file")

	if skipExecution {
		log.Printf("Skipping playbook execution (file-only mode)\n")
		return nil
	}

	// Run playbook on target servers
	containerPath := "/work/internal/ansible/vars/" + fileName
	// Use target-server from event to limit which servers to run on
	limit := strings.Join(event.TargetServer, ",")
	// Use relative path since working_dir in container is /work
	out, errOut, err := ansible.RunPlaybookCreateUser(ctx, "internal/ansible/playbooks/create_user.yml", containerPath, limit)

	log.Printf("CREATE result - STDOUT: %s\n", out)
	if errOut != "" {
		log.Printf("CREATE result - STDERR: %s\n", errOut)
	}

	// Check for "no hosts matched" or host not found errors (even if playbook didn't fail)
	combinedOutput := strings.ToLower(out + errOut)
	if strings.Contains(combinedOutput, "no hosts matched") ||
		strings.Contains(combinedOutput, "could not match supplied host pattern") ||
		strings.Contains(combinedOutput, "provided hosts list is empty") {
		errorMsg := fmt.Sprintf("Target server(s) not found in inventory. Please register the server(s) first: %v", event.TargetServer)
		log.Printf("CREATE USER ERROR: %s\n", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Only check for errors if the playbook actually failed
	if err != nil {
		// Check if user already exists by examining the output
		usernameLower := strings.ToLower(userVars.Username)

		// Check for specific "user already exists" error patterns (only in error context)
		userExistsPatterns := []string{
			"user '" + usernameLower + "' already exists on",
			"useradd: user '" + usernameLower + "' already exists",
			"failed: user '" + usernameLower + "' already exists",
		}

		for _, pattern := range userExistsPatterns {
			if strings.Contains(combinedOutput, pattern) {
				errorMsg := fmt.Sprintf("User '%s' already exists on target server(s)", userVars.Username)
				log.Printf("CREATE USER ERROR: %s\n", errorMsg)
				// Return error - main handler will send response
				return fmt.Errorf(errorMsg)
			}
		}

		// Check for other specific failure patterns (avoid matching "failed=0" which means success)
		if strings.Contains(combinedOutput, "fatal:") ||
			strings.Contains(combinedOutput, "error!") ||
			strings.Contains(combinedOutput, "unreachable=") {
			errorMsg := fmt.Sprintf("Failed to create user '%s': %v", userVars.Username, err)
			log.Printf("CREATE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		log.Printf("CREATE result - ERROR: %v\n", err)
		return err
	}

	// Success - clear password from vars file for security
	clearPasswordFromVarsFile(hostVarsPath, "plain_password")

	// Send success response
	successMsg := fmt.Sprintf("Successfully created user '%s' on target server(s)", userVars.Username)
	log.Printf("CREATE USER SUCCESS: %s\n", successMsg)
	SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
	return nil
}

// HandleDeleteUser deletes users from target servers
func HandleDeleteUser(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing DELETE USER event\n")

	// Parse payload to get username
	var deleteUserVars struct {
		Username    string `json:"username"`
		RemoveHome  bool   `json:"remove_home,omitempty"`
		ForceRemove bool   `json:"force_remove,omitempty"`
	}

	if len(event.Payload) == 0 {
		return fmt.Errorf("payload is required for delete_user playbook")
	}

	if err := json.Unmarshal(event.Payload, &deleteUserVars); err != nil {
		return fmt.Errorf("failed to parse delete_user payload: %w", err)
	}

	if deleteUserVars.Username == "" {
		return fmt.Errorf("username is required in payload")
	}

	// Save vars file using fixed template name
	varsDir := filepath.Join("internal", "ansible", "vars")
	if err := os.MkdirAll(varsDir, 0755); err != nil {
		return err
	}

	// Use fixed template file name: job-delete-users.json
	fileName := "job-delete-users.json"
	hostVarsPath := filepath.Join(varsDir, fileName)

	b, err := json.MarshalIndent(deleteUserVars, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(hostVarsPath, b, 0644); err != nil {
		return err
	}

	log.Printf("Created/Updated vars file: %s\n", hostVarsPath)

	// Check if we should skip playbook execution (file-only mode)
	skipExecution := len(event.TargetServer) == 0 || strings.Contains(event.Playbook, "_file")

	if skipExecution {
		log.Printf("Skipping playbook execution (file-only mode)\n")
		return nil
	}

	// Run playbook on target servers
	containerPath := "/work/internal/ansible/vars/" + fileName
	limit := strings.Join(event.TargetServer, ",")
	// Use relative path since working_dir in container is /work
	out, errOut, err := ansible.RunPlaybookDeleteUser(ctx, "internal/ansible/playbooks/delete_user.yml", containerPath, limit)

	log.Printf("DELETE USER result - STDOUT: %s\n", out)
	if errOut != "" {
		log.Printf("DELETE USER result - STDERR: %s\n", errOut)
	}

	// Check for "no hosts matched" or host not found errors (even if playbook didn't fail)
	combinedOutput := strings.ToLower(out + errOut)
	if strings.Contains(combinedOutput, "no hosts matched") ||
		strings.Contains(combinedOutput, "could not match supplied host pattern") ||
		strings.Contains(combinedOutput, "provided hosts list is empty") {
		errorMsg := fmt.Sprintf("Target server(s) not found in inventory. Please register the server(s) first: %v", event.TargetServer)
		log.Printf("DELETE USER ERROR: %s\n", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Only check for errors if the playbook actually failed
	if err != nil {
		// Check if user not found or other errors by examining the output
		usernameLower := strings.ToLower(deleteUserVars.Username)

		// Check for "user not found" or "does not exist" patterns (only in error context)
		userNotFoundPatterns := []string{
			"user '" + usernameLower + "' does not exist on",
			"userdel: user '" + usernameLower + "' does not exist",
			"failed: user '" + usernameLower + "' does not exist",
			"id: '" + usernameLower + "': no such user",
		}

		for _, pattern := range userNotFoundPatterns {
			if strings.Contains(combinedOutput, pattern) {
				errorMsg := fmt.Sprintf("User '%s' does not exist on target server(s)", deleteUserVars.Username)
				log.Printf("DELETE USER ERROR: %s\n", errorMsg)
				// Return error - main handler will send response
				return fmt.Errorf(errorMsg)
			}
		}

		// Check for other specific failure patterns (avoid matching "failed=0" which means success)
		if strings.Contains(combinedOutput, "fatal:") ||
			strings.Contains(combinedOutput, "error!") ||
			strings.Contains(combinedOutput, "unreachable=") {
			errorMsg := fmt.Sprintf("Failed to delete user '%s': %v", deleteUserVars.Username, err)
			log.Printf("DELETE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		log.Printf("DELETE USER result - ERROR: %v\n", err)
		return err
	}

	// Success - clear sensitive data from vars file for security
	clearPasswordFromVarsFile(hostVarsPath, "")

	// Send success response
	successMsg := fmt.Sprintf("Successfully deleted user '%s' from target server(s)", deleteUserVars.Username)
	log.Printf("DELETE USER SUCCESS: %s\n", successMsg)
	SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
	return nil
}
