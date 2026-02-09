package models

type CreateSaleRequest struct {
	CarID     int     `json:"car_id" binding:"required"`
	BuyerID   int     `json:"buyer_id" binding:"required"`
	SalePrice float64 `json:"sale_price" binding:"required"`
}
