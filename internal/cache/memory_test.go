package cache

import (
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
	store := NewInMemoryStore(50 * time.Millisecond)
	defer store.Stop()

	_ = store.Set("temp", "value", 50*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	store.mu.RLock()
	_, exists := store.items["temp"]
	store.mu.RUnlock()

	if exists {
		t.Fatal("expected expired key to be removed by background cleanup")
	}
}
