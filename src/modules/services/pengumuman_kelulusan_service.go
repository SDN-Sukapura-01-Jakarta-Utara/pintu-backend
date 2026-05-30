package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type PengumumanKelulusanService interface {
	ConfigurePengumuman(req *dtos.PengumumanKelulusanConfigRequest, fotoKepsek *multipart.FileHeader, ttdKepsek *multipart.FileHeader, userID uint) (*dtos.PengumumanKelulusanResponse, error)
	GetPengumuman() (*dtos.PengumumanKelulusanResponse, error)
	GetSettingPengumumanPublic() (*dtos.PengumumanKelulusanResponse, error)
}

type PengumumanKelulusanServiceImpl struct {
	repository repositories.PengumumanKelulusanRepository
	r2Storage  *utils.R2Storage
}

// NewPengumumanKelulusanService creates a new PengumumanKelulusan service
func NewPengumumanKelulusanService(repository repositories.PengumumanKelulusanRepository) PengumumanKelulusanService {
	return &PengumumanKelulusanServiceImpl{
		repository: repository,
		r2Storage:  utils.NewR2Storage(),
	}
}

// ConfigurePengumuman creates or updates pengumuman kelulusan configuration
func (s *PengumumanKelulusanServiceImpl) ConfigurePengumuman(req *dtos.PengumumanKelulusanConfigRequest, fotoKepsek *multipart.FileHeader, ttdKepsek *multipart.FileHeader, userID uint) (*dtos.PengumumanKelulusanResponse, error) {
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

		oldFotoKepsek := existing.FotoKepsek
		oldTtdKepsek := existing.TtdKepsek

		existing.SambutanKelulusan = req.SambutanKelulusan
		existing.TanggalPengumumanNilai = tanggalNilai
		existing.TanggalPengumumanKelulusan = tanggalKelulusan
		existing.NamaKepsek = req.NamaKepsek
		existing.UpdatedByID = &userID

		// Handle foto_kepsek deletion if requested
		if req.DeleteFotoKepsek {
			if oldFotoKepsek != "" {
				_ = s.r2Storage.DeleteFile(oldFotoKepsek)
			}
			existing.FotoKepsek = ""
		}

		// Handle foto_kepsek upload if provided
		if fotoKepsek != nil {
			uploadedPath, err := s.r2Storage.UploadFile(fotoKepsek, "pengumuman-kelulusan")
			if err != nil {
				return nil, fmt.Errorf("gagal upload foto kepsek: %s", err.Error())
			}

			// Delete old file if exists
			if oldFotoKepsek != "" && oldFotoKepsek != uploadedPath {
				_ = s.r2Storage.DeleteFile(oldFotoKepsek)
			}

			existing.FotoKepsek = uploadedPath
		}

		// Handle ttd_kepsek deletion if requested
		if req.DeleteTtdKepsek {
			if oldTtdKepsek != "" {
				_ = s.r2Storage.DeleteFile(oldTtdKepsek)
			}
			existing.TtdKepsek = ""
		}

		// Handle ttd_kepsek upload if provided
		if ttdKepsek != nil {
			uploadedPath, err := s.r2Storage.UploadFile(ttdKepsek, "pengumuman-kelulusan")
			if err != nil {
				return nil, fmt.Errorf("gagal upload ttd kepsek: %s", err.Error())
			}

			// Delete old file if exists
			if oldTtdKepsek != "" && oldTtdKepsek != uploadedPath {
				_ = s.r2Storage.DeleteFile(oldTtdKepsek)
			}

			existing.TtdKepsek = uploadedPath
		}

		if err := s.repository.Update(existing); err != nil {
			return nil, errors.New("gagal mengupdate konfigurasi pengumuman")
		}

		return s.mapToResponse(existing), nil
	}

	// Create new record
	var fotoKepsekPath, ttdKepsekPath string

	// Handle foto_kepsek upload if provided
	if fotoKepsek != nil {
		uploadedPath, err := s.r2Storage.UploadFile(fotoKepsek, "pengumuman-kelulusan")
		if err != nil {
			return nil, fmt.Errorf("gagal upload foto kepsek: %s", err.Error())
		}
		fotoKepsekPath = uploadedPath
	}

	// Handle ttd_kepsek upload if provided
	if ttdKepsek != nil {
		uploadedPath, err := s.r2Storage.UploadFile(ttdKepsek, "pengumuman-kelulusan")
		if err != nil {
			// Delete foto_kepsek if ttd_kepsek upload fails
			if fotoKepsekPath != "" {
				s.r2Storage.DeleteFile(fotoKepsekPath)
			}
			return nil, fmt.Errorf("gagal upload ttd kepsek: %s", err.Error())
		}
		ttdKepsekPath = uploadedPath
	}

	pengumuman := &models.PengumumanKelulusan{
		SambutanKelulusan:          req.SambutanKelulusan,
		TanggalPengumumanNilai:     tanggalNilai,
		TanggalPengumumanKelulusan: tanggalKelulusan,
		FotoKepsek:                 fotoKepsekPath,
		TtdKepsek:                  ttdKepsekPath,
		NamaKepsek:                 req.NamaKepsek,
		CreatedByID:                &userID,
		UpdatedByID:                &userID,
	}

	if err := s.repository.Create(pengumuman); err != nil {
		// Delete uploaded files if save to DB failed
		if fotoKepsekPath != "" {
			s.r2Storage.DeleteFile(fotoKepsekPath)
		}
		if ttdKepsekPath != "" {
			s.r2Storage.DeleteFile(ttdKepsekPath)
		}
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

// GetSettingPengumumanPublic retrieves the pengumuman kelulusan record with ID 1 (public API)
func (s *PengumumanKelulusanServiceImpl) GetSettingPengumumanPublic() (*dtos.PengumumanKelulusanResponse, error) {
	data, err := s.repository.GetByID(1)
	if err != nil {
		return nil, errors.New("data pengumuman tidak ditemukan")
	}

	return s.mapToResponse(data), nil
}

// mapToResponse maps PengumumanKelulusan model to PengumumanKelulusanResponse DTO
func (s *PengumumanKelulusanServiceImpl) mapToResponse(data *models.PengumumanKelulusan) *dtos.PengumumanKelulusanResponse {
	// Generate full URLs for files
	fotoKepsekURL := s.r2Storage.GetPublicURL(data.FotoKepsek)
	ttdKepsekURL := s.r2Storage.GetPublicURL(data.TtdKepsek)

	response := &dtos.PengumumanKelulusanResponse{
		ID:                         data.ID,
		SambutanKelulusan:          data.SambutanKelulusan,
		TanggalPengumumanNilai:     data.TanggalPengumumanNilai.Format("2006-01-02 15:04:05"),
		TanggalPengumumanKelulusan: data.TanggalPengumumanKelulusan.Format("2006-01-02 15:04:05"),
		FotoKepsek:                 fotoKepsekURL,
		TtdKepsek:                  ttdKepsekURL,
		NamaKepsek:                 data.NamaKepsek,
		CreatedAt:                  data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                  data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedByID:                data.CreatedByID,
		UpdatedByID:                data.UpdatedByID,
	}

	return response
}
