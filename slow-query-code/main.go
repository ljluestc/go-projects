package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/rocketlaunchr/remember-go"
    "github.com/rocketlaunchr/remember-go/memory"
    "github.com/rocketlaunchr/remember-go/redis"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
    "github.com/gomodule/redigo/redis"
)

// Key struct for generating cache keys
type Key struct {
    Search string `json:"search"`
    Page   int    `json:"page"`
}

// Result struct for query results
type Result struct {
    Title string `json:"title"`
}

// SlowRetrieve function type from remember-go
type SlowRetrieve func(ctx context.Context) (interface{}, error)

// Initialize SQLite database for demo
func initDB() *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    _, err = db.Exec(`
        CREATE TABLE books (title TEXT);
        INSERT INTO books (title) VALUES ('Golang Basics'), ('Golang Advanced'), ('Python Basics');
    `)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

// slowQuery simulates a slow database query
func slowQuery(db *sql.DB) SlowRetrieve {
    return func(ctx context.Context) (interface{}, error) {
        results := []Result{}
        key := ctx.Value("key").(Key)
        stmt := `SELECT title FROM books WHERE title LIKE ? ORDER BY title LIMIT ?, 20`
        rows, err := db.QueryContext(ctx, stmt, "%"+key.Search+"%", (key.Page-1)*20)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        for rows.Next() {
            var title string
            if err := rows.Scan(&title); err != nil {
                return nil, err
            }
            results = append(results, Result{Title: title})
        }
        // Simulate slowness
        time.Sleep(100 * time.Millisecond)
        return results, nil
    }
}

func main() {
    // Initialize database
    db := initDB()
    defer db.Close()

    // Initialize storage drivers
    memoryStore := memory.NewMemoryStore(10 * time.Minute)
    redisStore := redis.NewRedisStore(&redis.Pool{
        Dial: func() (redis.Conn, error) {
            return redis.Dial("tcp", "localhost:6379")
        },
    })

    // Example query
    key := Key{Search: "Golang", Page: 1}
    ctx := context.WithValue(context.Background(), "key", key)
    cacheKey := remember.CreateKeyStruct(key)
    exp := 10 * time.Minute

    // Using in-memory cache
    results, found, err := remember.Cache(ctx, memoryStore, cacheKey, exp, slowQuery(db), remember.Options{GobRegister: false})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("In-Memory Cache - Found: %v, Results: %+v\n", found, results.([]Result)))

    // Using Redis cache
    results, found, err = remember.Cache(ctx, redisStore, cacheKey, exp, slowQuery(db), remember.Options{GobRegister: false})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Redis Cache - Found: %v, Results: %+v\n", found, results.([]Result))
}