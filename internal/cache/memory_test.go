package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	store := NewInMemoryStore(0)

	err := store.Set("name", "bruce", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	value, ok := store.Get("name")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "bruce" {
		t.Fatalf("expected value bruce, got %s", value)
	}
}

func TestDelete(t *testing.T) {
	store := NewInMemoryStore(0)

	_ = store.Set("name", "bruce", 0)

	deleted := store.Delete("name")
	if !deleted {
		t.Fatal("expected delete to return true")
	}

	_, ok := store.Get("name")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestTTLExpirationOnGet(t *testing.T) {
	store := NewInMemoryStore(0)

	_ = store.Set("temp", "value", 100*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	_, ok := store.Get("temp")
	if ok {
		t.Fatal("expected key to expire")
	}
}

func TestBackgroundCleanup(t *testing.T) {
	store := NewInMemoryStore(10 * time.Millisecond)
	defer store.Stop()

	const key = "temp"
	_ = store.Set(key, "value", 20*time.Millisecond)

	if !waitUntil(500*time.Millisecond, func() bool { return !itemExists(store, key) }) {
		t.Fatal("expected expired key to be removed by background cleanup")
	}
}

func TestNewInMemoryStoreWithShards(t *testing.T) {
	store := NewInMemoryStoreWithShards(4, 0)

	if store.shardCount != 4 {
		t.Fatalf("expected 4 shards, got %d", store.shardCount)
	}
	if len(store.shards) != 4 {
		t.Fatalf("expected 4 initialized shards, got %d", len(store.shards))
	}
	for i, sh := range store.shards {
		if sh == nil || sh.items == nil {
			t.Fatalf("expected shard %d to be initialized", i)
		}
	}
}

func TestMultipleKeysAcrossShards(t *testing.T) {
	store := NewInMemoryStoreWithShards(16, 0)
	values := make(map[string]string)
	usedShards := make(map[*shard]struct{})

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		values[key] = value
		usedShards[store.getShard(key)] = struct{}{}
		if err := store.Set(key, value, 0); err != nil {
			t.Fatalf("set %q: %v", key, err)
		}
	}

	if len(usedShards) < 2 {
		t.Fatal("expected test keys to be distributed across multiple shards")
	}
	for key, want := range values {
		got, ok := store.Get(key)
		if !ok || got != want {
			t.Fatalf("get %q: got (%q, %t), want (%q, true)", key, got, ok, want)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewInMemoryStoreWithShards(16, 0)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			if err := store.Set(key, value, 0); err != nil {
				t.Errorf("set %q: %v", key, err)
				return
			}
			if got, ok := store.Get(key); !ok || got != value {
				t.Errorf("get %q: got (%q, %t), want (%q, true)", key, got, ok, value)
				return
			}
			if !store.Delete(key) {
				t.Errorf("delete %q: expected true", key)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < workers; i++ {
		key := fmt.Sprintf("key-%d", i)
		if _, ok := store.Get(key); ok {
			t.Fatalf("expected %q to be deleted", key)
		}
	}
}

func itemExists(store *InMemoryStore, key string) bool {
	sh := store.getShard(key)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	_, exists := sh.items[key]
	return exists
}

func waitUntil(timeout time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return condition()
}
