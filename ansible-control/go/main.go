package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"test-kafka/internal/kafka"
	"test-kafka/internal/models"
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
	v := models.SudoBlacklistVars{
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

func runPlaybookCreateUser(ctx context.Context, playbookPath string, extraVarsFile string) (string, string, error) {
	args := []string{
		"compose", "exec", "-T", "ansible",
		"ansible-playbook", playbookPath,
		"-i", "/work/inventory/hosts.ini",
		"-e", "@" + extraVarsFile,
		"-l", "target3", // or "--limit", "target1"	 // Limit the playbook to only run on target1
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

func main() {

	if err := kafka.RunConsumer("localhost:29092", "test"); err != nil {
		log.Fatal(err)
	}
}
