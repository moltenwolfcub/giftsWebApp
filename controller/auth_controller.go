package controller

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type authController struct {
	db *sql.DB
}

func (c *authController) register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, "auth_register.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, nil)
}

func (c *authController) registerSubmit(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
