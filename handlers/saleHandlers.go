package handlers

import (
	"net/http"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	carRepo  *repositories.CarRepository
	saleRepo *repositories.SaleRepository
}

func NewSaleHandler(carRepo *repositories.CarRepository, saleRepo *repositories.SaleRepository) *SaleHandler {
	return &SaleHandler{carRepo: carRepo, saleRepo: saleRepo}
}

func (handler *SaleHandler) CreateSale(c *gin.Context) {
	sellerID := c.GetInt("user_id")

	var req models.CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), req.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	sale := models.Sale{
		CarID:     req.CarID,
		SellerID:  sellerID,
		BuyerID:   req.BuyerID,
		SalePrice: req.SalePrice,
	}

	if err := handler.saleRepo.Create(c.Request.Context(), sale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = handler.saleRepo.MarkCarSold(c.Request.Context(), req.CarID)

	c.JSON(http.StatusCreated, gin.H{"message": "car sold"})
}

func (handler *SaleHandler) GetAllSales(c *gin.Context) {
	sales, err := handler.saleRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sales)
}
