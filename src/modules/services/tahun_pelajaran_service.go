package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// TahunPelajaranService handles business logic for TahunPelajaran
type TahunPelajaranService interface {
	Create(req *dtos.TahunPelajaranCreateRequest, userID uint) (*dtos.TahunPelajaranResponse, error)
	GetByID(id uint) (*dtos.TahunPelajaranResponse, error)
	GetAll(limit int, offset int) (*dtos.TahunPelajaranListResponse, error)
	GetAllWithFilter(params repositories.GetTahunPelajaranParams) (*dtos.TahunPelajaranListWithPaginationResponse, error)
	Update(req *dtos.TahunPelajaranUpdateRequest, userID uint) (*dtos.TahunPelajaranResponse, error)
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
func (s *TahunPelajaranServiceImpl) Create(req *dtos.TahunPelajaranCreateRequest, userID uint) (*dtos.TahunPelajaranResponse, error) {
	// Check if tahun_pelajaran already exists
	_, err := s.repository.GetByTahunPelajaran(req.TahunPelajaran)
	if err == nil {
		return nil, errors.New("tahun pelajaran sudah ada")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// If status is active, set all others to inactive
	if status == "active" {
		if err := s.repository.UpdateAllStatusToInactive(); err != nil {
			return nil, err
		}
	}

	data := &models.TahunPelajaran{
		TahunPelajaran: req.TahunPelajaran,
		Status:         status,
		CreatedByID:    &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves TahunPelajaran by ID
func (s *TahunPelajaranServiceImpl) GetByID(id uint) (*dtos.TahunPelajaranResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all TahunPelajaran
func (s *TahunPelajaranServiceImpl) GetAll(limit int, offset int) (*dtos.TahunPelajaranListResponse, error) {
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
	responses := make([]dtos.TahunPelajaranResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.TahunPelajaranListResponse{
		Data: responses,
	}, nil
}

// Update updates TahunPelajaran
func (s *TahunPelajaranServiceImpl) Update(req *dtos.TahunPelajaranUpdateRequest, userID uint) (*dtos.TahunPelajaranResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.TahunPelajaran != nil {
		existing.TahunPelajaran = *req.TahunPelajaran
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	// If status is being set to active, set all others to inactive
	if req.Status != nil && *req.Status == "active" {
		if err := s.repository.UpdateAllStatusToInactive(); err != nil {
			return nil, err
		}
	}

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes TahunPelajaran by ID
func (s *TahunPelajaranServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// GetAllWithFilter retrieves TahunPelajaran with filters and pagination
func (s *TahunPelajaranServiceImpl) GetAllWithFilter(params repositories.GetTahunPelajaranParams) (*dtos.TahunPelajaranListWithPaginationResponse, error) {
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
	responses := make([]dtos.TahunPelajaranResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.TahunPelajaranListWithPaginationResponse{
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
func (s *TahunPelajaranServiceImpl) mapToResponse(data *models.TahunPelajaran) *dtos.TahunPelajaranResponse {
	return &dtos.TahunPelajaranResponse{
		ID:             data.ID,
		TahunPelajaran: data.TahunPelajaran,
		Status:         data.Status,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
		CreatedByID:    data.CreatedByID,
		UpdatedByID:    data.UpdatedByID,
	}
}
