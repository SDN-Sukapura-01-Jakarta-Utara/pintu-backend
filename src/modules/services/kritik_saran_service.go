package services

import (
	"errors"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// KritikSaranService handles business logic for KritikSaran
type KritikSaranService interface {
	Create(req *dtos.KritikSaranCreateRequest) (*dtos.KritikSaranResponse, error)
	GetByID(id uint) (*dtos.KritikSaranResponse, error)
	GetAll(limit int, offset int) (*dtos.KritikSaranListResponse, error)
	GetAllWithFilter(params repositories.GetKritikSaranParams) (*dtos.KritikSaranListResponse, error)
	Delete(id uint) error
}

type KritikSaranServiceImpl struct {
	repository repositories.KritikSaranRepository
}

// NewKritikSaranService creates a new KritikSaran service
func NewKritikSaranService(repository repositories.KritikSaranRepository) KritikSaranService {
	return &KritikSaranServiceImpl{repository: repository}
}

// Create creates a new KritikSaran
func (s *KritikSaranServiceImpl) Create(req *dtos.KritikSaranCreateRequest) (*dtos.KritikSaranResponse, error) {
	// Create kritik saran record
	data := &models.KritikSaran{
		Nama:        req.Nama,
		KritikSaran: req.KritikSaran,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves KritikSaran by ID
func (s *KritikSaranServiceImpl) GetByID(id uint) (*dtos.KritikSaranResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all KritikSaran with pagination
func (s *KritikSaranServiceImpl) GetAll(limit int, offset int) (*dtos.KritikSaranListResponse, error) {
	// Set default limit and offset
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	data, total, err := s.repository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.KritikSaranResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.KritikSaranListResponse{
		Data:  responses,
		Total: total,
	}, nil
}

// GetAllWithFilter retrieves KritikSaran with filters and pagination
func (s *KritikSaranServiceImpl) GetAllWithFilter(params repositories.GetKritikSaranParams) (*dtos.KritikSaranListResponse, error) {
	// Validate and set default limit and offset
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.KritikSaranResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.KritikSaranListResponse{
		Data:  responses,
		Total: total,
	}, nil
}

// Delete deletes KritikSaran by ID
func (s *KritikSaranServiceImpl) Delete(id uint) error {
	// Check if exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("kritik saran not found")
	}

	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *KritikSaranServiceImpl) mapToResponse(data *models.KritikSaran) *dtos.KritikSaranResponse {
	return &dtos.KritikSaranResponse{
		ID:          data.ID,
		Nama:        data.Nama,
		KritikSaran: data.KritikSaran,
		CreatedAt:   data.CreatedAt,
	}
}
