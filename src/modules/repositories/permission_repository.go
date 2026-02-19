package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// PermissionRepository handles data operations for Permission
type PermissionRepository interface {
	Create(data *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	GetAll(limit, offset int) ([]models.Permission, int64, error)
	GetAllWithFilter(params GetPermissionsParams) ([]models.Permission, int64, error)
	GetByGroupName(groupName string) ([]models.Permission, error)
	GetBySystem(system string) ([]models.Permission, error)
	Update(data *models.Permission) error
	Delete(id uint) error
}

// GetPermissionsFilter represents filters for getting permissions
type GetPermissionsFilter struct {
	Name     string
	GroupName string
	SystemID uint
	Status   string
}

// GetPermissionsParams represents parameters for getting permissions with filters
type GetPermissionsParams struct {
	Filter GetPermissionsFilter
	Limit  int
	Offset int
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
	if err := r.db.Preload("System").First(&data, id).Error; err != nil {
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

	if err := r.db.Preload("System").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves permissions with filters and pagination
func (r *PermissionRepositoryImpl) GetAllWithFilter(params GetPermissionsParams) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	query := r.db.Preload("System")

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Filter.Name)+"%")
	}
	if params.Filter.GroupName != "" {
		query = query.Where("LOWER(group_name) LIKE ?", "%"+strings.ToLower(params.Filter.GroupName)+"%")
	}
	if params.Filter.SystemID > 0 {
		query = query.Where("system_id = ?", params.Filter.SystemID)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data with pagination
	if err := query.Limit(params.Limit).Offset(params.Offset).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetByGroupName retrieves permissions by group name
func (r *PermissionRepositoryImpl) GetByGroupName(groupName string) ([]models.Permission, error) {
	var data []models.Permission
	if err := r.db.Where("group_name = ?", groupName).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetBySystem retrieves permissions by system ID
func (r *PermissionRepositoryImpl) GetBySystem(systemID string) ([]models.Permission, error) {
	var data []models.Permission
	if err := r.db.Where("system_id = ?", systemID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Permission record
func (r *PermissionRepositoryImpl) Update(data *models.Permission) error {
	result := r.db.Model(&models.Permission{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"name":           data.Name,
		"description":    data.Description,
		"group_name":     data.GroupName,
		"system_id":      data.SystemID,
		"status":         data.Status,
		"updated_by_id":  data.UpdatedByID,
		"updated_at":     data.UpdatedAt,
	})
	return result.Error
}

// Delete deletes Permission record by ID
func (r *PermissionRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}
