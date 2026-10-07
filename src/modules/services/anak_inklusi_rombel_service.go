package services

import (
	"fmt"
	"math"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type AnakInklusiRombelService interface {
	SyncByTahunPelajaran(tahunPelajaranID uint, userID uint) (*dtos.AnakInklusiRombelSyncResponse, error)
	GetAllWithFilter(params repositories.GetAnakInklusiRombelParams) (*dtos.AnakInklusiRombelListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.AnakInklusiRombelResponse, error)
	Update(id uint, req *dtos.AnakInklusiRombelUpdateRequest, userID uint) (*dtos.AnakInklusiRombelResponse, error)
	Delete(id uint) error
}

type anakInklusiRombelService struct {
	anakInklusiRombelRepo repositories.AnakInklusiRombelRepository
	anakInklusiRepo       repositories.AnakInklusiRepository
	r2Storage             *utils.R2Storage
}

func NewAnakInklusiRombelService(
	anakInklusiRombelRepo repositories.AnakInklusiRombelRepository,
	anakInklusiRepo repositories.AnakInklusiRepository,
	r2Storage *utils.R2Storage,
) AnakInklusiRombelService {
	return &anakInklusiRombelService{
		anakInklusiRombelRepo: anakInklusiRombelRepo,
		anakInklusiRepo:       anakInklusiRepo,
		r2Storage:             r2Storage,
	}
}

func (s *anakInklusiRombelService) SyncByTahunPelajaran(tahunPelajaranID uint, userID uint) (*dtos.AnakInklusiRombelSyncResponse, error) {
	// Get all peserta_didik_rombel for the specified tahun_pelajaran
	pesertaDidikRombelList, err := s.anakInklusiRombelRepo.GetPesertaDidikRombelByTahunPelajaran(tahunPelajaranID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data peserta didik rombel: %w", err)
	}

	var details []dtos.AnakInklusiRombelSyncDetail
	totalSynced := 0
	totalSkipped := 0

	// Loop through each peserta_didik_rombel
	for _, pdr := range pesertaDidikRombelList {
		// Check if this peserta_didik exists in anak_inklusi
		anakInklusi, err := s.anakInklusiRepo.GetByPesertaDidikID(pdr.PesertaDidikID)
		
		if err != nil || anakInklusi == nil {
			// Skip if not found in anak_inklusi (not an inklusi student)
			continue
		}

		// Check if already exists in anak_inklusi_rombel
		exists, err := s.anakInklusiRombelRepo.CheckExists(anakInklusi.ID, pdr.ID)
		if err != nil {
			details = append(details, dtos.AnakInklusiRombelSyncDetail{
				PesertaDidikID:       pdr.PesertaDidikID,
				PesertaDidikNama:     getPesertaDidikNama(pdr),
				PesertaDidikRombelID: pdr.ID,
				AnakInklusiID:        anakInklusi.ID,
				Status:               "skipped",
				Message:              fmt.Sprintf("Error checking existence: %v", err),
			})
			totalSkipped++
			continue
		}

		if exists {
			// Skip if already exists
			details = append(details, dtos.AnakInklusiRombelSyncDetail{
				PesertaDidikID:       pdr.PesertaDidikID,
				PesertaDidikNama:     getPesertaDidikNama(pdr),
				PesertaDidikRombelID: pdr.ID,
				AnakInklusiID:        anakInklusi.ID,
				Status:               "skipped",
				Message:              "Data sudah ada",
			})
			totalSkipped++
			continue
		}

		// Create new anak_inklusi_rombel entry
		newEntry := &models.AnakInklusiRombel{
			AnakInklusiID:        anakInklusi.ID,
			PesertaDidikRombelID: pdr.ID,
			CreatedByID:          &userID,
			UpdatedByID:          &userID,
		}

		if err := s.anakInklusiRombelRepo.Create(newEntry); err != nil {
			details = append(details, dtos.AnakInklusiRombelSyncDetail{
				PesertaDidikID:       pdr.PesertaDidikID,
				PesertaDidikNama:     getPesertaDidikNama(pdr),
				PesertaDidikRombelID: pdr.ID,
				AnakInklusiID:        anakInklusi.ID,
				Status:               "skipped",
				Message:              fmt.Sprintf("Error saat menyimpan: %v", err),
			})
			totalSkipped++
			continue
		}

		// Success
		details = append(details, dtos.AnakInklusiRombelSyncDetail{
			PesertaDidikID:       pdr.PesertaDidikID,
			PesertaDidikNama:     getPesertaDidikNama(pdr),
			PesertaDidikRombelID: pdr.ID,
			AnakInklusiID:        anakInklusi.ID,
			Status:               "synced",
			Message:              "Berhasil disinkronkan",
		})
		totalSynced++
	}

	return &dtos.AnakInklusiRombelSyncResponse{
		TotalProcessed: len(details),
		TotalSynced:    totalSynced,
		TotalSkipped:   totalSkipped,
		Details:        details,
	}, nil
}

func getPesertaDidikNama(pdr models.PesertaDidikRombel) string {
	if pdr.PesertaDidik != nil {
		return pdr.PesertaDidik.Nama
	}
	return ""
}

func (s *anakInklusiRombelService) GetAllWithFilter(params repositories.GetAnakInklusiRombelParams) (*dtos.AnakInklusiRombelListWithPaginationResponse, error) {
	anakInklusiRombelList, total, err := s.anakInklusiRombelRepo.GetAllWithFilter(params)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data anak inklusi rombel: %w", err)
	}

	// Convert to response DTOs
	var responses []dtos.AnakInklusiRombelResponse
	for _, anakInklusiRombel := range anakInklusiRombelList {
		responses = append(responses, *s.toResponseDTO(&anakInklusiRombel))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	page := (params.Offset / params.Limit) + 1

	return &dtos.AnakInklusiRombelListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationResponse{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       page,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *anakInklusiRombelService) GetByID(id uint) (*dtos.AnakInklusiRombelResponse, error) {
	anakInklusiRombel, err := s.anakInklusiRombelRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("data anak inklusi rombel tidak ditemukan")
	}

	return s.toResponseDTO(anakInklusiRombel), nil
}

func (s *anakInklusiRombelService) Update(id uint, req *dtos.AnakInklusiRombelUpdateRequest, userID uint) (*dtos.AnakInklusiRombelResponse, error) {
	// Get existing data
	existing, err := s.anakInklusiRombelRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("data anak inklusi rombel tidak ditemukan")
	}

	// Update only allowed fields
	existing.GuruPendampingKhususID = req.GuruPendampingKhususID
	existing.GuruKelasID = req.GuruKelasID
	existing.Catatan = req.Catatan
	existing.UpdatedByID = &userID

	// Save to database
	if err := s.anakInklusiRombelRepo.Update(existing); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data: %w", err)
	}

	// Get updated data with relations
	updated, err := s.anakInklusiRombelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	return s.toResponseDTO(updated), nil
}

