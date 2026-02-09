package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// BidangStudiService handles business logic for BidangStudi
type BidangStudiService interface {
	Create(req *dtos.BidangStudiCreateRequest, userID uint) (*dtos.BidangStudiResponse, error)
	GetByID(id uint) (*dtos.BidangStudiResponse, error)
	GetAll(limit int, offset int) (*dtos.BidangStudiListResponse, error)
	Update(req *dtos.BidangStudiUpdateRequest, userID uint) (*dtos.BidangStudiResponse, error)
	Delete(id uint) error
}

type BidangStudiServiceImpl struct {
	repository repositories.BidangStudiRepository
}

// NewBidangStudiService creates a new BidangStudi service
func NewBidangStudiService(repository repositories.BidangStudiRepository) BidangStudiService {
	return &BidangStudiServiceImpl{repository: repository}
}

// Create creates a new BidangStudi
func (s *BidangStudiServiceImpl) Create(req *dtos.BidangStudiCreateRequest, userID uint) (*dtos.BidangStudiResponse, error) {
	// Check if name already exists
	_, err := s.repository.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("bidang studi dengan nama ini sudah ada")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.BidangStudi{
		Name:        req.Name,
		Status:      status,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves BidangStudi by ID
func (s *BidangStudiServiceImpl) GetByID(id uint) (*dtos.BidangStudiResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all BidangStudi
func (s *BidangStudiServiceImpl) GetAll(limit int, offset int) (*dtos.BidangStudiListResponse, error) {
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
	responses := make([]dtos.BidangStudiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.BidangStudiListResponse{
		Data: responses,
	}, nil
}

// Update updates BidangStudi
func (s *BidangStudiServiceImpl) Update(req *dtos.BidangStudiUpdateRequest, userID uint) (*dtos.BidangStudiResponse, error) {
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

// Delete deletes BidangStudi by ID
func (s *BidangStudiServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *BidangStudiServiceImpl) mapToResponse(data *models.BidangStudi) *dtos.BidangStudiResponse {
	return &dtos.BidangStudiResponse{
		ID:          data.ID,
		Name:        data.Name,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
