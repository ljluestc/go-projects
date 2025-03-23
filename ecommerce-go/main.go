package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/juicedata/juicefs/ecommerce-go/db"
    "github.com/juicedata/juicefs/ecommerce-go/handlers"
)

func main() {
    app := fiber.New()

    // Initialize database
    if err := db.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Routes
    app.Get("/products", handlers.GetProducts)
    app.Post("/cart/add", handlers.AddToCart)
    app.Get("/cart/:userID", handlers.GetCart)
    app.Post("/order/create", handlers.CreateOrder)

    // Start server
    log.Fatal(app.Listen(":8080"))
}