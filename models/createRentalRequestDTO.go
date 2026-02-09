package models

type CreateRentalRequest struct {
	CarID      int     `json:"car_id" binding:"required"`
	StartDate  string  `json:"start_date" binding:"required"`
	EndDate    string  `json:"end_date" binding:"required"`
	DailyPrice float64 `json:"daily_price" binding:"required"`
}
