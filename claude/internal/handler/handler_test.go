package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/model"
)

// --- Test helpers ---

func testConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Env:          "testing",
			Port:         "8080",
			URL:          "http://localhost:8080",
			APIDomain:    "localhost",
			UserDomain:   "http://localhost:3000",
			IsProduction: false,
		},
		JWT: config.JWTConfig{
			Secret: "test-jwt-secret-key-12345",
			TTL:    60 * time.Minute,
		},
		Session: config.SessionConfig{
			Secret: "test-session-secret",
		},
	}
}

func TestToUserResource(t *testing.T) {
	t.Run("user without phone", func(t *testing.T) {
		now := time.Now()
		user := &model.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     nil,
			CreatedAt: now,
			UpdatedAt: now,
		}
		res := toUserResource(user)
		if res.ID != 1 {
			t.Fatalf("expected ID=1, got %d", res.ID)
		}
		if res.Name != "John Doe" {
			t.Fatalf("expected name=John Doe, got %s", res.Name)
		}
		if res.Phone != "" {
			t.Fatalf("expected empty phone, got %q", res.Phone)
		}
		if res.EmailVerified {
			t.Fatal("expected email_verified=false")
		}
	})

	t.Run("user with phone", func(t *testing.T) {
		phone := "+2348012345678"
		now := time.Now()
		user := &model.User{
			ID:        2,
			Name:      "Jane",
			Email:     "jane@example.com",
			Phone:     &phone,
			CreatedAt: now,
			UpdatedAt: now,
		}
		res := toUserResource(user)
		if res.Phone != "+2348012345678" {
			t.Fatalf("expected phone=%s, got %s", phone, res.Phone)
		}
	})

	t.Run("user with verified email", func(t *testing.T) {
		now := time.Now()
		verified := now.Add(-1 * time.Hour)
		user := &model.User{
			ID:              3,
			Name:            "Verified User",
			Email:           "verified@example.com",
			EmailVerifiedAt: &verified,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		res := toUserResource(user)
		if !res.EmailVerified {
			t.Fatal("expected email_verified=true")
		}
	})

	t.Run("user with profile", func(t *testing.T) {
		now := time.Now()
		user := &model.User{
			ID:        4,
			Name:      "Profile User",
			Email:     "profile@example.com",
			CreatedAt: now,
			UpdatedAt: now,
			Profile: &model.Profile{
				ID:     1,
				Gender: "FEMALE",
				Age:    25,
			},
		}
		res := toUserResource(user)
		if res.Profile == nil {
			t.Fatal("expected profile to be present")
		}
		if res.Profile.Gender != "FEMALE" {
			t.Fatalf("expected gender=FEMALE, got %s", res.Profile.Gender)
		}
	})
}

func TestToClinicResource(t *testing.T) {
	t.Run("basic clinic", func(t *testing.T) {
		lat := 6.5244
		lng := 3.3792
		now := time.Now()
		clinic := &model.Clinic{
			ID:        1,
			Name:      "Lagos Clinic",
			Address:   "123 Main St, Lagos",
			Latitude:  &lat,
			Longitude: &lng,
			CreatedAt: now,
		}
		res := toClinicResource(clinic)
		if res.Name != "Lagos Clinic" {
			t.Fatalf("expected name=Lagos Clinic, got %s", res.Name)
		}
		if *res.Latitude != 6.5244 {
			t.Fatalf("expected lat=6.5244, got %f", *res.Latitude)
		}
		if res.Mode != "" {
			t.Fatal("expected empty mode for non-distance query")
		}
	})

	t.Run("clinic with distance info", func(t *testing.T) {
		lat := 6.5244
		lng := 3.3792
		now := time.Now()
		clinic := &model.Clinic{
			ID:             2,
			Name:           "Abuja Clinic",
			Address:        "456 Ring Rd, Abuja",
			Latitude:       &lat,
			Longitude:      &lng,
			Mode:           "km",
			Radius:         10,
			SearchRadius:   "10km",
			ActualDistance: 3.45,
			Distance:       "3.45km",
			CreatedAt:      now,
		}
		res := toClinicResource(clinic)
		if res.Mode != "km" {
			t.Fatalf("expected mode=km, got %s", res.Mode)
		}
		if res.Distance != "3.45km" {
			t.Fatalf("expected distance=3.45km, got %s", res.Distance)
		}
	})

	t.Run("clinic with locations count", func(t *testing.T) {
		now := time.Now()
		count := int64(5)
		clinic := &model.Clinic{
			ID:             3,
			Name:           "Test",
			Address:        "Test",
			LocationsCount: &count,
			CreatedAt:      now,
		}
		res := toClinicResource(clinic)
		if res.LocationsCount == nil || *res.LocationsCount != 5 {
			t.Fatalf("expected locations_count=5, got %v", res.LocationsCount)
		}
	})
}

func TestHealthCheckHandler(t *testing.T) {
	// We test the HealthCheck handler directly since it doesn't need a DB.
	cfg := testConfig()
	h := &GlobalHandler{refRepo: nil, cfg: cfg}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/global/", nil)
	rec := httptest.NewRecorder()
	h.HealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "Hello, World!" {
		t.Fatalf("expected 'Hello, World!', got %q", rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %s", ct)
	}
}

func TestRegisterValidation(t *testing.T) {
	// Test that the register handler validates requests properly.
	// We can't fully test without a DB, but we can test JSON parsing and validation.

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register",
			bytes.NewBufferString("not json"))
		rec := httptest.NewRecorder()
		h.Register(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("missing required fields returns 422", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		body, _ := json.Marshal(dto.RegisterRequest{})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register",
			bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		h.Register(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}

		var resp dto.ValidationErrorResponse
		json.NewDecoder(rec.Body).Decode(&resp)

		if _, ok := resp.Errors["name"]; !ok {
			t.Error("expected validation error for name")
		}
		if _, ok := resp.Errors["email"]; !ok {
			t.Error("expected validation error for email")
		}
		if _, ok := resp.Errors["password"]; !ok {
			t.Error("expected validation error for password")
		}
	})
}

func TestLoginValidation(t *testing.T) {
	t.Run("invalid JSON returns 400", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login",
			bytes.NewBufferString("{bad"))
		rec := httptest.NewRecorder()
		h.Login(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("missing fields returns 422", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		body, _ := json.Marshal(dto.LoginRequest{})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login",
			bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		h.Login(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}
	})
}

func TestRefreshHandler(t *testing.T) {
	t.Run("missing Authorization returns 401", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/refresh", nil)
		rec := httptest.NewRecorder()
		h.Refresh(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("malformed Authorization returns 401", func(t *testing.T) {
		cfg := testConfig()
		h := NewAuthHandler(nil, nil, nil, cfg)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/refresh", nil)
		req.Header.Set("Authorization", "JustOneWord")
		rec := httptest.NewRecorder()
		h.Refresh(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}

func TestLogoutHandler(t *testing.T) {
	cfg := testConfig()
	h := NewAuthHandler(nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/logout", nil)
	rec := httptest.NewRecorder()
	h.Logout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp dto.SuccessResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Message != "Successfully logged out." {
		t.Fatalf("expected 'Successfully logged out.', got %q", resp.Message)
	}
}

func TestEmailSuccess(t *testing.T) {
	cfg := testConfig()
	h := NewAuthHandler(nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/email/success", nil)
	rec := httptest.NewRecorder()
	h.EmailSuccess(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetUserRequiresAuth(t *testing.T) {
	// Without user ID in context, GetUser should return 401.
	cfg := testConfig()
	h := NewAuthHandler(nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	rec := httptest.NewRecorder()
	h.GetUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
