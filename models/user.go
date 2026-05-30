package models

import (
	"crypto/rand"
	"database/sql"
	"log"

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

func UsernameTaken(db *sql.DB, username string) (bool, error) {
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

const (
	saltLength uint   = 16
	iterations uint32 = 4
	memory     uint32 = 64 * 1024
	threads    uint8  = 2
	keyLength  uint32 = 32
)

func HashPassword(password string) (hash, salt []byte, err error) {
	salt, err = genSalt()
	if err != nil {
		log.Printf("Error generating salt for Registering password: %v", err)
		return nil, nil, err
	}

	hash = argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLength)

	return hash, salt, nil
}

func genSalt() ([]byte, error) {
	b := make([]byte, saltLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
