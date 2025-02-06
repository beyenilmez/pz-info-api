package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gorcon/rcon"
)

var (
	rconConn *rcon.Conn
	rconMu   sync.Mutex
)

// ensureRCONConnection checks if rconConn is nil; if so, tries to connect.
func ensureRCONConnection() error {
	rconMu.Lock()
	defer rconMu.Unlock()

	if rconConn != nil {
		return nil
	}

	socket := os.Getenv("RCON_SOCKET")
	password := os.Getenv("RCON_PASSWORD")
	conn, err := rcon.Dial(socket, password)
	if err != nil {
		return fmt.Errorf("failed to connect to RCON: %w", err)
	}

	rconConn = conn
	log.Println("Connected to RCON")
	return nil
}

// executeRCONCommand is a concurrency-safe way to run an RCON command.
func executeRCONCommand(cmd string) (string, error) {
	// Make sure we're connected
	if err := ensureRCONConnection(); err != nil {
		return "", err
	}

	rconMu.Lock()
	defer rconMu.Unlock()

	resp, err := rconConn.Execute(cmd)
	if err != nil {
		// Attempt one reconnect on error
		log.Printf("RCON error: %v. Reconnecting...", err)
		rconConn.Close()
		rconConn = nil // force reconnect next time

		// reconnect once
		if err2 := ensureRCONConnection(); err2 != nil {
			return "", fmt.Errorf("failed reconnect: %v", err2)
		}
		// re-try the command once
		resp, err = rconConn.Execute(cmd)
		if err != nil {
			return "", err
		}
	}
	return resp, nil
}

// GetPlayers runs the `players` command and parses the output.
func GetPlayers() ([]string, error) {
	resp, err := executeRCONCommand("players")
	if err != nil {
		return nil, err
	}

	// Example response:
	// Players connected (2):
	// - player1
	// - player2

	lines := strings.Split(resp, "\n")
	players := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			// Trim off the '-'
			name := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			players = append(players, name)
		}
	}
	return players, nil
}
