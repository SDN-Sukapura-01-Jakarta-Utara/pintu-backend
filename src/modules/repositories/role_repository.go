package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// RoleRepository handles data operations for Role
type RoleRepository interface {
	Create(data *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetAll() ([]models.Role, error)
	GetAllWithFilter(params GetRolesParams) ([]models.Role, int64, error)
	Update(data *models.Role) error
	Delete(id uint) error
	CreateRolePermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]models.RolePermission, error)
	DeleteRolePermissions(roleID uint) error
	GetPermissionsByIDs(ids []uint) ([]models.Permission, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

// GetRolesFilter represents filters for getting roles
type GetRolesFilter struct {
	Name     string
	SystemID uint
	Status   string
}

// GetRolesParams represents parameters for getting roles with filters
type GetRolesParams struct {
	Filter GetRolesFilter
	Limit  int
	Offset int
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
	if err := r.db.Preload("System").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Role records sorted by created_at DESC
func (r *RoleRepositoryImpl) GetAll() ([]models.Role, error) {
	var data []models.Role
	if err := r.db.Preload("System").Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllWithFilter retrieves roles with filters and pagination
func (r *RoleRepositoryImpl) GetAllWithFilter(params GetRolesParams) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.Preload("System")

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Filter.Name)+"%")
	}
	if params.Filter.SystemID > 0 {
		query = query.Where("system_id = ?", params.Filter.SystemID)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data with pagination and sorting by created_at DESC (newest first)
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Update updates Role record
func (r *RoleRepositoryImpl) Update(data *models.Role) error {
	result := r.db.Model(&models.Role{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"name":           data.Name,
		"description":    data.Description,
		"system_id":      data.SystemID,
		"status":         data.Status,
		"updated_by_id":  data.UpdatedByID,
		"updated_at":     data.UpdatedAt,
	})
	return result.Error
}

// Delete deletes Role record by ID
func (r *RoleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// CreateRolePermissions creates multiple role-permission associations
func (r *RoleRepositoryImpl) CreateRolePermissions(roleID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	
	rolePermissions := make([]models.RolePermission, len(permissionIDs))
	for i, permissionID := range permissionIDs {
		rolePermissions[i] = models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
	}
	
	return r.db.CreateInBatches(rolePermissions, 100).Error
}

// GetRolePermissions retrieves all permissions for a role
func (r *RoleRepositoryImpl) GetRolePermissions(roleID uint) ([]models.RolePermission, error) {
	var permissions []models.RolePermission
	if err := r.db.Where("role_id = ?", roleID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// DeleteRolePermissions deletes all permissions for a role
func (r *RoleRepositoryImpl) DeleteRolePermissions(roleID uint) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

// GetPermissionsByIDs retrieves permissions by multiple IDs
func (r *RoleRepositoryImpl) GetPermissionsByIDs(ids []uint) ([]models.Permission, error) {
	var permissions []models.Permission
	if len(ids) == 0 {
		return permissions, nil
	}
	if err := r.db.Preload("System").Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
