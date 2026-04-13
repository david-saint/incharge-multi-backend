package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/model"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/validator"
)

// AlgorithmHandler handles algorithm endpoints.
type AlgorithmHandler struct {
	algoRepo *repository.AlgorithmRepo
	cfg      *config.Config
}

// NewAlgorithmHandler creates a new AlgorithmHandler.
func NewAlgorithmHandler(algoRepo *repository.AlgorithmRepo, cfg *config.Config) *AlgorithmHandler {
	return &AlgorithmHandler{algoRepo: algoRepo, cfg: cfg}
}

// ListAlgorithms handles GET /algo — publicly accessible.
func (h *AlgorithmHandler) ListAlgorithms(w http.ResponseWriter, r *http.Request) {
	algos, err := h.algoRepo.ListAll()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch algorithms", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteJSON(w, http.StatusOK, algos)
}

// CreateAlgorithm handles POST /algo.
func (h *AlgorithmHandler) CreateAlgorithm(w http.ResponseWriter, r *http.Request) {
	var req dto.AlgorithmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	algo := &model.Algorithm{
		Text:       req.Text,
		Delay:      req.Delay,
		Active:     req.Active,
		OnPositive: req.OnPositive,
		OnNegative: req.OnNegative,
		NextMove:   req.NextMove,
		Series:     req.Series,
	}

	if algo.Active == "" {
		algo.Active = "N"
	}

	setNullString(&algo.ActionType, req.ActionType)
	setNullString(&algo.Positive, req.Positive)
	setNullString(&algo.Negative, req.Negative)
	setNullString(&algo.TempPlan, req.TempPlan)
	setNullString(&algo.TempPlanDirP, req.TempPlanDirP)
	setNullString(&algo.TempPlanDirN, req.TempPlanDirN)
	setNullString(&algo.ConditionalFactor, req.ConditionalFactor)
	setNullString(&algo.ConditionalOperator, req.ConditionalOperator)
	setNullString(&algo.ConditionalValue, req.ConditionalValue)
	setNullString(&algo.StateValue, req.StateValue)
	setNullString(&algo.Label, req.Label)
	setNullString(&algo.ProgestogenPossible, req.ProgestogenPossible)
	setNullString(&algo.ProgestogenPossibleDir, req.ProgestogenPossibleDir)

	if err := h.algoRepo.Create(algo); err != nil {
		dto.WriteServerError(w, "Failed to create algorithm", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, algo)
}

// UpdateAlgorithm handles PUT /algo/{id}.
func (h *AlgorithmHandler) UpdateAlgorithm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid algorithm ID"})
		return
	}

	algo, err := h.algoRepo.FindByID(uint(id))
	if err != nil {
		dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Algorithm not found"})
		return
	}

	var req dto.AlgorithmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	algo.Text = req.Text
	algo.Delay = req.Delay
	algo.OnPositive = req.OnPositive
	algo.OnNegative = req.OnNegative
	algo.NextMove = req.NextMove
	algo.Series = req.Series

	if req.Active != "" {
		algo.Active = req.Active
	}

	setNullString(&algo.ActionType, req.ActionType)
	setNullString(&algo.Positive, req.Positive)
	setNullString(&algo.Negative, req.Negative)
	setNullString(&algo.TempPlan, req.TempPlan)
	setNullString(&algo.TempPlanDirP, req.TempPlanDirP)
	setNullString(&algo.TempPlanDirN, req.TempPlanDirN)
	setNullString(&algo.ConditionalFactor, req.ConditionalFactor)
	setNullString(&algo.ConditionalOperator, req.ConditionalOperator)
	setNullString(&algo.ConditionalValue, req.ConditionalValue)
	setNullString(&algo.StateValue, req.StateValue)
	setNullString(&algo.Label, req.Label)
	setNullString(&algo.ProgestogenPossible, req.ProgestogenPossible)
	setNullString(&algo.ProgestogenPossibleDir, req.ProgestogenPossibleDir)

	if err := h.algoRepo.Update(algo); err != nil {
		dto.WriteServerError(w, "Failed to update algorithm", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteJSON(w, http.StatusOK, algo)
}

func setNullString(ns *model.NullString, val *string) {
	if val != nil {
		ns.String = *val
		ns.Valid = true
	} else {
		ns.Valid = false
	}
}
