package dto

// --- Request DTOs ---

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty,phone_ng_us"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// PasswordEmailRequest is the payload for password reset email.
type PasswordEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetRequest is the payload for password reset.
type PasswordResetRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Token                string `json:"token" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

// ProfileSaveRequest is the payload for creating/updating a profile.
type ProfileSaveRequest struct {
	Gender          string   `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	Age             *uint    `json:"age" validate:"omitempty"`
	DOB             string   `json:"dob" validate:"omitempty"`
	Address         string   `json:"address"`
	MaritalStatus   *string  `json:"marital_status" validate:"omitempty,oneof=SINGLE RELATIONSHIP"`
	Height          *uint    `json:"height" validate:"omitempty"`
	Weight          *float64 `json:"weight" validate:"omitempty"`
	EducationLevel  *uint    `json:"education_level" validate:"omitempty"`
	Occupation      string   `json:"occupation"`
	Children        *uint    `json:"children" validate:"omitempty"`
	Reason          *uint    `json:"reason" validate:"omitempty"`
	SexuallyActive  *bool    `json:"sexually_active"`
	PregnancyStatus *bool    `json:"pregnancy_status"`
	Religion        *string  `json:"religion" validate:"omitempty,oneof=CHRISTIANITY ISLAM OTHER"`
	ReligionSect    *string  `json:"religion_sect" validate:"omitempty,oneof=CATHOLIC PENTECOSTAL OTHER"`
}

// AlgorithmPlanRequest is the payload for storing a contraceptive plan.
type AlgorithmPlanRequest struct {
	Plan string `json:"plan" validate:"required"`
}

// ClinicRequest is the payload for creating/updating a clinic.
type ClinicRequest struct {
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	AddedByID uint    `json:"added_by_id" validate:"required"`
}

// AdminCreateRequest is the payload for creating a new admin.
type AdminCreateRequest struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"omitempty"`
	Password  string `json:"password" validate:"required,min=6"`
	Verified  string `json:"verified" validate:"omitempty,oneof=Y N"`
	UserType  string `json:"userType" validate:"required,oneof=Super Sub"`
}

// AdminUpdateRequest is the payload for updating an admin.
type AdminUpdateRequest struct {
	Firstname   *string `json:"firstname"`
	Lastname    *string `json:"lastname"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Verified    *string `json:"verified" validate:"omitempty,oneof=Y N"`
	UserType    *string `json:"userType" validate:"omitempty,oneof=Super Sub"`
	AccessToken *string `json:"accessToken"`
}

// AdminLoginRequest is the payload for admin login.
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AlgorithmCreateRequest is the payload for creating/updating an algorithm step.
type AlgorithmCreateRequest struct {
	Text                   string  `json:"text" validate:"required"`
	ActionType             *string `json:"actionType" validate:"omitempty,oneof=bool input date"`
	Positive               *string `json:"positive"`
	Negative               *string `json:"negative"`
	OnPositive             *uint   `json:"onPositive"`
	OnNegative             *uint   `json:"onNegative"`
	NextMove               *uint   `json:"nextMove"`
	TempPlan               *string `json:"tempPlan"`
	TempPlanDirP           *string `json:"tempPlanDirP"`
	TempPlanDirN           *string `json:"tempPlanDirN"`
	ConditionalFactor      *string `json:"conditionalFactor"`
	ConditionalOperator    *string `json:"conditionalOperator"`
	ConditionalValue       *string `json:"conditionalValue"`
	StateValue             *string `json:"stateValue"`
	Label                  *string `json:"label"`
	ProgestogenPossible    *string `json:"progestogenPossible" validate:"omitempty,oneof=true false"`
	ProgestogenPossibleDir *string `json:"progestogenPossibleDir" validate:"omitempty,oneof=positive negative"`
	Delay                  int     `json:"delay"`
	Series                 *int    `json:"series"`
	Active                 string  `json:"active" validate:"omitempty,oneof=Y N"`
}
