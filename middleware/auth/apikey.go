package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// APIKeyMiddleware verifies the API key in the `Authorization` header.
func APIKeyMiddleware(next http.Handler) http.Handler {
	apiKey, err := getOrGenerateAPIKey()
	fmt.Printf("\n=== MCPilot-Pair API KEY GENERATED ===\n%s\n=== COPY THIS KEY ===\n\n", apiKey)
	if err != nil {
		panic("API key not configured")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("invalid bearer token: %s", authHeader)
			http.Error(w, "Invalid authorization format: expected 'Bearer <api_key>'", http.StatusUnauthorized)
			return
		}

		if parts[1] != apiKey {
			log.Printf("failed to verify %s != %s", parts[1], apiKey)
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getOrGenerateAPIKey reads the API key from the XDG-compliant config directory or generates a new one.
func getOrGenerateAPIKey() (string, error) {
	// XDG-compliant config directory: ~/.config/mcpilot-pair/
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine config directory: %v", err)
	}
	apiKeyDir := filepath.Join(configDir, "mcpilot-pair")
	apiKeyFile := filepath.Join(apiKeyDir, "api-key.txt")

	if err := os.MkdirAll(apiKeyDir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %v", err)
	}

	apiKey, err := os.ReadFile(apiKeyFile)
	if err == nil {
		return strings.TrimSpace(string(apiKey)), nil
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("could not generate API key: %v", err)
	}
	apiKeyBytes := []byte(base64.StdEncoding.EncodeToString(key))

	if err := os.WriteFile(apiKeyFile, apiKeyBytes, 0600); err != nil {
		return "", fmt.Errorf("could not save API key: %v", err)
	}

	return strings.TrimSpace(string(apiKeyBytes)), nil
}
