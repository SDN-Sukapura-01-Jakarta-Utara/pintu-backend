package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// RoleService handles business logic for Role
type RoleService interface {
	Create(data *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Update(data *models.Role) error
	Delete(id uint) error
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

// GetByID retrieves Role by ID
func (s *RoleServiceImpl) GetByID(id uint) (*models.Role, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all Role
func (s *RoleServiceImpl) GetAll() ([]models.Role, error) {
	return s.repository.GetAll()
}

// Update updates Role
func (s *RoleServiceImpl) Update(data *models.Role) error {
	return s.repository.Update(data)
}

// Delete deletes Role by ID
func (s *RoleServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
