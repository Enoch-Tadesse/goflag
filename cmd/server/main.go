package main

import (
	"github.com/Enoch-Tadesse/goflag/config"
	"github.com/Enoch-Tadesse/goflag/pkg/db"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariable()
	db.ConnectToDb()
	db.MigrateTables()
}

func main() {
	router := gin.Default()

	router.Run(":8080")
}
