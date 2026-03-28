package cache

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		items: make(map[string]Item),
	}
}

func (s *InMemoryStore) Set(key string, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := Item{
		Value: value,
	}

	if ttl > 0 {
		item.HasExpiry = true
		item.ExpiresAt = time.Now().Add(ttl)
	}

	s.items[key] = item
	return nil
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	item, ok := s.items[key]
	s.mu.RUnlock()

	if !ok {
		return "", false
	}

	if item.HasExpiry && time.Now().After(item.ExpiresAt) {
		s.mu.Lock()
		delete(s.items, key)
		s.mu.Unlock()
		return "", false
	}

	return item.Value, true
}

func (s *InMemoryStore) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[key]; !ok {
		return false
	}

	delete(s.items, key)
	return true
}
