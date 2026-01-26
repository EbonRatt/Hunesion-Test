package ansible

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunPlaybook runs a basic Ansible playbook
func RunPlaybook(ctx context.Context, playbookPath string) (string, string, error) {
	args := []string{"compose", "exec", "-T", "ansible", "ansible-playbook", playbookPath}
	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunPlaybookCreateUser runs the create_user playbook
func RunPlaybookCreateUser(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	// Use "run" instead of "exec" - "run" creates a new container if service is not running
	// "exec" requires the service to already be running
	args := []string{
		"compose", "run", "--rm", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "internal/inventory/files/hosts.ini",
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

// RunPlaybookChangePassword runs the change_password playbook
func RunPlaybookChangePassword(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	// Use "run" instead of "exec" - "run" creates a new container if service is not running
	args := []string{
		"compose", "run", "--rm", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "internal/inventory/files/hosts.ini",
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

// RunPlaybookDeleteUser runs the delete_user playbook
func RunPlaybookDeleteUser(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	// Use "run" instead of "exec" - "run" creates a new container if service is not running
	args := []string{
		"compose", "run", "--rm", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "internal/inventory/files/hosts.ini",
		"-e", "@" + extraVarsFile,
	}
	if strings.TrimSpace(limit) != "" {
		args = append(args, "-l", limit)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Set working directory to project root (where docker-compose.yml should be)
	wd, wdErr := os.Getwd()
	if wdErr == nil {
		if strings.Contains(wd, "ansible-control") {
			projectRoot := filepath.Join(wd, "..", "..")
			if absPath, absErr := filepath.Abs(projectRoot); absErr == nil {
				cmd.Dir = absPath
				log.Printf("Running docker compose from: %s\n", absPath)
			}
		} else {
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

// RunPlaybookEnableUser runs the enable_user playbook
func RunPlaybookEnableUser(ctx context.Context, playbookPath string, extraVarsFile string, limit string) (string, string, error) {
	// Use "run" instead of "exec" - "run" creates a new container if service is not running
	args := []string{
		"compose", "run", "--rm", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "internal/inventory/files/hosts.ini",
		"-e", "@" + extraVarsFile,
	}
	if strings.TrimSpace(limit) != "" {
		args = append(args, "-l", limit)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Set working directory to project root (where docker-compose.yml should be)
	wd, wdErr := os.Getwd()
	if wdErr == nil {
		if strings.Contains(wd, "ansible-control") {
			projectRoot := filepath.Join(wd, "..", "..")
			if absPath, absErr := filepath.Abs(projectRoot); absErr == nil {
				cmd.Dir = absPath
				log.Printf("Running docker compose from: %s\n", absPath)
			}
		} else {
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
