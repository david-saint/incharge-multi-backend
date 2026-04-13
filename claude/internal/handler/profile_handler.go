package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/middleware"
	"github.com/incharge/server/internal/model"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/validator"
	"gorm.io/gorm"
)

// ProfileHandler handles profile endpoints.
type ProfileHandler struct {
	profileRepo *repository.ProfileRepo
	cfg         *config.Config
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(profileRepo *repository.ProfileRepo, cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{profileRepo: profileRepo, cfg: cfg}
}

// SaveProfile handles POST /api/v1/user/profile/ — create or update profile.
func (h *ProfileHandler) SaveProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	var req dto.ProfileSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	// Validate conditional: religion_sect required if religion == CHRISTIANITY.
	if req.Religion != nil && *req.Religion == "CHRISTIANITY" {
		if req.ReligionSect == nil || *req.ReligionSect == "" {
			dto.WriteValidationError(w, map[string][]string{
				"religion_sect": {"The religion sect field is required when religion is CHRISTIANITY."},
			})
			return
		}
	}

	// Build profile model.
	profile := &model.Profile{
		Gender: req.Gender,
	}

	// Age
	if req.Age != nil {
		profile.Age = *req.Age
	}

	// Date of birth
	if req.DOB != "" {
		dob, err := time.Parse("2006-01-02", req.DOB)
		if err != nil {
			dob, err = time.Parse(time.RFC3339, req.DOB)
			if err != nil {
				dob = time.Now()
			}
		}
		profile.DateOfBirth = dob
	} else {
		profile.DateOfBirth = time.Now()
	}

	// Address
	profile.Address = req.Address

	// Marital status
	if req.MaritalStatus != nil {
		profile.MaritalStatus = *req.MaritalStatus
	} else {
		profile.MaritalStatus = "SINGLE"
	}

	// Height / weight
	profile.Height = req.Height
	profile.Weight = req.Weight

	// Education level
	if req.EducationLevel != nil {
		profile.EducationLevelID = req.EducationLevel
	} else {
		defaultEdu := uint(14)
		profile.EducationLevelID = &defaultEdu
	}

	// Occupation
	if req.Occupation != "" {
		profile.Occupation = sql.NullString{String: req.Occupation, Valid: true}
	}

	// Children
	if req.Children != nil {
		profile.NumberOfChildren = req.Children
	} else {
		zero := uint(0)
		profile.NumberOfChildren = &zero
	}

	// Reason
	if req.Reason != nil {
		profile.ContraceptionReasonID = req.Reason
	} else {
		defaultReason := uint(3)
		profile.ContraceptionReasonID = &defaultReason
	}

	// Sexually active
	if req.SexuallyActive != nil {
		profile.SexuallyActive = *req.SexuallyActive
	}

	// Pregnancy status
	if req.PregnancyStatus != nil {
		profile.PregnancyStatus = *req.PregnancyStatus
	}

	// Religion
	if req.Religion != nil {
		profile.Religion = sql.NullString{String: *req.Religion, Valid: true}
	} else {
		profile.Religion = sql.NullString{String: "OTHER", Valid: true}
	}

	// Religion sect
	if req.ReligionSect != nil {
		profile.ReligionSect = sql.NullString{String: *req.ReligionSect, Valid: true}
	}

	if err := h.profileRepo.Upsert(userID, profile); err != nil {
		dto.WriteServerError(w, "Failed to save profile", err, h.cfg.App.IsProduction)
		return
	}

	// Reload with relationships.
	saved, err := h.profileRepo.FindByUserID(userID, "reason", "educationLevel")
	if err != nil {
		dto.WriteServerError(w, "Failed to load profile", err, h.cfg.App.IsProduction)
		return
	}

	resource := toProfileResource(saved)
	dto.WriteJSON(w, http.StatusCreated, dto.SuccessResponse{
		Status:  "success",
		Message: "Profile saved successfully",
		Data:    resource,
	})
}

// GetProfile handles GET /api/v1/user/profile/.
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	// Parse ?with= parameter.
	var withs []string
	if withParam := r.URL.Query().Get("with"); withParam != "" {
		withs = strings.Split(withParam, ",")
	}

	profile, err := h.profileRepo.FindByUserID(userID, withs...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Profile not found"})
			return
		}
		dto.WriteServerError(w, "Failed to fetch profile", err, h.cfg.App.IsProduction)
		return
	}

	resource := toProfileResource(profile)
	dto.WriteJSON(w, http.StatusOK, resource)
}

// StoreAlgorithmPlan handles POST /api/v1/user/profile/algorithm.
func (h *ProfileHandler) StoreAlgorithmPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	var req dto.AlgorithmPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	profile, err := h.profileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Profile not found. Please create a profile first."})
			return
		}
		dto.WriteServerError(w, "Failed to fetch profile", err, h.cfg.App.IsProduction)
		return
	}

	// Merge contraceptive_plan into existing meta.
	meta := profile.GetMeta()
	meta["contraceptive_plan"] = req.Plan
	if err := profile.SetMeta(meta); err != nil {
		dto.WriteServerError(w, "Failed to update meta", err, h.cfg.App.IsProduction)
		return
	}

	if err := h.profileRepo.UpdateMeta(profile.ID, profile.Meta); err != nil {
		dto.WriteServerError(w, "Failed to save plan", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "Contraceptive plan saved successfully", req.Plan)
}
