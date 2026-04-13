package dto

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

// --- Response DTOs ---

// SuccessResponse is the standard success response envelope.
type SuccessResponse struct {
	Status  interface{} `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard error response.
type ErrorResponse struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

// ValidationErrorResponse is a 422 validation error response.
type ValidationErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

// TokenResponse is the JWT token response.
type TokenResponse struct {
	Token string `json:"token"`
}

// FaqContentResponse is the response for FAQ content.
type FaqContentResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

// UserResource is the JSON representation of a user.
type UserResource struct {
	ID             uint             `json:"id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	Phone          string           `json:"phone"`
	EmailVerified  bool             `json:"email_verified"`
	PhoneConfirmed bool             `json:"phone_confirmed"`
	Profile        *ProfileResource `json:"profile,omitempty"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
}

// ProfileResource is the JSON representation of a profile.
type ProfileResource struct {
	ID              uint           `json:"id"`
	Age             uint           `json:"age"`
	Gender          string         `json:"gender"`
	DateOfBirth     string         `json:"date_of_birth"`
	Address         string         `json:"address"`
	Latitude        *float64       `json:"latitude"`
	Longitude       *float64       `json:"longitude"`
	MaritalStatus   string         `json:"marital_status"`
	Height          *uint          `json:"height"`
	Weight          *float64       `json:"weight"`
	Occupation      string         `json:"occupation"`
	Children        *uint          `json:"children"`
	SexuallyActive  bool           `json:"sexually_active"`
	PregnancyStatus bool           `json:"pregnancy_status"`
	Religion        *string        `json:"religion"`
	ReligionSect    *string        `json:"religion_sect"`
	Reason          *NamedResource `json:"reason,omitempty"`
	EducationLevel  *NamedResource `json:"education_level,omitempty"`
	User            *UserResource  `json:"user,omitempty"`
}

// NamedResource is a generic {id, name} resource.
type NamedResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ReasonResource wraps a contraception reason with optional profiles.
type ReasonResource struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

// ClinicResource is the JSON representation of a clinic.
type ClinicResource struct {
	ID             uint               `json:"id"`
	Name           string             `json:"name"`
	Address        string             `json:"address"`
	Latitude       *float64           `json:"latitude"`
	Longitude      *float64           `json:"longitude"`
	CreatedAt      string             `json:"created_at"`
	Mode           string             `json:"mode,omitempty"`
	Radius         float64            `json:"radius,omitempty"`
	SearchRadius   string             `json:"search_radius,omitempty"`
	ActualDistance float64            `json:"actual_distance,omitempty"`
	Distance       string             `json:"distance,omitempty"`
	Locations      []LocationResource `json:"locations,omitempty"`
	LocationsCount *int64             `json:"locations_count,omitempty"`
}

// LocationResource is the JSON representation of a location.
type LocationResource struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	StateID   uint     `json:"state_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// PaginatedResponse is the standard pagination envelope.
type PaginatedResponse struct {
	Data         interface{} `json:"data"`
	CurrentPage  int         `json:"current_page"`
	PerPage      int         `json:"per_page"`
	Total        int64       `json:"total"`
	LastPage     int         `json:"last_page"`
	From         int         `json:"from"`
	To           int         `json:"to"`
	FirstPageURL string      `json:"first_page_url"`
	LastPageURL  string      `json:"last_page_url"`
	NextPageURL  *string     `json:"next_page_url"`
	PrevPageURL  *string     `json:"prev_page_url"`
	Path         string      `json:"path"`
}

// NewPaginatedResponse builds a pagination envelope.
func NewPaginatedResponse(data interface{}, total int64, page, perPage int, path string) PaginatedResponse {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	from := 0
	to := 0
	if total > 0 {
		from = (page-1)*perPage + 1
		to = page * perPage
		if int64(to) > total {
			to = int(total)
		}
	}

	firstURL := path + "?page=1"
	lastURL := path + "?page=" + strconv.Itoa(lastPage)

	var nextURL *string
	if page < lastPage {
		u := path + "?page=" + strconv.Itoa(page+1)
		nextURL = &u
	}
	var prevURL *string
	if page > 1 {
		u := path + "?page=" + strconv.Itoa(page-1)
		prevURL = &u
	}

	return PaginatedResponse{
		Data:         data,
		CurrentPage:  page,
		PerPage:      perPage,
		Total:        total,
		LastPage:     lastPage,
		From:         from,
		To:           to,
		FirstPageURL: firstURL,
		LastPageURL:  lastURL,
		NextPageURL:  nextURL,
		PrevPageURL:  prevURL,
		Path:         path,
	}
}

// JSON helper functions

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// WriteSuccess writes a success response.
func WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// WriteCreated writes a 201 created response.
func WriteCreated(w http.ResponseWriter, status string, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// WriteValidationError writes a 422 validation error response.
func WriteValidationError(w http.ResponseWriter, errors map[string][]string) {
	WriteJSON(w, http.StatusUnprocessableEntity, ValidationErrorResponse{Errors: errors})
}

// WriteAuthError writes a 401 auth error response.
func WriteAuthError(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": message})
}

// WriteServerError writes a 500 error response.
func WriteServerError(w http.ResponseWriter, message string, err error, isProduction bool) {
	resp := ErrorResponse{
		Status:  false,
		Message: message,
	}
	if !isProduction && err != nil {
		resp.Error = err.Error()
	}
	WriteJSON(w, http.StatusInternalServerError, resp)
}
