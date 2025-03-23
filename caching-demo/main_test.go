package main

import (
    "context"
    "testing"
    "time"

    "github.com/rocketlaunchr/remember-go/memory"
)

func TestFetchBooks(t *testing.T) {
    db := initDB()
    defer db.Close()

    ms := memory.NewMemoryStore(10 * time.Minute)
    ctx := context.Background()
    key := Key{Search: "golang", Page: 1}
    exp := 10 * time.Second

    // First call: should hit database
    start := time.Now()
    results, found, err := fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        t.Fatalf("fetchBooks failed: %v", err)
    }
    if found {
        t.Error("Expected first call to not be found in cache")
    }
    if len(results) != 4 { // Expect 4 "golang" titles
        t.Errorf("Expected 4 results, got %d", len(results))
    }
    duration := time.Since(start)
    if duration < 2*time.Second {
        t.Errorf("Expected slow query (>2s), took %v", duration)
    }

    // Second call: should hit cache
    start = time.Now()
    results, found, err = fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        t.Fatalf("fetchBooks failed: %v", err)
    }
    if !found {
        t.Error("Expected second call to be found in cache")
    }
    if len(results) != 4 {
        t.Errorf("Expected 4 results, got %d", len(results))
    }
    duration = time.Since(start)
    if duration > 10*time.Millisecond {
        t.Errorf("Expected fast query (<10ms), took %v", duration)
    }
}

func TestCacheExpiration(t *testing.T) {
    db := initDB()
    defer db.Close()

    ms := memory.NewMemoryStore(10 * time.Minute)
    ctx := context.Background()
    key := Key{Search: "golang", Page: 1}
    exp := 1 * time.Second

    // First call
    _, found, err := fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        t.Fatalf("fetchBooks failed: %v", err)
    }
    if found {
        t.Error("Expected first call to not be found in cache")
    }

    // Wait for cache to expire
    time.Sleep(2 * time.Second)

    // Second call after expiration
    _, found, err = fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        t.Fatalf("fetchBooks failed: %v", err)
    }
    if found {
        t.Error("Expected cache to expire, but data was found")
    }
}