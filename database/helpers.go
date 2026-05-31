package database

import (
	"database/sql"
	"embed"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

const databasePath = "database"

func connectDB(loc string) *sql.DB {
	db, err := sql.Open("sqlite", loc)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// defer db.Close()
	log.Println("Connected to database")

	return db
}
func ConnectDB() *sql.DB {
	return connectDB(filepath.Join(databasePath, "dev.db"))
}
func ConnectTestDB() *sql.DB {
	return connectDB(":memory:")
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateDB(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	log.Println("Checking for database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}
	log.Println("Database migrations up to date")
}
