package model

import (
	"database/sql"

	"github.com/g-harel/edelweiss/internal/database"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is a user.
type User struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Hash     string `json:"-"`
}

// Users interface is implemented by the "users" struct.
type Users interface {
	Add(email, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
	ChangeVerified(email string, value bool) error
	ChangeHash(email, password, newPassword string) error
}

// Users holds the users table queries.
type users struct {
	db *sql.DB
}

func (u *users) Add(e, p string) (*User, error) {
	err := database.ValidateEmail(e)
	if err != nil {
		return nil, err
	}

	err = database.ValidatePassword(p)
	if err != nil {
		return nil, err
	}

	h, err := hashPassword(p)
	if err != nil {
		return nil, err
	}

	stmt, err := u.db.Prepare(`
		INSERT INTO users (uuid, email, hash)
		VALUES            ($1, $2, $3)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	id := uuid.NewV4().String()
	_, err = stmt.Exec(id, e, h)
	if err != nil {
		return nil, err
	}

	return &User{
		UUID:     id,
		Email:    e,
		Hash:     h,
		Verified: false,
	}, nil
}

func (u *users) Authenticate(e, p string) (*User, error) {
	stmt, err := u.db.Prepare(`
		SELECT uuid, hash, verified
		FROM   users
		WHERE  email=$1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	usr := &User{
		Email: e,
	}
	err = stmt.QueryRow(e).Scan(&usr.UUID, &usr.Hash, &usr.Verified)
	if err != nil {
		return nil, err
	}

	return usr, checkPassword(usr.Hash, p)
}

func (u *users) ChangeHash(e, p, n string) error {
	err := database.ValidatePassword(p)
	if err != nil {
		return err
	}

	usr, err := u.Authenticate(e, p)
	if err != nil {
		return err
	}

	stmt, err := u.db.Prepare(`
		UPDATE users
		SET    hash=$1
		WHERE  uuid=$2;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	h, err := hashPassword(n)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(h, usr.UUID)
	return err
}

func (u *users) ChangeVerified(e string, v bool) error {
	stmt, err := u.db.Prepare(`
		UPDATE users
		SET    verified=$1
		WHERE  email=$2;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(v, e)
	return err
}

func hashPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), 12)
	return string(h), err
}

func checkPassword(h, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}
