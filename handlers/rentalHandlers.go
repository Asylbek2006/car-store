package handlers

import (
	"net/http"
	"strconv"
	"time"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	carRepo    *repositories.CarRepository
	rentalRepo *repositories.RentalRepository
}

func NewRentalHandler(carRepo *repositories.CarRepository, rentalRepo *repositories.RentalRepository) *RentalHandler {
	return &RentalHandler{carRepo: carRepo, rentalRepo: rentalRepo}
}

func (handler *RentalHandler) CreateRental(c *gin.Context) {
	renterID := c.GetInt("user_id")

	var req models.CreateRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (YYYY-MM-DD)"})
		return
	}

	if !endDate.After(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), req.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId == renterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot rent your own car"})
		return
	}

	days := int(endDate.Sub(startDate).Hours() / 24)
	if days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rental duration"})
		return
	}

	total := float64(days) * req.DailyPrice

	rental := models.Rental{
		CarID:      req.CarID,
		OwnerID:    car.OwnerId,
		RenterID:   renterID,
		StartDate:  startDate,
		EndDate:    endDate,
		DailyPrice: req.DailyPrice,
		TotalPrice: total,
		Status:     "active",
	}

	if err := handler.rentalRepo.Create(c.Request.Context(), rental); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = handler.rentalRepo.MarkCarRented(c.Request.Context(), req.CarID)

	c.JSON(http.StatusCreated, gin.H{"message": "car rented"})
}

func (handler *RentalHandler) CompleteRental(c *gin.Context) {
	userID := c.GetInt("user_id")
	rentalID, _ := strconv.Atoi(c.Param("id"))

	rental, err := handler.rentalRepo.GetByID(c.Request.Context(), rentalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rental not found"})
		return
	}

	if rental.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your rental"})
		return
	}

	_ = handler.rentalRepo.UpdateStatus(c.Request.Context(), rentalID, "completed")
	_ = handler.carRepo.MarkAvailable(c.Request.Context(), rental.CarID)

	c.JSON(http.StatusOK, gin.H{"message": "rental completed"})
}

func (handler *RentalHandler) CancelRental(c *gin.Context) {
	userID := c.GetInt("user_id")
	rentalID, _ := strconv.Atoi(c.Param("id"))

	rental, err := handler.rentalRepo.GetByID(c.Request.Context(), rentalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rental not found"})
		return
	}

	if rental.OwnerID != userID && rental.RenterID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	_ = handler.rentalRepo.UpdateStatus(c.Request.Context(), rentalID, "cancelled")
	_ = handler.carRepo.MarkAvailable(c.Request.Context(), rental.CarID)

	c.JSON(http.StatusOK, gin.H{"message": "rental cancelled"})
}

func (handler *RentalHandler) GetAllRentals(c *gin.Context) {
	rentals, err := handler.rentalRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rentals)
}
