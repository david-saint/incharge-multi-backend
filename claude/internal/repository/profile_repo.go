package repository

import (
	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// ProfileRepo handles profile database operations.
type ProfileRepo struct {
	db *gorm.DB
}

// NewProfileRepo creates a new ProfileRepo.
func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

// FindByUserID finds a profile by its user ID.
func (r *ProfileRepo) FindByUserID(userID uint, withs ...string) (*model.Profile, error) {
	q := r.db.Where("user_id = ?", userID)
	for _, w := range withs {
		switch w {
		case "user":
			q = q.Preload("User")
		case "reason":
			q = q.Preload("ContraceptionReason")
		case "educationLevel":
			q = q.Preload("EducationLevel")
		}
	}
	var profile model.Profile
	if err := q.First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// Upsert creates or updates a profile for the given user ID.
func (r *ProfileRepo) Upsert(userID uint, profile *model.Profile) error {
	var existing model.Profile
	err := r.db.Where("user_id = ?", userID).First(&existing).Error
	if err == nil {
		// Update existing.
		profile.ID = existing.ID
		profile.UserID = userID
		return r.db.Save(profile).Error
	}
	// Create new.
	profile.UserID = userID
	return r.db.Create(profile).Error
}

// UpdateMeta updates just the meta JSON field.
func (r *ProfileRepo) UpdateMeta(profileID uint, meta interface{}) error {
	return r.db.Model(&model.Profile{}).Where("id = ?", profileID).Update("meta", meta).Error
}
