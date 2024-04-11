package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache()
	key := "testKey"
	value := "testValue"
	expiration := int64(5) // 5 seconds

	cache.Set(key, value, expiration)

	retrievedValue, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find value for key %s", key)
	}

	if retrievedValue != value {
		t.Errorf("Expected value %s, got %s", value, retrievedValue)
	}

	// Test expiration
	time.Sleep(6 * time.Second)
	_, foundAfterExpiration := cache.Get(key)
	if foundAfterExpiration {
		t.Errorf("Expected not to find value for key %s after expiration", key)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache()
	key := "testKey"
	value := "testValue"

	cache.Set(key, value, 0)
	cache.Delete(key)

	_, found := cache.Get(key)
	if found {
		t.Errorf("Expected not to find value for key %s after deletion", key)
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after clear, got %d", cache.Size())
	}
}

func TestCache_Size(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	size := cache.Size()
	if size != 2 {
		t.Errorf("Expected cache size to be 2, got %d", size)
	}
}
