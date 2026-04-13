package model

import (
	"time"

	"gorm.io/gorm"
)

// State represents the states table.
type State struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Latitude  *float64       `gorm:"type:decimal(10,7);null" json:"latitude"`
	Longitude *float64       `gorm:"type:decimal(10,7);null" json:"longitude"`
	Meta      *string        `gorm:"type:json;null" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TableName overrides the table name.
func (State) TableName() string { return "states" }

// Country represents the countries table.
type Country struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name.
func (Country) TableName() string { return "countries" }