func (s *anakInklusiRombelService) Delete(id uint) error {
	// Check if data exists
	existing, err := s.anakInklusiRombelRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("data anak inklusi rombel tidak ditemukan")
	}

	// Delete from database (soft delete using GORM)
	if err := s.anakInklusiRombelRepo.Delete(existing.ID); err != nil {
		return fmt.Errorf("gagal menghapus data: %w", err)
	}

	return nil
}

func (s *anakInklusiRombelService) toResponseDTO(anakInklusiRombel *models.AnakInklusiRombel) *dtos.AnakInklusiRombelResponse {
	resp := &dtos.AnakInklusiRombelResponse{
		ID:                     anakInklusiRombel.ID,
		AnakInklusiID:          anakInklusiRombel.AnakInklusiID,
		PesertaDidikRombelID:   anakInklusiRombel.PesertaDidikRombelID,
		GuruPendampingKhususID: anakInklusiRombel.GuruPendampingKhususID,
		GuruKelasID:            anakInklusiRombel.GuruKelasID,
		Catatan:                anakInklusiRombel.Catatan,
		CreatedAt:              anakInklusiRombel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:              anakInklusiRombel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Include GuruPendampingKhusus data if preloaded
	if anakInklusiRombel.GuruPendampingKhusus != nil {
		resp.GuruPendampingKhusus = &dtos.KepegawaianSimple{
			ID:       anakInklusiRombel.GuruPendampingKhusus.ID,
			Nama:     anakInklusiRombel.GuruPendampingKhusus.Nama,
			NIP:      anakInklusiRombel.GuruPendampingKhusus.NIP,
			Foto:     anakInklusiRombel.GuruPendampingKhusus.Foto,
			FotoURL:  s.r2Storage.GetPublicURL(anakInklusiRombel.GuruPendampingKhusus.Foto),
			Kategori: anakInklusiRombel.GuruPendampingKhusus.Kategori,
			Jabatan:  anakInklusiRombel.GuruPendampingKhusus.Jabatan,
		}
	}

	// Include GuruKelas data if preloaded
	if anakInklusiRombel.GuruKelas != nil {
		resp.GuruKelas = &dtos.KepegawaianSimple{
			ID:       anakInklusiRombel.GuruKelas.ID,
			Nama:     anakInklusiRombel.GuruKelas.Nama,
			NIP:      anakInklusiRombel.GuruKelas.NIP,
			Foto:     anakInklusiRombel.GuruKelas.Foto,
			FotoURL:  s.r2Storage.GetPublicURL(anakInklusiRombel.GuruKelas.Foto),
			Kategori: anakInklusiRombel.GuruKelas.Kategori,
			Jabatan:  anakInklusiRombel.GuruKelas.Jabatan,
		}
	}

	// Include PesertaDidikRombel data if preloaded
	if anakInklusiRombel.PesertaDidikRombel != nil {
		resp.PesertaDidikRombel = &dtos.PesertaDidikRombelSimple{
			ID:               anakInklusiRombel.PesertaDidikRombel.ID,
			PesertaDidikID:   anakInklusiRombel.PesertaDidikRombel.PesertaDidikID,
			RombelID:         anakInklusiRombel.PesertaDidikRombel.RombelID,
			TahunPelajaranID: anakInklusiRombel.PesertaDidikRombel.TahunPelajaranID,
		}

		// Include Rombel data if preloaded
		if anakInklusiRombel.PesertaDidikRombel.Rombel != nil {
			resp.PesertaDidikRombel.Rombel = &dtos.RombelSimple{
				ID:     anakInklusiRombel.PesertaDidikRombel.Rombel.ID,
				Name:   anakInklusiRombel.PesertaDidikRombel.Rombel.Name,
				Status: anakInklusiRombel.PesertaDidikRombel.Rombel.Status,
			}
		}

		// Include TahunPelajaran data if preloaded
		if anakInklusiRombel.PesertaDidikRombel.TahunPelajaran != nil {
			resp.PesertaDidikRombel.TahunPelajaran = &dtos.TahunPelajaranSimple{
				ID:             anakInklusiRombel.PesertaDidikRombel.TahunPelajaran.ID,
				TahunPelajaran: anakInklusiRombel.PesertaDidikRombel.TahunPelajaran.TahunPelajaran,
				Status:         anakInklusiRombel.PesertaDidikRombel.TahunPelajaran.Status,
			}
		}
	}

	// Include AnakInklusi data if preloaded
	if anakInklusiRombel.AnakInklusi != nil {
		anakInklusiResp := &dtos.AnakInklusiResponse{
			ID:                  anakInklusiRombel.AnakInklusi.ID,
			PesertaDidikID:      anakInklusiRombel.AnakInklusi.PesertaDidikID,
			JenisHambatan:       anakInklusiRombel.AnakInklusi.JenisHambatan,
			FileSuratDokter:     anakInklusiRombel.AnakInklusi.FileSuratDokter,
			FileSuratDokterURL:  s.r2Storage.GetPublicURL(anakInklusiRombel.AnakInklusi.FileSuratDokter),
			Status:              anakInklusiRombel.AnakInklusi.Status,
			Catatan:             anakInklusiRombel.AnakInklusi.Catatan,
			CreatedAt:           anakInklusiRombel.AnakInklusi.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:           anakInklusiRombel.AnakInklusi.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Handle nullable date fields
		if anakInklusiRombel.AnakInklusi.TanggalIdentifikasi != nil {
			formatted := anakInklusiRombel.AnakInklusi.TanggalIdentifikasi.Format("2006-01-02")
			anakInklusiResp.TanggalIdentifikasi = &formatted
		}
		if anakInklusiRombel.AnakInklusi.TanggalDiagnosa != nil {
			formatted := anakInklusiRombel.AnakInklusi.TanggalDiagnosa.Format("2006-01-02")
			anakInklusiResp.TanggalDiagnosa = &formatted
		}
		if anakInklusiRombel.AnakInklusi.TanggalKadaluarsaSurat != nil {
			formatted := anakInklusiRombel.AnakInklusi.TanggalKadaluarsaSurat.Format("2006-01-02")
			anakInklusiResp.TanggalKadaluarsaSurat = &formatted
		}

		// Include PesertaDidik data if preloaded
		if anakInklusiRombel.AnakInklusi.PesertaDidik != nil {
			anakInklusiResp.PesertaDidik = &dtos.PesertaDidikSimple{
				ID:       anakInklusiRombel.AnakInklusi.PesertaDidik.ID,
				NIS:      anakInklusiRombel.AnakInklusi.PesertaDidik.NIS,
				NISN:     anakInklusiRombel.AnakInklusi.PesertaDidik.NISN,
				Nama:     anakInklusiRombel.AnakInklusi.PesertaDidik.Nama,
				Photo:    anakInklusiRombel.AnakInklusi.PesertaDidik.Photo,
				PhotoURL: s.r2Storage.GetPublicURL(anakInklusiRombel.AnakInklusi.PesertaDidik.Photo),
				NamaAyah: anakInklusiRombel.AnakInklusi.PesertaDidik.NamaAyah,
				NamaIbu:  anakInklusiRombel.AnakInklusi.PesertaDidik.NamaIbu,
				Status:   anakInklusiRombel.AnakInklusi.PesertaDidik.Status,
			}
		}

		resp.AnakInklusi = anakInklusiResp
	}

	return resp
}
