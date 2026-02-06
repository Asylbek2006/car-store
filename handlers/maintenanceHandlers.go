package handlers

import (
	"net/http"
	"strconv"
	"time"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	carRepo         *repositories.CarRepository
	maintenanceRepo *repositories.MaintenanceRepository
}

func NewMaintenanceHandler(carRepo *repositories.CarRepository, maintenanceRepo *repositories.MaintenanceRepository) *MaintenanceHandler {
	return &MaintenanceHandler{carRepo: carRepo, maintenanceRepo: maintenanceRepo}
}

func (handler *MaintenanceHandler) CreateMaintenance(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, _ := strconv.Atoi(c.Param("id"))

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil || car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access"})
		return
	}

	var req models.MaintenanceRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.CarID = carID
	_ = handler.maintenanceRepo.Create(c.Request.Context(), req)

	c.JSON(http.StatusCreated, gin.H{"message": "maintenance added"})
}

func (handler *MaintenanceHandler) GetMaintenance(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")

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

	if role != "admin" && car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access"})
		return
	}

	list, err := handler.maintenanceRepo.GetByCar(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *MaintenanceHandler) PatchMaintenance(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid maintenance id"})
		return
	}

	var req models.UpdateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	m, err := handler.maintenanceRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "maintenance record not found"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), m.CarID)
	if err != nil {
		c.JSON(404, gin.H{"error": "car not found"})
		return
	}
	if role != "admin" && car.OwnerId != userID {
		c.JSON(403, gin.H{"error": "no access"})
		return
	}

	if req.ServiceType != nil {
		m.ServiceType = *req.ServiceType
	}
	if req.ServiceDate != nil {
		date, err := time.Parse("2006-01-02", *req.ServiceDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid service_date format (YYYY-MM-DD)"})
			return
		}
		m.ServiceDate = date
	}
	if req.Mileage != nil {
		m.Mileage = *req.Mileage
	}
	if req.Cost != nil {
		m.Cost = *req.Cost
	}
	if req.Notes != nil {
		m.Notes = *req.Notes
	}

	if err := handler.maintenanceRepo.Update(c.Request.Context(), m); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, m)
}
