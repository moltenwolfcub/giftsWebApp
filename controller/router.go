package controller

import (
	"database/sql"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/config"
)

const templatesPath = "templates"

func BuildRouter(database *sql.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	authController := &authController{db: database, cfg: cfg}
	mux.HandleFunc("GET /auth/register/", authController.register)
	mux.HandleFunc("POST /auth/register/", authController.registerSubmit)
	mux.HandleFunc("GET /auth/login/", authController.login)
	mux.HandleFunc("POST /auth/login/", authController.loginSubmit)

	wishlistController := &wishlistController{db: database}
	mux.HandleFunc("GET /wishlist/", wishlistController.index)
	// mux.Handle("GET /wishlist/", authController.AuthRequired(wishlistController.index))
	mux.HandleFunc("GET /wishlist/add_item", wishlistController.addItem)
	mux.HandleFunc("POST /wishlist/add_item", wishlistController.addItemSubmit)
	mux.HandleFunc("GET /wishlist/edit_item/{id}", wishlistController.editItem)
	mux.HandleFunc("POST /wishlist/edit_item/{id}", wishlistController.editItemSubmit)
	mux.HandleFunc("POST /wishlist/delete_item/{id}", wishlistController.deleteItem)

	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "index.html", nil)
}
