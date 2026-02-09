package models

type CreateExpenseRequest struct {
	ExpenseType string  `json:"expense_type" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	ExpenseDate string  `json:"expense_date" binding:"required"`
	Description string  `json:"description"`
}
