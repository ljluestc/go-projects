package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "sync"
)

type Store struct {
    data map[string]string
    mu   sync.RWMutex
}

func NewStore() *Store {
    return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

func (s *Store) Get(key string) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data[key]
}

func handleClient(conn net.Conn, store *Store) {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        parts := strings.Fields(scanner.Text())
        switch strings.ToUpper(parts[0]) {
        case "SET":
            if len(parts) == 3 {
                store.Set(parts[1], parts[2])
                fmt.Fprintln(conn, "+OK")
            }
        case "GET":
            if len(parts) == 2 {
                value := store.Get(parts[1])
                if value == "" {
                    fmt.Fprintln(conn, "$-1")
                } else {
                    fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
                }
            }
        }
    }
}

func main() {
    store := NewStore()
    listener, err := net.Listen("tcp", ":6379")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Println("Redis Server running on :6379")

    // Persistence (simplified)
    go func() {
        for range time.Tick(10 * time.Second) {
            store.mu.RLock()
            data, _ := json.Marshal(store.data)
            store.mu.RUnlock()
            os.WriteFile("dump.json", data, 0644)
        }
    }()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleClient(conn, store)
    }
}