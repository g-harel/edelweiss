package users

import (
	"fmt"
	"regexp"
	"errors"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// User is a user
type User struct {
	ID       int    `json:"-"`
	DomainID int    `json:"domain_id"`
	Email    string `json:"email"`
	Hash     string `json:"-"`
}

func hashPassword(password string) ([]byte, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func validateEmail(email string) error {
	matched, err := regexp.Match("^.+@.+\\..+$", []byte(email))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("Invalid email address")
	}

	return nil
}

// Add adds a new user to the database.
// The new user's id is added to the User struct.
func Add(db *sql.DB, u *User) error {
	stmt, err := db.Prepare(`
		INSERT INTO users (domain_id, email, hash)
			VALUES ($1, $2, $3)
			RETURNING (id)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hash, err := hashPassword(u.Hash)
	if err != nil {
		return err
	}

	err = validateEmail(u.Email)
	if err != nil {
		return err
	}

	row := stmt.QueryRow(u.DomainID, u.Email, hash)
	err = row.Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll reads all users from the database
func ReadAll(db *sql.DB) []User {
	rows, err := db.Query(`
		SELECT * FROM users
	`)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var res []User

	for rows.Next() {
		user := new(User)
		rows.Scan(&user.ID, &user.DomainID, &user.Email, &user.Hash)
		res = append(res, *user)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	return res
}
