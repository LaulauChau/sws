package cache

import (
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache[[]string](1 * time.Hour)

	// Test empty cache
	data, ok := cache.Get()
	if ok {
		t.Error("Expected cache miss for empty cache")
	}
	if data != nil {
		t.Error("Expected nil data for empty cache")
	}

	// Test setting and getting cache
	testData := []string{"test1", "test2", "test3"}

	cache.Set(testData)
	data, ok = cache.Get()
	if !ok {
		t.Error("Expected cache hit after setting data")
	}
	if len(data) != len(testData) {
		t.Errorf("Expected %d items, got %d", len(testData), len(data))
	}

	// Test cache expiration
	cache = NewCache[[]string](1 * time.Millisecond)
	cache.Set(testData)
	time.Sleep(2 * time.Millisecond)
	data, ok = cache.Get()
	if ok {
		t.Error("Expected cache miss after expiration")
	}
	if data != nil {
		t.Error("Expected nil data after expiration")
	}
}

// Test with a struct type
type testStruct struct {
	ID   int
	Name string
}

func TestCacheWithStruct(t *testing.T) {
	cache := NewCache[testStruct](1 * time.Hour)

	test := testStruct{
		ID:   1,
		Name: "test",
	}

	cache.Set(test)
	data, ok := cache.Get()
	if !ok {
		t.Error("Expected cache hit after setting struct")
	}
	if data.ID != test.ID || data.Name != test.Name {
		t.Error("Cached struct data doesn't match original")
	}
}

func TestCache_Basic(t *testing.T) {
	t.Run("new cache is empty", func(t *testing.T) {
		cache := NewCache[string](time.Hour)
		_, ok := cache.Get()
		if ok {
			t.Error("new cache should be empty")
		}
	})

	t.Run("can set and get value", func(t *testing.T) {
		cache := NewCache[string](time.Hour)
		cache.Set("test value")

		value, ok := cache.Get()
		if !ok {
			t.Error("should get value after setting")
		}
		if value != "test value" {
			t.Errorf("got %q, want %q", value, "test value")
		}
	})

	t.Run("expires after timeout", func(t *testing.T) {
		cache := NewCache[int](100 * time.Millisecond)
		cache.Set(42)

		time.Sleep(200 * time.Millisecond)

		_, ok := cache.Get()
		if ok {
			t.Error("cache should expire after timeout")
		}
	})
}

func TestCache_Concurrent(t *testing.T) {
	cache := NewCache[int](time.Hour)
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrently write values
	for i := 0; i < goroutines; i++ {
		go func(val int) {
			defer wg.Done()
			cache.Set(val)
		}(i)
	}

	// Concurrently read values while writing
	for i := 0; i < goroutines; i++ {
		go func() {
			_, _ = cache.Get()
		}()
	}

	wg.Wait()
}

func TestCache_Types(t *testing.T) {
	t.Run("works with struct", func(t *testing.T) {
		type person struct {
			Name string
			Age  int
		}

		cache := NewCache[person](time.Hour)
		p := person{Name: "Alice", Age: 30}
		cache.Set(p)

		got, ok := cache.Get()
		if !ok {
			t.Fatal("should get value")
		}
		if got != p {
			t.Errorf("got %+v, want %+v", got, p)
		}
	})

	t.Run("works with slice", func(t *testing.T) {
		cache := NewCache[[]int](time.Hour)
		data := []int{1, 2, 3}
		cache.Set(data)

		got, ok := cache.Get()
		if !ok {
			t.Fatal("should get value")
		}
		if len(got) != len(data) {
			t.Errorf("got len %d, want %d", len(got), len(data))
		}
		for i := range data {
			if got[i] != data[i] {
				t.Errorf("at index %d: got %d, want %d", i, got[i], data[i])
			}
		}
	})
}
