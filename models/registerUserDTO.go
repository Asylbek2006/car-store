package models

type RegisterUserRequest struct {
	Full_name string `json:"full_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
