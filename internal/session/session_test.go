package session

import (
	"testing"
)

func TestSessions(t *testing.T) {
	st, err := newStore()
	if err != nil {
		t.Fatalf("could not create new session store: %v", err)
	}

	se := Session{
		id:    "abc-def",
		store: st,
		errHandler: func(err error) {
			t.Errorf("store operation failed: %v", err)
		},
	}

	t.Run("Get", func(t *testing.T) {
		key := "key"
		value := "value"

		err := st.set(se.id, key, value)
		if err != nil {
			t.Fatalf("could not set value in store: %v", err)
		}

		res := se.Get(key)

		if res != value {
			t.Fatalf("returned value (%v) does not match \"%v\"", res, value)
		}
	})

	t.Run("Set", func(t *testing.T) {
		key := "key"
		value := "value"

		se.Set(key, value)

		res, err := st.get(se.id, key)
		if err != nil {
			t.Fatalf("could not get value from store: %v", err)
		}

		if res != value {
			t.Fatalf("assigned value (%v) does not match actual \"%v\"", value, res)
		}
	})
}
