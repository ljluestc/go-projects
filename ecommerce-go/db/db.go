package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/juicedata/juicefs/ecommerce-go/models"
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./ecommerce.db")
    if err != nil {
        return err
    }

    // Create tables
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            price REAL,
            type TEXT,
            description TEXT,
            image TEXT,
            stock INTEGER
        );
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id TEXT,
            total REAL,
            status TEXT,
            currency TEXT,
            shipping REAL,
            discount REAL
        );
    `)
    if err != nil {
        return err
    }

    // Seed initial products
    _, err = DB.Exec(`
        INSERT OR IGNORE INTO products (id, name, price, type, description, image, stock) VALUES
        (1, 'Hot Air Balloons (Pack of 2)', 235.99, 'physical', 'Glide through the air!', 'balloons.jpg', 10),
        (2, 'Flood Survival Guide (PDF)', 19.99, 'digital', 'Learn to survive floods.', 'guide.pdf', 100);
    `)
    return err
}