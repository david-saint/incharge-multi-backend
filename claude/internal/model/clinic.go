package model

import (
	"time"

	"gorm.io/gorm"
)

// Clinic represents the clinics table.
type Clinic struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Address   string         `gorm:"type:text;not null" json:"address"`
	Latitude  *float64       `gorm:"type:decimal(10,7);null" json:"latitude"`
	Longitude *float64       `gorm:"type:decimal(10,7);null" json:"longitude"`
	AddedByID uint           `gorm:"not null" json:"-"`
	Meta      *string        `gorm:"type:json;null" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	Locations []Location `gorm:"many2many:locatables;foreignKey:ID;joinForeignKey:LocatableID;References:ID;joinReferences:LocationID" json:"locations,omitempty"`

	// Computed fields (not stored, populated by queries)
	Mode           string  `gorm:"-" json:"mode,omitempty"`
	Radius         float64 `gorm:"-" json:"radius,omitempty"`
	SearchRadius   string  `gorm:"-" json:"search_radius,omitempty"`
	ActualDistance float64 `gorm:"-" json:"actual_distance,omitempty"`
	Distance       string  `gorm:"-" json:"distance,omitempty"`
	DistanceRaw    float64 `gorm:"-" json:"-"`
	LocationsCount *int64  `gorm:"-" json:"locations_count,omitempty"`
}

// TableName overrides the table name.
func (Clinic) TableName() string { return "clinics" }

// Location represents the locations table.
type Location struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	StateID   uint           `gorm:"not null" json:"state_id"`
	CountryID uint           `gorm:"not null" json:"country_id"`
	Latitude  *float64       `gorm:"type:decimal(10,7);null" json:"latitude"`
	Longitude *float64       `gorm:"type:decimal(10,7);null" json:"longitude"`
	Meta      *string        `gorm:"type:json;null" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	State   *State   `gorm:"foreignKey:StateID" json:"state,omitempty"`
	Country *Country `gorm:"foreignKey:CountryID" json:"country,omitempty"`
}

// TableName overrides the table name.
func (Location) TableName() string { return "locations" }

// Locatable represents the polymorphic pivot table.
type Locatable struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LocationID    uint      `gorm:"not null" json:"location_id"`
	LocatableID   uint      `gorm:"not null" json:"locatable_id"`
	LocatableType string    `gorm:"type:varchar(255);not null" json:"locatable_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName overrides the table name.
func (Locatable) TableName() string { return "locatables" }
