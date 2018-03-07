package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

var cooKey = "session"
var contextKey = "session"

var lifespan = time.Hour * 24 * 14

// Manager manages sessions.
type Manager struct {
	store Store
}

// NewManager creates a new session manager.
func NewManager(s Store) *Manager {
	return &Manager{s}
}

// Load fetches or creates the session from the request context.
func (m *Manager) Load(c *gin.Context) (*Session, error) {
	val, exists := c.Get(contextKey)
	if exists {
		s, ok := val.(*Session)
		if !ok {
			return nil, fmt.Errorf("session could not be created from context value")
		}
		return s, nil
	}

	var s *Session

	ck, err := c.Request.Cookie(cooKey)
	if err != nil {
		s, err = m.createSession(c.Writer)
	} else {
		s, err = m.findSession(ck.Value)
		if err != nil {
			s, err = m.createSession(c.Writer)
		}
	}
	if err != nil {
		return nil, err
	}

	c.Set(contextKey, s)

	return s, nil
}

// Refresh erases the current session and starts a blank one.
// A session should always exist, so it is replace instead of
// just being destroyed.
func (m *Manager) Refresh(c *gin.Context) (*Session, error) {
	s, err := m.createSession(c.Writer)
	if err != nil {
		return nil, err
	}

	c.Set(contextKey, s)

	// Session data also deleted from store if found. An error in
	// fetching the cookie key is not a deal-breaker since the store
	// data will eventually expire.
	ck, err := c.Request.Cookie(cooKey)
	if err == nil {
		m.store.delete(ck.Value)
	}

	return s, nil
}

func (m *Manager) createSession(w http.ResponseWriter) (*Session, error) {
	sessionID := uuid.NewV4().String()

	err := m.store.create(sessionID)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cooKey,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(lifespan.Seconds()),
	})

	return &Session{
		id:    sessionID,
		store: m.store,
	}, nil
}

func (m *Manager) findSession(id string) (*Session, error) {
	err := m.store.touch(id)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:    id,
		store: m.store,
	}, nil
}
