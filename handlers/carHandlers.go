package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	carRepo *repositories.CarRepository
}

func NewCarHandlers(carRepo *repositories.CarRepository) *CarHandler {
	return &CarHandler{carRepo: carRepo}
}

func (handler *CarHandler) CreateCar(c *gin.Context) {
	var request models.CreateCarRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid json body"))
		return
	}

	userIdValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewApiError("unauthorized"))
		return
	}

	ownerId := userIdValue.(int)

	car := models.Car{
		OwnerId: ownerId,
		Brand:   request.Brand,
		Model:   request.Model,
		Year:    request.Year,
		Price:   request.Price,
		Status:  request.Status,
	}

	carId, err := handler.carRepo.Create(c.Request.Context(), car)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not create car"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"car_id": carId})
}
