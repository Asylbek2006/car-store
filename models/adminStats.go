package models

type AdminStatsResponse struct {
	Users        int     `json:"users"`
	Cars         int     `json:"cars"`
	CarsForSale  int     `json:"cars_for_sale"`
	CarsForRent  int     `json:"cars_for_rent"`
	TotalSales   int     `json:"total_sales"`
	TotalRentals int     `json:"total_rentals"`
	TotalRevenue float64 `json:"total_revenue"`
	TopCar       *Car    `json:"top_car,omitempty"`
}
