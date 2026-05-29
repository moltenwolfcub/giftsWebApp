package main

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	"github.com/moltenwolfcub/giftsWebApp/controller"

	_ "modernc.org/sqlite"
)

const databasePath = "database"

var db *sql.DB

func main() {
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

	controller.BuildRouter(db)
	log.Fatal(http.ListenAndServe(":8040", nil))
}
