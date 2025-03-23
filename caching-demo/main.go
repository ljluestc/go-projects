package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/rocketlaunchr/remember-go"
    "github.com/rocketlaunchr/remember-go/memory"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Key struct for generating cache keys
type Key struct {
    Search string
    Page   int `json:"page"`
}

// Result struct for query results
type Result struct {
    Title string
}

// initDB initializes an in-memory SQLite database with sample data
func initDB() *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }

    createTable := `
        CREATE TABLE books (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL
        );
    `
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }

    insertData := `
        INSERT INTO books (title) VALUES
        ('Golang Basics'),
        ('Advanced Golang'),
        ('Golang Web Development'),
        ('Python for Beginners'),
        ('Golang Concurrency');
    `
    _, err = db.Exec(insertData)
    if err != nil {
        log.Fatalf("Failed to insert data: %v", err)
    }

    return db
}

// slowQuery creates a SlowRetrieve function for database queries
func slowQuery(db *sql.DB, key Key) remember.SlowRetrieve {
    return func(ctx context.Context) (interface{}, error) {
        // Simulate slow query
        time.Sleep(2 * time.Second)

        results := []Result{}
        stmt := `
            SELECT title
            FROM books
            WHERE title LIKE ?
            ORDER BY title
            LIMIT ?, 20
        `
        rows, err := db.QueryContext(ctx, stmt, "%"+key.Search+"%", (key.Page-1)*20)
        if err != nil {
            return nil, fmt.Errorf("query failed: %v", err)
        }
        defer rows.Close()

        for rows.Next() {
            var title string
            if err := rows.Scan(&title); err != nil {
                return nil, fmt.Errorf("scan failed: %v", err)
            }
            results = append(results, Result{Title: title})
        }
        return results, nil
    }
}

// fetchBooks retrieves books, using cache if available
func fetchBooks(ctx context.Context, db *sql.DB, ms *memory.MemoryStore, key Key, exp time.Duration) ([]Result, bool, error) {
    cacheKey := remember.CreateKeyStruct(key)
    results, found, err := remember.Cache(ctx, ms, cacheKey, exp, slowQuery(db, key), remember.Options{GobRegister: false})
    if err != nil {
        return nil, false, err
    }
    return results.([]Result), found, nil
}

func main() {
    db := initDB()
    defer db.Close()

    ms := memory.NewMemoryStore(10 * time.Minute)
    ctx := context.Background()
    key := Key{Search: "golang", Page: 1}
    exp := 10 * time.Minute

    // First call: hits the database
    fmt.Println("First query (should be slow):")
    start := time.Now()
    results, found, err := fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Found in cache: %v\n", found)
    for _, r := range results {
        fmt.Println(r.Title)
    }
    fmt.Printf("Time taken: %v\n\n", time.Since(start))

    // Second call: hits the cache
    fmt.Println("Second query (should be fast):")
    start = time.Now()
    results, found, err = fetchBooks(ctx, db, ms, key, exp)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Found in cache: %v\n", found)
    for _, r := range results {
        fmt.Println(r.Title)
    }
    fmt.Printf("Time taken: %v\n", time.Since(start))
}