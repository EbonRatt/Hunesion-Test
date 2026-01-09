package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Executor handles remote command execution
type Executor struct {
	client *ssh.Client
}

// NewExecutor creates a new SSH executor
func NewExecutor(host, user, password, keyPath string) (*Executor, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Use password authentication if provided
	if password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(password),
		}
	} else if keyPath != "" {
		// Use key-based authentication
		key, err := readPrivateKey(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(key),
		}
	} else {
		return nil, fmt.Errorf("either password or SSH key must be provided")
	}

	// Connect to the remote server
	addr := host
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%s:22", addr)
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	return &Executor{client: client}, nil
}

// Execute runs a command on the remote host
func (e *Executor) Execute(command string) (string, error) {
	session, err := e.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// Close closes the SSH connection
func (e *Executor) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

// readPrivateKey reads a private key from file path
func readPrivateKey(keyPath string) (ssh.Signer, error) {
	// For simplicity, we'll use a basic approach
	// In production, you'd read from file system
	// This is a placeholder - you'll need to implement actual key reading
	return nil, fmt.Errorf("key-based auth not fully implemented - use password for now")
}

// CreateUser executes the user creation command on remote host
func CreateUser(host, sshUser, sshPass, sshKey, username, userID string) error {
	executor, err := NewExecutor(host, sshUser, sshPass, sshKey)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}
	defer executor.Close()

	// Check if user already exists
	checkCmd := fmt.Sprintf("id -u %s 2>/dev/null", username)
	output, _ := executor.Execute(checkCmd)
	if strings.TrimSpace(output) != "" {
		return fmt.Errorf("user %s already exists", username)
	}

	// Create user with specified user ID
	createCmd := fmt.Sprintf("sudo useradd -m -u %s %s", userID, username)
	output, err = executor.Execute(createCmd)
	if err != nil {
		return fmt.Errorf("failed to create user: %w, output: %s", err, output)
	}

	// Set a default password (optional)
	setPassCmd := fmt.Sprintf("echo '%s:%s' | sudo chpasswd", username, "TempPass123!")
	executor.Execute(setPassCmd)

	fmt.Printf("Successfully created user %s (ID: %s) on %s\n", username, userID, host)
	return nil
}
