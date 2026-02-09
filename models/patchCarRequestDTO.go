package models

type PatchCarRequest struct {
	Brand  *string  `json:"brand"`
	Model  *string  `json:"model"`
	Year   *int     `json:"year"`
	Price  *float64 `json:"price"`
	Status *string  `json:"status"`
}
