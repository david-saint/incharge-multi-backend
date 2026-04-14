package models

import (
	"time"

	"gorm.io/gorm"
)

// User Model
type User struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"size:255;not null" json:"name"`
	Email            string         `gorm:"size:255;unique;not null" json:"email"`
	EmailVerifiedAt  *time.Time     `json:"email_verified_at,omitempty"`
	Phone            string         `gorm:"size:50;unique;not null" json:"phone"`
	PhoneConfirmedAt *time.Time     `json:"phone_confirmed_at,omitempty"`
	Password         string         `gorm:"size:255;not null" json:"-"`
	RememberToken    *string        `gorm:"size:100" json:"-"`
	Profile          *Profile       `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// Profile Model
type Profile struct {
	ID                    uint                 `gorm:"primaryKey" json:"id"`
	UserID                uint                 `gorm:"index;not null" json:"-"`
	Age                   uint                 `gorm:"not null" json:"age"`
	Gender                string               `gorm:"type:enum('MALE', 'FEMALE', 'OTHER');not null" json:"gender"`
	DateOfBirth           time.Time            `gorm:"not null" json:"date_of_birth"`
	Address               string               `gorm:"type:text" json:"address"`
	Latitude              *float64             `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude             *float64             `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	MaritalStatus         string               `gorm:"type:enum('SINGLE', 'RELATIONSHIP');not null" json:"marital_status"`
	Height                *uint                `json:"height,omitempty"`
	Weight                *float64             `gorm:"type:decimal(10,2)" json:"weight,omitempty"`
	EducationLevelID      *uint                `gorm:"index" json:"education_level_id,omitempty"`
	EducationLevel        *EducationLevel      `gorm:"foreignKey:EducationLevelID" json:"education_level,omitempty"`
	Occupation            *string              `gorm:"size:255" json:"occupation,omitempty"`
	NumberOfChildren      *uint                `json:"children,omitempty"`
	ContraceptionReasonID uint                 `gorm:"index;not null" json:"reason_id,omitempty"`
	Reason                *ContraceptionReason `gorm:"foreignKey:ContraceptionReasonID" json:"reason,omitempty"`
	SexuallyActive        bool                 `gorm:"not null" json:"sexually_active"`
	PregnancyStatus       bool                 `gorm:"not null" json:"pregnancy_status"`
	Religion              string               `gorm:"type:enum('CHRISTIANITY', 'ISLAM', 'OTHER');not null" json:"religion"`
	ReligionSect          *string              `gorm:"type:enum('CATHOLIC', 'PENTECOSTAL', 'OTHER')" json:"religion_sect,omitempty"`
	Meta                  *string              `gorm:"type:json" json:"meta,omitempty"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

// Clinic Model
type Clinic struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Address   string         `gorm:"type:text;not null" json:"address"`
	Latitude  *float64       `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude *float64       `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	AddedByID uint           `gorm:"not null" json:"added_by_id"`
	Meta      *string        `gorm:"type:json" json:"meta,omitempty"`
	Locations []Location     `gorm:"many2many:locatables;foreignKey:ID;joinForeignKey:LocatableID;references:ID;joinReferences:LocationID" json:"locations,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Computed Distance Fields
	Distance       *string  `gorm:"-" json:"distance,omitempty"`
	ActualDistance *float64 `gorm:"-" json:"actual_distance,omitempty"`
}

// Location Model
type Location struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	StateID   uint           `gorm:"index;not null" json:"state_id"`
	CountryID uint           `gorm:"index;not null" json:"country_id"`
	State     *State         `gorm:"foreignKey:StateID" json:"state,omitempty"`
	Country   *Country       `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	Latitude  *float64       `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude *float64       `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	Meta      *string        `gorm:"type:json" json:"meta,omitempty"`
	Clinics   []Clinic       `gorm:"many2many:locatables;foreignKey:ID;joinForeignKey:LocationID;references:ID;joinReferences:LocatableID" json:"clinics,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Locatable Polymorphic Pivot Model
type Locatable struct {
	ID            uint      `gorm:"primaryKey"`
	LocationID    uint      `gorm:"index"`
	LocatableID   uint      `gorm:"index"`
	LocatableType string    `gorm:"size:255;index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// State Model
type State struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Slug      string         `gorm:"size:255;unique;not null" json:"slug"`
	Latitude  *float64       `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude *float64       `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	Meta      *string        `gorm:"type:json" json:"meta,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Country Model
type Country struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Code      string    `gorm:"size:10;unique;not null" json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Contraception Reason Model
type ContraceptionReason struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Value     string         `gorm:"size:255;not null" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Education Level Model
type EducationLevel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FAQ Group Model
type FaqGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Faq       *Faq      `gorm:"foreignKey:FaqGroupID" json:"faq,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FAQ Model
type Faq struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	FaqGroupID uint           `gorm:"unique;not null" json:"-"` // One-to-one with FAQ Group
	Content    *string        `gorm:"type:json" json:"content"` // supersedes old text columns
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Algorithm Model
type Algorithm struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	Text                   string         `gorm:"type:text;not null" json:"text"`
	ActionType             *string        `gorm:"type:enum('bool', 'input', 'date')" json:"actionType,omitempty"`
	Positive               *string        `gorm:"size:255" json:"positive,omitempty"`
	Negative               *string        `gorm:"size:255" json:"negative,omitempty"`
	OnPositive             *uint          `gorm:"index" json:"onPositive,omitempty"`
	OnNegative             *uint          `gorm:"index" json:"onNegative,omitempty"`
	NextMove               *uint          `gorm:"index" json:"nextMove,omitempty"`
	TempPlan               *string        `gorm:"size:255" json:"tempPlan,omitempty"`
	TempPlanDirP           *string        `gorm:"size:255" json:"tempPlanDirP,omitempty"`
	TempPlanDirN           *string        `gorm:"size:255" json:"tempPlanDirN,omitempty"`
	ConditionalFactor      *string        `gorm:"size:255" json:"conditionalFactor,omitempty"`
	ConditionalOperator    *string        `gorm:"size:50" json:"conditionalOperator,omitempty"`
	ConditionalValue       *string        `gorm:"size:255" json:"conditionalValue,omitempty"`
	StateValue             *string        `gorm:"size:255" json:"stateValue,omitempty"`
	Label                  *string        `gorm:"size:255" json:"label,omitempty"`
	ProgestogenPossible    *string        `gorm:"type:enum('true', 'false')" json:"progestogenPossible,omitempty"`
	ProgestogenPossibleDir *string        `gorm:"type:enum('positive', 'negative')" json:"progestogenPossibleDir,omitempty"`
	Delay                  int            `gorm:"not null" json:"delay"`
	Series                 *int           `json:"series,omitempty"`
	Active                 string         `gorm:"type:enum('Y', 'N');default:'N';not null" json:"active"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// Admin Model
type Admin struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Firstname     string         `gorm:"size:255;not null" json:"firstname"`
	Lastname      string         `gorm:"size:255;not null" json:"lastname"`
	Phone         *string        `gorm:"size:50" json:"phone,omitempty"`
	Email         string         `gorm:"size:255;unique;not null" json:"email"`
	Verified      string         `gorm:"type:enum('Y', 'N');default:'N';not null" json:"verified"`
	UserType      string         `gorm:"type:enum('Super', 'Sub');not null" json:"userType"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	AccessToken   *string        `gorm:"type:text" json:"accessToken,omitempty"`
	RememberToken *string        `gorm:"size:100" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// Password Reset Model
type PasswordReset struct {
	Email     string    `gorm:"size:255;index;not null" json:"email"`
	Token     string    `gorm:"size:255;not null" json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
