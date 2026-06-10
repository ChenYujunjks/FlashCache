// internal/cache/memory.go	内存存储实现
package cache

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu              sync.RWMutex // 可以直接自定义
	items           map[string]Item
	cleanupInterval time.Duration
	stopCh          chan struct{}
}

func NewInMemoryStore(cleanupInterval time.Duration) *InMemoryStore {
	store := &InMemoryStore{
		items:           make(map[string]Item),
		cleanupInterval: cleanupInterval,
		stopCh:          make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go store.startCleanup()
	}

	return store
}

func (s *InMemoryStore) Set(key string, value string, ttl time.Duration) error {
	//所有访问 items 的代码，都约定先拿 mu 这把锁。
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

func (s *InMemoryStore) startCleanup() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.deleteExpiredItems()
		case <-s.stopCh:
			return
		}
	}
}

func (s *InMemoryStore) deleteExpiredItems() {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, item := range s.items {
		if item.HasExpiry && now.After(item.ExpiresAt) {
			delete(s.items, key)
			println("cleanup tick executed")
		}
	}
}

func (s *InMemoryStore) Stop() {
	close(s.stopCh)
}
