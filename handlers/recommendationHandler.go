package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"car-management-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	carRepo *repositories.CarRepository
}

func NewRecommendationHandler(carRepo *repositories.CarRepository) *RecommendationHandler {
	return &RecommendationHandler{carRepo: carRepo}
}

func (handler *RecommendationHandler) GetQuestions(c *gin.Context) {
	questions := []models.RecommendationQuestion{
		{
			ID:       1,
			Question: "Who do you usually drive with?",
			Options:  []string{"Friends", "Family", "Colleagues", "Alone"},
		},
		{
			ID:       2,
			Question: "Where do you drive most of the time?",
			Options:  []string{"City", "Highway", "Off-road", "Mixed"},
		},
		{
			ID:       3,
			Question: "What is your preferred fuel type?",
			Options:  []string{"Petrol", "Diesel", "Hybrid", "Electric"},
		},
		{
			ID:       4,
			Question: "What is your budget range?",
			Options:  []string{"Low", "Medium", "High", "Premium"},
		},
		{
			ID:       5,
			Question: "What is your driving style?",
			Options:  []string{"Calm", "Balanced", "Sporty", "Aggressive"},
		},
		{
			ID:       6,
			Question: "How important is comfort to you?",
			Options:  []string{"Low", "Medium", "High", "Very High"},
		},
		{
			ID:       7,
			Question: "What car size do you prefer?",
			Options:  []string{"Small", "Medium", "Large", "Doesn't matter"},
		},
		{
			ID:       8,
			Question: "What type of roads do you mostly use?",
			Options:  []string{"Good roads", "Bad roads", "Mountains", "Mixed"},
		},
		{
			ID:       9,
			Question: "How important is fuel efficiency?",
			Options:  []string{"Low", "Medium", "High", "Very High"},
		},
		{
			ID:       10,
			Question: "What matters most to you?",
			Options:  []string{"Price", "Safety", "Performance", "Technology"},
		},
	}

	c.JSON(http.StatusOK, questions)
}

func (handler *RecommendationHandler) GetRecommendation(c *gin.Context) {
	var req models.RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	carType := services.CalculateCarType(req.Answers)

	cars, err := handler.carRepo.GetRecommendedCars(
		c.Request.Context(),
		carType,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"recommended_type": carType,
		"cars":             cars,
	})
}
