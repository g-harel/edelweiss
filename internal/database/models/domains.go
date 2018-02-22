package models

// Domain is a domain.
type Domain struct {
	ID   int    `json:"-"`
	Name string `json:"email"`
	Data string `json:"data"`
}
