package services

import (
	"errors"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

type PengumumanKelulusanService interface {
	ConfigurePengumuman(req *dtos.PengumumanKelulusanConfigRequest, userID uint) (*dtos.PengumumanKelulusanResponse, error)
	GetPengumuman() (*dtos.PengumumanKelulusanResponse, error)
}

type PengumumanKelulusanServiceImpl struct {
	repository repositories.PengumumanKelulusanRepository
}

// NewPengumumanKelulusanService creates a new PengumumanKelulusan service
func NewPengumumanKelulusanService(repository repositories.PengumumanKelulusanRepository) PengumumanKelulusanService {
	return &PengumumanKelulusanServiceImpl{
		repository: repository,
	}
}

// ConfigurePengumuman creates or updates pengumuman kelulusan configuration
func (s *PengumumanKelulusanServiceImpl) ConfigurePengumuman(req *dtos.PengumumanKelulusanConfigRequest, userID uint) (*dtos.PengumumanKelulusanResponse, error) {
	// Parse tanggal_pengumuman_nilai (YYYY-MM-DD HH:MM:SS format)
	tanggalNilai, err := time.Parse("2006-01-02 15:04:05", req.TanggalPengumumanNilai)
	if err != nil {
		return nil, errors.New("format tanggal_pengumuman_nilai tidak valid, gunakan YYYY-MM-DD HH:MM:SS")
	}

	// Parse tanggal_pengumuman_kelulusan (YYYY-MM-DD HH:MM:SS format)
	tanggalKelulusan, err := time.Parse("2006-01-02 15:04:05", req.TanggalPengumumanKelulusan)
	if err != nil {
		return nil, errors.New("format tanggal_pengumuman_kelulusan tidak valid, gunakan YYYY-MM-DD HH:MM:SS")
	}

	// Check if ID is provided (update) or not (create)
	if req.ID != nil && *req.ID > 0 {
		// Update existing record
		existing, err := s.repository.GetByID(*req.ID)
		if err != nil {
			return nil, errors.New("data pengumuman tidak ditemukan")
		}

		existing.SambutanKelulusan = req.SambutanKelulusan
		existing.TanggalPengumumanNilai = tanggalNilai
		existing.TanggalPengumumanKelulusan = tanggalKelulusan
		existing.UpdatedByID = &userID

		if err := s.repository.Update(existing); err != nil {
			return nil, errors.New("gagal mengupdate konfigurasi pengumuman")
		}

		return s.mapToResponse(existing), nil
	}

	// Create new record
	pengumuman := &models.PengumumanKelulusan{
		SambutanKelulusan:          req.SambutanKelulusan,
		TanggalPengumumanNilai:     tanggalNilai,
		TanggalPengumumanKelulusan: tanggalKelulusan,
		CreatedByID:                &userID,
		UpdatedByID:                &userID,
	}

	if err := s.repository.Create(pengumuman); err != nil {
		return nil, errors.New("gagal menyimpan konfigurasi pengumuman")
	}

	return s.mapToResponse(pengumuman), nil
}

// GetPengumuman retrieves the pengumuman kelulusan record with ID 1
func (s *PengumumanKelulusanServiceImpl) GetPengumuman() (*dtos.PengumumanKelulusanResponse, error) {
	data, err := s.repository.GetFirst()
	if err != nil {
		return nil, errors.New("data pengumuman tidak ditemukan")
	}

	return s.mapToResponse(data), nil
}

// mapToResponse maps PengumumanKelulusan model to PengumumanKelulusanResponse DTO
func (s *PengumumanKelulusanServiceImpl) mapToResponse(data *models.PengumumanKelulusan) *dtos.PengumumanKelulusanResponse {
	response := &dtos.PengumumanKelulusanResponse{
		ID:                         data.ID,
		SambutanKelulusan:          data.SambutanKelulusan,
		TanggalPengumumanNilai:     data.TanggalPengumumanNilai.Format("2006-01-02 15:04:05"),
		TanggalPengumumanKelulusan: data.TanggalPengumumanKelulusan.Format("2006-01-02 15:04:05"),
		CreatedAt:                  data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                  data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedByID:                data.CreatedByID,
		UpdatedByID:                data.UpdatedByID,
	}

	return response
}
