package models

type UserResponse struct {
	User_id      int    `json:"user_id"`
	Full_name    string `json:"full_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
