package models

// User is a user.
type User struct {
	ID       int    `json:"-"`
	DomainID int    `json:"domain_id"`
	Email    string `json:"email"`
	Hash     string `json:"-"`
}
