package models

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Type        string  `json:"type"` // "physical" or "digital"
    Description string  `json:"description"`
    Image       string  `json:"image"`
    Stock       int     `json:"stock"`
}