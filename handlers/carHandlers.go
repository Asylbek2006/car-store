package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"net/http"
	"strconv"

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

func (handler *CarHandler) UpdateCar(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, _ := strconv.Atoi(c.Param("id"))

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	var req models.UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	car.Brand = req.Brand
	car.Model = req.Model
	car.Year = req.Year
	car.Price = req.Price
	car.Status = req.Status

	if err := handler.carRepo.Update(c.Request.Context(), car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "car updated"})
}

func (handler *CarHandler) DeleteCar(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, _ := strconv.Atoi(c.Param("id"))

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	if err := handler.carRepo.Delete(c.Request.Context(), carID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "car deleted"})
}

func (handler *CarHandler) GetAllCars(c *gin.Context) {
	cars, err := handler.carRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (handler *CarHandler) PatchCar(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	carID, _ := strconv.Atoi(c.Param("id"))

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil {
		c.JSON(404, gin.H{"error": "car not found"})
		return
	}

	if role != "admin" && car.OwnerId != userID {
		c.JSON(403, gin.H{"error": "no access"})
		return
	}

	var req models.PatchCarRequest
	c.ShouldBindJSON(&req)

	if req.Brand != nil {
		car.Brand = *req.Brand
	}
	if req.Model != nil {
		car.Model = *req.Model
	}
	if req.Year != nil {
		car.Year = *req.Year
	}
	if req.Price != nil {
		car.Price = *req.Price
	}
	if req.Status != nil {
		car.Status = *req.Status
	}

	handler.carRepo.UpdateFields(c.Request.Context(), car)
	c.JSON(200, car)
}

func (handler *CarHandler) BuyCar(c *gin.Context) {
	// 1. Читаем JSON из тела запроса
	var req models.BuyCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON или отсутствует car_id"})
		return
	}

	// Проверяем, что ID пришел
	if req.CarID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "car_id обязателен"})
		return
	}

	userID := c.GetInt("user_id")

	// 2. Вызываем репозиторий (передаем ID из JSON)
	err := handler.carRepo.BuyCarTransaction(c.Request.Context(), userID, req.CarID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Покупка успешно завершена!",
	})
}
