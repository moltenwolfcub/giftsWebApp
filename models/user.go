package models

import (
	"database/sql"
	"log"
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
