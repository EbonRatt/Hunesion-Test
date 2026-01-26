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

	// Validate username format (no spaces, not empty)
	if userVars.Username == "" {
		return fmt.Errorf("username is required and cannot be empty")
	}
	if strings.Contains(userVars.Username, " ") {
		return fmt.Errorf("username cannot contain spaces: '%s'", userVars.Username)
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

		// Check for connection/unreachable errors
		if strings.Contains(combinedOutput, "unreachable") ||
			strings.Contains(combinedOutput, "connection refused") ||
			strings.Contains(combinedOutput, "timed out") ||
			strings.Contains(combinedOutput, "no route to host") {
			errorMsg := fmt.Sprintf("Cannot connect to target server(s). Please check server connectivity and SSH configuration: %v", event.TargetServer)
			log.Printf("CREATE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		// Check for permission/authentication errors
		if strings.Contains(combinedOutput, "permission denied") ||
			strings.Contains(combinedOutput, "authentication failed") ||
			strings.Contains(combinedOutput, "access denied") {
			errorMsg := fmt.Sprintf("Permission denied. Please check SSH credentials and sudo/become configuration for server(s): %v", event.TargetServer)
			log.Printf("CREATE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		// Check for other specific failure patterns (avoid matching "failed=0" which means success)
		if strings.Contains(combinedOutput, "fatal:") ||
			strings.Contains(combinedOutput, "error!") {
			// Try to extract more specific error message from output
			errorMsg := fmt.Sprintf("Failed to create user '%s' on target server(s)", userVars.Username)

			// Look for more specific error messages in the output
			if strings.Contains(combinedOutput, "invalid username") {
				errorMsg = fmt.Sprintf("Invalid username format: '%s'", userVars.Username)
			} else if strings.Contains(combinedOutput, "invalid characters") {
				errorMsg = fmt.Sprintf("Username contains invalid characters: '%s'", userVars.Username)
			}

			log.Printf("CREATE USER ERROR: %s\n", errorMsg)
			log.Printf("Full output for debugging: %s\n", combinedOutput)
			return fmt.Errorf(errorMsg)
		}

		log.Printf("CREATE result - ERROR: %v\n", err)
		log.Printf("Full output for debugging: %s\n", combinedOutput)
		return fmt.Errorf("Failed to create user '%s': %v", userVars.Username, err)
	}

	// Success - clear password from vars file for security
	clearPasswordFromVarsFile(hostVarsPath, "plain_password")

	// Send success response
	successMsg := fmt.Sprintf("Successfully created user '%s' on target server(s)", userVars.Username)
	log.Printf("CREATE USER SUCCESS: %s\n", successMsg)
	SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
	return nil
}

// DeleteUserVars represents variables for the delete_user playbook
type DeleteUserVars struct {
	Username    string `json:"username"`
	RemoveHome  bool   `json:"remove_home,omitempty"`
	ForceRemove bool   `json:"force_remove,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle boolean fields as both string and bool.
// This is necessary because Java Map<String, String> sends all values as strings.
func (d *DeleteUserVars) UnmarshalJSON(data []byte) error {
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

	// Helper function to get bool value (handles both string and bool)
	// This is critical for Java Map<String, String> which sends booleans as strings
	getBool := func(key string, defaultValue bool) bool {
		if val, ok := raw[key]; ok && val != nil {
			switch v := val.(type) {
			case bool:
				return v
			case string:
				// Handle string format (from Java Map<String, String>)
				if v == "true" || v == "True" || v == "TRUE" {
					return true
				}
				if v == "false" || v == "False" || v == "FALSE" {
					return false
				}
			}
		}
		return defaultValue
	}

	// Extract all fields
	d.Username = getString("username")
	d.RemoveHome = getBool("remove_home", false)
	d.ForceRemove = getBool("force_remove", false)

	return nil
}

// HandleDeleteUser disables users on target servers (soft delete)
func HandleDeleteUser(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing DELETE USER event (soft delete - disable)\n")

	// Parse payload to get username
	var deleteUserVars DeleteUserVars

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

	log.Printf("DISABLE USER result - STDOUT: %s\n", out)
	if errOut != "" {
		log.Printf("DISABLE USER result - STDERR: %s\n", errOut)
	}

	// Check for "no hosts matched" or host not found errors (even if playbook didn't fail)
	combinedOutput := strings.ToLower(out + errOut)
	if strings.Contains(combinedOutput, "no hosts matched") ||
		strings.Contains(combinedOutput, "could not match supplied host pattern") ||
		strings.Contains(combinedOutput, "provided hosts list is empty") {
		errorMsg := fmt.Sprintf("Target server(s) not found in inventory. Please register the server(s) first: %v", event.TargetServer)
		log.Printf("DISABLE USER ERROR: %s\n", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Only check for errors if the playbook actually failed
	if err != nil {
		// Check if user not found or other errors by examining the output
		usernameLower := strings.ToLower(deleteUserVars.Username)

		// Check for "user not found" or "does not exist" patterns (only in error context)
		userNotFoundPatterns := []string{
			"user '" + usernameLower + "' does not exist on",
			"failed: user '" + usernameLower + "' does not exist",
			"id: '" + usernameLower + "': no such user",
		}

		for _, pattern := range userNotFoundPatterns {
			if strings.Contains(combinedOutput, pattern) {
				errorMsg := fmt.Sprintf("User '%s' does not exist on target server(s)", deleteUserVars.Username)
				log.Printf("DISABLE USER ERROR: %s\n", errorMsg)
				// Return error - main handler will send response
				return fmt.Errorf(errorMsg)
			}
		}

		// Check for other specific failure patterns (avoid matching "failed=0" which means success)
		if strings.Contains(combinedOutput, "fatal:") ||
			strings.Contains(combinedOutput, "error!") ||
			strings.Contains(combinedOutput, "unreachable=") {
			errorMsg := fmt.Sprintf("Failed to disable user '%s': %v", deleteUserVars.Username, err)
			log.Printf("DISABLE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		log.Printf("DISABLE USER result - ERROR: %v\n", err)
		return err
	}

	// Success - clear sensitive data from vars file for security
	clearPasswordFromVarsFile(hostVarsPath, "")

	// Send success response
	successMsg := fmt.Sprintf("Successfully disabled user '%s' on target server(s) (soft delete)", deleteUserVars.Username)
	log.Printf("DISABLE USER SUCCESS: %s\n", successMsg)
	SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
	return nil
}

// HandleEnableUser enables users on target servers (re-enable after soft delete)
func HandleEnableUser(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Processing ENABLE USER event\n")

	// Parse payload to get username and optional shell
	var enableUserVars struct {
		Username string `json:"username"`
		Shell    string `json:"shell,omitempty"` // Optional shell to restore (default: /bin/bash)
	}

	if len(event.Payload) == 0 {
		return fmt.Errorf("payload is required for enable_user playbook")
	}

	if err := json.Unmarshal(event.Payload, &enableUserVars); err != nil {
		return fmt.Errorf("failed to parse enable_user payload: %w", err)
	}

	if enableUserVars.Username == "" {
		return fmt.Errorf("username is required in payload")
	}

	// Set default shell if not provided
	if enableUserVars.Shell == "" {
		enableUserVars.Shell = "/bin/bash"
	}

	// Save vars file
	varsDir := filepath.Join("internal", "ansible", "vars")
	if err := os.MkdirAll(varsDir, 0755); err != nil {
		return err
	}

	fileName := "job-enable-users.json"
	hostVarsPath := filepath.Join(varsDir, fileName)

	b, err := json.MarshalIndent(enableUserVars, "", "  ")
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
	out, errOut, err := ansible.RunPlaybookEnableUser(ctx, "internal/ansible/playbooks/enable_user.yml", containerPath, limit)

	log.Printf("ENABLE USER result - STDOUT: %s\n", out)
	if errOut != "" {
		log.Printf("ENABLE USER result - STDERR: %s\n", errOut)
	}

	// Check for "no hosts matched" or host not found errors (even if playbook didn't fail)
	combinedOutput := strings.ToLower(out + errOut)
	if strings.Contains(combinedOutput, "no hosts matched") ||
		strings.Contains(combinedOutput, "could not match supplied host pattern") ||
		strings.Contains(combinedOutput, "provided hosts list is empty") {
		errorMsg := fmt.Sprintf("Target server(s) not found in inventory. Please register the server(s) first: %v", event.TargetServer)
		log.Printf("ENABLE USER ERROR: %s\n", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Only check for errors if the playbook actually failed
	if err != nil {
		// Check if user not found or other errors by examining the output
		usernameLower := strings.ToLower(enableUserVars.Username)

		// Check for "user not found" or "does not exist" patterns (only in error context)
		userNotFoundPatterns := []string{
			"user '" + usernameLower + "' does not exist on",
			"failed: user '" + usernameLower + "' does not exist",
			"id: '" + usernameLower + "': no such user",
		}

		for _, pattern := range userNotFoundPatterns {
			if strings.Contains(combinedOutput, pattern) {
				errorMsg := fmt.Sprintf("User '%s' does not exist on target server(s)", enableUserVars.Username)
				log.Printf("ENABLE USER ERROR: %s\n", errorMsg)
				return fmt.Errorf(errorMsg)
			}
		}

		// Check for other specific failure patterns (avoid matching "failed=0" which means success)
		if strings.Contains(combinedOutput, "fatal:") ||
			strings.Contains(combinedOutput, "error!") ||
			strings.Contains(combinedOutput, "unreachable=") {
			errorMsg := fmt.Sprintf("Failed to enable user '%s': %v", enableUserVars.Username, err)
			log.Printf("ENABLE USER ERROR: %s\n", errorMsg)
			return fmt.Errorf(errorMsg)
		}

		log.Printf("ENABLE USER result - ERROR: %v\n", err)
		return err
	}

	// Success - clear sensitive data from vars file for security
	clearPasswordFromVarsFile(hostVarsPath, "")

	// Send success response
	successMsg := fmt.Sprintf("Successfully enabled user '%s' on target server(s)", enableUserVars.Username)
	log.Printf("ENABLE USER SUCCESS: %s\n", successMsg)
	SendResponse(ctx, producer, responseTopic, event, "success", successMsg, "")
	return nil
}
