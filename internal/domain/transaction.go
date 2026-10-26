package domain

import "time"

type Transaction struct {
    ID         int       `json:"id"`
    Amount     float64   `json:"amount"`
    CategoryID int       `json:"category_id"`
    CreatedAt  time.Time `json:"created_at"`
}