package db

import (
	"fmt"
	"github.com/Enoch-Tadesse/goflag/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDb() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	log.Printf("dsn: %s", dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Unable to connect to DB, check if DB exist. Error : %v", err)
	}
	log.Println("Successfully connected to DB")
}

func MigrateTables() {
	tables := []any{
		&models.User{},
		&models.Flag{},
		&models.Segment{},
		&models.FlagSegment{},
		&models.SegmentRule{},
		&models.FlagUser{},
	}
	for _, table := range tables {
		if err := DB.AutoMigrate(table); err != nil {
			log.Fatalf("Error migrating %T: %v\n", table, err)
		} else {
			log.Printf("Successfully migrated model:%T\n", table)
		}
	}
}
