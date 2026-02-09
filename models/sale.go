package models

import "time"

type Sale struct {
	SaleID    int       `json:"sale_id"`
	CarID     int       `json:"car_id"`
	SellerID  int       `json:"seller_id"`
	BuyerID   int       `json:"buyer_id"`
	SalePrice float64   `json:"sale_price"`
	SaleDate  time.Time `json:"sale_date"`
}
