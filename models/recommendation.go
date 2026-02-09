package models

type RecommendationQuestion struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type RecommendationRequest struct {
	Answers map[string]string `json:"answers"`
}
