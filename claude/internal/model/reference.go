package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ContraceptionReason represents the contraception_reasons table.
type ContraceptionReason struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Value     string         `gorm:"type:varchar(255);not null" json:"value"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TableName overrides the table name.
func (ContraceptionReason) TableName() string { return "contraception_reasons" }

// EducationLevel represents the education_levels table.
type EducationLevel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name.
func (EducationLevel) TableName() string { return "education_levels" }

// FaqGroup represents the faq_groups table.
type FaqGroup struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships (one-to-one with FAQ)
	Faq *Faq `gorm:"foreignKey:FaqGroupID" json:"faq,omitempty"`
}

// TableName overrides the table name.
func (FaqGroup) TableName() string { return "faq_groups" }

// Faq represents the faqs table.
type Faq struct {
	ID         uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	FaqGroupID uint             `gorm:"uniqueIndex;not null" json:"faq_group_id"`
	Content    *json.RawMessage `gorm:"type:json;null" json:"content"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// TableName overrides the table name.
func (Faq) TableName() string { return "faqs" }
