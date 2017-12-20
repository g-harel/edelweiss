package session

import (
	"testing"
)

func TestManager(t *testing.T) {
	/* m, err := NewManager()
	if err != nil {
		t.Fatalf("could not create new manager: %s", err)
	}

	t.Run("Start", func(t *testing.T) {
		w := httptest.NewRecorder()
		s, err := m.Start(w)
		if err != nil {
			t.Fatalf("error starting new session: %v", err)
		}

		// check that a session cookie was set in response
		var match bool
		for _, c := range w.Result().Cookies() {
			if c.Name == cookey && c.MaxAge == int(lifespan.Seconds()) {
				match = true
			}
		}
		if match != true {
			t.Errorf("cookie was not properly set")
		}

		// check that map was created in redis with the correct timeout
		ttl, err := m.client.TTL(s.ID).Result()
		if err != nil {
			t.Fatalf("session storage was not created: %v", err)
		}
		if ttl > lifespan {
			t.Errorf("session timeout does not match %v vs. %v", ttl, lifespan)
		}
	})

	t.Run("Load", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)

		sessionID := "session-id"

		r.AddCookie(&http.Cookie{
			Name:   cookey,
			Value:  sessionID,
			MaxAge: int(lifespan.Seconds()),
		})
		m.client.Set(sessionID, "", lifespan+time.Minute)

		// check that correct session is loaded
		s, err := m.Load(r)
		if err != nil {
			t.Fatalf("could not load session: %v", err)
		}
		if s.ID != sessionID || s.client != m.client {
			t.Errorf("created invalid session")
		}

		// check that redis key expiration was extended
		ttl, err := m.client.TTL(sessionID).Result()
		if err != nil {
			t.Fatalf("%s", err)
		}
		if ttl > lifespan {
			t.Errorf("session timeout does not match %v vs. %v", ttl, lifespan)
		}
	}) */
}
