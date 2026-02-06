package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// PermissionService handles business logic for Permission
type PermissionService interface {
	Create(data *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	GetAll(limit, offset int) ([]models.Permission, int64, error)
	GetByGroupName(groupName string) ([]models.Permission, error)
	GetBySystem(system string) ([]models.Permission, error)
	Update(data *models.Permission) error
	Delete(id uint) error
}

type PermissionServiceImpl struct {
	repository repositories.PermissionRepository
}

// NewPermissionService creates a new Permission service
func NewPermissionService(repository repositories.PermissionRepository) PermissionService {
	return &PermissionServiceImpl{repository: repository}
}

// Create creates a new Permission
func (s *PermissionServiceImpl) Create(data *models.Permission) error {
	return s.repository.Create(data)
}

// GetByID retrieves Permission by ID
func (s *PermissionServiceImpl) GetByID(id uint) (*models.Permission, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all Permissions
func (s *PermissionServiceImpl) GetAll(limit, offset int) ([]models.Permission, int64, error) {
	return s.repository.GetAll(limit, offset)
}

// GetByGroupName retrieves permissions by group name
func (s *PermissionServiceImpl) GetByGroupName(groupName string) ([]models.Permission, error) {
	return s.repository.GetByGroupName(groupName)
}

// GetBySystem retrieves permissions by system
func (s *PermissionServiceImpl) GetBySystem(system string) ([]models.Permission, error) {
	return s.repository.GetBySystem(system)
}

// Update updates Permission
func (s *PermissionServiceImpl) Update(data *models.Permission) error {
	return s.repository.Update(data)
}

// Delete deletes Permission by ID
func (s *PermissionServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
