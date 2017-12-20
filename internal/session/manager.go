package session

import (
	"net/http"

	"github.com/satori/go.uuid"
)

var cookey = "edelweiss-session-id"

// Manager facilitates creating and fetching user sessions.
type Manager struct {
	store *store
}

// NewManager creates a new session manager.
func NewManager() (*Manager, error) {
	s, err := newStore()
	if err != nil {
		return nil, err
	}

	return &Manager{s}, nil
}

// Start starts a new session.
func (m *Manager) Start(w http.ResponseWriter) (*Session, error) {
	sessionID := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:   cookey,
		Value:  sessionID,
		MaxAge: int(lifespan.Seconds()),
	})

	session := &Session{
		ID:    sessionID,
		store: m.store,
	}
	return session, nil
}

// Load finds the existing session.
func (m *Manager) Load(r *http.Request) (*Session, error) {
	c, err := r.Cookie(cookey)
	if err != nil {
		return nil, err
	}

	sessionID := c.Value

	session := &Session{
		ID:    sessionID,
		store: m.store,
	}
	return session, nil
}
