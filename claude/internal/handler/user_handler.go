package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/repository"
)

// UserHandler handles user management endpoints (admin-side).
type UserHandler struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userRepo *repository.UserRepo, cfg *config.Config) *UserHandler {
	return &UserHandler{userRepo: userRepo, cfg: cfg}
}

// ListUsers handles GET /api/v1/user/users/ — paginated user list.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	users, total, err := h.userRepo.ListWithProfile(page, 50)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch users", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.UserResource, len(users))
	for i, u := range users {
		resources[i] = toUserResource(&u)
	}

	path := h.cfg.App.URL + r.URL.Path
	resp := dto.NewPaginatedResponse(resources, total, page, 50, path)
	dto.WriteJSON(w, http.StatusOK, resp)
}

// ListDeletedUsers handles GET /api/v1/user/users/deletedUser.
func (h *UserHandler) ListDeletedUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.ListDeleted()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch deleted users", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.UserResource, len(users))
	for i, u := range users {
		resources[i] = toUserResource(&u)
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// RestoreUser handles PUT /api/v1/user/users/update/{user_id}.
func (h *UserHandler) RestoreUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}

	if err := h.userRepo.Restore(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to restore user", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "User restored successfully", nil)
}

// DeleteUser handles DELETE /api/v1/user/users/deleteUser/{user_id}.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}

	if err := h.userRepo.SoftDelete(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to delete user", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "User deleted successfully", nil)
}
