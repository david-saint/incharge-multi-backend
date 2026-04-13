package handlers

import (
	"net/http"

	"incharge/internal/config"
	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Cfg config.Config
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	// Check unique constraints
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ValidationErrorResponse(c, map[string][]string{"email": {"The email has already been taken."}})
		return
	}

	if req.Phone != "" {
		if err := database.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			utils.ValidationErrorResponse(c, map[string][]string{"phone": {"The phone has already been taken."}})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	newUser := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// TODO: Fire registration event -> sends verification email

	utils.SuccessResponse(c, 201, "User registered successfully", newUser)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"email": {"These credentials do not match our records."}})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"email": {"These credentials do not match our records."}})
		return
	}

	token, err := utils.GenerateJWT(user.ID, h.Cfg.JWTSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Simple JWT logout relies on client-side deletion since we don't have a token blacklist
	utils.SuccessResponse(c, 200, "Successfully logged out.", nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Issue a new token for the currently authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission Denied"})
		return
	}

	token, err := utils.GenerateJWT(userID.(uint), h.Cfg.JWTSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token")
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Permission Denied")
		return
	}

	utils.SuccessResponse(c, 200, "User retrieved successfully", user)
}
