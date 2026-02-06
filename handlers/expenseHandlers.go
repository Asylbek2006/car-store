package handlers

import (
	"net/http"
	"strconv"
	"time"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	carRepo     *repositories.CarRepository
	expenseRepo *repositories.ExpenseRepository
}

func NewExpenseHandler(carRepo *repositories.CarRepository, expenseRepo *repositories.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{carRepo: carRepo, expenseRepo: expenseRepo}
}

func (handler *ExpenseHandler) CreateExpense(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil || car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	var req models.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense_date format (YYYY-MM-DD)"})
		return
	}

	expense := models.Expense{
		CarID:       carID,
		ExpenseType: req.ExpenseType,
		Amount:      req.Amount,
		ExpenseDate: expenseDate,
		Description: req.Description,
	}

	if err := handler.expenseRepo.Create(c.Request.Context(), expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "expense added"})
}

func (handler *ExpenseHandler) GetExpense(c *gin.Context) {
	userID := c.GetInt("user_id")
	carID, _ := strconv.Atoi(c.Param("id"))

	car, err := handler.carRepo.GetByID(c.Request.Context(), carID)
	if err != nil || car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	list, _ := handler.expenseRepo.GetByCar(c.Request.Context(), carID)
	c.JSON(http.StatusOK, list)
}

func (handler *ExpenseHandler) PatchExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.UpdateExpenseRequest
	c.ShouldBindJSON(&req)

	expense, _ := handler.expenseRepo.GetByID(c.Request.Context(), id)

	if req.ExpenseType != nil {
		expense.ExpenseType = *req.ExpenseType
	}
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.ExpenseDate != nil {
		expense.ExpenseDate, _ = time.Parse("2006-01-02", *req.ExpenseDate)
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}

	handler.expenseRepo.Update(c.Request.Context(), expense)
	c.JSON(200, expense)
}
