package main

import (
	"context" // For cancellation
	"encoding/json"
	"fmt" // For printing to console
	"log" // For logging errors
	"os"  // For OS signals
	"os/signal"
	"syscall"

	"bytes"
	"os/exec"
	"path/filepath"

	"bufio"
	"strings"

	"github.com/segmentio/kafka-go"
	// Kafka library
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

func runPlaybookSudoBlacklist(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
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

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func applySudoCommandBlacklist(ctx context.Context, username string, blacklistCommands []string, limit string) (string, string, error) {
	v := SudoBlacklistVars{
		Username:          username,
		BlacklistCommands: blacklistCommands,
	}

	projectRoot := filepath.Join("..", "..", "ansible")
	varsDir := filepath.Join(projectRoot, "vars")
	if err := os.MkdirAll(varsDir, 0755); err != nil {
		return "", "", err
	}

	fileName := "job-sudo-blacklist-" + username + ".json"
	hostVarsPath := filepath.Join(varsDir, fileName)

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", "", err
	}
	if err := os.WriteFile(hostVarsPath, b, 0644); err != nil {
		return "", "", err
	}

	containerPath := "/work/vars/" + fileName
	return runPlaybookSudoBlacklist(ctx, "/work/playbooks/sudo_blacklist.yml", containerPath, limit)
}

func KafkaConsumer() {
	// ============================================
	// STEP 1: Declare variables
	// ============================================

	// broker: Kafka server address (host:port)
	// localhost = your computer, 9092 = Kafka port
	broker := "localhost:29092"

	// topic: The Kafka topic name (like a channel/queue name)
	topic := "test"

	// Print what we're connecting to
	fmt.Println("Starting Kafka Consumer...")
	fmt.Printf("Broker: %s\n", broker)
	fmt.Printf("Topic: %s\n", topic)
	fmt.Println("Waiting for messages... (Press Ctrl+C to stop)")
	fmt.Println()

	// ============================================
	// STEP 2: Create Kafka connection
	// ============================================

	// kafka.NewReader: Creates a Kafka consumer/reader
	// ReaderConfig{}: Configuration object (settings)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker}, // []string = slice of strings (list)
		Topic:   topic,            // Which topic to read from
	})

	// defer: "Do this at the end, no matter what"
	// reader.Close(): Closes the connection (cleanup)
	defer reader.Close()

	// ============================================
	// STEP 3: Set up context for cancellation
	// ============================================

	// context.Background(): Base context (starting point)
	// We'll cancel this when user presses Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Always cancel when done (cleanup)

	// Handle Ctrl+C signal (graceful shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown
	go func() {
		<-sigChan // Wait for Ctrl+C
		fmt.Println("\nStopping consumer...")
		cancel() // Cancel the context to stop reading
	}()

	// ============================================
	// STEP 4: Consume messages in a loop
	// ============================================

	// Loop forever to read messages
	for {
		// ReadMessage: Reads one message from Kafkaø
		// ctx: Context for cancellation (stops when Ctrl+C pressed)
		msg, err := reader.ReadMessage(ctx)

		// Check if context was cancelled (user pressed Ctrl+C)
		// Check both the error and the context itself
		if err == context.Canceled || ctx.Err() == context.Canceled {
			fmt.Println("Consumer stopped.")
			return // Exit the function
		}

		// Check for other errors
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			// Don't continue if context is done - exit instead
			if ctx.Err() != nil {
				fmt.Println("Consumer stopped due to context cancellation.")
				return
			}
			continue // Skip this message, try next one
		}

		// ============================================
		// STEP 5: Display the message
		// ============================================

		// Print message details
		fmt.Println("--- Message Received ---")
		fmt.Printf("Topic:     %s\n", msg.Topic)
		fmt.Printf("Partition: %d\n", msg.Partition)
		fmt.Printf("Offset:    %d\n", msg.Offset)
		fmt.Printf("Key:       %s\n", string(msg.Key))
		fmt.Printf("Value:     %s\n", string(msg.Value))
		fmt.Println("------------------------")
		fmt.Println()
	}
}

