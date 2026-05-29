package database

import (
	"database/sql"
	"log"
	"path/filepath"
)

const databasePath = "database"

func LoadDatabase() *sql.DB {
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
