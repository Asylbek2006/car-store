package models

type CreateFuelRequest struct {
	FillDate string  `json:"fill_date" binding:"required"`
	Liters   float64 `json:"liters" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Mileage  int     `json:"mileage"`
}
