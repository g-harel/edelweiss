package session

import "testing"

func TestStore(t *testing.T) {
	s, err := newStore()
	if err != nil {
		t.Fatalf("could not create new store: %v", err)
	}

	t.Run("create", func(t *testing.T) {
		sessionID := "1234-abcd"

		err := s.create(sessionID)
		if err != nil {
			t.Fatalf("create operation failed: %v", err)
		}

		ttl, err := s.client.TTL(sessionID).Result()
		if err != nil {
			t.Fatalf("error reading session expiry: %v", err)
		}

		if ttl == -1 {
			t.Fatal("session expiry was not set")
		}
		if ttl == -2 {
			t.Fatal("session was not created")
		}
	})

	t.Run("delete", func(t *testing.T) {
		sessionID := "abcd-1234"

		_, err := s.client.HMSet(sessionID, map[string]interface{}{"test": 0}).Result()
		if err != nil {
			t.Fatalf("could not write to client: %v", err)
		}

		err = s.delete(sessionID)
		if err != nil {
			t.Fatalf("delete operation failed: %v", err)
		}

		res, err := s.client.Exists(sessionID).Result()
		if err != nil {
			t.Fatalf("could not read from client: %v", err)
		}

		if res > 0 {
			t.Fatal("session was not deleted")
		}

		err = s.delete("other-id")
		if err != nil {
			t.Fatalf("deleting a nonexistent session should not fail: %v", err)
		}
	})

	t.Run("touch", func(t *testing.T) {
		sessionID := "efgh-5678"

		_, err := s.client.HMSet(sessionID, map[string]interface{}{"test": 0}).Result()
		if err != nil {
			t.Fatalf("could not write to client: %v", err)
		}

		err = s.touch(sessionID)
		if err != nil {
			t.Fatalf("touch operation failed: %v", err)
		}

		ttl, err := s.client.TTL(sessionID).Result()
		if err != nil {
			t.Fatalf("error reading session expiry: %v", err)
		}
		if ttl == -1 {
			t.Fatal("session expiry was not set")
		}
	})

	t.Run("get", func(t *testing.T) {
		sessionID := "5678-efgh"
		value := "test"

		_, err := s.client.HMSet(sessionID, map[string]interface{}{
			"test": value,
		}).Result()
		if err != nil {
			t.Fatalf("could not write to client: %v", err)
		}

		res, err := s.get(sessionID, "test")
		if err != nil {
			t.Fatalf("get operation failed: %v", err)
		}

		if res != value {
			t.Fatal("hash key value does not match")
		}
	})

	t.Run("set", func(t *testing.T) {
		sessionID := "aceg-0248"
		value := "test"

		err := s.set(sessionID, "test", value)
		if err != nil {
			t.Fatalf("set operation failed: %v", err)
		}

		res, err := s.client.HGet(sessionID, "test").Result()
		if err != nil {
			t.Fatalf("could not read from hash key")
		}

		if res != value {
			t.Fatal("hash key value does not match")
		}
	})
}
