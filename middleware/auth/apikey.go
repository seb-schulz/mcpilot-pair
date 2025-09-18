package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// APIKeyMiddleware verifies the API key in the `Authorization` header.
func APIKeyMiddleware(next http.Handler) http.Handler {
	apiKey, err := getOrGenerateAPIKey()
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
			http.Error(w, "Invalid authorization format: expected 'Bearer <api_key>'", http.StatusUnauthorized)
			return
		}

		if parts[1] != apiKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getOrGenerateAPIKey reads the API key from the file or generates a new one.
func getOrGenerateAPIKey() (string, error) {
	apiKeyFile := ".mcpilot-pair-api-key.txt"

	apiKey, err := os.ReadFile(apiKeyFile)
	if err == nil {
		return string(apiKey), nil
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("could not generate API key: %v", err)
	}
	apiKeyBytes := []byte(base64.StdEncoding.EncodeToString(key))

	if err := os.WriteFile(apiKeyFile, apiKeyBytes, 0600); err != nil {
		return "", fmt.Errorf("could not save API key: %v", err)
	}

	apiKeyStr := string(apiKeyBytes)
	fmt.Printf("\n=== MCPilot-Pair API KEY GENERATED ===\n%s\n=== COPY THIS KEY ===\n\n", apiKeyStr)
	return apiKeyStr, nil
}
