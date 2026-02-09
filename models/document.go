package models

import "time"

type Document struct {
	DocumentID   int       `json:"document_id"`
	CarID        int       `json:"car_id"`
	DocumentType string    `json:"document_type"`
	ExpiryDate   time.Time `json:"expiry_date"`
	FileURL      string    `json:"file_url"`
}
