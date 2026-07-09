package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// KonfigurasiMutasiSiswaService handles business logic for Konfigurasi Mutasi Siswa
type KonfigurasiMutasiSiswaService interface {
	UpsertSetting(req *dtos.KonfigurasiMutasiSiswaRequest, file *multipart.FileHeader, userID uint) (*dtos.KonfigurasiMutasiSiswaResponse, error)
	GetSetting() (*dtos.KonfigurasiMutasiSiswaResponse, error)
}

type KonfigurasiMutasiSiswaServiceImpl struct {
	repository repositories.KonfigurasiMutasiSiswaRepository
	r2Storage  *utils.R2Storage
}

// NewKonfigurasiMutasiSiswaService creates a new Konfigurasi Mutasi Siswa service
func NewKonfigurasiMutasiSiswaService(repository repositories.KonfigurasiMutasiSiswaRepository, r2Storage *utils.R2Storage) KonfigurasiMutasiSiswaService {
	return &KonfigurasiMutasiSiswaServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// UpsertSetting creates or updates Konfigurasi Mutasi Siswa with ID = 1
func (s *KonfigurasiMutasiSiswaServiceImpl) UpsertSetting(req *dtos.KonfigurasiMutasiSiswaRequest, file *multipart.FileHeader, userID uint) (*dtos.KonfigurasiMutasiSiswaResponse, error) {
	// Parse tanggal
	tanggalBuka, err := time.Parse("2006-01-02", req.TanggalBukaPendaftaran)
	if err != nil {
		return nil, fmt.Errorf("format tanggal buka pendaftaran tidak valid, gunakan YYYY-MM-DD")
	}

	tanggalTutup, err := time.Parse("2006-01-02", req.TanggalTutupPendaftaran)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tutup pendaftaran tidak valid, gunakan YYYY-MM-DD")
	}

	// Check if record with ID = 1 exists
	existing, err := s.repository.GetByID(1)

	// Handle grup WA
	var grupWA *string
	if req.GrupWA != "" {
		grupWA = &req.GrupWA
	}

	var templateSPTJMPath *string
	var oldTemplatePath *string

	if err != nil {
		// Record not found, create new one with ID = 1
		
		// Upload template SPTJM if provided
		if file != nil {
			path, err := s.r2Storage.UploadFile(file, "mutasi-siswa/template-sptjm")
			if err != nil {
				return nil, fmt.Errorf("gagal upload template SPTJM: %w", err)
			}
			templateSPTJMPath = &path
		}

		data := &models.KonfigurasiMutasiSiswa{
			ID:                      1,
			TanggalBukaPendaftaran:  tanggalBuka,
			TanggalTutupPendaftaran: tanggalTutup,
			NamaKepalaSekolah:       req.NamaKepalaSekolah,
			NIPKepalaSekolah:        req.NIPKepalaSekolah,
			NamaKetuaPanitia:        req.NamaKetuaPanitia,
			NIPKetuaPanitia:         req.NIPKetuaPanitia,
			TemplateSPTJM:           templateSPTJMPath,
			GrupWA:                  grupWA,
			CreatedByID:             &userID,
			UpdatedByID:             &userID,
		}

		if err := s.repository.Create(data); err != nil {
			// If database save fails, delete uploaded file
			if templateSPTJMPath != nil && *templateSPTJMPath != "" {
				_ = s.r2Storage.DeleteFile(*templateSPTJMPath)
			}
			return nil, err
		}

		return s.mapToResponse(data), nil
	}

	// Record exists, update it
	oldTemplatePath = existing.TemplateSPTJM

	existing.TanggalBukaPendaftaran = tanggalBuka
	existing.TanggalTutupPendaftaran = tanggalTutup
	existing.NamaKepalaSekolah = req.NamaKepalaSekolah
	existing.NIPKepalaSekolah = req.NIPKepalaSekolah
	existing.NamaKetuaPanitia = req.NamaKetuaPanitia
	existing.NIPKetuaPanitia = req.NIPKetuaPanitia
	existing.GrupWA = grupWA
	existing.UpdatedByID = &userID

	// Update template SPTJM if provided
	if file != nil {
		// Upload new template
		path, err := s.r2Storage.UploadFile(file, "mutasi-siswa/template-sptjm")
		if err != nil {
			return nil, fmt.Errorf("gagal upload template SPTJM: %w", err)
		}

		// Delete old template if exists
		if oldTemplatePath != nil && *oldTemplatePath != "" {
			_ = s.r2Storage.DeleteFile(*oldTemplatePath)
		}

		existing.TemplateSPTJM = &path
	}

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// GetSetting retrieves konfigurasi mutasi siswa (auth required)
func (s *KonfigurasiMutasiSiswaServiceImpl) GetSetting() (*dtos.KonfigurasiMutasiSiswaResponse, error) {
	// Try to get record with ID = 1
	data, err := s.repository.GetByID(1)

	if err != nil {
		// Record not found, return nil
		return nil, nil
	}

	// Return full setting data
	return s.mapToResponse(data), nil
}

// mapToResponse maps KonfigurasiMutasiSiswa model to response DTO
func (s *KonfigurasiMutasiSiswaServiceImpl) mapToResponse(data *models.KonfigurasiMutasiSiswa) *dtos.KonfigurasiMutasiSiswaResponse {
	// Convert template SPTJM path to public URL
	var templateURL *string
	if data.TemplateSPTJM != nil && *data.TemplateSPTJM != "" {
		url := s.r2Storage.GetPublicURL(*data.TemplateSPTJM)
		templateURL = &url
	}

	return &dtos.KonfigurasiMutasiSiswaResponse{
		ID:                      data.ID,
		TanggalBukaPendaftaran:  data.TanggalBukaPendaftaran.Format("2006-01-02"),
		TanggalTutupPendaftaran: data.TanggalTutupPendaftaran.Format("2006-01-02"),
		NamaKepalaSekolah:       data.NamaKepalaSekolah,
		NIPKepalaSekolah:        data.NIPKepalaSekolah,
		NamaKetuaPanitia:        data.NamaKetuaPanitia,
		NIPKetuaPanitia:         data.NIPKetuaPanitia,
		TemplateSPTJM:           templateURL,
		GrupWA:                  data.GrupWA,
		CreatedAt:               data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:               data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
