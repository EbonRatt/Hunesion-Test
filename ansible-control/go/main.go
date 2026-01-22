package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"test-kafka/internal/kafka"
	"test-kafka/internal/models"

	kafkago "github.com/segmentio/kafka-go"
)

func runPlaybook(ctx context.Context, playbookPath string) (string, string, error) {
	args := []string{"compose", "exec", "-T", "ansible", "ansible-playbook", playbookPath}
	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func runPlaybookCreateUser(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	args := []string{
		"compose", "exec", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "/work/inventory/hosts.ini",
		"-e", "@" + extraVarsFile,
	}
	if strings.TrimSpace(limit) != "" {
		args = append(args, "-l", limit)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Set working directory to project root (where docker-compose.yml should be)
	// Get absolute path to ensure docker compose can find the file
	wd, wdErr := os.Getwd()
	if wdErr == nil {
		// If we're in ansible-control/go/, go up 2 levels
		if strings.Contains(wd, "ansible-control") {
			projectRoot := filepath.Join(wd, "..", "..")
			if absPath, absErr := filepath.Abs(projectRoot); absErr == nil {
				cmd.Dir = absPath
				log.Printf("Running docker compose from: %s\n", absPath)
			}
		} else {
			// Already in project root or different structure
			cmd.Dir = wd
			log.Printf("Running docker compose from: %s\n", wd)
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func runPlaybookChangePassword(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	args := []string{
		"compose", "exec", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "/work/inventory/hosts.ini",
		"-e", "@" + extraVarsFile,
	}
	if strings.TrimSpace(limit) != "" {
		args = append(args, "-l", limit)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Set working directory to project root (where docker-compose.yml should be)
	// Get absolute path to ensure docker compose can find the file
	wd, wdErr := os.Getwd()
	if wdErr == nil {
		// If we're in ansible-control/go/, go up 2 levels
		if strings.Contains(wd, "ansible-control") {
			projectRoot := filepath.Join(wd, "..", "..")
			if absPath, absErr := filepath.Abs(projectRoot); absErr == nil {
				cmd.Dir = absPath
				log.Printf("Running docker compose from: %s\n", absPath)
			}
		} else {
			// Already in project root or different structure
			cmd.Dir = wd
			log.Printf("Running docker compose from: %s\n", wd)
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func upsertHostInGroup(filePath, group, alias, hostLine string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var out []string
	grpHeader := "[" + group + "]"
	inGroup := false
	found := false

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		// section handling
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			// leaving group: if we didn't find alias, insert before next section
			if inGroup && !found {
				out = append(out, hostLine)
				found = true
			}
			inGroup = (trim == grpHeader)
			out = append(out, line)
			continue
		}

		// inside group: replace if alias exists
		if inGroup {
			if strings.HasPrefix(trim, alias+" ") || trim == alias {
				out = append(out, hostLine)
				found = true
				continue // skip old line
			}
		}

		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		return err
	}

	// group not found: append group at end
	if !strings.Contains(strings.Join(out, "\n"), grpHeader) {
		out = append(out, "", grpHeader, hostLine, "")
	} else if inGroup && !found {
		// file ended while still in group
		out = append(out, hostLine)
	}

	return os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644)
}

// func main() {

// 	// Create user on target1

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	defer cancel()

// 	// Get the project root (ansible directory)
// 	// Assuming main.go is in ansible-control/go/
// 	// Project root is 2 levels up: ../../ansible
// 	projectRoot := filepath.Join("..", "..", "ansible")
// 	varsDir := filepath.Join(projectRoot, "vars")
// 	hostVarsPath := filepath.Join(varsDir, "job-123.json")

// 	// Create directory
// 	_ = os.MkdirAll(varsDir, 0755)

// 	v := Vars{
// 		Username:      "hello",
// 		UserShell:     "/bin/bash",
// 		UserGroups:    []string{"sudo"},
// 		PlainPassword: "800",
// 	}
// 	b, _ := json.MarshalIndent(v, "", "  ")
// 	_ = os.WriteFile(hostVarsPath, b, 0644)

// 	// Inside container: /work/vars/job-123.json
// 	containerPath := "/work/vars/job-123.json"
// 	out, errOut, err := runPlaybookCreateUser(ctx, "/work/playbooks/create_user.yml", containerPath)

// 	fmt.Println("STDOUT:\n", out)
// 	fmt.Println("STDERR:\n", errOut)
// 	fmt.Println("ERR:\n", err)

// 	// Add host to group
// 	// invPath := filepath.Join("..", "..", "ansible", "inventory", "hosts.ini")

// 	// hostLine := "target3 ansible_host=192.168.230.128 ansible_user=ebon ansible_password=123 ansible_port=22"
// 	// if err := upsertHostInGroup(invPath, "targets", "target3", hostLine); err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println("Updated inventory:", invPath)

// 	// Define Path From Project Root
// 	// projectRoot := filepath.Join("..", "..", "ansible")
// 	// varsDir := filepath.Join(projectRoot, "vars")
// 	// hostVarsPath := filepath.Join(varsDir, "job-change-password.json")

// 	// userPassword := PasswordVars{
// 	// 	Username:      "ebon",
// 	// 	PlainPassword: "800",
// 	// }

// 	// Create directory
// 	// b, _ := json.MarshalIndent(userPassword, "", "  ")
// 	// _ = os.WriteFile(hostVarsPath, b, 0644)

// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	// defer cancel()

// 	// containerPath := "/work/vars/job-change-password.json"
// 	// out, errOut, err := runPlaybookChangePassword(ctx, "/work/playbooks/change_password.yml", containerPath, "target3")

// 	// fmt.Println("STDOUT:\n", out)
// 	// fmt.Println("STDERR:\n", errOut)
// 	// fmt.Println("ERR:\n", err)

// 	// Apply sudo blacklist
// 	// ctx := context.Background()
// 	// out, errOut, err := applySudoCommandBlacklist(
// 	// 	ctx,
// 	// 	"ebon",
// 	// 	[]string{"/usr/bin/rm", "/usr/bin/cat", "/usr/bin/touch", "/usr/bin/cp"},
// 	// 	"target3", // limit host; use "" to apply to all in group
// 	// )
// 	// fmt.Println(out, errOut, err)
// }

// AnsibleEventHandler processes Kafka messages and routes them to Ansible playbooks
// based on the event-type field in the JSON
func AnsibleEventHandler(ctx context.Context, msg kafkago.Message) error {
	// Step 1: Check if message is valid JSON
	var event models.KafkaEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		// Not valid JSON - treat as plain text message
		log.Printf("Received plain text message: %s\n", string(msg.Value))
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

	// Step 4: Route based on event-type
	switch strings.ToUpper(event.EventType) {
	case "CREATE":
		return handleCreateEvent(ctx, &event)
	case "UPDATE":
		return handleUpdateEvent(ctx, &event)
	case "DELETE":
		return handleDeleteEvent(ctx, &event)
	default:
		log.Printf("Unknown event-type: %s\n", event.EventType)
		return fmt.Errorf("unknown event-type: %s", event.EventType)
	}
}

// handleCreateEvent processes "Create" events
func handleCreateEvent(ctx context.Context, event *models.KafkaEvent) error {
	log.Printf("Processing CREATE event for playbook: %s\n", event.Playbook)

	// Handle adding servers to hosts file
	if event.Playbook == "add_server" {
		return handleAddServer(ctx, event)
	}

	// Example: If playbook is "create_user" or "create_user_file", parse payload and create user
	if event.Playbook == "create_user" || event.Playbook == "create_user_file" {
		var userVars models.Vars
		if err := json.Unmarshal(event.Payload, &userVars); err != nil {
			return fmt.Errorf("failed to parse create_user payload: %w", err)
		}

		// Save vars file using fixed template name
		projectRoot := filepath.Join("..", "..", "ansible")
		varsDir := filepath.Join(projectRoot, "vars")
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
		containerPath := "/work/vars/" + fileName
		// Use target-server from event to limit which servers to run on
		limit := strings.Join(event.TargetServer, ",")
		out, errOut, err := runPlaybookCreateUser(ctx, "/work/playbooks/create_user.yml", containerPath, limit)

		log.Printf("CREATE result - STDOUT: %s\n", out)
		if errOut != "" {
			log.Printf("CREATE result - STDERR: %s\n", errOut)
		}
		if err != nil {
			log.Printf("CREATE result - ERROR: %v\n", err)
			return err
		}
		return nil
	}

	// Add more playbook handlers here...
	log.Printf("Unhandled playbook for CREATE: %s\n", event.Playbook)
	return nil
}

// handleAddServer adds target servers to the Ansible hosts file
func handleAddServer(ctx context.Context, event *models.KafkaEvent) error {
	log.Printf("Adding servers to Ansible hosts file: %v\n", event.TargetServer)

	// Get the hosts file path
	projectRoot := filepath.Join("..", "..", "ansible")
	invPath := filepath.Join(projectRoot, "inventory", "hosts.ini")

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

		// Use group from serverInfo if specified, otherwise use default
		group := groupName
		if serverInfo.Group != "" {
			group = serverInfo.Group
		}

		// Add/update the host in the hosts file
		log.Printf("Adding server to hosts file: %s in group [%s]\n", serverInfo.Alias, group)
		if err := upsertHostInGroup(invPath, group, serverInfo.Alias, hostLine); err != nil {
			log.Printf("Error adding server %s to hosts file: %v\n", serverInfo.Alias, err)
			return fmt.Errorf("failed to add server %s: %w", serverInfo.Alias, err)
		}

		log.Printf("Successfully added server %s to hosts file\n", serverInfo.Alias)
	}

	return nil
}

// handleUpdateEvent processes "Update" events
func handleUpdateEvent(ctx context.Context, event *models.KafkaEvent) error {
	log.Printf("Processing UPDATE event for playbook: %s\n", event.Playbook)

	// Example: If playbook is "change_password", update password
	if event.Playbook == "change_password" {
		var passwordVars models.PasswordVars
		if err := json.Unmarshal(event.Payload, &passwordVars); err != nil {
			return fmt.Errorf("failed to parse change_password payload: %w", err)
		}

		// Save vars file
		projectRoot := filepath.Join("..", "..", "ansible")
		varsDir := filepath.Join(projectRoot, "vars")
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
		containerPath := "/work/vars/" + fileName
		limit := strings.Join(event.TargetServer, ",")
		out, errOut, err := runPlaybookChangePassword(ctx, "/work/playbooks/change_password.yml", containerPath, limit)

		log.Printf("UPDATE result - STDOUT: %s\n", out)
		if errOut != "" {
			log.Printf("UPDATE result - STDERR: %s\n", errOut)
		}
		if err != nil {
			log.Printf("UPDATE result - ERROR: %v\n", err)
			return err
		}
		return nil
	}

	log.Printf("Unhandled playbook for UPDATE: %s\n", event.Playbook)
	return nil
}

// handleDeleteEvent processes "Delete" events
func handleDeleteEvent(ctx context.Context, event *models.KafkaEvent) error {
	log.Printf("Processing DELETE event for playbook: %s\n", event.Playbook)
	// Add your delete logic here
	// For example: delete user, remove sudo blacklist, etc.
	return nil
}

func main() {
	// Create Kafka consumer with custom handler
	config := kafka.Config{
		Broker: "localhost:29092",
		Topic:  "test",
	}

	consumer := kafka.NewConsumer(config, AnsibleEventHandler)

	// Start consuming (blocks until stopped)
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
