package repository

import (
	"errors"
	"time"

	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// UserRepo handles user database operations.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user.
func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID (excluding soft-deleted).
func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email.
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// EmailExists checks if an email is already taken.
func (r *UserRepo) EmailExists(email string) bool {
	var count int64
	r.db.Unscoped().Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// PhoneExists checks if a phone number is already taken.
// Returns false for empty strings (no phone provided).
func (r *UserRepo) PhoneExists(phone string) bool {
	if phone == "" {
		return false
	}
	var count int64
	r.db.Unscoped().Model(&model.User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// VerifyEmail sets the email_verified_at timestamp.
func (r *UserRepo) VerifyEmail(id uint) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("email_verified_at", &now).Error
}

// UpdatePassword updates a user's password.
func (r *UserRepo) UpdatePassword(id uint, hash string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", hash).Error
}

// SoftDelete soft-deletes a user.
func (r *UserRepo) SoftDelete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Restore restores a soft-deleted user.
func (r *UserRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// ListWithProfile returns paginated users with profiles eager-loaded.
func (r *UserRepo) ListWithProfile(page, perPage int) ([]model.User, int64, error) {
	var total int64
	r.db.Model(&model.User{}).Count(&total)

	var users []model.User
	offset := (page - 1) * perPage
	err := r.db.Preload("Profile").Preload("Profile.EducationLevel").Preload("Profile.ContraceptionReason").
		Order("id desc").Offset(offset).Limit(perPage).Find(&users).Error
	return users, total, err
}

// ListDeleted returns all soft-deleted users.
func (r *UserRepo) ListDeleted() ([]model.User, error) {
	var users []model.User
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").
		Preload("Profile").Find(&users).Error
	return users, err
}

// --- Password Reset ---

// CreatePasswordReset stores a password reset token.
func (r *UserRepo) CreatePasswordReset(email, token string) error {
	// Delete any existing tokens for this email.
	r.db.Where("email = ?", email).Delete(&model.PasswordReset{})
	now := time.Now()
	return r.db.Create(&model.PasswordReset{
		Email:     email,
		Token:     token,
		CreatedAt: &now,
	}).Error
}

// FindPasswordReset finds a valid (non-expired) password reset token.
func (r *UserRepo) FindPasswordReset(email, token string) (*model.PasswordReset, error) {
	var pr model.PasswordReset
	cutoff := time.Now().Add(-60 * time.Minute)
	err := r.db.Where("email = ? AND token = ? AND created_at > ?", email, token, cutoff).First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pr, nil
}

// DeletePasswordResets removes all reset tokens for an email.
func (r *UserRepo) DeletePasswordResets(email string) error {
	return r.db.Where("email = ?", email).Delete(&model.PasswordReset{}).Error
}
