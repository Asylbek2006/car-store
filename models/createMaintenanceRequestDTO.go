package models

type CreateMaintenanceRequest struct {
	ServiceType string  `json:"service_type" binding:"required"`
	ServiceDate string  `json:"service_date" binding:"required"`
	Mileage     int     `json:"mileage"`
	Cost        float64 `json:"cost"`
	Notes       string  `json:"notes"`
}
