package controller

import (
	"database/sql"
	"net/http"
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

	authController := &authController{db: database}
	mux.HandleFunc("GET /auth/register/", authController.register)
	mux.HandleFunc("POST /auth/register/", authController.registerSubmit)

	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "index.html", nil)
}
