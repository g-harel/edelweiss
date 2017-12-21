package session

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestManager(t *testing.T) {
	m, _ := NewManager()

	t.Run("Middleware", func(t *testing.T) {
		w := httptest.NewRecorder()
		sessionID := "1234-abcd"

		r1 := httptest.NewRequest("GET", "/", nil)
		r1.AddCookie(&http.Cookie{
			Name:  cookey,
			Value: sessionID,
		})
		m.store.set(sessionID, "id", sessionID)
		h1 := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			s := m.Load(r)
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
		m.Middleware(h1)(w, r1, httprouter.Params{})

		r2 := httptest.NewRequest("GET", "/", nil)
		h2 := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			s := m.Load(r)
			if s == nil {
				t.Errorf("session should be created when not in cookies")
			}
			val, err := m.store.client.Exists(sessionID).Result()
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
		m.Middleware(h2)(w, r2, httprouter.Params{})
	})

	t.Run("Load", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)

		s := m.Load(r)
		if s != nil {
			t.Errorf("load should return nil when no session is set")
		}

		s1 := &Session{
			id:    "abcd-1234",
			store: m.store,
		}
		r = r.WithContext(context.WithValue(r.Context(), &cookey, s1))
		s2 := m.Load(r)
		if s1 != s2 {
			t.Fatalf("load should fetch the stored session")
		}
	})
}
