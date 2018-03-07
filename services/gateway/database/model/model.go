package model

import "database/sql"

// Model represents the database model and contains table queries.
type Model struct {
	Users Users
}

// New creates a new model instance.
func New(db *sql.DB) *Model {
	return &Model{
		Users: &users{db},
	}
}
