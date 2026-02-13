package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type RecommendationHandler struct {
	carRepo *repositories.CarRepository
	genAI   *genai.GenerativeModel
}

func NewRecommendationHandler(carRepo *repositories.CarRepository) *RecommendationHandler {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	// ✅ FIX: Use 'gemini-1.5-flash'. It works perfectly with library v0.19.0
	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	return &RecommendationHandler{
		carRepo: carRepo,
		genAI:   model,
	}
}

// GetQuestions returns the full list of 10 questions
func (handler *RecommendationHandler) GetQuestions(c *gin.Context) {
	questions := []models.RecommendationQuestion{
		{
			ID:       1,
			Question: "Who do you usually drive with?",
			Options:  []string{"Just me", "Me and a partner", "Family with kids", "Clients or Colleagues"},
		},
		{
			ID:       2,
			Question: "Where do you drive most of the time?",
			Options:  []string{"City streets (Traffic)", "Highways (Long distance)", "Rural/Off-road", "Mixed commute"},
		},
		{
			ID:       3,
			Question: "What is your preferred fuel type?",
			Options:  []string{"Petrol (Gasoline)", "Diesel", "Hybrid", "Electric (EV)"},
		},
		{
			ID:       4,
			Question: "What is your budget range?",
			Options:  []string{"Budget-Friendly", "Mid-Range", "High-End", "Luxury / No Limit"},
		},
		{
			ID:       5,
			Question: "How would you describe your driving style?",
			Options:  []string{"Calm & Relaxed", "Standard & Balanced", "Sporty & Fast", "Adventurous"},
		},
		{
			ID:       6,
			Question: "How important is technology and features?",
			Options:  []string{"Not important", "Nice to have", "Important", "Must have latest tech"},
		},
		{
			ID:       7,
			Question: "What car size do you prefer?",
			Options:  []string{"Compact / Small", "Sedan / Medium", "SUV / Large", "Truck / Utility"},
		},
		{
			ID:       8,
			Question: "How much cargo space (trunk) do you need?",
			Options:  []string{"Minimal (Gym bag)", "Average (Groceries)", "Large (Strollers/Luggage)", "Huge (Furniture/Equipment)"},
		},
		{
			ID:       9,
			Question: "What is your top priority?",
			Options:  []string{"Reliability", "Fuel Efficiency", "Performance", "Safety"},
		},
		{
			ID:       10,
			Question: "Do you care about brand prestige?",
			Options:  []string{"No, just a car", "A little bit", "Yes, I like known brands", "Yes, status is key"},
		},
	}

	c.JSON(200, questions)
}

func (handler *RecommendationHandler) GetRecommendation(c *gin.Context) {
	var req models.RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 1. Build the Prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Recommend the single best car model (Make, Model, Year) for a user with these preferences:\n")

	for qID, answer := range req.Answers {
		promptBuilder.WriteString(fmt.Sprintf("- Question %s: %s\n", qID, answer))
	}

	promptBuilder.WriteString("\nReturn a JSON object with keys: 'car_name' (e.g. 2024 Toyota Camry), 'reasoning' (short sales pitch under 30 words), and 'search_query' (best 3 words to search images for this car).")

	// 2. Call Gemini
	resp, err := handler.genAI.GenerateContent(c.Request.Context(), genai.Text(promptBuilder.String()))
	if err != nil {
		c.JSON(500, gin.H{"error": "AI generation failed: " + err.Error()})
		return
	}

	// 3. Return the AI's answer
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]

		if txt, ok := part.(genai.Text); ok {
			c.Data(200, "application/json", []byte(txt))
			return
		}
	}

	c.JSON(500, gin.H{"error": "No response from AI"})
}
