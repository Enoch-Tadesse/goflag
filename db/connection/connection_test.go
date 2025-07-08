package connection

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Enoch-Tadesse/goflag/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func TestDBConnectionAndMigrations(t *testing.T) {
	keys := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_NAME",
		"DB_PORT",
	}

	// test env variables loading
	config.LoadEnvVariable(".env.test")

	for _, key := range keys {
		if os.Getenv(key) == "" {
			t.Fatalf("Failed to load .env key %s", key)
		}
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// test db connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open DB Connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping DB: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, "../../db/migrations"); err != nil {
		t.Fatalf("Goose migration failed: %v", err)
	}

	tables := []string{
		"users",
		"segments",
		"flags",
		"flag_segments",
		"flag_users",
		"segment_rules",
	}

	// check individual table existence
	for _, table := range tables {
		if !tableExists(db, dbName, table) {
			t.Fatalf("Expected %s table not found", table)
		}
	}
}

// tableExists checks whether a table with the given name exists
// in the specified database schema.
func tableExists(db *sql.DB, dbName string, table string) bool {
	row := db.QueryRow(`
		SELECT Count(*) FROM information_schema.tables
		WHERE table_schema = ? AND table_name = ?
	`, dbName, table)

	var count int
	_ = row.Scan(&count)
	return count == 1
}

// columnExists verifies whether a specific column exists in a given table
// within the specified database schema. This is useful for validating
// column-level changes introduced via SQL migrations.
func columnExists(db *sql.DB, dbName, table, column string) bool {
	row := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ? AND column_name = ?`,
		dbName, table, column)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count == 1
}
