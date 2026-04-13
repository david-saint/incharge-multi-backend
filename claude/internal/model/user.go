package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// User represents the users table.
type User struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Email            string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	EmailVerifiedAt  *time.Time     `gorm:"type:timestamp;null" json:"-"`
	Phone            *string        `gorm:"type:varchar(50);uniqueIndex" json:"phone"`
	PhoneConfirmedAt *time.Time     `gorm:"type:timestamp;null" json:"-"`
	Password         string         `gorm:"type:varchar(255);not null" json:"-"`
	RememberToken    sql.NullString `gorm:"type:varchar(100);null" json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	// Relationships
	Profile *Profile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// TableName overrides the table name.
func (User) TableName() string { return "users" }

// IsEmailVerified returns whether the user's email is verified.
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsPhoneConfirmed returns whether the user's phone is confirmed.
func (u *User) IsPhoneConfirmed() bool {
	return u.PhoneConfirmedAt != nil
}
