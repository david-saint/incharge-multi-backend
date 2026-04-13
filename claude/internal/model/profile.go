package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Profile represents the profiles table.
type Profile struct {
	ID                    uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                uint             `gorm:"uniqueIndex;not null" json:"-"`
	Age                   uint             `gorm:"type:int unsigned;default:0" json:"age"`
	Gender                string           `gorm:"type:enum('MALE','FEMALE','OTHER');not null" json:"gender"`
	DateOfBirth           time.Time        `gorm:"type:datetime" json:"date_of_birth"`
	Address               string           `gorm:"type:text" json:"address"`
	Latitude              *float64         `gorm:"type:decimal(10,7);null" json:"latitude"`
	Longitude             *float64         `gorm:"type:decimal(10,7);null" json:"longitude"`
	MaritalStatus         string           `gorm:"type:enum('SINGLE','RELATIONSHIP');default:'SINGLE'" json:"marital_status"`
	Height                *uint            `gorm:"type:int unsigned;null" json:"height"`
	Weight                *float64         `gorm:"type:decimal(10,2);null" json:"weight"`
	EducationLevelID      *uint            `gorm:"null" json:"-"`
	Occupation            sql.NullString   `gorm:"type:varchar(255);null" json:"occupation"`
	NumberOfChildren      *uint            `gorm:"type:int unsigned;null" json:"children"`
	ContraceptionReasonID *uint            `gorm:"null" json:"-"`
	SexuallyActive        bool             `gorm:"type:tinyint(1);default:0" json:"sexually_active"`
	PregnancyStatus       bool             `gorm:"type:tinyint(1);default:0" json:"pregnancy_status"`
	Religion              sql.NullString   `gorm:"type:enum('CHRISTIANITY','ISLAM','OTHER');null" json:"religion"`
	ReligionSect          sql.NullString   `gorm:"type:enum('CATHOLIC','PENTECOSTAL','OTHER');null" json:"religion_sect"`
	Meta                  *json.RawMessage `gorm:"type:json;null" json:"-"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`

	// Relationships
	User                *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	EducationLevel      *EducationLevel      `gorm:"foreignKey:EducationLevelID" json:"education_level,omitempty"`
	ContraceptionReason *ContraceptionReason `gorm:"foreignKey:ContraceptionReasonID" json:"reason,omitempty"`
}

// TableName overrides the table name.
func (Profile) TableName() string { return "profiles" }

// GetMeta parses the meta JSON field into a map.
func (p *Profile) GetMeta() map[string]interface{} {
	if p.Meta == nil {
		return make(map[string]interface{})
	}
	var m map[string]interface{}
	if err := json.Unmarshal(*p.Meta, &m); err != nil {
		return make(map[string]interface{})
	}
	return m
}

// SetMeta serializes a map into the meta JSON field.
func (p *Profile) SetMeta(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	raw := json.RawMessage(b)
	p.Meta = &raw
	return nil
}
