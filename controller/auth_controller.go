package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/models"
)

type authController struct {
	db *sql.DB
}

func (c *authController) register(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "auth_register.html", nil)
}

func (c *authController) registerSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	userTaken, err := models.UsernameTaken(c.db, username)
	if err != nil {
		log.Printf("Error checking for duplicate username on register: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Print(username, userTaken)

	if userTaken {
		// TODO: add already submitted form data along with
		// hint complaining about duplicate username
		http.Redirect(w, r, "/auth/register", http.StatusFound)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
