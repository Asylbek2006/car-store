package models

type UpdateExpenseRequest struct {
	ExpenseType *string  `json:"expense_type"`
	Amount      *float64 `json:"amount"`
	ExpenseDate *string  `json:"expense_date"`
	Description *string  `json:"description"`
}
