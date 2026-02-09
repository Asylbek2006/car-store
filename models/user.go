package models

import "time"

type User struct {
	User_id      int       `json:"user_id"`
	Full_name    string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Created_at   time.Time `json:"created_at"`
	Role         string    `json:"role"`
	Balance      float64   `json:"balance" gorm:"default:1000000"`
}
