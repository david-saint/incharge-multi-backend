package dto

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestNewPaginatedResponse(t *testing.T) {
	t.Run("basic pagination", func(t *testing.T) {
		data := []string{"a", "b", "c"}
		resp := NewPaginatedResponse(data, 100, 2, 10, "/api/items")

		if resp.CurrentPage != 2 {
			t.Fatalf("CurrentPage: expected 2, got %d", resp.CurrentPage)
		}
		if resp.PerPage != 10 {
			t.Fatalf("PerPage: expected 10, got %d", resp.PerPage)
		}
		if resp.Total != 100 {
			t.Fatalf("Total: expected 100, got %d", resp.Total)
		}
		if resp.LastPage != 10 {
			t.Fatalf("LastPage: expected 10, got %d", resp.LastPage)
		}
		if resp.From != 11 {
			t.Fatalf("From: expected 11, got %d", resp.From)
		}
		if resp.To != 20 {
			t.Fatalf("To: expected 20, got %d", resp.To)
		}
		if resp.FirstPageURL != "/api/items?page=1" {
			t.Fatalf("FirstPageURL: expected /api/items?page=1, got %s", resp.FirstPageURL)
		}
		if resp.LastPageURL != "/api/items?page=10" {
			t.Fatalf("LastPageURL: expected /api/items?page=10, got %s", resp.LastPageURL)
		}
		if resp.NextPageURL == nil || *resp.NextPageURL != "/api/items?page=3" {
			t.Fatalf("NextPageURL: expected /api/items?page=3, got %v", resp.NextPageURL)
		}
		if resp.PrevPageURL == nil || *resp.PrevPageURL != "/api/items?page=1" {
			t.Fatalf("PrevPageURL: expected /api/items?page=1, got %v", resp.PrevPageURL)
		}
	})

	t.Run("first page has no prev URL", func(t *testing.T) {
		resp := NewPaginatedResponse(nil, 50, 1, 10, "/items")
		if resp.PrevPageURL != nil {
			t.Fatalf("PrevPageURL should be nil on first page, got %v", resp.PrevPageURL)
		}
		if resp.NextPageURL == nil {
			t.Fatal("NextPageURL should not be nil when more pages exist")
		}
	})

	t.Run("last page has no next URL", func(t *testing.T) {
		resp := NewPaginatedResponse(nil, 50, 5, 10, "/items")
		if resp.NextPageURL != nil {
			t.Fatalf("NextPageURL should be nil on last page, got %v", resp.NextPageURL)
		}
		if resp.PrevPageURL == nil {
			t.Fatal("PrevPageURL should not be nil on last page")
		}
	})

	t.Run("empty result set", func(t *testing.T) {
		resp := NewPaginatedResponse([]string{}, 0, 1, 10, "/items")
		if resp.From != 0 {
			t.Fatalf("From: expected 0 for empty result, got %d", resp.From)
		}
		if resp.To != 0 {
			t.Fatalf("To: expected 0 for empty result, got %d", resp.To)
		}
		if resp.LastPage != 1 {
			t.Fatalf("LastPage: expected 1 (minimum), got %d", resp.LastPage)
		}
	})

	t.Run("To is capped at total", func(t *testing.T) {
		resp := NewPaginatedResponse(nil, 25, 3, 10, "/items")
		if resp.To != 25 {
			t.Fatalf("To: expected 25 (capped at total), got %d", resp.To)
		}
	})
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, 201, map[string]string{"key": "value"})

	if rec.Code != 201 {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected content-type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["key"] != "value" {
		t.Fatalf("expected key=value, got %v", body)
	}
}

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccess(rec, "OK", map[string]int{"count": 5})

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body SuccessResponse
	json.NewDecoder(rec.Body).Decode(&body)
	if body.Message != "OK" {
		t.Fatalf("expected message=OK, got %s", body.Message)
	}
}

func TestWriteValidationError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteValidationError(rec, map[string][]string{
		"email": {"The email field is required."},
	})

	if rec.Code != 422 {
		t.Fatalf("expected 422, got %d", rec.Code)
	}

	var body ValidationErrorResponse
	json.NewDecoder(rec.Body).Decode(&body)
	if len(body.Errors["email"]) != 1 {
		t.Fatalf("expected 1 email error, got %d", len(body.Errors["email"]))
	}
}

func TestWriteAuthError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteAuthError(rec, "Permission Denied")

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestWriteServerError(t *testing.T) {
	t.Run("production hides error details", func(t *testing.T) {
		rec := httptest.NewRecorder()
		WriteServerError(rec, "Something went wrong", nil, true)
		if rec.Code != 500 {
			t.Fatalf("expected 500, got %d", rec.Code)
		}
		var body ErrorResponse
		json.NewDecoder(rec.Body).Decode(&body)
		if body.Error != "" {
			t.Fatalf("production should hide error details, got %q", body.Error)
		}
	})

	t.Run("non-production shows error details", func(t *testing.T) {
		rec := httptest.NewRecorder()
		WriteServerError(rec, "Fail", &testErr{}, false)
		var body ErrorResponse
		json.NewDecoder(rec.Body).Decode(&body)
		if body.Error != "test error" {
			t.Fatalf("expected error detail, got %q", body.Error)
		}
	})
}

type testErr struct{}

func (e *testErr) Error() string { return "test error" }
