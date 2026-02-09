package models

type UpdateMaintenanceRequest struct {
	ServiceType *string  `json:"service_type"`
	ServiceDate *string  `json:"service_date"`
	Mileage     *int     `json:"mileage"`
	Cost        *float64 `json:"cost"`
	Notes       *string  `json:"notes"`
}
