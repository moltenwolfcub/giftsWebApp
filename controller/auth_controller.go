package controller

import (
	"database/sql"
	"net/http"
)

type authController struct {
	db *sql.DB
}

func (c *authController) register(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "auth_register.html", nil)
}

func (c *authController) registerSubmit(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
