package controller

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/models"
)

type contextKey string

const sessionContextKey contextKey = "session"

type authController struct {
	db  *sql.DB
	cfg *config.Config
}

func (c *authController) register(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "auth_register.html", nil)
}

func (c *authController) registerSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	if len(username) == 0 {
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}

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

	if len(password) == 0 {
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}

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

	hexHash := hex.EncodeToString(hash)
	hexSalt := hex.EncodeToString(salt)

	_, err = c.db.Exec("INSERT INTO users (username, password, salt) VALUES (?, ?, ?);", username, hexHash, hexSalt)
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

	var hexDBPassword, hexSalt string
	c.db.QueryRow("SELECT password, salt FROM users WHERE username=?", username).Scan(&hexDBPassword, &hexSalt)

	dbPassword, err := hex.DecodeString(hexDBPassword)
	if err != nil {
		log.Printf("Error decoding password hash: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	salt, err := hex.DecodeString(hexSalt)
	if err != nil {
		log.Printf("Error decoding salt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	password := r.FormValue("password")
	hash := models.HashPasswordWithSalt(password, c.cfg, salt)

	if subtle.ConstantTimeCompare(dbPassword, hash) != 1 {
		// TODO: add already submitted form data along with
		// hint complaining about invalid credentials (NOT INVALID USERNAME)
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	// build session object
	// Generate JWT
	// set cookie
	// cookie := http.Cookie{}
	// http.SetCookie(w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *authController) AuthRequired(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read cookie
		// parse and validate JWT
		// store session object in context
		var session *string = nil
		if session == nil {
			// http.Redirect(/login)
			http.Error(w, "Unauthenticated", 403)
			return
		}

		userID := "something generated from cookie"

		ctx := context.WithValue(r.Context(), sessionContextKey, userID)
		handler(w, r.WithContext(ctx))
	})
}

func extractUserID(ctx context.Context) string {
	userID := ctx.Value(sessionContextKey)
	if userID == nil {
		return ""
	}
	return userID.(string)
}
