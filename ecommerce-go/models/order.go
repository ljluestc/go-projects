package models

type Order struct {
    ID         int     `json:"id"`
    UserID     string  `json:"user_id"`
    Total      float64 `json:"total"`
    Status     string  `json:"status"`
    Currency   string  `json:"currency"`
    Shipping   float64 `json:"shipping"`
    Discount   float64 `json:"discount"`
}