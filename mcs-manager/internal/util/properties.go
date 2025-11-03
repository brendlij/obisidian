package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseProperties reads and parses a Java properties file
func ParseProperties(filePath string) (map[string]string, error) {
	properties := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open properties file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			properties[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading properties file: %w", err)
	}

	return properties, nil
}

// SaveProperties writes properties to a Java properties file
func SaveProperties(filePath string, properties map[string]string) error {
	// First, read existing file to preserve comments and order
	var lines []string
	var existingKeys map[string]bool

	if file, err := os.Open(filePath); err == nil {
		defer file.Close()
		existingKeys = make(map[string]bool)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			// Keep comments and empty lines
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				lines = append(lines, line)
				continue
			}

			// For properties, check if we have a new value
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				existingKeys[key] = true

				// If we have a new value, update it
				if newValue, exists := properties[key]; exists {
					lines = append(lines, fmt.Sprintf("%s=%s", key, newValue))
				} else {
					lines = append(lines, line)
				}
			}
		}
	} else {
		existingKeys = make(map[string]bool)
	}

	// Add new properties that weren't in the file
	for key, value := range properties {
		if !existingKeys[key] {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Write back to file
	output := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filePath, []byte(output), 0o644); err != nil {
		return fmt.Errorf("failed to write properties file: %w", err)
	}

	return nil
}

// GetServerPropertiesPath returns the path to server.properties
func GetServerPropertiesPath(serverPath string) string {
	return filepath.Join(serverPath, "server.properties")
}
