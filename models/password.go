package models

import "time"

// Database Model
type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
	Used      bool
}

// JSON Request: Step 1 (Request Link)
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// JSON Request: Step 2 (Reset Password)
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
