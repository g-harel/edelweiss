package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestManager(t *testing.T) {
	m, _ := NewManager()

	t.Run("Middleware", func(t *testing.T) {
		w := httptest.NewRecorder()

		router := gin.Default()
		router.Use(m.Middleware)

		sessionID := "1234-abcd"
		req1, _ := http.NewRequest("GET", "/req1", nil)
		req1.AddCookie(&http.Cookie{
			Name:  cookey,
			Value: sessionID,
		})
		err := m.store.set(sessionID, "id", sessionID)
		if err != nil {
			t.Errorf("could not set value in session store")
		}
		router.GET("/req1", func(c *gin.Context) {
			s := m.Load(c)
			if s.id != sessionID {
				t.Errorf("session should be read from request cookies")
			}
			val, err := m.store.client.TTL(sessionID).Result()
			if err != nil {
				t.Errorf("could not read ttl from session store client")
			}
			if val < 0 {
				t.Errorf("existing session ttl should be reset")
			}
		})
		router.ServeHTTP(w, req1)

		req2, _ := http.NewRequest("GET", "/req2", nil)
		router.GET("/req2", func(c *gin.Context) {
			s := m.Load(c)
			if s == nil {
				t.Errorf("session should be created when not in cookies")
			}
			val, err := m.store.client.Exists(s.id).Result()
			if err != nil {
				t.Errorf("could not check session store client")
			}
			if val < 1 {
				t.Errorf("session not created in store")
			}
			header := w.Header().Get("Set-Cookie")
			if !strings.Contains(header, s.id) {
				t.Errorf("set-cookie header was not added to response")
			}
		})
		router.ServeHTTP(w, req2)
	})

	t.Run("Load", func(t *testing.T) {
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		s1 := m.Load(c)
		if s1 != nil {
			t.Errorf("load should return nil when no session is set")
		}

		s2 := &Session{
			id:    "abcd-1234",
			store: m.store,
		}
		c.Set("session", s2)
		s3 := m.Load(c)
		if s2 != s3 {
			t.Fatalf("load should fetch the stored session")
		}
	})
}
