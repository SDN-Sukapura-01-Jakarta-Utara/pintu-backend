package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// SystemService handles business logic for System
type SystemService interface {
	Create(data *models.System) error
	GetByID(id uint) (*models.System, error)
	GetAll() ([]models.System, error)
	GetAllWithFilter(params repositories.GetSystemsParams) ([]models.System, int64, error)
	Update(data *models.System) error
	Delete(id uint) error
}

type SystemServiceImpl struct {
	repository repositories.SystemRepository
}

// NewSystemService creates a new System service
func NewSystemService(repository repositories.SystemRepository) SystemService {
	return &SystemServiceImpl{repository: repository}
}

// Create creates a new System
func (s *SystemServiceImpl) Create(data *models.System) error {
	return s.repository.Create(data)
}

// GetByID retrieves System by ID
func (s *SystemServiceImpl) GetByID(id uint) (*models.System, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all System
func (s *SystemServiceImpl) GetAll() ([]models.System, error) {
	return s.repository.GetAll()
}

// GetAllWithFilter retrieves systems with filters and pagination
func (s *SystemServiceImpl) GetAllWithFilter(params repositories.GetSystemsParams) ([]models.System, int64, error) {
	return s.repository.GetAllWithFilter(params)
}

// Update updates System
func (s *SystemServiceImpl) Update(data *models.System) error {
	return s.repository.Update(data)
}

// Delete deletes System by ID
func (s *SystemServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
