package controller

import (
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/models"
)

type authController struct {
	db  *sql.DB
	cfg *config.Config
}

func (c *authController) register(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "auth_register.html", nil)
}

func (c *authController) registerSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	userTaken, err := models.UsernameExists(c.db, username)
	if err != nil {
		log.Printf("Error checking for duplicate username on register: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userTaken {
		// TODO: add already submitted form data along with
		// hint complaining about duplicate username
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	if password != confirmPassword {
		// TODO: add already submitted form data along with
		// hint complaining about non-identical passwords
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}
	//TODO: ensure password isn't too long
	//TODO: maybe force passwords to be kinda secure

	hash, salt, err := models.HashPassword(password, c.cfg)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.db.Exec("INSERT INTO users (username, password, salt) VALUES (?, ?, ?);", username, hash, salt)
	if err != nil {
		log.Printf("Error registering user to database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//TODO: sign the user in after registering them

	//TODO: redirect user to their wishlist once its user protected
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *authController) login(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "auth_login.html", nil)
}

func (c *authController) loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	userExists, err := models.UsernameExists(c.db, username)
	if err != nil {
		log.Printf("Error looking up username: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !userExists {
		// TODO: add already submitted form data along with
		// hint complaining about invalid credentials (NOT INVALID USERNAME)
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	var dbPassword, salt []byte
	c.db.QueryRow("SELECT password, salt FROM users WHERE username=?", username).Scan(&dbPassword, &salt)

	password := r.FormValue("password")
	hash := models.HasPasswordWithSalt(password, c.cfg, salt)

	if subtle.ConstantTimeCompare(dbPassword, hash) != 1 {
		// TODO: add already submitted form data along with
		// hint complaining about invalid credentials (NOT INVALID USERNAME)
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
