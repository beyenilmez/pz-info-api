package server

import (
	"log"
	"strconv"
	"strings"
)

// Options holds data parsed from the `showoptions` RCON command.
type Options struct {
	MaxPlayers        int    `json:"MaxPlayers"`
	PublicDescription string `json:"PublicDescription"`
	PublicName        string `json:"PublicName"`
}

// GetOptions retrieves server options by calling `showoptions` via RCON.
func GetOptions() (Options, error) {
	var opts Options

	resp, err := executeRCONCommand("showoptions")
	if err != nil {
		return opts, err
	}

	lines := strings.Split(resp, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "* ") {
			continue
		}

		parts := strings.SplitN(strings.TrimPrefix(line, "* "), "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "MaxPlayers":
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("Warning: could not parse MaxPlayers: %v", err)
				continue
			}
			opts.MaxPlayers = i
		case "PublicDescription":
			opts.PublicDescription = val
		case "PublicName":
			opts.PublicName = val
		}
	}
	return opts, nil
}
