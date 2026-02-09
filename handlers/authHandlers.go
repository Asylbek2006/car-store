package handlers

import (
	"car-management-system/config"
	"car-management-system/models"
	"car-management-system/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepository
}

func NewAuthHandler(authRepo *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

func (handler *AuthHandler) SignUp(c *gin.Context) {
	var request models.RegisterUserRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not bind json body"))
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not hash password"))
		return
	}
	user := models.User{
		Full_name:    request.Full_name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := handler.authRepo.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not create an user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (handler *AuthHandler) GetAll(c *gin.Context) {
	users, err := handler.authRepo.FindAll(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (handler *AuthHandler) SignIn(c *gin.Context) {
	var request models.SignInRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not bind json body"))
		return
	}

	// This now expects a *models.User, not an int
	user, err := handler.authRepo.FindByEmailHash(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("invalid email or password"))
		return
	}

	// Now we can access user.PasswordHash because user is a struct
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("invalid email or password"))
		return
	}

	claims := models.JwtClaims{
		UserID: user.User_id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.User_id),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not sign JWT"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func (handler *AuthHandler) SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully exited",
	})
}

func (handler *AuthHandler) GetUserBalance(c *gin.Context) {
	userID := c.GetInt("user_id")

	balance, err := handler.authRepo.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not get user balance"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
