package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// StrukturOrganisasiService handles business logic for StrukturOrganisasi
type StrukturOrganisasiService interface {
	Create(req *dtos.StrukturOrganisasiCreateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error)
	GetByID(id uint) (*dtos.StrukturOrganisasiResponse, error)
	GetAll(limit int, offset int) (*dtos.StrukturOrganisasiListResponse, error)
	GetAllWithFilter(params repositories.GetStrukturOrganisasiParams) (*dtos.StrukturOrganisasiListWithPaginationResponse, error)
	Update(req *dtos.StrukturOrganisasiUpdateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error)
	Delete(id uint) error
}

type StrukturOrganisasiServiceImpl struct {
	repository repositories.StrukturOrganisasiRepository
	pegawaiRepo repositories.KepegawaianRepository
}

// NewStrukturOrganisasiService creates a new StrukturOrganisasi service
func NewStrukturOrganisasiService(repository repositories.StrukturOrganisasiRepository, pegawaiRepo repositories.KepegawaianRepository) StrukturOrganisasiService {
	return &StrukturOrganisasiServiceImpl{
		repository: repository,
		pegawaiRepo: pegawaiRepo,
	}
}

// Create creates a new StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) Create(req *dtos.StrukturOrganisasiCreateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error) {
	// Validate that either pegawai_id or nama_non_pegawai is provided
	if req.PegawaiID == nil && req.NamaNonPegawai == "" {
		return nil, errors.New("either pegawai_id or nama_non_pegawai must be provided")
	}

	// If pegawai_id is provided, validate it exists
	if req.PegawaiID != nil {
		_, err := s.pegawaiRepo.GetByID(*req.PegawaiID)
		if err != nil {
			return nil, errors.New("pegawai not found")
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.StrukturOrganisasi{
		PegawaiID:         req.PegawaiID,
		NamaNonPegawai:    req.NamaNonPegawai,
		JabatanNonPegawai: req.JabatanNonPegawai,
		Urutan:            req.Urutan,
		Relasi:            req.Relasi,
		Status:            status,
		CreatedByID:       &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves StrukturOrganisasi by ID
func (s *StrukturOrganisasiServiceImpl) GetByID(id uint) (*dtos.StrukturOrganisasiResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("struktur organisasi not found")
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) GetAll(limit int, offset int) (*dtos.StrukturOrganisasiListResponse, error) {
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
	responses := make([]dtos.StrukturOrganisasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.StrukturOrganisasiListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves StrukturOrganisasi with filters and pagination
func (s *StrukturOrganisasiServiceImpl) GetAllWithFilter(params repositories.GetStrukturOrganisasiParams) (*dtos.StrukturOrganisasiListWithPaginationResponse, error) {
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
	responses := make([]dtos.StrukturOrganisasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.StrukturOrganisasiListWithPaginationResponse{
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

// Update updates StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) Update(req *dtos.StrukturOrganisasiUpdateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("struktur organisasi not found")
	}

	// Validate if pegawai_id is being updated (explicitly set in request)
	if req.PegawaiIDSet {
		if req.PegawaiID != nil {
			_, err := s.pegawaiRepo.GetByID(*req.PegawaiID)
			if err != nil {
				return nil, errors.New("pegawai not found")
			}
		}
		existing.PegawaiID = req.PegawaiID
	}

	// Update nama_non_pegawai if provided
	if req.NamaNonPegawai != nil {
		existing.NamaNonPegawai = *req.NamaNonPegawai
	}

	// Update jabatan_non_pegawai if provided
	if req.JabatanNonPegawai != nil {
		existing.JabatanNonPegawai = *req.JabatanNonPegawai
	}

	// Update urutan if provided
	if req.Urutan != nil {
		existing.Urutan = *req.Urutan
	}

	// Update relasi if provided
	if req.Relasi != nil {
		existing.Relasi = *req.Relasi
	}

	// Update status if provided
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes StrukturOrganisasi by ID
func (s *StrukturOrganisasiServiceImpl) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("struktur organisasi not found")
	}
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *StrukturOrganisasiServiceImpl) mapToResponse(data *models.StrukturOrganisasi) *dtos.StrukturOrganisasiResponse {
	response := &dtos.StrukturOrganisasiResponse{
		ID:                data.ID,
		PegawaiID:         data.PegawaiID,
		NamaNonPegawai:    data.NamaNonPegawai,
		JabatanNonPegawai: data.JabatanNonPegawai,
		Urutan:            data.Urutan,
		Relasi:            data.Relasi,
		Status:            data.Status,
		CreatedAt:         data.CreatedAt,
		UpdatedAt:         data.UpdatedAt,
		CreatedByID:       data.CreatedByID,
		UpdatedByID:       data.UpdatedByID,
	}

	// Add Pegawai data if available
	if data.Pegawai != nil {
		response.Pegawai = &dtos.PegawaiSimpleResponse{
			ID:      data.Pegawai.ID,
			Nama:    data.Pegawai.Nama,
			Jabatan: data.Pegawai.Jabatan,
			Status:  data.Pegawai.Status,
		}
	}

	return response
}
