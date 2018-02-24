package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// User is a user.
type User struct {
	UUID  int    `json:"uuid"`
	Email string `json:"email"`
	Hash  string `json:"-"`
}

// Users holds the users table queries.
type Users struct {
	db *sql.DB
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
