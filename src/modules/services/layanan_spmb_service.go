package services

import (
	"fmt"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// LayananSPMBService handles business logic for Layanan SPMB
type LayananSPMBService interface {
	CreatePublic(req *dtos.LayananSPMBCreateRequest) (*dtos.LayananSPMBResponse, error)
	GetAllWithFilter(req *dtos.LayananSPMBGetAllRequest) (*dtos.LayananSPMBListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.LayananSPMBResponse, error)
	UpdateStatus(req *dtos.LayananSPMBUpdateStatusRequest) (*dtos.LayananSPMBResponse, error)
	DeleteLayananSPMB(id uint) error
}

type LayananSPMBServiceImpl struct {
	repository repositories.LayananSPMBRepository
}

// NewLayananSPMBService creates a new Layanan SPMB service
func NewLayananSPMBService(repository repositories.LayananSPMBRepository) LayananSPMBService {
	return &LayananSPMBServiceImpl{
		repository: repository,
	}
}

// CreatePublic creates a new Layanan SPMB from public form
func (s *LayananSPMBServiceImpl) CreatePublic(req *dtos.LayananSPMBCreateRequest) (*dtos.LayananSPMBResponse, error) {
	// Create layanan_spmb record
	data := &models.LayananSPMB{
		NamaOrangTua:     req.NamaOrangTua,
		NomorTelepon:     req.NomorTelepon,
		Alamat:           req.Alamat,
		NamaLengkapMurid: req.NamaLengkapMurid,
		Keperluan:        req.Keperluan,
		Status:           "pending",
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// mapToResponse maps LayananSPMB model to response DTO
func (s *LayananSPMBServiceImpl) mapToResponse(data *models.LayananSPMB) *dtos.LayananSPMBResponse {
	return &dtos.LayananSPMBResponse{
		ID:               data.ID,
		NamaOrangTua:     data.NamaOrangTua,
		NomorTelepon:     data.NomorTelepon,
		Alamat:           data.Alamat,
		NamaLengkapMurid: data.NamaLengkapMurid,
		Keperluan:        data.Keperluan,
		TanggalLaporan:   data.TanggalLaporan.Format("2006-01-02 15:04:05"),
		Status:           data.Status,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// GetAllWithFilter retrieves all Layanan SPMB with filters and pagination
func (s *LayananSPMBServiceImpl) GetAllWithFilter(req *dtos.LayananSPMBGetAllRequest) (*dtos.LayananSPMBListWithPaginationResponse, error) {
	// Set default pagination
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Build filter params
	params := repositories.GetLayananSPMBParams{
		Filter: repositories.GetLayananSPMBFilter{
			StartDate:    req.Search.StartDate,
			EndDate:      req.Search.EndDate,
			NamaOrangTua: req.Search.NamaOrangTua,
			NamaMurid:    req.Search.NamaMurid,
			Status:       req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	}

	// Get data from repository
	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	var responses []dtos.LayananSPMBResponse
	for _, item := range data {
		responses = append(responses, *s.mapToResponse(&item))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dtos.LayananSPMBListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationMeta{
			Limit:      limit,
			Offset:     offset,
			Page:       page,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetByID retrieves Layanan SPMB by ID
func (s *LayananSPMBServiceImpl) GetByID(id uint) (*dtos.LayananSPMBResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", id)
	}

	return s.mapToResponse(data), nil
}

// UpdateStatus updates status of Layanan SPMB
func (s *LayananSPMBServiceImpl) UpdateStatus(req *dtos.LayananSPMBUpdateStatusRequest) (*dtos.LayananSPMBResponse, error) {
	// Get existing data
	data, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", req.ID)
	}

	// Validate status value
	if req.Status != "pending" && req.Status != "selesai" {
		return nil, fmt.Errorf("status harus 'pending' atau 'selesai'")
	}

	// Update status
	data.Status = req.Status

	// Save to database
	if err := s.repository.Update(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// DeleteLayananSPMB soft deletes Layanan SPMB by ID
func (s *LayananSPMBServiceImpl) DeleteLayananSPMB(id uint) error {
	// Check if record exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", id)
	}

	// Soft delete
	if err := s.repository.SoftDelete(id); err != nil {
		return err
	}

	return nil
}
