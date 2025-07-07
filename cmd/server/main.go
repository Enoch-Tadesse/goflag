package main

import (
	"log"
	"net/http"

	"github.com/Enoch-Tadesse/goflag/config"
	"github.com/Enoch-Tadesse/goflag/db/connection"
)

func setup() error {
	if err := config.LoadEnvVariable(".env"); err != nil {
		return err
	}
	if err := connection.ConnectToDB(); err != nil {
		return err
	}
	if err := connection.RunMigrations(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	http.ListenAndServe(":8080", nil)
}
