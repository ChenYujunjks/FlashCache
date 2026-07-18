package cache

import (
	"hash/fnv"
	"sync"
	"time"
)

const defaultShardCount = 16

type shard struct {
	mu    sync.RWMutex
	items map[string]Item
}

type InMemoryStore struct {
	shards          []*shard
	shardCount      uint32
	cleanupInterval time.Duration
	stopCh          chan struct{}
}

func NewInMemoryStore(cleanupInterval time.Duration) *InMemoryStore {
	return NewInMemoryStoreWithShards(defaultShardCount, cleanupInterval)
}

func NewInMemoryStoreWithShards(shardCount int, cleanupInterval time.Duration) *InMemoryStore {
	if shardCount <= 0 {
		shardCount = defaultShardCount
	}

	store := &InMemoryStore{
		shards:          make([]*shard, shardCount),
		shardCount:      uint32(shardCount),
		cleanupInterval: cleanupInterval,
		stopCh:          make(chan struct{}),
	}

	for i := 0; i < shardCount; i++ {
		store.shards[i] = &shard{
			items: make(map[string]Item),
		}
	}

	if cleanupInterval > 0 {
		go store.startCleanup()
	}

	return store
}

func (s *InMemoryStore) Set(key string, value string, ttl time.Duration) error {
	sh := s.getShard(key)

	sh.mu.Lock()
	defer sh.mu.Unlock()

	item := Item{
		Value: value,
	}

	if ttl > 0 {
		item.HasExpiry = true
		item.ExpiresAt = time.Now().Add(ttl)
	}

	sh.items[key] = item
	return nil
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	sh := s.getShard(key)

	sh.mu.RLock()
	item, ok := sh.items[key]
	sh.mu.RUnlock()

	if !ok {
		return "", false
	}

	if item.HasExpiry && time.Now().After(item.ExpiresAt) {
		sh.mu.Lock()
		delete(sh.items, key)
		sh.mu.Unlock()
		return "", false
	}

	return item.Value, true
}

func (s *InMemoryStore) Delete(key string) bool {
	sh := s.getShard(key)

	sh.mu.Lock()
	defer sh.mu.Unlock()

	if _, ok := sh.items[key]; !ok {
		return false
	}

	delete(sh.items, key)
	return true
}

func (s *InMemoryStore) getShard(key string) *shard {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	index := hash.Sum32() % s.shardCount
	return s.shards[index]
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

	for _, sh := range s.shards {
		sh.mu.Lock()

		for key, item := range sh.items {
			if item.HasExpiry && now.After(item.ExpiresAt) {
				delete(sh.items, key)
			}
		}

		sh.mu.Unlock()
	}
}

func (s *InMemoryStore) Stop() {
	close(s.stopCh)
}
