package models

import (
	"crypto/rand"
	"database/sql"
	"log"

	"github.com/moltenwolfcub/giftsWebApp/config"
	"golang.org/x/crypto/argon2"
)

type User struct {
	Id       int
	Username string
	password string
	salt     string
}

func LoadUser(db *sql.DB, id int) (*User, error) {
	u := User{}
	err := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&u.Id, &u.Username, &u.password, &u.salt)
	if err != nil {
		log.Printf("Error finding user [%d] from database: %v", id, err)
		return nil, err
	}
	return &u, nil
}

func UsernameExists(db *sql.DB, username string) (bool, error) {
	var numFound int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username).Scan(&numFound)
	if err != nil {
		log.Printf("Error checking username [%s] exists: %v", username, err)
		return true, err
	}

	if numFound > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func HashPassword(password string, cfg *config.Config) (hash, salt []byte, err error) {
	salt, err = genSalt(cfg.SaltLength)
	if err != nil {
		log.Printf("Error generating salt for Registering password: %v", err)
		return nil, nil, err
	}

	hash = HasPasswordWithSalt(password, cfg, salt)

	return hash, salt, nil
}

func HasPasswordWithSalt(password string, cfg *config.Config, salt []byte) []byte {
	hash := argon2.IDKey(append([]byte(password), cfg.Pepper...), salt, cfg.ArgonIterations, cfg.ArgonMem, cfg.ArgonThreads, cfg.HashLength)
	return hash
}

func genSalt(length uint) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
