package models

import "time"

type Expense struct {
	ExpenseID   int       `json:"expense_id"`
	CarID       int       `json:"car_id"`
	ExpenseType string    `json:"expense_type"`
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	Description string    `json:"description"`
}
