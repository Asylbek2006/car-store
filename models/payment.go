package models

import "time"

type Payment struct {
	PaymentID   int       `json:"payment_id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	PaymentType string    `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at"`
}
