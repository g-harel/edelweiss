package session

// Session contains identifying information about the user.
type Session struct {
	id    string
	store Store
}

// Get fetches session data.
func (s *Session) Get(key string) (string, error) {
	return s.store.get(s.id, key)
}

// Set modifies session data.
func (s *Session) Set(key, value string) error {
	return s.store.set(s.id, key, value)
}
