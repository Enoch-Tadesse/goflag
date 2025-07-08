package config

import (
	"os"
	"testing"
)

func TestLoadEnvVariable(t *testing.T) {
	// Load the test env file committed in the repo
	LoadEnvVariable(".env.test")

	// Expected keys and values (adjust as per your .env.test content)
	expectedVars := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}

	for _, key := range expectedVars {
		if val := os.Getenv(key); val == "" {
			t.Errorf("Expected environment variable %s to be set, but got empty", key)
		}
	}
}
