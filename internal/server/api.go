package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// ApiResponse is what we return via our HTTP endpoint.
type ApiResponse struct {
	Status            string   `json:"status"`
	Players           []string `json:"players"`
	OnlinePlayerCount int      `json:"onlinePlayerCount"`
	MaxPlayers        int      `json:"maxPlayers"`
	PlayersString     string   `json:"playersString"`
	PublicName        string   `json:"publicName"`
	PublicDescription string   `json:"publicDescription"`
}

var (
	cacheMu        sync.RWMutex
	cachedResponse ApiResponse
	cacheExpiry    time.Time
	cacheTTL       = loadCacheTTLFromEnv()
)

// loadCacheTTLFromEnv reads the environment variable `CACHE_TTL_SECONDS`,
// converts it to an integer, and returns a `time.Duration`. Defaults to 5 seconds if unset/invalid.
func loadCacheTTLFromEnv() time.Duration {
	godotenv.Load()

	const envVar = "CACHE_TTL_SECONDS"
	defaultSeconds := 5

	val := os.Getenv(envVar)
	if val == "" {
		// Fallback if no env var is set
		log.Printf("%s not set, defaulting to %d seconds\n", envVar, defaultSeconds)
		return time.Duration(defaultSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid %s=%q, defaulting to %d seconds\n", envVar, val, defaultSeconds)
		return time.Duration(defaultSeconds) * time.Second
	}

	ttl := time.Duration(seconds) * time.Second
	log.Printf("%s=%d seconds\n", envVar, seconds)
	return ttl
}

// APIHandler handles the root ("/") endpoint, returning cached data if valid,
// or fresh data if the cache is expired.
func APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	now := time.Now()

	// 1. Check if cached data is still valid
	cacheMu.RLock()
	if now.Before(cacheExpiry) {
		log.Println("Serving response from cache")
		response := cachedResponse
		cacheMu.RUnlock()

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error writing cached JSON response: %v", err)
		}
		return
	}
	cacheMu.RUnlock()

	// 2. Cache expired or not set; fetch fresh data.
	players, err := GetPlayers()
	if err != nil {
		log.Printf("Error getting players: %v", err)
		_ = json.NewEncoder(w).Encode(ApiResponse{Status: "Offline"})
		return
	}

	options, err := GetOptions()
	if err != nil {
		log.Printf("Error getting options: %v", err)
		_ = json.NewEncoder(w).Encode(ApiResponse{Status: "Offline"})
		return
	}

	// 3. Build a fresh response
	response := ApiResponse{
		Status:            "Online",
		Players:           players,
		OnlinePlayerCount: len(players),
		MaxPlayers:        options.MaxPlayers,
		PlayersString:     fmt.Sprintf("%d/%d", len(players), options.MaxPlayers),
		PublicName:        options.PublicName,
		PublicDescription: options.PublicDescription,
	}

	// 4. Store in cache (under write lock)
	cacheMu.Lock()
	cachedResponse = response
	cacheExpiry = now.Add(cacheTTL)
	cacheMu.Unlock()

	// 5. Send fresh response
	log.Println("Serving fresh response, cache updated")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
