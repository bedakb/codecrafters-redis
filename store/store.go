package store

import (
	"sync"
	"time"
)

type entry struct {
	data string
	ts   time.Time
}

// Store is a concurrent-safe key value store.
//
// Each key in the store could have a custom expiration time.
// The expiration mechanisim this implementation using is called lazy expiration which means that key will be kept
// in the memory and removed once accessed, in case the key expired.
type Store struct {
	values map[string]entry
	mu     sync.Mutex
}

// New creates a new instance of Store.
func New() *Store {
	return &Store{
		values: make(map[string]entry),
	}
}

// Get returns the actual value for the given key and value indicating the existence.
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if ok && !v.ts.IsZero() {
		if v.ts.Before(time.Now()) {
			delete(s.values, key)
			return "", false
		}
	}
	return v.data, ok
}

// Set sets the value for the given key.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = entry{
		data: value,
	}
}

// SetWithExpirity sets the value for the key and assigns given expirity.
func (s *Store) SetWithExpirity(key, value string, expirity time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = entry{
		data: value,
		ts:   time.Now().Add(expirity),
	}
}

// Len returns length of the store.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Lock()
	return len(s.values)
}
