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
	id         string
	store      sessionStorer
	errHandler errHandler
}

// Get fetches session data.
func (s *Session) Get(key string) string {
	val, err := s.store.get(s.id, key)
	if err != nil {
		s.errHandler(fmt.Errorf("error reading from store: %v", err))
		return ""
	}

	return val
}

// Set modifies session data.
func (s *Session) Set(key, value string) {
	err := s.store.set(s.id, key, value)
	if err != nil {
		s.errHandler(fmt.Errorf("error setting value in store: %v", err))
	}
}

// ID is a getter for the session's id property.
func (s *Session) ID() string {
	return s.id
}
