package controller

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

const templatesPath = "templates"

func BuildRouter(database *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	wishlistController := &wishlistController{db: database}
	mux.HandleFunc("GET /wishlist/", wishlistController.index)
	mux.HandleFunc("GET /wishlist/add_item", wishlistController.addItem)
	mux.HandleFunc("POST /wishlist/add_item", wishlistController.addItemSubmit)
	mux.HandleFunc("GET /wishlist/edit_item/{id}", wishlistController.editItem)
	mux.HandleFunc("POST /wishlist/edit_item/{id}", wishlistController.editItemSubmit)
	mux.HandleFunc("POST /wishlist/delete_item/{id}", wishlistController.deleteItem)

	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, "index.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, nil)
}
