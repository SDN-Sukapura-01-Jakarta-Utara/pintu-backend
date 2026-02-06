package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// PermissionRepository handles data operations for Permission
type PermissionRepository interface {
	Create(data *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	GetAll(limit, offset int) ([]models.Permission, int64, error)
	GetByGroupName(groupName string) ([]models.Permission, error)
	GetBySystem(system string) ([]models.Permission, error)
	Update(data *models.Permission) error
	Delete(id uint) error
}

type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new Permission repository
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// Create creates a new Permission record
func (r *PermissionRepositoryImpl) Create(data *models.Permission) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Permission by ID
func (r *PermissionRepositoryImpl) GetByID(id uint) (*models.Permission, error) {
	var data models.Permission
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Permission records with pagination
func (r *PermissionRepositoryImpl) GetAll(limit, offset int) ([]models.Permission, int64, error) {
	var data []models.Permission
	var total int64

	if err := r.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByGroupName retrieves permissions by group name
func (r *PermissionRepositoryImpl) GetByGroupName(groupName string) ([]models.Permission, error) {
	var data []models.Permission
	if err := r.db.Where("group_name = ?", groupName).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetBySystem retrieves permissions by system
func (r *PermissionRepositoryImpl) GetBySystem(system string) ([]models.Permission, error) {
	var data []models.Permission
	if err := r.db.Where("system = ?", system).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Permission record
func (r *PermissionRepositoryImpl) Update(data *models.Permission) error {
	return r.db.Save(data).Error
}

// Delete deletes Permission record by ID
func (r *PermissionRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}
