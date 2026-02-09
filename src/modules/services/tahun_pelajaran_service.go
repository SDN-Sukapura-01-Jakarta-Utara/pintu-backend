package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// TahunPelajaranService handles business logic for TahunPelajaran
type TahunPelajaranService interface {
	Create(data *models.TahunPelajaran) error
	GetByID(id uint) (*models.TahunPelajaran, error)
	GetAll() ([]models.TahunPelajaran, error)
	Update(data *models.TahunPelajaran) error
	Delete(id uint) error
}

type TahunPelajaranServiceImpl struct {
	repository repositories.TahunPelajaranRepository
}

// NewTahunPelajaranService creates a new TahunPelajaran service
func NewTahunPelajaranService(repository repositories.TahunPelajaranRepository) TahunPelajaranService {
	return &TahunPelajaranServiceImpl{repository: repository}
}

// Create creates a new TahunPelajaran
func (s *TahunPelajaranServiceImpl) Create(data *models.TahunPelajaran) error {
	return s.repository.Create(data)
}

// GetByID retrieves TahunPelajaran by ID
func (s *TahunPelajaranServiceImpl) GetByID(id uint) (*models.TahunPelajaran, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all TahunPelajaran
func (s *TahunPelajaranServiceImpl) GetAll() ([]models.TahunPelajaran, error) {
	return s.repository.GetAll()
}

// Update updates TahunPelajaran
func (s *TahunPelajaranServiceImpl) Update(data *models.TahunPelajaran) error {
	return s.repository.Update(data)
}

// Delete deletes TahunPelajaran by ID
func (s *TahunPelajaranServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
