package services

import (
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// SettingLayananSPMBService handles business logic for Setting Layanan SPMB
type SettingLayananSPMBService interface {
	UpsertSetting(req *dtos.SettingLayananSPMBRequest, userID uint) (*dtos.SettingLayananSPMBResponse, error)
	GetGrupWAPublic() (*dtos.GrupWASPMBResponse, error)
	GetSetting() (*dtos.SettingLayananSPMBResponse, error)
}

type SettingLayananSPMBServiceImpl struct {
	repository repositories.SettingLayananSPMBRepository
}

// NewSettingLayananSPMBService creates a new Setting Layanan SPMB service
func NewSettingLayananSPMBService(repository repositories.SettingLayananSPMBRepository) SettingLayananSPMBService {
	return &SettingLayananSPMBServiceImpl{
		repository: repository,
	}
}

// UpsertSetting creates or updates Setting Layanan SPMB with ID = 1
func (s *SettingLayananSPMBServiceImpl) UpsertSetting(req *dtos.SettingLayananSPMBRequest, userID uint) (*dtos.SettingLayananSPMBResponse, error) {
	// Check if record with ID = 1 exists
	existing, err := s.repository.GetByID(1)

	// Handle nullable fields
	var namaKepsek, nipKepsek, namaKetua, nipKetua, grupWA *string
	if req.NamaKepalaSekolah != "" {
		namaKepsek = &req.NamaKepalaSekolah
	}
	if req.NIPKepalaSekolah != "" {
		nipKepsek = &req.NIPKepalaSekolah
	}
	if req.NamaKetuaPanitia != "" {
		namaKetua = &req.NamaKetuaPanitia
	}
	if req.NIPKetuaPanitia != "" {
		nipKetua = &req.NIPKetuaPanitia
	}
	if req.GrupWA != "" {
		grupWA = &req.GrupWA
	}

	if err != nil {
		// Record not found, create new one with ID = 1
		data := &models.SettingLayananSPMB{
			ID:                1,
			NamaKepalaSekolah: namaKepsek,
			NIPKepalaSekolah:  nipKepsek,
			NamaKetuaPanitia:  namaKetua,
			NIPKetuaPanitia:   nipKetua,
			GrupWA:            grupWA,
			CreatedByID:       &userID,
			UpdatedByID:       &userID,
		}

		if err := s.repository.Create(data); err != nil {
			return nil, err
		}

		return s.mapToResponse(data), nil
	}

	// Record exists, update it
	existing.NamaKepalaSekolah = namaKepsek
	existing.NIPKepalaSekolah = nipKepsek
	existing.NamaKetuaPanitia = namaKetua
	existing.NIPKetuaPanitia = nipKetua
	existing.GrupWA = grupWA
	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// mapToResponse maps SettingLayananSPMB model to response DTO
func (s *SettingLayananSPMBServiceImpl) mapToResponse(data *models.SettingLayananSPMB) *dtos.SettingLayananSPMBResponse {
	return &dtos.SettingLayananSPMBResponse{
		ID:                data.ID,
		NamaKepalaSekolah: data.NamaKepalaSekolah,
		NIPKepalaSekolah:  data.NIPKepalaSekolah,
		NamaKetuaPanitia:  data.NamaKetuaPanitia,
		NIPKetuaPanitia:   data.NIPKetuaPanitia,
		GrupWA:            data.GrupWA,
		CreatedAt:         data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// GetGrupWAPublic retrieves grup WA link for public (no auth)
func (s *SettingLayananSPMBServiceImpl) GetGrupWAPublic() (*dtos.GrupWASPMBResponse, error) {
	// Try to get record with ID = 1
	data, err := s.repository.GetByID(1)
	
	if err != nil {
		// Record not found, return null grup_wa
		return &dtos.GrupWASPMBResponse{
			GrupWA: nil,
		}, nil
	}

	// Return grup_wa from existing record
	return &dtos.GrupWASPMBResponse{
		GrupWA: data.GrupWA,
	}, nil
}

// GetSetting retrieves setting layanan SPMB (auth required)
func (s *SettingLayananSPMBServiceImpl) GetSetting() (*dtos.SettingLayananSPMBResponse, error) {
	// Try to get record with ID = 1
	data, err := s.repository.GetByID(1)
	
	if err != nil {
		// Record not found, return empty/null values
		return &dtos.SettingLayananSPMBResponse{
			ID:                1,
			NamaKepalaSekolah: nil,
			NIPKepalaSekolah:  nil,
			NamaKetuaPanitia:  nil,
			NIPKetuaPanitia:   nil,
			GrupWA:            nil,
			CreatedAt:         "",
			UpdatedAt:         "",
		}, nil
	}

	// Return full setting data
	return s.mapToResponse(data), nil
}
