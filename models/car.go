package models

import "time"

type Car struct {
	CarId     int       `json:"car_id"`
	OwnerId   int       `json:"owner_id"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Year      int       `json:"year"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
