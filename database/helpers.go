package database

import (
	"database/sql"
	"embed"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const databasePath = "database"

func LoadDatabase(migrate bool) *sql.DB {
	if migrate {
		goose.SetBaseFS(embedMigrations)
		// handle Goose
	}

	db, err := sql.Open("sqlite", filepath.Join(databasePath, "dev.db"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	log.Println("Connected to database")

	return db
}
