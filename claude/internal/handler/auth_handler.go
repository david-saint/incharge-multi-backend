package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/middleware"
	"github.com/incharge/server/internal/model"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/service"
	"github.com/incharge/server/internal/validator"
	"gorm.io/gorm"
)

// AuthHandler handles user authentication endpoints.
type AuthHandler struct {
	userRepo     *repository.UserRepo
	authService  *service.AuthService
	emailService *service.EmailService
	cfg          *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(
	userRepo *repository.UserRepo,
	authService *service.AuthService,
	emailService *service.EmailService,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		authService:  authService,
		emailService: emailService,
		cfg:          cfg,
	}
}

// Register handles POST /api/v1/user/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	// Check email uniqueness.
	if h.userRepo.EmailExists(req.Email) {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"The email has already been taken."},
		})
		return
	}

	// Check phone uniqueness.
	if req.Phone != "" && h.userRepo.PhoneExists(req.Phone) {
		dto.WriteValidationError(w, map[string][]string{
			"phone": {"The phone has already been taken."},
		})
		return
	}

	hash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		dto.WriteServerError(w, "Failed to hash password", err, h.cfg.App.IsProduction)
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    strings.ToLower(req.Email),
		Password: hash,
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := h.userRepo.Create(user); err != nil {
		dto.WriteServerError(w, "Failed to create user", err, h.cfg.App.IsProduction)
		return
	}

	// Send verification email (synchronous per spec).
	go func() {
		verifyURL := middleware.GenerateSignedURL(
			h.cfg.App.URL+"/api/v1/user/email/verify",
			user.ID,
			h.cfg.JWT.Secret,
			4320, // 72 hours
		)
		if err := h.emailService.SendVerificationEmail(user.Email, user.Name, verifyURL); err != nil {
			slog.Error("failed to send verification email", "error", err, "user_id", user.ID)
		}
	}()

	resource := toUserResource(user)
	dto.WriteJSON(w, http.StatusCreated, dto.SuccessResponse{
		Status:  true,
		Message: "User registered successfully",
		Data:    resource,
	})
}

// Login handles POST /api/v1/user/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	user, err := h.userRepo.FindByEmail(strings.ToLower(req.Email))
	if err != nil {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"These credentials do not match our records."},
		})
		return
	}

	if !h.authService.CheckPassword(user.Password, req.Password) {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"These credentials do not match our records."},
		})
		return
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		dto.WriteServerError(w, "Failed to generate token", err, h.cfg.App.IsProduction)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	dto.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
}

// Logout handles POST /api/v1/user/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a production system, you'd add the token to a blocklist.
	// For now, we simply return success as the client should discard the token.
	dto.WriteSuccess(w, "Successfully logged out.", nil)
}

// Refresh handles GET /api/v1/user/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	userID, err := h.authService.ValidateToken(parts[1])
	if err != nil {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		dto.WriteServerError(w, "Failed to generate token", err, h.cfg.App.IsProduction)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	dto.WriteJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
}

// PasswordEmail handles POST /api/v1/user/password/email.
func (h *AuthHandler) PasswordEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	user, err := h.userRepo.FindByEmail(strings.ToLower(req.Email))
	if err != nil {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"We can't find a user with that email address."},
		})
		return
	}

	token, err := service.GenerateRandomToken(32)
	if err != nil {
		dto.WriteServerError(w, "Failed to generate token", err, h.cfg.App.IsProduction)
		return
	}

	if err := h.userRepo.CreatePasswordReset(user.Email, token); err != nil {
		dto.WriteServerError(w, "Failed to create reset token", err, h.cfg.App.IsProduction)
		return
	}

	resetURL := h.cfg.App.UserDomain + "/reset-password/" + token
	go func() {
		if err := h.emailService.SendPasswordResetEmail(user.Email, resetURL); err != nil {
			slog.Error("failed to send password reset email", "error", err)
		}
	}()

	dto.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "We have emailed your password reset link!",
	})
}

