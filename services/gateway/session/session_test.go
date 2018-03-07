package session

import (
	"testing"
)

func TestSessions(t *testing.T) {
	st := &MockStore{}

	se := Session{
		id:    "abc-def",
		store: st,
	}

	t.Run("Get", func(t *testing.T) {
		key := "key"
		value := "value"

		st.set(se.id, key, value)

		res, _ := se.Get(key)

		if res != value {
			t.Fatalf("returned value (%v) does not match \"%v\"", res, value)
		}
	})

	t.Run("Set", func(t *testing.T) {
		key := "key"
		value := "value"

		se.Set(key, value)

		res, _ := st.get(se.id, key)

		if res != value {
			t.Fatalf("assigned value (%v) does not match actual \"%v\"", value, res)
		}
	})
}
