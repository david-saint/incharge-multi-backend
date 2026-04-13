package routes

import (
	"incharge/internal/api/handlers"
	"incharge/internal/api/middleware"
	"incharge/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) {
	adminHandler := &handlers.AdminHandler{}

	// Admin Web Panel Routes (Top Level, Not API v1)
	r.POST("/login", adminHandler.Login)
	r.GET("/logout", adminHandler.Logout)
	r.GET("/algo", adminHandler.ListAlgorithms)
	r.POST("/admin", adminHandler.CreateAdmin) // Create Admin

	adminProtected := r.Group("/")
	adminProtected.Use(middleware.IsAdmin())
	{
		adminProtected.GET("/getAdminDet", adminHandler.GetAdminDet)
		adminProtected.GET("/allAdmins", func(c *gin.Context) {})
		adminProtected.PUT("/updateAdmin/:admin_id", func(c *gin.Context) {})
		adminProtected.POST("/algo", func(c *gin.Context) {})
		adminProtected.PUT("/algo/:id", func(c *gin.Context) {})
	}

	v1 := r.Group("/api/v1")

	// Global / Public
	globalHandler := &handlers.GlobalHandler{}
	global := v1.Group("/global")
	{
		global.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": true, "message": "Hello, World!", "data": nil})
		})
		global.GET("/contraception-reasons", globalHandler.GetContraceptionReasons)
		global.GET("/contraception-reasons/:id", globalHandler.GetContraceptionReason)
		global.GET("/education-levels", globalHandler.GetEducationLevels)
		global.GET("/faq-groups", globalHandler.GetFaqGroups)
		global.GET("/faq-groups/:id", globalHandler.GetFaqGroup)
	}

	authHandler := &handlers.AuthHandler{Cfg: cfg}

	// User Endpoints
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/register", authHandler.Register)
		userGroup.POST("/login", authHandler.Login)
		userGroup.POST("/password/email", func(c *gin.Context) {}) // Placeholder
		userGroup.POST("/password/reset", func(c *gin.Context) {}) // Placeholder

		// Authenticated Routes
		authRequired := userGroup.Group("")
		authRequired.Use(middleware.Auth(cfg))
		{
			authRequired.POST("/logout", authHandler.Logout)
			authRequired.GET("/refresh", authHandler.Refresh)
			authRequired.GET("/email/resend", func(c *gin.Context) {}) // Placeholder

			// User Context required
			userRequired := authRequired.Group("")
			userRequired.Use(middleware.User())
			{
				userRequired.GET("/", authHandler.GetUser)

				profileHandler := &handlers.ProfileHandler{}
				profileGroup := userRequired.Group("/profile")
				{
					profileGroup.GET("/", profileHandler.GetProfile)
					profileGroup.POST("/", profileHandler.CreateOrUpdateProfile)
					profileGroup.POST("/algorithm", profileHandler.StoreAlgorithmPlan)
				}

				clinicHandler := &handlers.ClinicHandler{}
				clinicGroup := userRequired.Group("/clinics")
				{
					clinicGroup.GET("/", clinicHandler.ListClinics)
				}
			}
		}
	}
}
