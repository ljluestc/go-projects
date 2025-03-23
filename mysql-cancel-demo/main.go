package main

import (
    "context"
    "fmt"
    "time"

    mysqlgo "github.com/rocketlaunchr/mysql-go"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := mysqlgo.Open("user:password@tcp(localhost:3306)/testdb")
    if err != nil {
        fmt.Println("Error opening DB:", err)
        return
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    rows, err := db.QueryContext(ctx, "SELECT SLEEP(5)")
    if err != nil {
        fmt.Println("Query error:", err)
        return
    }
    defer rows.Close()

    fmt.Println("Query executed (should cancel due to timeout)")
}