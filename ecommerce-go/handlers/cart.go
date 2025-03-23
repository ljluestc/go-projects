package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/juicedata/juicefs/ecommerce-go/models"
)

var carts = make(map[string]models.Cart) // In-memory cart storage

func AddToCart(c *fiber.Ctx) error {
    var item models.CartItem
    if err := c.BodyParser(&item); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    userID := c.FormValue("user_id")
    cart, exists := carts[userID]
    if !exists {
        cart = models.Cart{UserID: userID}
    }

    cart.Items = append(cart.Items, item)
    carts[userID] = cart
    return c.JSON(cart)
}

func GetCart(c *fiber.Ctx) error {
    userID := c.Params("userID")
    cart, exists := carts[userID]
    if !exists {
        return c.JSON(models.Cart{UserID: userID})
    }
    return c.JSON(cart)
}