package main

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/rocketlaunchr/remember-go"
    "github.com/rocketlaunchr/remember-go/memory"
    "github.com/rocketlaunchr/remember-go/redis"
    "github.com/gomodule/redigo/redis"
    "github.com/stretchr/testify/assert"
)

func TestSlowQueryCache(t *testing.T) {
    // Initialize database
    db := initDB()
    defer db.Close()

    // Test cases
    tests := []struct {
        name       string
        store      remember.StorageDriver
        key        Key
        wantFound  bool
        wantTitles []string
    }{
        {
            name:      "In-Memory Cache - First Run",
            store:     memory.NewMemoryStore(1 * time.Second),
            key:       Key{Search: "Golang", Page: 1},
            wantFound: false,
            wantTitles: []string{"Golang Advanced", "Golang Basics"},
        },
        {
            name:      "Redis Cache - First Run",
            store:     redis.NewRedisStore(&redis.Pool{Dial: func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") }}),
            key:       Key{Search: "Golang", Page: 1},
            wantFound: false,
            wantTitles: []string{"Golang Advanced", "Golang Basics"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.WithValue(context.Background(), "key", tt.key)
            cacheKey := remember.CreateKeyStruct(tt.key)
            exp := 1 * time.Second

            // First run - should not be found in cache
            results, found, err := remember.Cache(ctx, tt.store, cacheKey, exp, slowQuery(db), remember.Options{GobRegister: false})
            assert.NoError(t, err)
            assert.Equal(t, tt.wantFound, found)

            res := results.([]Result)
            var titles []string
            for _, r := range res {
                titles = append(titles, r.Title)
            }
            assert.Equal(t, tt.wantTitles, titles)

            // Second run - should be found in cache
            results, found, err = remember.Cache(ctx, tt.store, cacheKey, exp, slowQuery(db), remember.Options{GobRegister: false})
            assert.NoError(t, err)
            assert.True(t, found)

            res = results.([]Result)
            titles = []string{}
            for _, r := range res {
                titles = append(titles, r.Title)
            }
            assert.Equal(t, tt.wantTitles, titles)

            // Wait for cache to expire and check again
            time.Sleep(2 * time.Second)
            results, found, err = remember.Cache(ctx, tt.store, cacheKey, exp, slowQuery(db), remember.Options{GobRegister: false})
            assert.NoError(t, err)
            assert.False(t, found) // Cache should have expired
        })
    }
}