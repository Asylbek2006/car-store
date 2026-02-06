package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetHandler struct {
	authRepo  *repositories.AuthRepository
	resetRepo *repositories.PasswordResetRepository
}

func NewPasswordResetHandler(authRepo *repositories.AuthRepository, resetRepo *repositories.PasswordResetRepository) *PasswordResetHandler {
	return &PasswordResetHandler{
		authRepo:  authRepo,
		resetRepo: resetRepo,
	}
}

func (handler *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.authRepo.FindEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset link sent"})
		return
	}

	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	expiresAt := time.Now().Add(15 * time.Minute)

	handler.resetRepo.Create(c.Request.Context(), user.User_id, token, expiresAt)

	c.JSON(http.StatusOK, gin.H{
		"reset_token": token,
		"expires_at":  expiresAt,
	})
}

func (handler *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := handler.resetRepo.GetByToken(c.Request.Context(), req.Token)
	if err != nil || pr.Used || time.Now().After(pr.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	handler.authRepo.UpdatePassword(c.Request.Context(), pr.UserID, string(hash))
	handler.resetRepo.MarkUsed(c.Request.Context(), req.Token)

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}
