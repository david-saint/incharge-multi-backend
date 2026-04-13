package repository

import (
	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// AdminRepo handles admin database operations.
type AdminRepo struct {
	db *gorm.DB
}

// NewAdminRepo creates a new AdminRepo.
func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// Create inserts a new admin.
func (r *AdminRepo) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

// FindByID finds an admin by ID.
func (r *AdminRepo) FindByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByEmail finds an admin by email.
func (r *AdminRepo) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// SuperAdminExists checks if any Super admin exists.
func (r *AdminRepo) SuperAdminExists() bool {
	var count int64
	r.db.Model(&model.Admin{}).Where("userType = ?", "Super").Count(&count)
	return count > 0
}

// ListPaginated returns paginated admins ordered by verified DESC.
func (r *AdminRepo) ListPaginated(page, perPage int) ([]model.Admin, int64, error) {
	var total int64
	r.db.Model(&model.Admin{}).Count(&total)

	var admins []model.Admin
	offset := (page - 1) * perPage
	err := r.db.Order("verified DESC, id ASC").Offset(offset).Limit(perPage).Find(&admins).Error
	return admins, total, err
}

// Update updates an admin record.
func (r *AdminRepo) Update(admin *model.Admin) error {
	return r.db.Save(admin).Error
}
