package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/juicedata/juicefs/ecommerce-go/db"
    "github.com/juicedata/juicefs/ecommerce-go/models"
)

func GetProducts(c *fiber.Ctx) error {
    rows, err := db.DB.Query("SELECT id, name, price, type, description, image, stock FROM products")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Type, &p.Description, &p.Image, &p.Stock); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        products = append(products, p)
    }
    return c.JSON(products)
}