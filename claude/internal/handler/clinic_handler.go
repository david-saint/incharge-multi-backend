package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/model"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/validator"
)

// ClinicHandler handles clinic endpoints.
type ClinicHandler struct {
	clinicRepo *repository.ClinicRepo
	cfg        *config.Config
}

// NewClinicHandler creates a new ClinicHandler.
func NewClinicHandler(clinicRepo *repository.ClinicRepo, cfg *config.Config) *ClinicHandler {
	return &ClinicHandler{clinicRepo: clinicRepo, cfg: cfg}
}

// ListClinics handles GET /api/v1/user/clinics/ — full query capabilities.
func (h *ClinicHandler) ListClinics(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := repository.ClinicQueryParams{
		Search:      q.Get("search"),
		Mode:        q.Get("mode"),
		Sort:        q.Get("sort"),
		With:        q.Get("with"),
		WithTrashed: q.Get("withTrashed") == "true" || q.Get("withTrashed") == "1",
		OnlyTrashed: q.Get("onlyTrashed") == "true" || q.Get("onlyTrashed") == "1",
		WithCount:   q.Get("withCount"),
	}

	if params.Mode == "" {
		params.Mode = "km"
	}

	// Parse lat/lng.
	if latStr := q.Get("latitude"); latStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err == nil {
			params.Latitude = &lat
		}
	}
	if lngStr := q.Get("longitude"); lngStr != "" {
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err == nil {
			params.Longitude = &lng
		}
	}

	// Parse radius.
	if radiusStr := q.Get("radius"); radiusStr != "" {
		radius, err := strconv.ParseFloat(radiusStr, 64)
		if err == nil {
			params.Radius = radius
		}
	}
	if params.Radius <= 0 {
		params.Radius = 10
	}

	// Parse pagination.
	params.Page, _ = strconv.Atoi(q.Get("page"))
	params.PerPage, _ = strconv.Atoi(q.Get("per_page"))
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	clinics, total, err := h.clinicRepo.List(params)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch clinics", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.ClinicResource, len(clinics))
	for i, c := range clinics {
		resources[i] = toClinicResource(&c)
	}

	// If pagination requested, return paginated response.
	if params.Page > 0 {
		path := h.cfg.App.URL + r.URL.Path
		resp := dto.NewPaginatedResponse(resources, total, params.Page, params.PerPage, path)
		dto.WriteJSON(w, http.StatusOK, resp)
		return
	}

	// Otherwise return all.
	dto.WriteJSON(w, http.StatusOK, resources)
}

// GetClinicsPaginated handles GET /api/v1/user/clinics/getClinics — simple paginated list.
func (h *ClinicHandler) GetClinicsPaginated(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	clinics, total, err := h.clinicRepo.ListPaginated(page, 50)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch clinics", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.ClinicResource, len(clinics))
	for i, c := range clinics {
		resources[i] = toClinicResource(&c)
	}

	path := h.cfg.App.URL + r.URL.Path
	resp := dto.NewPaginatedResponse(resources, total, page, 50, path)
	dto.WriteJSON(w, http.StatusOK, resp)
}

// GetDeletedClinics handles GET /api/v1/user/clinics/deletedClinics.
func (h *ClinicHandler) GetDeletedClinics(w http.ResponseWriter, r *http.Request) {
	clinics, err := h.clinicRepo.ListDeleted()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch deleted clinics", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.ClinicResource, len(clinics))
	for i, c := range clinics {
		resources[i] = toClinicResource(&c)
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// AddClinic handles POST /api/v1/user/clinics/addClinic.
func (h *ClinicHandler) AddClinic(w http.ResponseWriter, r *http.Request) {
	var req dto.ClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	clinic := &model.Clinic{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  &req.Latitude,
		Longitude: &req.Longitude,
		AddedByID: req.AddedByID,
	}

	if err := h.clinicRepo.Create(clinic); err != nil {
		dto.WriteServerError(w, "Failed to create clinic", err, h.cfg.App.IsProduction)
		return
	}

	resource := toClinicResource(clinic)
	dto.WriteJSON(w, http.StatusCreated, dto.SuccessResponse{
		Status:  true,
		Message: "Clinic created successfully",
		Data:    resource,
	})
}

// UpdateClinic handles PUT /api/v1/user/clinics/update/{clinic_id}.
func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "clinic_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}

	clinic, err := h.clinicRepo.FindByID(uint(id))
	if err != nil {
		dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Clinic not found"})
		return
	}

	var req dto.ClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	clinic.Name = req.Name
	clinic.Address = req.Address
	clinic.Latitude = &req.Latitude
	clinic.Longitude = &req.Longitude

	if err := h.clinicRepo.Update(clinic); err != nil {
		dto.WriteServerError(w, "Failed to update clinic", err, h.cfg.App.IsProduction)
		return
	}

	resource := toClinicResource(clinic)
	dto.WriteSuccess(w, "Clinic updated successfully", resource)
}

// RevertDeleteClinic handles PUT /api/v1/user/clinics/revertDelete/{clinic_id}.
func (h *ClinicHandler) RevertDeleteClinic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "clinic_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}

	if err := h.clinicRepo.Restore(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to restore clinic", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "Clinic restored successfully", nil)
}

// DeleteClinic handles DELETE /api/v1/user/clinics/deleteClinic/{clinic_id}.
func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "clinic_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}

	if err := h.clinicRepo.SoftDelete(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to delete clinic", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "Clinic deleted successfully", nil)
}

// --- Helper ---

func toClinicResource(c *model.Clinic) dto.ClinicResource {
	r := dto.ClinicResource{
		ID:        c.ID,
		Name:      c.Name,
		Address:   c.Address,
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}

	if c.Mode != "" {
		r.Mode = c.Mode
		r.Radius = c.Radius
		r.SearchRadius = c.SearchRadius
		r.ActualDistance = c.ActualDistance
		r.Distance = c.Distance
	}

	if c.Locations != nil {
		locs := make([]dto.LocationResource, len(c.Locations))
		for i, l := range c.Locations {
			locs[i] = dto.LocationResource{
				ID:        l.ID,
				Name:      l.Name,
				StateID:   l.StateID,
				Latitude:  l.Latitude,
				Longitude: l.Longitude,
			}
		}
		r.Locations = locs
	}

	r.LocationsCount = c.LocationsCount

	return r
}
