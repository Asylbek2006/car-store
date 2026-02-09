package models

import "time"

type FuelLog struct {
	FuelID   int       `json:"fuel_id"`
	CarID    int       `json:"car_id"`
	FillDate time.Time `json:"fill_date"`
	Liters   float64   `json:"liters"`
	Price    float64   `json:"price"`
	Mileage  int       `json:"mileage"`
}
