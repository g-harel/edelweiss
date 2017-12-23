package session

import (
	"fmt"
)

type sessionStorer interface {
	get(id, key string) (string, error)
	set(id, key, value string) error
}

// Session contains identifying information about the user.
type Session struct {
	id    string
	store sessionStorer
}

// Get fetches session data.
func (s *Session) Get(key string) (string, error) {
	val, err := s.store.get(s.id, key)
	if err != nil {
		return val, fmt.Errorf("error reading from store: %v", err)
	}

	return val, nil
}

// Set modifies session data.
func (s *Session) Set(key, value string) error {
	err := s.store.set(s.id, key, value)
	if err != nil {
		return fmt.Errorf("error setting value in store: %v", err)
	}

	return nil
}

// ID is a getter for the session's id property.
func (s *Session) ID() string {
	return s.id
}
