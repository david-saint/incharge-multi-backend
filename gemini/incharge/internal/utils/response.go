package utils

import (
	"github.com/gin-gonic/gin"
)

// Standard Response Wrappers
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  false,
		"message": message,
	})
}

func ValidationErrorResponse(c *gin.Context, errors map[string][]string) {
	c.JSON(422, gin.H{
		"errors": errors,
	})
}
