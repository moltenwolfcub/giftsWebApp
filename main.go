package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const templatesPath = "templates"
const databasePath = "database"

type Wishlist struct {
	Items []string
}

var dummyWishlist = Wishlist{[]string{"Chocolate", "Sweets", "Guitar"}}

var db *sql.DB

func main() {
	db, err := sql.Open("sqlite", filepath.Join(databasePath, "dev.db"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wishlist/", handleWishlist)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home\nHello World")
}

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "wishlist.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, dummyWishlist)
}
