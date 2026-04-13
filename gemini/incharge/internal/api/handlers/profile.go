package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{}

type ProfileRequest struct {
	Gender          string     `json:"gender" binding:"required,oneof=MALE FEMALE OTHER"`
	Age             *uint      `json:"age"`
	Dob             *time.Time `json:"dob"`
	Address         *string    `json:"address"`
	MaritalStatus   *string    `json:"marital_status" binding:"omitempty,oneof=SINGLE RELATIONSHIP"`
	Height          *uint      `json:"height"`
	Weight          *float64   `json:"weight"`
	EducationLevel  *uint      `json:"education_level"`
	Occupation      *string    `json:"occupation"`
	Children        *uint      `json:"children"`
	Reason          *uint      `json:"reason"`
	SexuallyActive  *bool      `json:"sexually_active"`
	PregnancyStatus *bool      `json:"pregnancy_status"`
	Religion        *string    `json:"religion" binding:"omitempty,oneof=CHRISTIANITY ISLAM OTHER"`
	ReligionSect    *string    `json:"religion_sect" binding:"omitempty,oneof=CATHOLIC PENTECOSTAL OTHER"`
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Permission Denied")
		return
	}

	var profile models.Profile
	query := database.DB.Where("user_id = ?", userID)

	with := c.Query("with")
	if with != "" {
		if utils.Contains(with, "user") {
			query = query.Preload("User")
		}
		if utils.Contains(with, "reason") {
			query = query.Preload("Reason")
		}
		if utils.Contains(with, "educationLevel") {
			query = query.Preload("EducationLevel")
		}
	}

	if err := query.First(&profile).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Profile not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved", profile)
}

func (h *ProfileHandler) CreateOrUpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Permission Denied")
		return
	}

	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	// Apply Defaults
	age := uint(0)
	if req.Age != nil {
		age = *req.Age
	}

	dob := time.Now()
	if req.Dob != nil {
		dob = *req.Dob
	}

	address := ""
	if req.Address != nil {
		address = *req.Address
	}

	maritalStatus := "SINGLE"
	if req.MaritalStatus != nil {
		maritalStatus = *req.MaritalStatus
	}

	educationLevelID := uint(14)
	if req.EducationLevel != nil {
		educationLevelID = *req.EducationLevel
	}

	children := uint(0)
	if req.Children != nil {
		children = *req.Children
	}

	reasonID := uint(3)
	if req.Reason != nil {
		reasonID = *req.Reason
	}

	sexuallyActive := false
	if req.SexuallyActive != nil {
		sexuallyActive = *req.SexuallyActive
	}

	pregnancyStatus := false
	if req.PregnancyStatus != nil {
		pregnancyStatus = *req.PregnancyStatus
	}

	religion := "OTHER"
	if req.Religion != nil {
		religion = *req.Religion
	}

	if religion == "CHRISTIANITY" && req.ReligionSect == nil {
		utils.ValidationErrorResponse(c, map[string][]string{"religion_sect": {"The religion sect field is required when religion is CHRISTIANITY."}})
		return
	}

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		// Create
		profile = models.Profile{
			UserID:                userID.(uint),
			Gender:                req.Gender,
			Age:                   age,
			DateOfBirth:           dob,
			Address:               address,
			MaritalStatus:         maritalStatus,
			Height:                req.Height,
			Weight:                req.Weight,
			EducationLevelID:      &educationLevelID,
			Occupation:            req.Occupation,
			NumberOfChildren:      &children,
			ContraceptionReasonID: reasonID,
			SexuallyActive:        sexuallyActive,
			PregnancyStatus:       pregnancyStatus,
			Religion:              religion,
			ReligionSect:          req.ReligionSect,
		}
		if err := database.DB.Create(&profile).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create profile")
			return
		}
		utils.SuccessResponse(c, 201, "Profile created successfully", profile)
		return
	}

	// Update
	profile.Gender = req.Gender
	profile.Age = age
	profile.DateOfBirth = dob
	profile.Address = address
	profile.MaritalStatus = maritalStatus
	profile.Height = req.Height
	profile.Weight = req.Weight
	profile.EducationLevelID = &educationLevelID
	profile.Occupation = req.Occupation
	profile.NumberOfChildren = &children
	profile.ContraceptionReasonID = reasonID
	profile.SexuallyActive = sexuallyActive
	profile.PregnancyStatus = pregnancyStatus
	profile.Religion = religion
	profile.ReligionSect = req.ReligionSect

	if err := database.DB.Save(&profile).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.SuccessResponse(c, 200, "Profile updated successfully", profile)
}

type AlgorithmRequest struct {
	Plan string `json:"plan" binding:"required"`
}

func (h *ProfileHandler) StoreAlgorithmPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Permission Denied")
		return
	}

	var req AlgorithmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, map[string][]string{"fields": {err.Error()}})
		return
	}

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Profile not found")
		return
	}

	meta := make(map[string]interface{})
	if profile.Meta != nil {
		json.Unmarshal([]byte(*profile.Meta), &meta)
	}

	meta["contraceptive_plan"] = req.Plan
	metaJSON, _ := json.Marshal(meta)
	metaStr := string(metaJSON)
	profile.Meta = &metaStr

	if err := database.DB.Save(&profile).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to store plan")
		return
	}

	utils.SuccessResponse(c, 200, "Plan stored successfully", req.Plan)
}
