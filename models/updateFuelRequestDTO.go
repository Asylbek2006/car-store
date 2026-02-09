package models

type UpdateFuelRequest struct {
	FillDate *string  `json:"fill_date"`
	Liters   *float64 `json:"liters"`
	Price    *float64 `json:"price"`
	Mileage  *int     `json:"mileage"`
}
