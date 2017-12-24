package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var cookey = "session"

var lifespan = time.Hour * 24 * 14

type errHandler func(error)

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

// Close closes the connection to the store.
func (m *Manager) Close() error {
	return m.store.close()
}

// Middleware creates a middlware function which adds session info to the request context.
func (m *Manager) Middleware(c *gin.Context) {
	e := errHandler(func(err error) {
		handleErr(c, err)
	})
	s := getSession(m, c.Request, c.Writer, e)
	c.Set("session", s)
	c.Next()
}

// Load fetches the session from the request context.
func (m *Manager) Load(c *gin.Context) *Session {
	val, exists := c.Get("session")
	if exists == false {
		handleErr(c, fmt.Errorf("session was not found in context"))
		return nil
	}
	s, ok := val.(*Session)
	if ok != true {
		handleErr(c, fmt.Errorf("session could not be created from context value"))
		return nil
	}
	return s
}

func handleErr(c *gin.Context, err error) {
	c.AbortWithError(500, err)
}

// find or create a session for the current request.
func getSession(m *Manager, r *http.Request, w http.ResponseWriter, e errHandler) *Session {
	cookie, foundErr := r.Cookie(cookey)
	if foundErr != nil {
		return createSession(m, w, e)
	}

	return findSession(m, cookie.Value, e)
}

func createSession(m *Manager, w http.ResponseWriter, e errHandler) *Session {
	sessionID := uuid.NewV4().String()

	err := m.store.create(sessionID)
	if err != nil {
		e(err)
	}

	cookie := &http.Cookie{
		Name:     cookey,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(lifespan.Seconds()),
	}
	http.SetCookie(w, cookie)

	return &Session{
		id:         sessionID,
		store:      m.store,
		errHandler: e,
	}
}

func findSession(m *Manager, id string, e errHandler) *Session {
	err := m.store.touch(id)
	if err != nil {
		e(err)
	}

	return &Session{
		id:         id,
		store:      m.store,
		errHandler: e,
	}
}
