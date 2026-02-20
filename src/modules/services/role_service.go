package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// RoleService handles business logic for Role
type RoleService interface {
	Create(data *models.Role) error
	CreateWithPermissions(data *models.Role, permissionIDs []uint) error
	GetByID(id uint) (*models.Role, error)
	GetAll() ([]models.Role, error)
	GetAllWithFilter(params repositories.GetRolesParams) ([]models.Role, int64, error)
	GetRoleWithPermissions(id uint) (*models.Role, []models.RolePermission, error)
	GetRoleWithPermissionDetails(id uint) (*models.Role, []models.Permission, error)
	Update(data *models.Role) error
	Delete(id uint) error
	AssignPermissions(roleID uint, permissionIDs []uint) error
}

type RoleServiceImpl struct {
	repository repositories.RoleRepository
}

// NewRoleService creates a new Role service
func NewRoleService(repository repositories.RoleRepository) RoleService {
	return &RoleServiceImpl{repository: repository}
}

// Create creates a new Role
func (s *RoleServiceImpl) Create(data *models.Role) error {
	return s.repository.Create(data)
}

// CreateWithPermissions creates a new Role with permissions
func (s *RoleServiceImpl) CreateWithPermissions(data *models.Role, permissionIDs []uint) error {
	if err := s.repository.Create(data); err != nil {
		return err
	}
	
	// Only assign permissions if there are any
	if len(permissionIDs) > 0 {
		return s.AssignPermissions(data.ID, permissionIDs)
	}
	
	return nil
}

// GetByID retrieves Role by ID
func (s *RoleServiceImpl) GetByID(id uint) (*models.Role, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all Role
func (s *RoleServiceImpl) GetAll() ([]models.Role, error) {
	return s.repository.GetAll()
}

// GetAllWithFilter retrieves roles with filters and pagination
func (s *RoleServiceImpl) GetAllWithFilter(params repositories.GetRolesParams) ([]models.Role, int64, error) {
	return s.repository.GetAllWithFilter(params)
}

// Update updates Role
func (s *RoleServiceImpl) Update(data *models.Role) error {
	return s.repository.Update(data)
}

// GetRoleWithPermissions retrieves role with its permissions
func (s *RoleServiceImpl) GetRoleWithPermissions(id uint) (*models.Role, []models.RolePermission, error) {
	role, err := s.repository.GetByID(id)
	if err != nil {
		return nil, nil, err
	}
	
	permissions, err := s.repository.GetRolePermissions(id)
	if err != nil {
		return nil, nil, err
	}
	
	return role, permissions, nil
}

// GetRoleWithPermissionDetails retrieves role with full permission details
func (s *RoleServiceImpl) GetRoleWithPermissionDetails(id uint) (*models.Role, []models.Permission, error) {
	role, err := s.repository.GetByID(id)
	if err != nil {
		return nil, nil, err
	}
	
	// Get role permission IDs
	rolePermissions, err := s.repository.GetRolePermissions(id)
	if err != nil {
		return nil, nil, err
	}
	
	// Extract permission IDs
	permissionIDs := make([]uint, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissionIDs[i] = rp.PermissionID
	}
	
	// Get full permission details if there are any
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		permissions, err = s.repository.GetPermissionsByIDs(permissionIDs)
		if err != nil {
			return nil, nil, err
		}
	}
	
	return role, permissions, nil
}

// AssignPermissions assigns permissions to a role
func (s *RoleServiceImpl) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// Clear existing permissions first
	if err := s.repository.DeleteRolePermissions(roleID); err != nil {
		return err
	}
	
	// Assign new permissions
	return s.repository.CreateRolePermissions(roleID, permissionIDs)
}

// Delete deletes Role by ID
func (s *RoleServiceImpl) Delete(id uint) error {
	// Delete role permissions first
	if err := s.repository.DeleteRolePermissions(id); err != nil {
		return err
	}
	
	// Then delete the role
	return s.repository.Delete(id)
}
