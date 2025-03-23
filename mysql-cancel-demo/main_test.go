package main

import (
    "context"
    "testing"
    "time"

    mysqlgo "github.com/rocketlaunchr/mysql-go"
    "github.com/stretchr/testify/assert"
)

func TestQueryCancellation(t *testing.T) {
    db, err := mysqlgo.Open("user:password@tcp(localhost:3306)/testdb")
    assert.NoError(t, err)
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    _, err = db.QueryContext(ctx, "SELECT SLEEP(3)")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "context deadline exceeded")
}