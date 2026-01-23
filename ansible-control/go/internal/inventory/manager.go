package inventory

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// UpsertHostInGroup adds or updates a host in a specific group in the hosts file
func UpsertHostInGroup(filePath, group, alias, hostLine string) error {
	// Get absolute path to ensure we're writing to the correct file
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	log.Printf("upsertHostInGroup: filePath=%s, absPath=%s, group=%s, alias=%s\n", filePath, absPath, group, alias)

	f, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", absPath, err)
	}
	defer f.Close()

	var out []string
	grpHeader := "[" + group + "]"
	inGroup := false
	found := false
	groupAdded := false

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		// section handling
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			// leaving group: if we didn't find alias, append it before leaving
			if inGroup && !found {
				log.Printf("Leaving group [%s], appending %s before next section\n", group, alias)
				out = append(out, hostLine)
				found = true
			}
			inGroup = (trim == grpHeader)
			if inGroup {
				groupAdded = true
				log.Printf("Found group header: %s\n", grpHeader)
			}
			out = append(out, line)
			continue
		}

		// inside group: replace if alias exists
		if inGroup {
			if strings.HasPrefix(trim, alias+" ") || trim == alias {
				log.Printf("Found existing alias %s, replacing\n", alias)
				out = append(out, hostLine)
				found = true
				continue // skip old line
			}
		}

		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// If group doesn't exist, create it at the end
	if !groupAdded {
		log.Printf("Group [%s] not found, creating at end\n", group)
		if len(out) > 0 && out[len(out)-1] != "" {
			out = append(out, "")
		}
		out = append(out, grpHeader, hostLine)
	} else if inGroup && !found {
		// Group exists, we're still in it, and alias wasn't found - append at end of group
		log.Printf("Appending %s to end of group [%s]\n", alias, group)
		out = append(out, hostLine)
	}

	result := strings.Join(out, "\n")
	log.Printf("Writing to file %s:\n%s\n", absPath, result)

	if err := os.WriteFile(absPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Successfully wrote to %s\n", absPath)
	return nil
}

// RemoveHostFromGroup removes a host from a specific group in the hosts file
func RemoveHostFromGroup(filePath, group, alias string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var out []string
	grpHeader := "[" + group + "]"
	inGroup := false
	removed := false

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		// section handling
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inGroup = (trim == grpHeader)
			out = append(out, line)
			continue
		}

		// inside group: skip if alias matches
		if inGroup {
			if strings.HasPrefix(trim, alias+" ") || trim == alias {
				removed = true
				continue // skip this line (remove it)
			}
		}

		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		return err
	}

	if removed {
		log.Printf("Removed host %s from group [%s]\n", alias, group)
	} else {
		log.Printf("Host %s not found in group [%s]\n", alias, group)
	}

	return os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644)
}
