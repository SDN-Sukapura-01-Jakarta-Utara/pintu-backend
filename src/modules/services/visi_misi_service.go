package services

import (
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// VisiMisiService handles business logic for VisiMisi
type VisiMisiService interface {
	Create(req *dtos.VisiMisiCreateRequest, userID uint) (*dtos.VisiMisiResponse, error)
	GetByID(id uint) (*dtos.VisiMisiResponse, error)
	GetAll(limit int, offset int) (*dtos.VisiMisiListResponse, error)
	Update(req *dtos.VisiMisiUpdateRequest, userID uint) (*dtos.VisiMisiResponse, error)
	Delete(id uint) error
}

type VisiMisiServiceImpl struct {
	repository repositories.VisiMisiRepository
}

// NewVisiMisiService creates a new VisiMisi service
func NewVisiMisiService(repository repositories.VisiMisiRepository) VisiMisiService {
	return &VisiMisiServiceImpl{repository: repository}
}

// Create creates a new VisiMisi
func (s *VisiMisiServiceImpl) Create(req *dtos.VisiMisiCreateRequest, userID uint) (*dtos.VisiMisiResponse, error) {
	data := &models.VisiMisi{
		Visi:        req.Visi,
		Misi:        req.Misi,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves VisiMisi by ID
func (s *VisiMisiServiceImpl) GetByID(id uint) (*dtos.VisiMisiResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all VisiMisi
func (s *VisiMisiServiceImpl) GetAll(limit int, offset int) (*dtos.VisiMisiListResponse, error) {
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
	responses := make([]dtos.VisiMisiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.VisiMisiListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// Update updates VisiMisi
func (s *VisiMisiServiceImpl) Update(req *dtos.VisiMisiUpdateRequest, userID uint) (*dtos.VisiMisiResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Visi != nil {
		existing.Visi = *req.Visi
	}
	if req.Misi != nil {
		existing.Misi = *req.Misi
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes VisiMisi by ID
func (s *VisiMisiServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *VisiMisiServiceImpl) mapToResponse(data *models.VisiMisi) *dtos.VisiMisiResponse {
	return &dtos.VisiMisiResponse{
		ID:          data.ID,
		Visi:        data.Visi,
		Misi:        data.Misi,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
