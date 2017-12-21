package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

var cookey = "session"

var lifespan = time.Hour * 24 * 14

// Manager manages sessions.
type Manager struct {
	store *store
}

// NewManager creates a new Manager.
func NewManager() (*Manager, error) {
	s, err := newStore()
	if err != nil {
		return nil, err
	}

	return &Manager{s}, nil
}

// Middleware creates a middlware function which adds session info to the request context.
func (m *Manager) Middleware(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s, err := getSession(m.store, r, w)
		if err != nil {
			fmt.Printf("error initalizing session from request: %v\n", err)
		}
		r = r.WithContext(context.WithValue(r.Context(), &cookey, s))
		next(w, r, p)
	})
}

// Load loads the session from the request context.
// Returns a value of nil if the session was not found.
func (m *Manager) Load(r *http.Request) *Session {
	val, _ := r.Context().Value(&cookey).(*Session)
	return val
}

// Close closes the connection to the store.
func (m *Manager) Close() error {
	return m.store.close()
}

// find or create a session for the current request.
func getSession(s *store, r *http.Request, w http.ResponseWriter) (*Session, error) {
	c, err := r.Cookie(cookey)
	if err != nil {
		return createSession(s, w)
	}

	return findSession(s, c.Value, r)
}

func createSession(s *store, w http.ResponseWriter) (*Session, error) {
	sessionID := uuid.NewV4().String()

	c := &http.Cookie{
		Name:   cookey,
		Value:  sessionID,
		MaxAge: int(lifespan.Seconds()),
	}

	err := s.create(sessionID)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, c)
	return &Session{
		id:    sessionID,
		store: s,
	}, nil
}

func findSession(s *store, id string, r *http.Request) (*Session, error) {
	err := s.touch(id)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:    id,
		store: s,
	}, nil
}