// Variables for the create_user playbook
// Variables are passed to the playbook via the --extra-vars flag
type Vars struct {
	Username      string   `json:"username"`
	UserShell     string   `json:"user_shell,omitempty"`
	UserGroups    []string `json:"user_groups,omitempty"`
	PlainPassword string   `json:"plain_password,omitempty"`
}

type PasswordVars struct {
	Username      string `json:"username"`
	PlainPassword string `json:"plain_password"`
}

type SudoBlacklistVars struct {
	Username          string   `json:"username"`
	BlacklistCommands []string `json:"blacklist_commands"`
}

func runPlaybookCreateUser(ctx context.Context, playbookPath string, extraVarsFile string) (string, string, error) {
	args := []string{
		"compose", "exec", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "/work/inventory/hosts.ini",
		"-e", "@" + extraVarsFile,
		"-l", "target1", // or "--limit", "target1"	 // Limit the playbook to only run on target1
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

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

func main() {

	// Create user on target1

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer cancel()

	// // Get the project root (ansible directory)
	// // Assuming main.go is in ansible-control/go/
	// // Project root is 2 levels up: ../../ansible
	// projectRoot := filepath.Join("..", "..", "ansible")
	// varsDir := filepath.Join(projectRoot, "vars")
	// hostVarsPath := filepath.Join(varsDir, "job-123.json")

	// // Create directory
	// _ = os.MkdirAll(varsDir, 0755)

	// v := Vars{
	// 	Username:      "ebon",
	// 	UserShell:     "/bin/bash",
	// 	UserGroups:    []string{"sudo"},
	// 	PlainPassword: "password",
	// }
	// b, _ := json.MarshalIndent(v, "", "  ")
	// _ = os.WriteFile(hostVarsPath, b, 0644)

	// // Inside container: /work/vars/job-123.json
	// containerPath := "/work/vars/job-123.json"
	// out, errOut, err := runPlaybookCreateUser(ctx, "/work/playbooks/create_user.yml", containerPath)

	// fmt.Println("STDOUT:\n", out)
	// fmt.Println("STDERR:\n", errOut)
	// fmt.Println("ERR:\n", err)

	// Add host to group
	// invPath := filepath.Join("..", "..", "ansible", "inventory", "hosts.ini")

	// hostLine := "target3 ansible_host=192.168.230.128 ansible_user=ebon ansible_password=123 ansible_port=22"
	// if err := upsertHostInGroup(invPath, "targets", "target3", hostLine); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Updated inventory:", invPath)

	// Define Path From Project Root
	// projectRoot := filepath.Join("..", "..", "ansible")
	// varsDir := filepath.Join(projectRoot, "vars")
	// hostVarsPath := filepath.Join(varsDir, "job-change-password.json")

	// userPassword := PasswordVars{
	// 	Username:      "ebon",
	// 	PlainPassword: "600",
	// }

	// // Create directory
	// b, _ := json.MarshalIndent(userPassword, "", "  ")
	// _ = os.WriteFile(hostVarsPath, b, 0644)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer cancel()

	// containerPath := "/work/vars/job-change-password.json"
	// out, errOut, err := runPlaybookChangePassword(ctx, "/work/playbooks/change_password.yml", containerPath, "target3")

	// fmt.Println("STDOUT:\n", out)
	// fmt.Println("STDERR:\n", errOut)
	// fmt.Println("ERR:\n", err)

	// Apply sudo blacklist
	ctx := context.Background()
	out, errOut, err := applySudoCommandBlacklist(
		ctx,
		"ebon",
		[]string{"/usr/bin/rm", "/usr/bin/cat", "/usr/bin/touch", "/usr/bin/cp"},
		"target3", // limit host; use "" to apply to all in group
	)
	fmt.Println(out, errOut, err)
}
