package middleware

import (
	"net/http"
	"strings"

	"incharge/internal/config"
	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-gonic/gin"
)

// Auth Middleware: Verifies the JWT Token
func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission Denied"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission Denied"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := utils.ValidateJWT(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission Denied"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission Denied"})
			c.Abort()
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Next()
	}
}

// User Middleware: Verifies the authenticated entity is a User and loads them
func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access this resource"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access this resource"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// Global throttle mock (to be fleshed out with redis or rate-limiter package)
func Throttle(limit int, window string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a production system, implement a token bucket algorithm
		// using Redis or an in-memory cache to enforce 120/min
		// For the scope of this implementation, we register the middleware
		c.Next()
	}
}
