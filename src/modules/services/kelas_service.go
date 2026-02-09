package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// KelasService handles business logic for Kelas
type KelasService interface {
	Create(req *dtos.KelasCreateRequest, userID uint) (*dtos.KelasResponse, error)
	GetByID(id uint) (*dtos.KelasResponse, error)
	GetAll(limit int, offset int) (*dtos.KelasListResponse, error)
	Update(req *dtos.KelasUpdateRequest, userID uint) (*dtos.KelasResponse, error)
	Delete(id uint) error
}

type KelasServiceImpl struct {
	repository repositories.KelasRepository
}

// NewKelasService creates a new Kelas service
func NewKelasService(repository repositories.KelasRepository) KelasService {
	return &KelasServiceImpl{repository: repository}
}

// Create creates a new Kelas
func (s *KelasServiceImpl) Create(req *dtos.KelasCreateRequest, userID uint) (*dtos.KelasResponse, error) {
	// Check if name already exists
	_, err := s.repository.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("kelas dengan nama ini sudah ada")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.Kelas{
		Name:        req.Name,
		Status:      status,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Kelas by ID
func (s *KelasServiceImpl) GetByID(id uint) (*dtos.KelasResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Kelas
func (s *KelasServiceImpl) GetAll(limit int, offset int) (*dtos.KelasListResponse, error) {
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

	data, _, err := s.repository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.KelasResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.KelasListResponse{
		Data: responses,
	}, nil
}

// Update updates Kelas
func (s *KelasServiceImpl) Update(req *dtos.KelasUpdateRequest, userID uint) (*dtos.KelasResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Kelas by ID
func (s *KelasServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *KelasServiceImpl) mapToResponse(data *models.Kelas) *dtos.KelasResponse {
	return &dtos.KelasResponse{
		ID:          data.ID,
		Name:        data.Name,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
