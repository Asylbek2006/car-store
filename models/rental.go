package models

import "time"

type Rental struct {
	RentalID   int       `json:"rental_id"`
	CarID      int       `json:"car_id"`
	OwnerID    int       `json:"owner_id"`
	RenterID   int       `json:"renter_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	DailyPrice float64   `json:"daily_price"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
}
