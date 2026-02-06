package handlers

import (
	"net/http"
	"strconv"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewRepo *repositories.ReviewRepository
	carRepo    *repositories.CarRepository
}

func NewReviewHandler(carRepo *repositories.CarRepository, reviewRepo *repositories.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{
		carRepo:    carRepo,
		reviewRepo: reviewRepo,
	}
}

func (handler *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	if _, err := handler.carRepo.GetByID(c.Request.Context(), carID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	exists, err := handler.reviewRepo.ExistsByUser(c.Request.Context(), carID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you already reviewed this car"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := models.Review{
		CarID:   carID,
		UserID:  userID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	if err := handler.reviewRepo.Create(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review added"})
}

func (handler *ReviewHandler) GetByCar(c *gin.Context) {
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	reviews, err := handler.reviewRepo.GetByCar(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (handler *ReviewHandler) PatchReview(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	role := c.GetString("role")

	review, _ := handler.reviewRepo.GetByID(c.Request.Context(), id)

	if role != "admin" && review.UserID != userID {
		c.JSON(403, gin.H{"error": "no access"})
		return
	}

	var req models.UpdateReviewRequest
	c.ShouldBindJSON(&req)

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	handler.reviewRepo.Update(c.Request.Context(), review)
	c.JSON(200, review)
}
