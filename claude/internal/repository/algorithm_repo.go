package repository

import (
	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// AlgorithmRepo handles algorithm database operations.
type AlgorithmRepo struct {
	db *gorm.DB
}

// NewAlgorithmRepo creates a new AlgorithmRepo.
func NewAlgorithmRepo(db *gorm.DB) *AlgorithmRepo {
	return &AlgorithmRepo{db: db}
}

// ListAll returns all algorithms ordered by active ASC, id ASC.
func (r *AlgorithmRepo) ListAll() ([]model.Algorithm, error) {
	var algorithms []model.Algorithm
	err := r.db.Order("active ASC, id ASC").Find(&algorithms).Error
	return algorithms, err
}

// Create creates a new algorithm step.
func (r *AlgorithmRepo) Create(algo *model.Algorithm) error {
	return r.db.Create(algo).Error
}

// Update updates an algorithm step.
func (r *AlgorithmRepo) Update(algo *model.Algorithm) error {
	return r.db.Save(algo).Error
}

// FindByID finds an algorithm step by ID.
func (r *AlgorithmRepo) FindByID(id uint) (*model.Algorithm, error) {
	var algo model.Algorithm
	if err := r.db.First(&algo, id).Error; err != nil {
		return nil, err
	}
	return &algo, nil
}
