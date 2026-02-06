package handlers

import (
	"net/http"
	"strconv"
	"time"

	"car-management-system/models"
	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	carRepo      *repositories.CarRepository
	documentRepo *repositories.DocumentRepository
}

func NewDocumentHandler(carRepo *repositories.CarRepository, documentRepo *repositories.DocumentRepository) *DocumentHandler {
	return &DocumentHandler{
		carRepo:      carRepo,
		documentRepo: documentRepo,
	}
}

func (handler *DocumentHandler) CreateDocument(c *gin.Context) {
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

	var req models.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiry_date format (YYYY-MM-DD)"})
		return
	}

	doc := models.Document{
		CarID:        carID,
		DocumentType: req.DocumentType,
		ExpiryDate:   expiryDate,
		FileURL:      req.FileURL,
	}

	if err := handler.documentRepo.Create(c.Request.Context(), doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "document added"})
}

func (handler *DocumentHandler) GetByCar(c *gin.Context) {
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

	list, err := handler.documentRepo.GetByCar(c.Request.Context(), carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *DocumentHandler) DeleteDocument(c *gin.Context) {
	userID := c.GetInt("user_id")

	documentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	doc, err := handler.documentRepo.GetByID(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	car, err := handler.carRepo.GetByID(c.Request.Context(), doc.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	if car.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your car"})
		return
	}

	if err := handler.documentRepo.Delete(c.Request.Context(), documentID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}

func (handler *DocumentHandler) PatchDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.UpdateDocumentRequest
	c.ShouldBindJSON(&req)

	document, _ := handler.documentRepo.GetByID(c.Request.Context(), id)

	if req.DocumentType != nil {
		document.DocumentType = *req.DocumentType
	}
	if req.ExpiryDate != nil {
		document.ExpiryDate, _ = time.Parse("2006-01-02", *req.ExpiryDate)
	}
	if req.FileURL != nil {
		document.FileURL = *req.FileURL
	}

	handler.documentRepo.Update(c.Request.Context(), document)
	c.JSON(200, document)
}
