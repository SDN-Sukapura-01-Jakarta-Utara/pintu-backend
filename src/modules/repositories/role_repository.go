package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// RoleRepository handles data operations for Role
type RoleRepository interface {
	Create(data *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Update(data *models.Role) error
	Delete(id uint) error
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new Role repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create creates a new Role record
func (r *RoleRepositoryImpl) Create(data *models.Role) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Role by ID
func (r *RoleRepositoryImpl) GetByID(id uint) (*models.Role, error) {
	var data models.Role
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Role records
func (r *RoleRepositoryImpl) GetAll() ([]models.Role, error) {
	var data []models.Role
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Role record
func (r *RoleRepositoryImpl) Update(data *models.Role) error {
	return r.db.Save(data).Error
}

// Delete deletes Role record by ID
func (r *RoleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
