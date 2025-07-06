package main

import (
	"net/http"

	"github.com/Enoch-Tadesse/goflag/config"
	"github.com/Enoch-Tadesse/goflag/db/connection"
)

func init() {
	config.LoadEnvVariable()
	connection.ConnectToDB()
	connection.RunMigrations()
}

func main() {

	http.ListenAndServe(":8080", nil)
}
