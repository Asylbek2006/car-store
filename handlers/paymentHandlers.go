package handlers

import (
	"net/http"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentRepo *repositories.PaymentRepository
}

func NewPaymentHandler(paymentRepo *repositories.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{paymentRepo: paymentRepo}
}

func (handler *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := models.Payment{
		UserID:      userID,
		Amount:      req.Amount,
		PaymentType: req.PaymentType,
	}

	if err := handler.paymentRepo.Create(c.Request.Context(), payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "payment recorded"})
}

func (handler *PaymentHandler) GetMyPayments(c *gin.Context) {
	userID := c.GetInt("user_id")

	payments, err := handler.paymentRepo.GetMyPayments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (handler *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := handler.paymentRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}
