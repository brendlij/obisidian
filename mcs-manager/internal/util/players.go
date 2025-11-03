package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ReadPlayersFromLog attempts to parse player count from server logs
// Looks for multiple patterns that indicate player join/leave events
func ReadPlayersFromLog(logPath string) (current int, max int, err error) {
	file, err := os.Open(logPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Multiple patterns to match player count messages
	// Pattern 1: "there are X of a max of Y players online" (Vanilla)
	pattern1 := regexp.MustCompile(`there are (\d+) of a max of (\d+) players online`)
	// Pattern 2: "X players online" (simple format, may be at end of startup)
	pattern2 := regexp.MustCompile(`\[SERVER\].*?(\d+) players? online`)
	// Pattern 3: Player join messages (Fabric/Paper format)
	playerJoinPattern := regexp.MustCompile(`\]\s+\[Server thread/INFO\]:\s+(\w+)\s+joined the game`)
	playerLeavePattern := regexp.MustCompile(`\]\s+\[Server thread/INFO\]:\s+(\w+)\s+left the game`)

	scanner := bufio.NewScanner(file)
	currentPlayers := 0
	maxPlayers := 0
	foundExactCount := false
	playersOnline := make(map[string]bool) // Track who's online

	// Read the entire file to find the most recent player count
	for scanner.Scan() {
		line := scanner.Text()

		// Try pattern 1 (most reliable - exact count)
		if matches := pattern1.FindStringSubmatch(line); len(matches) == 3 {
			cur, _ := strconv.Atoi(matches[1])
			max, _ := strconv.Atoi(matches[2])
			currentPlayers = cur
			maxPlayers = max
			foundExactCount = true
			// Reset player tracking since we found exact count
			playersOnline = make(map[string]bool)
			continue
		}

		// Try pattern 2
		if matches := pattern2.FindStringSubmatch(line); len(matches) == 2 && !foundExactCount {
			cur, _ := strconv.Atoi(matches[1])
			currentPlayers = cur
			continue
		}

		// Count join events (even if we found exact count, keep tracking for updates)
		if matches := playerJoinPattern.FindStringSubmatch(line); len(matches) == 2 {
			playerName := matches[1]
			if !playersOnline[playerName] {
				playersOnline[playerName] = true
				if !foundExactCount {
					currentPlayers++
				}
			}
		}

		// Count leave events
		if matches := playerLeavePattern.FindStringSubmatch(line); len(matches) == 2 {
			playerName := matches[1]
			if playersOnline[playerName] {
				delete(playersOnline, playerName)
				if !foundExactCount && currentPlayers > 0 {
					currentPlayers--
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("error reading log file: %w", err)
	}

	// If we found exact count, use it; otherwise use tracked players
	if foundExactCount && maxPlayers > 0 {
		return currentPlayers, maxPlayers, nil
	}

	// Use tracked player count if available
	if !foundExactCount {
		currentPlayers = len(playersOnline)
	}

	// If we have current players but no max, this is still a valid result
	// The caller should handle filling in max from config if needed
	if currentPlayers > 0 || maxPlayers > 0 {
		return currentPlayers, maxPlayers, nil
	}

	return 0, 0, fmt.Errorf("no player count found in logs")
}

// ReadMaxPlayersFromConfig reads max-players from server.properties
func ReadMaxPlayersFromConfig(configPath string) (int, error) {
	props, err := ParseProperties(configPath)
	if err != nil {
		return 0, err
	}

	maxPlayersStr, exists := props["max-players"]
	if !exists {
		return 20, nil // default
	}

	maxPlayers, err := strconv.Atoi(strings.TrimSpace(maxPlayersStr))
	if err != nil {
		return 20, nil // default on parse error
	}

	return maxPlayers, nil
}
