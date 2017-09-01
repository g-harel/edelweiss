package users

import (
	"database/sql"
	"fmt"
)

// Test runs some mock actions.
func Test(db *sql.DB) {
	// adding users
	userList := []User{
		User{
			DomainID: 1,
			Email:    "email1@example.com",
			Hash:     "password123",
		},
		User{
			DomainID: 1,
			Email:    "email2@example.com",
			Hash:     "password123",
		},
		User{
			DomainID: 2,
			Email:    "email1@example.com",
			Hash:     "password123",
		},
	}
	for _, u := range userList {
		err := Add(db, &u)
		if err != nil {
			fmt.Println(err)
		}
	}

	// testing user funcs
	user := User{
		DomainID: 2,
		Email:    "email1@example.com",
	}
	err := Authenticate(db, &user, "password123")
	if err != nil {
		fmt.Println(err)
	}

	err = ChangePassword(db, user.ID, "123password")
	if err != nil {
		fmt.Println(err)
	}
	err = Authenticate(db, &user, "123password")
	if err != nil {
		fmt.Println(err)
	}
}
