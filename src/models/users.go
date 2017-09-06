package models

import (
	"database/sql"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// User is a user.
type User struct {
	ID       int    `json:"-"`
	DomainID int    `json:"domain_id"`
	Email    string `json:"email"`
	Hash     string `json:"-"`
}

// IUsers is the Users' database interface.
type IUsers interface {
	Add(email string, domainID int, password string) (int, error)
	Authenticate(email string, domainID int, password string) (int, error)
	ChangePassword(id int, password string) error
	ReadAll() ([]User, error)
}

// Users database
type Users struct {
	DB *sql.DB
}

// Add adds a new user to the database and returns the id.
func (u Users) Add(email string, domainID int, password string) (int, error) {
	stmt, err := u.DB.Prepare(`
		INSERT INTO users (domain_id, email, hash)
			VALUES ($1, $2, $3)
			RETURNING (id)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	hash, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	err = validateEmail(email)
	if err != nil {
		return 0, err
	}

	var id int
	row := stmt.QueryRow(domainID, email, hash)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Authenticate will authenticate a user with the database and return the id.
func (u Users) Authenticate(email string, domainID int, password string) (int, error) {
	stmt, err := u.DB.Prepare(`
		SELECT * FROM users
			WHERE email=$1 AND domain_id=$2
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(email, domainID)
	if err != nil {
		return 0, err
	}

	users, err := readUsers(rows)
	if err != nil {
		return 0, err
	}

	if len(users) < 1 || len(users) > 1 {
		return 0, errors.New("User not found")
	}
	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return 0, errors.New("Authentication failed")
	}

	return user.ID, nil
}

// ChangePassword will change the password of the user in the database
// Authentication must be done before this point.
func (u Users) ChangePassword(id int, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	stmt, err := u.DB.Prepare(`
		UPDATE users SET hash=$1 WHERE id=$2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash, id)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll reads all users from the database
func (u Users) ReadAll() ([]User, error) {
	rows, err := u.DB.Query(`
		SELECT * FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := readUsers(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readUsers(rows *sql.Rows) ([]User, error) {
	res := []User{}
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.DomainID, &user.Email, &user.Hash)
		if err != nil {
			continue
		}
		res = append(res, *user)
	}

	err := rows.Err()
	if err != nil {
		return res, err
	}

	return res, nil
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
