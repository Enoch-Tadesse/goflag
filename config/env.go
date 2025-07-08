package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariable(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return fmt.Errorf("unable to load dotenv file: %v ", err)
	}
	variables := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}
	for _, key := range variables {
		if value := os.Getenv(key); value == "" {
			return fmt.Errorf("environment variable %s is missing", key)
		}
	}
	return nil
}
