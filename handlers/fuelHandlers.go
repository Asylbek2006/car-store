package handlers

import (
	"net/http"
	"strconv"
	"time"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type FuelHandler struct {
	carRepo  *repositories.CarRepository
	fuelRepo *repositories.FuelRepository
}

func NewFuelHandler(carRepo *repositories.CarRepository, fuelRepo *repositories.FuelRepository) *FuelHandler {
	return &FuelHandler{
		carRepo:  carRepo,
		fuelRepo: fuelRepo,
	}
}

func (handler *FuelHandler) CreateFuel(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	var req models.CreateFuelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fillDate, err := time.Parse("2006-01-02", req.FillDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fill_date format (YYYY-MM-DD)"})
		return
	}

	fuel := models.FuelLog{
		CarID:    carID,
		FillDate: fillDate,
		Liters:   req.Liters,
		Price:    req.Price,
		Mileage:  req.Mileage,
	}

	if err := handler.fuelRepo.Create(c.Request.Context(), fuel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "fuel log added"})
}

func (handler *FuelHandler) GetFuel(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	fuelLogs, err := handler.fuelRepo.GetByCar(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fuelLogs)
}

func (handler *FuelHandler) PatchFuel(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fuel id"})
		return
	}

	var req models.UpdateFuelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fuel, err := handler.fuelRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fuel log not found"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), fuel.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}
	if role != "admin" && car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access"})
		return
	}

	if req.FillDate != nil {
		date, err := time.Parse("2006-01-02", *req.FillDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fill_date format (YYYY-MM-DD)"})
			return
		}
		fuel.FillDate = date
	}
	if req.Liters != nil {
		fuel.Liters = *req.Liters
	}
	if req.Price != nil {
		fuel.Price = *req.Price
	}
	if req.Mileage != nil {
		fuel.Mileage = *req.Mileage
	}

	if err := handler.fuelRepo.Update(c.Request.Context(), fuel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fuel)
}
