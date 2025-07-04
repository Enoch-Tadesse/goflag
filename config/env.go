package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariable() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Unable to load dotenv file: %v ", err)
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
			log.Fatalf("Environment variable %s is missing", key)
		}
	}
}