// PasswordReset handles POST /api/v1/user/password/reset.
func (h *AuthHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	pr, err := h.userRepo.FindPasswordReset(strings.ToLower(req.Email), req.Token)
	if err != nil {
		dto.WriteServerError(w, "Failed to verify token", err, h.cfg.App.IsProduction)
		return
	}
	if pr == nil {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"This password reset token is invalid."},
		})
		return
	}

	user, err := h.userRepo.FindByEmail(strings.ToLower(req.Email))
	if err != nil {
		dto.WriteValidationError(w, map[string][]string{
			"email": {"We can't find a user with that email address."},
		})
		return
	}

	hash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		dto.WriteServerError(w, "Failed to hash password", err, h.cfg.App.IsProduction)
		return
	}

	if err := h.userRepo.UpdatePassword(user.ID, hash); err != nil {
		dto.WriteServerError(w, "Failed to update password", err, h.cfg.App.IsProduction)
		return
	}

	h.userRepo.DeletePasswordResets(user.Email)

	dto.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Your password has been reset!",
	})
}

// VerifyEmail handles GET /api/v1/user/email/verify/{id}.
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}

	if err := h.userRepo.VerifyEmail(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to verify email", err, h.cfg.App.IsProduction)
		return
	}

	http.Redirect(w, r, h.cfg.App.UserDomain+"/email-verified", http.StatusFound)
}

// ResendVerification handles GET /api/v1/user/email/resend.
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	if user.IsEmailVerified() {
		dto.WriteSuccess(w, "Email already verified.", nil)
		return
	}

	verifyURL := middleware.GenerateSignedURL(
		h.cfg.App.URL+"/api/v1/user/email/verify",
		user.ID,
		h.cfg.JWT.Secret,
		4320,
	)

	if err := h.emailService.SendVerificationEmail(user.Email, user.Name, verifyURL); err != nil {
		slog.Error("failed to resend verification email", "error", err)
		dto.WriteServerError(w, "Failed to send verification email", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteSuccess(w, "Verification email sent.", nil)
}

// EmailSuccess handles GET /api/v1/user/email/success.
func (h *AuthHandler) EmailSuccess(w http.ResponseWriter, r *http.Request) {
	dto.WriteSuccess(w, "Email verified successfully.", nil)
}

// GetUser handles GET /api/v1/user/ — returns the authenticated user.
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}
		dto.WriteServerError(w, "Failed to get user", err, h.cfg.App.IsProduction)
		return
	}

	resource := toUserResource(user)
	dto.WriteJSON(w, http.StatusOK, resource)
}

// --- Helper functions ---

func toUserResource(u *model.User) dto.UserResource {
	var phone string
	if u.Phone != nil {
		phone = *u.Phone
	}
	r := dto.UserResource{
		ID:             u.ID,
		Name:           u.Name,
		Email:          u.Email,
		Phone:          phone,
		EmailVerified:  u.IsEmailVerified(),
		PhoneConfirmed: u.IsPhoneConfirmed(),
		CreatedAt:      u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      u.UpdatedAt.Format(time.RFC3339),
	}
	if u.Profile != nil {
		pr := toProfileResource(u.Profile)
		r.Profile = &pr
	}
	return r
}

func toProfileResource(p *model.Profile) dto.ProfileResource {
	r := dto.ProfileResource{
		ID:              p.ID,
		Age:             p.Age,
		Gender:          p.Gender,
		DateOfBirth:     p.DateOfBirth.Format(time.RFC3339),
		Address:         p.Address,
		Latitude:        p.Latitude,
		Longitude:       p.Longitude,
		MaritalStatus:   p.MaritalStatus,
		Height:          p.Height,
		Weight:          p.Weight,
		Children:        p.NumberOfChildren,
		SexuallyActive:  p.SexuallyActive,
		PregnancyStatus: p.PregnancyStatus,
	}

	if p.Occupation.Valid {
		r.Occupation = p.Occupation.String
	}
	if p.Religion.Valid {
		r.Religion = &p.Religion.String
	}
	if p.ReligionSect.Valid {
		r.ReligionSect = &p.ReligionSect.String
	}

	if p.ContraceptionReason != nil {
		r.Reason = &dto.NamedResource{
			ID:   p.ContraceptionReason.ID,
			Name: p.ContraceptionReason.Value,
		}
	}
	if p.EducationLevel != nil {
		r.EducationLevel = &dto.NamedResource{
			ID:   p.EducationLevel.ID,
			Name: p.EducationLevel.Name,
		}
	}
	if p.User != nil {
		ur := toUserResource(p.User)
		r.User = &ur
	}

	return r
}
