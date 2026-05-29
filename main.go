package main

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/moltenwolfcub/giftsWebApp/app"
	"github.com/moltenwolfcub/giftsWebApp/config"

	_ "modernc.org/sqlite"
)

const databasePath = "database"

var db *sql.DB

func main() {
	cfg := config.New()

	var err error
	db, err = sql.Open("sqlite", filepath.Join(databasePath, "dev.db"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Connected to database")

	app := app.NewAppServer(db, cfg)
	app.Run()
}
