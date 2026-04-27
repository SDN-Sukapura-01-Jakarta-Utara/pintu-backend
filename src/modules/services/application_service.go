package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// ApplicationService handles business logic for Application
type ApplicationService interface {
	Create(req *dtos.ApplicationCreateRequest, userID uint) (*dtos.ApplicationResponse, error)
	GetByID(id uint) (*dtos.ApplicationResponse, error)
	GetAll(limit int, offset int) (*dtos.ApplicationListResponse, error)
	GetAllWithFilter(params repositories.GetApplicationParams) (*dtos.ApplicationListWithPaginationResponse, error)
	Update(req *dtos.ApplicationUpdateRequest, userID uint) (*dtos.ApplicationResponse, error)
	Delete(id uint) error
}

type ApplicationServiceImpl struct {
	repository repositories.ApplicationRepository
}

// NewApplicationService creates a new Application service
func NewApplicationService(repository repositories.ApplicationRepository) ApplicationService {
	return &ApplicationServiceImpl{repository: repository}
}

// Create creates a new Application
func (s *ApplicationServiceImpl) Create(req *dtos.ApplicationCreateRequest, userID uint) (*dtos.ApplicationResponse, error) {
	// Check if nama already exists
	_, err := s.repository.GetByNama(req.Nama)
	if err == nil {
		return nil, errors.New("application dengan nama ini sudah ada")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// If show_in_jumbotron is true, unset all other applications
	if req.ShowInJumbotron {
		if err := s.repository.UnsetAllShowInJumbotron(); err != nil {
			return nil, err
		}
	}

	data := &models.Application{
		Nama:            req.Nama,
		Link:            req.Link,
		ShowInJumbotron: req.ShowInJumbotron,
		Status:          status,
		CreatedByID:     &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Application by ID
func (s *ApplicationServiceImpl) GetByID(id uint) (*dtos.ApplicationResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Application
func (s *ApplicationServiceImpl) GetAll(limit int, offset int) (*dtos.ApplicationListResponse, error) {
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
	responses := make([]dtos.ApplicationResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.ApplicationListResponse{
		Data: responses,
	}, nil
}

// Update updates Application
func (s *ApplicationServiceImpl) Update(req *dtos.ApplicationUpdateRequest, userID uint) (*dtos.ApplicationResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Nama != nil {
		existing.Nama = *req.Nama
	}
	if req.Link != nil {
		existing.Link = *req.Link
	}
	if req.ShowInJumbotron != nil {
		// If show_in_jumbotron is being set to true, unset all other applications
		if *req.ShowInJumbotron {
			if err := s.repository.UnsetAllShowInJumbotron(); err != nil {
				return nil, err
			}
		}
		existing.ShowInJumbotron = *req.ShowInJumbotron
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

// Delete deletes Application by ID
func (s *ApplicationServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// GetAllWithFilter retrieves Application with filters and pagination
func (s *ApplicationServiceImpl) GetAllWithFilter(params repositories.GetApplicationParams) (*dtos.ApplicationListWithPaginationResponse, error) {
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
	responses := make([]dtos.ApplicationResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.ApplicationListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationInfo{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       (params.Offset / params.Limit) + 1,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// mapToResponse maps model to DTO response
func (s *ApplicationServiceImpl) mapToResponse(data *models.Application) *dtos.ApplicationResponse {
	return &dtos.ApplicationResponse{
		ID:              data.ID,
		Nama:            data.Nama,
		Link:            data.Link,
		ShowInJumbotron: data.ShowInJumbotron,
		Status:          data.Status,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
		CreatedByID:     data.CreatedByID,
		UpdatedByID:     data.UpdatedByID,
	}
}
