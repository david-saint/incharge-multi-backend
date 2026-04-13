package model

import "time"

// PasswordReset represents the password_resets table.
type PasswordReset struct {
	Email     string     `gorm:"type:varchar(255);index;not null" json:"email"`
	Token     string     `gorm:"type:varchar(255);not null" json:"token"`
	CreatedAt *time.Time `gorm:"type:timestamp;null" json:"created_at"`
}

// TableName overrides the table name.
func (PasswordReset) TableName() string { return "password_resets" }
