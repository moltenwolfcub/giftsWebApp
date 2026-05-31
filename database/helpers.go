package database

import (
	"database/sql"
	"embed"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

const databasePath = "database"

func connectDB(loc string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", loc)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// defer db.Close()
	log.Println("Connected to database")

	return db, nil
}
func ConnectDB() (*sql.DB, error) {
	return connectDB(filepath.Join(databasePath, "dev.db"))
}
func ConnectTestDB() (*sql.DB, error) {
	return connectDB(":memory:")
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateDB(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	log.Println("Checking for database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	log.Println("Database migrations up to date")

	return nil
}
