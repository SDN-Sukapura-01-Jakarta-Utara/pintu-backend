package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// EkstrakurikulerService handles business logic for Ekstrakurikuler
type EkstrakurikulerService interface {
	Create(req *dtos.EkstrakurikulerCreateRequest, userID uint) (*dtos.EkstrakurikulerResponse, error)
	GetByID(id uint) (*dtos.EkstrakurikulerResponse, error)
	GetAll(limit int, offset int) (*dtos.EkstrakurikulerListResponse, error)
	Update(req *dtos.EkstrakurikulerUpdateRequest, userID uint) (*dtos.EkstrakurikulerResponse, error)
	Delete(id uint) error
}

type EkstrakurikulerServiceImpl struct {
	repository      repositories.EkstrakurikulerRepository
	kelasRepository repositories.KelasRepository
}

// NewEkstrakurikulerService creates a new Ekstrakurikuler service
func NewEkstrakurikulerService(repository repositories.EkstrakurikulerRepository, kelasRepository repositories.KelasRepository) EkstrakurikulerService {
	return &EkstrakurikulerServiceImpl{repository: repository, kelasRepository: kelasRepository}
}

// Create creates a new Ekstrakurikuler
func (s *EkstrakurikulerServiceImpl) Create(req *dtos.EkstrakurikulerCreateRequest, userID uint) (*dtos.EkstrakurikulerResponse, error) {
	// Check if name already exists
	_, err := s.repository.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("ekstrakurikuler dengan nama ini sudah ada")
	}

	// Validate all kelas_ids exist and are active
	if len(req.KelasIDs) == 0 {
		return nil, errors.New("minimal harus ada 1 kelas")
	}

	for _, kelasID := range req.KelasIDs {
		kelas, err := s.kelasRepository.GetByID(kelasID)
		if err != nil {
			return nil, errors.New("kelas dengan id " + string(rune(kelasID)) + " tidak ditemukan")
		}
		if kelas.Status != "active" {
			return nil, errors.New("kelas dengan id " + string(rune(kelasID)) + " sudah tidak aktif")
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.Ekstrakurikuler{
		Name:        req.Name,
		KelasIDs:    models.KelasIDs(req.KelasIDs),
		Kategori:    req.Kategori,
		Status:      status,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Ekstrakurikuler by ID
func (s *EkstrakurikulerServiceImpl) GetByID(id uint) (*dtos.EkstrakurikulerResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Ekstrakurikuler
func (s *EkstrakurikulerServiceImpl) GetAll(limit int, offset int) (*dtos.EkstrakurikulerListResponse, error) {
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
	responses := make([]dtos.EkstrakurikulerResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.EkstrakurikulerListResponse{
		Data: responses,
	}, nil
}

// Update updates Ekstrakurikuler
func (s *EkstrakurikulerServiceImpl) Update(req *dtos.EkstrakurikulerUpdateRequest, userID uint) (*dtos.EkstrakurikulerResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Kategori != nil {
		existing.Kategori = *req.Kategori
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if len(req.KelasIDs) > 0 {
		// Validate all kelas_ids exist and are active
		for _, kelasID := range req.KelasIDs {
			kelas, err := s.kelasRepository.GetByID(kelasID)
			if err != nil {
				return nil, errors.New("kelas dengan id " + string(rune(kelasID)) + " tidak ditemukan")
			}
			if kelas.Status != "active" {
				return nil, errors.New("kelas dengan id " + string(rune(kelasID)) + " sudah tidak aktif")
			}
		}
		existing.KelasIDs = models.KelasIDs(req.KelasIDs)
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Ekstrakurikuler by ID
func (s *EkstrakurikulerServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *EkstrakurikulerServiceImpl) mapToResponse(data *models.Ekstrakurikuler) *dtos.EkstrakurikulerResponse {
	kelasIDs := []uint(data.KelasIDs)
	return &dtos.EkstrakurikulerResponse{
		ID:          data.ID,
		Name:        data.Name,
		KelasIDs:    kelasIDs,
		Kategori:    data.Kategori,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
