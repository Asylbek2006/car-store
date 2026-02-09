package models

type CreatePaymentRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	PaymentType string  `json:"payment_type" binding:"required"`
}
