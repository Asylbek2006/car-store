package handlers

import (
	"net/http"

	"car-management-system/repositories"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminRepo *repositories.AdminRepository
}

func NewAdminHandler(adminRepo *repositories.AdminRepository) *AdminHandler {
	return &AdminHandler{adminRepo: adminRepo}
}

func (handler *AdminHandler) GetStats(c *gin.Context) {
	stats, err := handler.adminRepo.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
