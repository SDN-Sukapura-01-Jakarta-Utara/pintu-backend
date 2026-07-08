package services

import (
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// KonfigurasiAbsensiService handles business logic for Konfigurasi Absensi
type KonfigurasiAbsensiService interface {
	UpsertKonfigurasi(req *dtos.KonfigurasiAbsensiRequest) (*dtos.KonfigurasiAbsensiResponse, error)
	GetKonfigurasi() (*dtos.KonfigurasiAbsensiResponse, error)
}

type KonfigurasiAbsensiServiceImpl struct {
	repository repositories.KonfigurasiAbsensiRepository
}

// NewKonfigurasiAbsensiService creates a new Konfigurasi Absensi service
func NewKonfigurasiAbsensiService(repository repositories.KonfigurasiAbsensiRepository) KonfigurasiAbsensiService {
	return &KonfigurasiAbsensiServiceImpl{
		repository: repository,
	}
}

// UpsertKonfigurasi creates or updates Konfigurasi Absensi with ID = 1
func (s *KonfigurasiAbsensiServiceImpl) UpsertKonfigurasi(req *dtos.KonfigurasiAbsensiRequest) (*dtos.KonfigurasiAbsensiResponse, error) {
	// Check if record with ID = 1 exists
	existing, err := s.repository.GetByID(1)

	// Handle nullable fields
	var namaKepsek, nipKepsek *string
	if req.NamaKepsek != "" {
		namaKepsek = &req.NamaKepsek
	}
	if req.NIPKepsek != "" {
		nipKepsek = &req.NIPKepsek
	}

	if err != nil {
		// Record not found, create new one with ID = 1
		data := &models.KonfigurasiAbsensi{
			ID:               1,
			JamDatangMulai:   req.JamDatangMulai,
			JamMaxDatang:     req.JamMaxDatang,
			JamDatangSelesai: req.JamDatangSelesai,
			JamPulangMulai:   req.JamPulangMulai,
			JamPulangSelesai: req.JamPulangSelesai,
			NamaKepsek:       namaKepsek,
			NIPKepsek:        nipKepsek,
		}

		if err := s.repository.Create(data); err != nil {
			return nil, err
		}

		return s.mapToResponse(data), nil
	}

	// Record exists, update it
	existing.JamDatangMulai = req.JamDatangMulai
	existing.JamMaxDatang = req.JamMaxDatang
	existing.JamDatangSelesai = req.JamDatangSelesai
	existing.JamPulangMulai = req.JamPulangMulai
	existing.JamPulangSelesai = req.JamPulangSelesai
	existing.NamaKepsek = namaKepsek
	existing.NIPKepsek = nipKepsek

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// GetKonfigurasi retrieves konfigurasi absensi with ID = 1
func (s *KonfigurasiAbsensiServiceImpl) GetKonfigurasi() (*dtos.KonfigurasiAbsensiResponse, error) {
	// Try to get record with ID = 1
	data, err := s.repository.GetByID(1)
	
	if err != nil {
		// Record not found, return nil
		return nil, nil
	}

	// Return konfigurasi data
	return s.mapToResponse(data), nil
}

// mapToResponse maps KonfigurasiAbsensi model to response DTO
func (s *KonfigurasiAbsensiServiceImpl) mapToResponse(data *models.KonfigurasiAbsensi) *dtos.KonfigurasiAbsensiResponse {
	return &dtos.KonfigurasiAbsensiResponse{
		ID:               data.ID,
		JamDatangMulai:   data.JamDatangMulai,
		JamMaxDatang:     data.JamMaxDatang,
		JamDatangSelesai: data.JamDatangSelesai,
		JamPulangMulai:   data.JamPulangMulai,
		JamPulangSelesai: data.JamPulangSelesai,
		NamaKepsek:       data.NamaKepsek,
		NIPKepsek:        data.NIPKepsek,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
