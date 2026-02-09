package models

import "time"

type Review struct {
	ReviewID  int       `json:"review_id"`
	CarID     int       `json:"car_id"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
