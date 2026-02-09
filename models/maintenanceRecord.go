package models

import "time"

type MaintenanceRecord struct {
	MaintenanceID int       `json:"maintenance_id"`
	CarID         int       `json:"car_id"`
	ServiceType   string    `json:"service_type"`
	ServiceDate   time.Time `json:"service_date"`
	Mileage       int       `json:"mileage"`
	Cost          float64   `json:"cost"`
	Notes         string    `json:"notes"`
}
