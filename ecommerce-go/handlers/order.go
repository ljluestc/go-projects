package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/juicedata/juicefs/ecommerce-go/db"
    "github.com/juicedata/juicefs/ecommerce-go/models"
)

func CreateOrder(c *fiber.Ctx) error {
    var order models.Order
    if err := c.BodyParser(&order); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    result, err := db.DB.Exec(
        "INSERT INTO orders (user_id, total, status, currency, shipping, discount) VALUES (?, ?, ?, ?, ?, ?)",
        order.UserID, order.Total, "pending", order.Currency, order.Shipping, order.Discount,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    id, _ := result.LastInsertId()
    order.ID = int(id)
    return c.JSON(order)
}