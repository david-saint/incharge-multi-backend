package repository

import (
	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// ReferenceRepo handles read operations for reference/lookup tables.
type ReferenceRepo struct {
	db *gorm.DB
}

// NewReferenceRepo creates a new ReferenceRepo.
func NewReferenceRepo(db *gorm.DB) *ReferenceRepo {
	return &ReferenceRepo{db: db}
}

// --- Contraception Reasons ---

// ListContraceptionReasons returns all contraception reasons.
func (r *ReferenceRepo) ListContraceptionReasons() ([]model.ContraceptionReason, error) {
	var reasons []model.ContraceptionReason
	err := r.db.Find(&reasons).Error
	return reasons, err
}

// FindContraceptionReason finds a contraception reason by ID.
func (r *ReferenceRepo) FindContraceptionReason(id uint) (*model.ContraceptionReason, error) {
	var reason model.ContraceptionReason
	if err := r.db.First(&reason, id).Error; err != nil {
		return nil, err
	}
	return &reason, nil
}

// --- Education Levels ---

// ListEducationLevels returns all education levels.
func (r *ReferenceRepo) ListEducationLevels() ([]model.EducationLevel, error) {
	var levels []model.EducationLevel
	err := r.db.Find(&levels).Error
	return levels, err
}

// --- FAQ Groups ---

// ListFaqGroups returns all FAQ groups.
func (r *ReferenceRepo) ListFaqGroups() ([]model.FaqGroup, error) {
	var groups []model.FaqGroup
	err := r.db.Find(&groups).Error
	return groups, err
}

// FindFaqGroupWithContent finds a FAQ group and its associated FAQ content.
func (r *ReferenceRepo) FindFaqGroupWithContent(id uint) (*model.FaqGroup, error) {
	var group model.FaqGroup
	if err := r.db.Preload("Faq").First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
