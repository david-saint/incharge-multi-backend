package handlers

import (
	"net/http"

	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct{}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if admin.Verified != "Y" {
		utils.ErrorResponse(c, http.StatusForbidden, "Admin not verified")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	session := sessions.Default(c)
	session.Set("admin_id", admin.ID)
	session.Save()

	utils.SuccessResponse(c, 200, "Login successful", admin)
}

func (h *AdminHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin")
}

type AdminCreateRequest struct {
	Firstname string  `json:"firstname" binding:"required"`
	Lastname  string  `json:"lastname" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     *string `json:"phone"`
	Password  string  `json:"password" binding:"required,min=6"`
	Verified  string  `json:"verified" binding:"required,oneof=Y N"`
	UserType  string  `json:"userType" binding:"required,oneof=Super Sub"`
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	var existing models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.ValidationErrorResponse(c, map[string][]string{"email": {"Email already taken"}})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	admin := models.Admin{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		Verified:  req.Verified,
		UserType:  req.UserType,
	}

	// Check if this is the first Super admin
	var count int64
	database.DB.Model(&models.Admin{}).Count(&count)

	if err := database.DB.Create(&admin).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create admin")
		return
	}

	if count == 0 && admin.UserType == "Super" {
		// Auto-login
		session := sessions.Default(c)
		session.Set("admin_id", admin.ID)
		session.Save()
	}

	utils.SuccessResponse(c, 201, "Admin created", admin)
}

func (h *AdminHandler) ListAlgorithms(c *gin.Context) {
	var algos []models.Algorithm
	// Spec: ordered by active ASC, id ASC
	database.DB.Order("active ASC, id ASC").Find(&algos)
	utils.SuccessResponse(c, 200, "Algorithms retrieved", algos)
}

func (h *AdminHandler) GetAdminDet(c *gin.Context) {
	admin, exists := c.Get("admin")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Not logged in")
		return
	}
	utils.SuccessResponse(c, 200, "Admin retrieved", admin)
}
