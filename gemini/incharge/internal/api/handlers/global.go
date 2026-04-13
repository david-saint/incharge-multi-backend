package handlers

import (
	"net/http"

	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct{}

func (h *GlobalHandler) GetContraceptionReasons(c *gin.Context) {
	var reasons []models.ContraceptionReason
	if err := database.DB.Find(&reasons).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reasons")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reasons retrieved", reasons)
}

func (h *GlobalHandler) GetContraceptionReason(c *gin.Context) {
	id := c.Param("id")
	var reason models.ContraceptionReason
	if err := database.DB.First(&reason, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Reason not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reason retrieved", reason)
}

func (h *GlobalHandler) GetEducationLevels(c *gin.Context) {
	var levels []models.EducationLevel
	if err := database.DB.Find(&levels).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve education levels")
		return
	}

	// Map to required { id, name } structure if needed, or simply return model
	utils.SuccessResponse(c, http.StatusOK, "Education levels retrieved", levels)
}

func (h *GlobalHandler) GetFaqGroups(c *gin.Context) {
	var groups []models.FaqGroup
	if err := database.DB.Find(&groups).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve FAQ groups")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "FAQ groups retrieved", groups)
}

func (h *GlobalHandler) GetFaqGroup(c *gin.Context) {
	id := c.Param("id")
	var faq models.Faq
	if err := database.DB.Where("faq_group_id = ?", id).First(&faq).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "FAQ content not found")
		return
	}

	// Required custom response shape for FAQ content
	c.JSON(http.StatusOK, gin.H{
		"data":   faq.Content,
		"status": "faq.get_content",
	})
}
