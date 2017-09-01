package users

import (
	"database/sql"
	"errors"
	"regexp"

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

func read(rows *sql.Rows) ([]User, error) {
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

// Authenticate will authenticate a user with the database.
// It will also add missing information to the User reference.
func Authenticate(db *sql.DB, u *User, password string) error {
	stmt, err := db.Prepare(`
		SELECT * FROM users
			WHERE email=$1 AND domain_id=$2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(u.Email, u.DomainID)
	if err != nil {
		return err
	}

	users, err := read(rows)
	if err != nil {
		return err
	}

	if len(users) < 1 || len(users) > 1 {
		return errors.New("User not found")
	}
	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return errors.New("Authentication failed")
	}

	u.ID = user.ID

	return nil
}

// ChangePassword will change the password of the user in the database
// Authentication is must be done before this point.
func ChangePassword(db *sql.DB, id int, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
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
func ReadAll(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`
		SELECT * FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := read(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}
