package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"test-kafka/internal/inventory"
	"test-kafka/internal/kafka"
	"test-kafka/internal/models"
)

// HandleAddServer adds target servers to the Ansible hosts file
func HandleAddServer(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	log.Printf("Adding servers to Ansible hosts file: %v\n", event.TargetServer)
	log.Printf("Payload: %s\n", string(event.Payload))

	// Get the hosts file path from internal structure
	invPath := filepath.Join("internal", "inventory", "files", "hosts.ini")

	// Get absolute path
	absInvPath, err := filepath.Abs(invPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Printf("Hosts file path: %s\n", absInvPath)

	// Default group name
	groupName := "targets"
	if event.Playbook == "add_server" {
		// Try to get group from payload if specified
		var serverInfo models.ServerInfo
		if err := json.Unmarshal(event.Payload, &serverInfo); err == nil && serverInfo.Group != "" {
			groupName = serverInfo.Group
		}
	}

	// Process each target server
	for _, serverAlias := range event.TargetServer {
		// Parse server info from payload or use event fields
		var serverInfo models.ServerInfo

		// Try to parse payload as ServerInfo
		if len(event.Payload) > 0 {
			if err := json.Unmarshal(event.Payload, &serverInfo); err != nil {
				log.Printf("Warning: Could not parse payload as ServerInfo, using event fields: %v\n", err)
				// Fallback: use event fields to construct server info (handle "null" strings)
				user := event.User
				if user == "null" || user == "" {
					user = ""
				}
				password := event.Password
				if password == "null" || password == "" {
					password = ""
				}
				serverInfo = models.ServerInfo{
					Alias:           serverAlias,
					AnsibleHost:     serverAlias, // Use alias as host if not specified
					AnsibleUser:     user,
					AnsiblePort:     22, // Default SSH port
					AnsiblePassword: password,
				}
			} else {
				// Payload parsed successfully - use target server alias (overrides payload alias)
				serverInfo.Alias = serverAlias
				// If ansible_host not specified, use alias
				if serverInfo.AnsibleHost == "" {
					serverInfo.AnsibleHost = serverAlias
				}
			}
		} else {
			// No payload, use event fields (handle "null" strings)
			user := event.User
			if user == "null" || user == "" {
				user = ""
			}
			password := event.Password
			if password == "null" || password == "" {
				password = ""
			}
			serverInfo = models.ServerInfo{
				Alias:           serverAlias,
				AnsibleHost:     serverAlias,
				AnsibleUser:     user,
				AnsiblePort:     22,
				AnsiblePassword: password,
			}
		}

		// Set default port if not provided
		if serverInfo.AnsiblePort == 0 {
			serverInfo.AnsiblePort = 22
		}

		// Build the host line for Ansible hosts.ini
		// Start with required fields
		hostLine := fmt.Sprintf("%s ansible_host=%s ansible_port=%d",
			serverInfo.Alias, serverInfo.AnsibleHost, serverInfo.AnsiblePort)

		// Add ansible_user only if provided
		if serverInfo.AnsibleUser != "" && serverInfo.AnsibleUser != "null" {
			hostLine += fmt.Sprintf(" ansible_user=%s", serverInfo.AnsibleUser)
		}

		// Add password if provided (skip if "null")
		if serverInfo.AnsiblePassword != "" && serverInfo.AnsiblePassword != "null" {
			hostLine += fmt.Sprintf(" ansible_password=%s", serverInfo.AnsiblePassword)
		}

		// Add become password if provided (skip if "null")
		if serverInfo.AnsibleBecomePassword != "" && serverInfo.AnsibleBecomePassword != "null" {
			hostLine += fmt.Sprintf(" ansible_become_password=%s", serverInfo.AnsibleBecomePassword)
		}

		// Add become method if provided (skip if "null")
		if serverInfo.AnsibleBecomeMethod != "" && serverInfo.AnsibleBecomeMethod != "null" {
			hostLine += fmt.Sprintf(" ansible_become_method=%s", serverInfo.AnsibleBecomeMethod)
		}

		// Use group from serverInfo if specified, otherwise use default
		group := groupName
		if serverInfo.Group != "" {
			group = serverInfo.Group
		}

		// Add/update the host in the hosts file
		log.Printf("Adding server to hosts file: %s in group [%s]\n", serverInfo.Alias, group)
		log.Printf("Host line: %s\n", hostLine)
		if err := inventory.UpsertHostInGroup(absInvPath, group, serverInfo.Alias, hostLine); err != nil {
			log.Printf("Error adding server %s to hosts file: %v\n", serverInfo.Alias, err)
			return fmt.Errorf("failed to add server %s: %w", serverInfo.Alias, err)
		}

		log.Printf("Successfully added server %s to hosts file\n", serverInfo.Alias)
	}

	// Send success response
	message := fmt.Sprintf("Successfully registered %d server(s) to Ansible hosts file", len(event.TargetServer))
	SendResponse(ctx, producer, responseTopic, event, "success", message, "")

	return nil
}

// HandleRemoveServer removes target servers from the Ansible hosts file
func HandleRemoveServer(ctx context.Context, event *models.KafkaEvent, producer *kafka.Producer, responseTopic string) error {
	// Get the hosts file path from internal structure
	invPath := filepath.Join("internal", "inventory", "files", "hosts.ini")

	// Get absolute path
	absInvPath, err := filepath.Abs(invPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Printf("Hosts file path: %s\n", absInvPath)

	// Parse payload to get alias and group
	var payloadData map[string]interface{}
	groupName := "targets"
	var aliasesToRemove []string

	if len(event.Payload) > 0 {
		if err := json.Unmarshal(event.Payload, &payloadData); err == nil {
			// Get group from payload
			if group, ok := payloadData["group"].(string); ok && group != "" {
				groupName = group
			}
			// Get alias from payload
			if alias, ok := payloadData["alias"].(string); ok && alias != "" {
				aliasesToRemove = append(aliasesToRemove, alias)
			}
		}
	}

	// Fall back to target-server if no alias in payload
	if len(aliasesToRemove) == 0 {
		if len(event.TargetServer) > 0 {
			aliasesToRemove = event.TargetServer
			log.Printf("Using target-server: %v\n", aliasesToRemove)
		} else {
			return fmt.Errorf("alias is required (either in payload or target-server)")
		}
	}

	log.Printf("Removing servers from Ansible hosts file: %v in group [%s]\n", aliasesToRemove, groupName)

	// Process each server to remove
	for _, serverAlias := range aliasesToRemove {
		log.Printf("Removing server %s from group [%s]\n", serverAlias, groupName)
		if err := inventory.RemoveHostFromGroup(absInvPath, groupName, serverAlias); err != nil {
			log.Printf("Error removing server %s from hosts file: %v\n", serverAlias, err)
			return fmt.Errorf("failed to remove server %s: %w", serverAlias, err)
		}
		log.Printf("Successfully removed server %s from hosts file\n", serverAlias)
	}

	// Send success response
	message := fmt.Sprintf("Successfully removed %d server(s) from Ansible hosts file", len(aliasesToRemove))
	SendResponse(ctx, producer, responseTopic, event, "success", message, "")

	return nil
}
