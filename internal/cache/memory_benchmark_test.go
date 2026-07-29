package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const benchmarkKeyCount = 4096

// globalLockStore preserves the pre-sharding design as a benchmark baseline.
// It intentionally lives in a _test.go file so it is not included in the
// production binary.
type globalLockStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func newGlobalLockStore() *globalLockStore {
	return &globalLockStore{items: make(map[string]Item)}
}

func (s *globalLockStore) Set(key, value string, ttl time.Duration) error {
	item := Item{Value: value}
	if ttl > 0 {
		item.HasExpiry = true
		item.ExpiresAt = time.Now().Add(ttl)
	}

	s.mu.Lock()
	s.items[key] = item
	s.mu.Unlock()
	return nil
}

func (s *globalLockStore) Get(key string) (string, bool) {
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

func (s *globalLockStore) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[key]; !ok {
		return false
	}
	delete(s.items, key)
	return true
}

func BenchmarkStoreParallelReadDistributed(b *testing.B) {
	keys := benchmarkKeys()
	benchmarkStores(b, func(b *testing.B, store Store) {
		for _, key := range keys {
			if err := store.Set(key, "value", 0); err != nil {
				b.Fatal(err)
			}
		}

		var worker atomic.Uint64
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			operation := worker.Add(1)
			for pb.Next() {
				key := keys[operation%uint64(len(keys))]
				if _, ok := store.Get(key); !ok {
					b.Fatalf("expected %q to exist", key)
				}
				operation++
			}
		})
	})
}

func BenchmarkStoreParallelWriteDistributed(b *testing.B) {
	keys := benchmarkKeys()
	benchmarkStores(b, func(b *testing.B, store Store) {
		var worker atomic.Uint64
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			operation := worker.Add(1)
			for pb.Next() {
				key := keys[operation%uint64(len(keys))]
				if err := store.Set(key, "value", 0); err != nil {
					b.Fatal(err)
				}
				operation++
			}
		})
	})
}

func BenchmarkStoreParallelMixedDistributed(b *testing.B) {
	keys := benchmarkKeys()
	benchmarkStores(b, func(b *testing.B, store Store) {
		for _, key := range keys {
			if err := store.Set(key, "value", 0); err != nil {
				b.Fatal(err)
			}
		}

		var worker atomic.Uint64
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			operation := worker.Add(1)
			for pb.Next() {
				key := keys[operation%uint64(len(keys))]
				if operation%10 < 8 {
					if _, ok := store.Get(key); !ok {
						b.Fatalf("expected %q to exist", key)
					}
					operation++
					continue
				}
				if err := store.Set(key, "value", 0); err != nil {
					b.Fatal(err)
				}
				operation++
			}
		})
	})
}

func BenchmarkStoreParallelReadHotKey(b *testing.B) {
	benchmarkStores(b, func(b *testing.B, store Store) {
		if err := store.Set("hot-key", "value", 0); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if _, ok := store.Get("hot-key"); !ok {
					b.Fatal("expected hot-key to exist")
				}
			}
		})
	})
}

func benchmarkStores(b *testing.B, benchmark func(*testing.B, Store)) {
	b.Helper()
	b.Run("GlobalLock", func(b *testing.B) {
		benchmark(b, newGlobalLockStore())
	})
	b.Run("ShardedLock16", func(b *testing.B) {
		benchmark(b, NewInMemoryStoreWithShards(16, 0))
	})
}

func benchmarkKeys() []string {
	keys := make([]string, benchmarkKeyCount)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}
