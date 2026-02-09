package handlers

import (
	"car-management-system/models"
	"car-management-system/repositories"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

// 1. Forgot Password
func (handler *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ИСПРАВЛЕНИЕ 1: Вызываем FindByEmail (который мы добавили в репозиторий)
	// Он возвращает int (userID)
	userID, err := handler.authRepo.FindByEmail(c.Request.Context(), req.Email)

	if err != nil {
		// Если email не найден, не выдаем ошибку пользователю (безопасность)
		c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset link sent"})
		return
	}

	// Генерируем токен
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	// ИСПРАВЛЕНИЕ 2: Убираем expiresAt из аргументов.
	// Репозиторий (Create) сам добавит +1 час к текущему времени.
	err = handler.resetRepo.Create(c.Request.Context(), userID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Логируем ссылку в консоль (СИМУЛЯЦИЯ EMAIL)
	fmt.Println("=====================================")
	fmt.Println("EMAIL TO:", req.Email)
	fmt.Printf("LINK: http://localhost:3000/reset-password.html?token=%s\n", token)
	fmt.Println("=====================================")

	c.JSON(http.StatusOK, gin.H{"message": "reset link sent (check server console)"})
}

// 2. Reset Password
func (handler *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Получаем данные токена
	pr, err := handler.resetRepo.GetByToken(c.Request.Context(), req.Token)

	// Проверяем валидность, использован ли он, и не истек ли срок
	if err != nil || pr.Used || time.Now().After(pr.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// 2. Хешируем новый пароль
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// 3. Обновляем пароль пользователя
	err = handler.authRepo.UpdatePassword(c.Request.Context(), pr.UserID, string(hash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update password"})
		return
	}

	// 4. Помечаем токен как использованный
	handler.resetRepo.MarkUsed(c.Request.Context(), req.Token)

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}
