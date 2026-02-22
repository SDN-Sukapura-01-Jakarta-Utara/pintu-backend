package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// RombelService handles business logic for Rombel
type RombelService interface {
	Create(req *dtos.RombelCreateRequest, userID uint) (*dtos.RombelResponse, error)
	GetByID(id uint) (*dtos.RombelResponse, error)
	GetAll(limit int, offset int) (*dtos.RombelListResponse, error)
	GetAllWithFilter(params repositories.GetRombelParams) (*dtos.RombelListWithPaginationResponse, error)
	Update(req *dtos.RombelUpdateRequest, userID uint) (*dtos.RombelResponse, error)
	Delete(id uint) error
}

type RombelServiceImpl struct {
	repository      repositories.RombelRepository
	kelasRepository repositories.KelasRepository
}

// NewRombelService creates a new Rombel service
func NewRombelService(repository repositories.RombelRepository, kelasRepository repositories.KelasRepository) RombelService {
	return &RombelServiceImpl{repository: repository, kelasRepository: kelasRepository}
}

// Create creates a new Rombel
func (s *RombelServiceImpl) Create(req *dtos.RombelCreateRequest, userID uint) (*dtos.RombelResponse, error) {
	// Check if name already exists
	_, err := s.repository.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("rombel dengan nama ini sudah ada")
	}

	// Validate kelas exists and is active
	kelas, err := s.kelasRepository.GetByID(req.KelasID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}
	if kelas.Status != "active" {
		return nil, errors.New("kelas sudah tidak aktif")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.Rombel{
		Name:        req.Name,
		Status:      status,
		KelasID:     req.KelasID,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Rombel by ID
func (s *RombelServiceImpl) GetByID(id uint) (*dtos.RombelResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Rombel
func (s *RombelServiceImpl) GetAll(limit int, offset int) (*dtos.RombelListResponse, error) {
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
	responses := make([]dtos.RombelResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.RombelListResponse{
		Data: responses,
	}, nil
}

// Update updates Rombel
func (s *RombelServiceImpl) Update(req *dtos.RombelUpdateRequest, userID uint) (*dtos.RombelResponse, error) {
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
	if req.KelasID != nil {
		// Validate kelas exists and is active
		kelas, err := s.kelasRepository.GetByID(*req.KelasID)
		if err != nil {
			return nil, errors.New("kelas tidak ditemukan")
		}
		if kelas.Status != "active" {
			return nil, errors.New("kelas sudah tidak aktif")
		}
		existing.KelasID = *req.KelasID
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Rombel by ID
func (s *RombelServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// GetAllWithFilter retrieves Rombel with filters and pagination
func (s *RombelServiceImpl) GetAllWithFilter(params repositories.GetRombelParams) (*dtos.RombelListWithPaginationResponse, error) {
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
	responses := make([]dtos.RombelResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.RombelListWithPaginationResponse{
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
func (s *RombelServiceImpl) mapToResponse(data *models.Rombel) *dtos.RombelResponse {
	var kelasDetail dtos.KelasDetail
	if data.Kelas != nil {
		kelasDetail = dtos.KelasDetail{
			ID:     data.Kelas.ID,
			Name:   data.Kelas.Name,
			Status: data.Kelas.Status,
		}
	}

	return &dtos.RombelResponse{
		ID:          data.ID,
		Name:        data.Name,
		Status:      data.Status,
		KelasID:     data.KelasID,
		Kelas:       kelasDetail,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
