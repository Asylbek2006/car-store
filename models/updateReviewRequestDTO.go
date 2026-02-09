package models

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}
