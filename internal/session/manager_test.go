package session

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestManager(t *testing.T) {
	m := NewManager(&MockStore{})

	t.Run("Load", func(t *testing.T) {
		t.Run("new session", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "http://example.com/", nil)

			s, err := m.Load(c)
			if err != nil {
				t.Errorf("error creating session")
			}

			_, exists := c.Get(contextKey)
			if !exists {
				t.Errorf("session not added to context after being created")
			}

			header := w.Header().Get("Set-Cookie")
			if !strings.Contains(header, s.id) {
				t.Errorf("set-cookie header was not added to response")
			}
		})

		t.Run("session from cookie", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", "http://example.com/", nil)

			_, err := m.Load(c)
			if err != nil {
				t.Errorf("error loading session from cookie")
			}

			_, exists := c.Get(contextKey)
			if !exists {
				t.Errorf("session not added to context after being created")
			}
		})

		t.Run("session from context", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", "http://example.com/", nil)

			original := &Session{
				id:    "abcd-1234",
				store: m.store,
			}
			c.Set(contextKey, original)

			s, err := m.Load(c)
			if err != nil {
				t.Fatalf("error loading session from context")
			}
			if original != s {
				t.Fatalf("stored session does not match original")
			}
		})
	})
}
