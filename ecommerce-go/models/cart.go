package models

type CartItem struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}

type Cart struct {
    UserID string      `json:"user_id"`
    Items  []CartItem `json:"items"`
}