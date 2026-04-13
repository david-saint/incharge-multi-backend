package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents the admins table.
type Admin struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Firstname     string         `gorm:"column:firstname;type:varchar(255);not null" json:"firstname"`
	Lastname      string         `gorm:"column:lastname;type:varchar(255);not null" json:"lastname"`
	Phone         NullString     `gorm:"column:phone;type:varchar(50);null" json:"phone"`
	Email         string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Verified      string         `gorm:"column:verified;type:enum('Y','N');default:'N'" json:"verified"`
	UserType      string         `gorm:"column:userType;type:enum('Super','Sub');not null" json:"userType"`
	Password      string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	AccessToken   NullString     `gorm:"column:accessToken;type:text;null" json:"accessToken,omitempty"`
	RememberToken NullString     `gorm:"column:remember_token;type:varchar(100);null" json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// TableName overrides the table name.
func (Admin) TableName() string { return "admins" }

// IsVerified checks if the admin is verified.
func (a *Admin) IsVerified() bool {
	return a.Verified == "Y"
}

// IsSuper checks if the admin is a super admin.
func (a *Admin) IsSuper() bool {
	return a.UserType == "Super"
}
