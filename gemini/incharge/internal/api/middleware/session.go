package middleware

import (
	"net/http"

	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IsLogged Middleware: Verifies session is active
func IsLogged() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		adminID := session.Get("admin_id")
		if adminID == nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Not logged in")
			c.Abort()
			return
		}

		var admin models.Admin
		if err := database.DB.First(&admin, adminID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Not logged in")
			c.Abort()
			return
		}

		c.Set("admin", admin)
		c.Next()
	}
}

// IsAdmin Middleware: Verifies web session is active AND admin's verified = 'Y'
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		adminID := session.Get("admin_id")
		if adminID == nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Not logged in")
			c.Abort()
			return
		}

		var admin models.Admin
		if err := database.DB.First(&admin, adminID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Not logged in")
			c.Abort()
			return
		}

		if admin.Verified != "Y" {
			utils.ErrorResponse(c, http.StatusForbidden, "Admin not verified")
			c.Abort()
			return
		}

		c.Set("admin", admin)
		c.Next()
	}
}
