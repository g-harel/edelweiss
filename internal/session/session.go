package session

import (
	"fmt"
	"time"
)

var lifespan = time.Hour * 24 * 14

type sessionStorer interface {
	get(ID, key string) (string, error)
	set(ID, key, value string) error
}

// Session contains identifying information about the user.
type Session struct {
	ID    string
	store sessionStorer
}

// Get fetches session data.
func (s *Session) Get(key string) (string, error) {
	val, err := s.store.get(s.ID, key)
	if err != nil {
		return val, fmt.Errorf("error reading from store: %v", err)
	}

	return val, nil
}

// Set modifies session data.
func (s *Session) Set(key, value string) error {
	err := s.store.set(s.ID, key, value)
	if err != nil {
		return fmt.Errorf("error setting value in store: %v", err)
	}

	return nil
}
