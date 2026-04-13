package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/repository"
	"gorm.io/gorm"
)

// GlobalHandler handles public/global endpoints.
type GlobalHandler struct {
	refRepo *repository.ReferenceRepo
	cfg     *config.Config
}

// NewGlobalHandler creates a new GlobalHandler.
func NewGlobalHandler(refRepo *repository.ReferenceRepo, cfg *config.Config) *GlobalHandler {
	return &GlobalHandler{refRepo: refRepo, cfg: cfg}
}

// HealthCheck handles GET /api/v1/ — returns "Hello, World!"
func (h *GlobalHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

// ListContraceptionReasons handles GET /api/v1/global/contraception-reasons.
func (h *GlobalHandler) ListContraceptionReasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := h.refRepo.ListContraceptionReasons()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch contraception reasons", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.ReasonResource, len(reasons))
	for i, reason := range reasons {
		resources[i] = dto.ReasonResource{
			ID:    reason.ID,
			Value: reason.Value,
		}
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// GetContraceptionReason handles GET /api/v1/global/contraception-reasons/{id}.
func (h *GlobalHandler) GetContraceptionReason(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
		return
	}

	reason, err := h.refRepo.FindContraceptionReason(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Not found"})
			return
		}
		dto.WriteServerError(w, "Failed to fetch reason", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.ReasonResource{
		ID:    reason.ID,
		Value: reason.Value,
	})
}

// ListEducationLevels handles GET /api/v1/global/education-levels.
func (h *GlobalHandler) ListEducationLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.refRepo.ListEducationLevels()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch education levels", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.NamedResource, len(levels))
	for i, level := range levels {
		resources[i] = dto.NamedResource{
			ID:   level.ID,
			Name: level.Name,
		}
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// ListFaqGroups handles GET /api/v1/global/faq-groups.
func (h *GlobalHandler) ListFaqGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.refRepo.ListFaqGroups()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch FAQ groups", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteJSON(w, http.StatusOK, groups)
}

// GetFaqGroup handles GET /api/v1/global/faq-groups/{id}.
func (h *GlobalHandler) GetFaqGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
		return
	}

	group, err := h.refRepo.FindFaqGroupWithContent(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Not found"})
			return
		}
		dto.WriteServerError(w, "Failed to fetch FAQ group", err, h.cfg.App.IsProduction)
		return
	}

	var content interface{}
	if group.Faq != nil && group.Faq.Content != nil {
		content = group.Faq.Content
	}

	dto.WriteJSON(w, http.StatusOK, dto.FaqContentResponse{
		Data:   content,
		Status: "faq.get_content",
	})
}
