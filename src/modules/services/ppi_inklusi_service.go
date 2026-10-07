package services

import (
	"fmt"
	"math"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type PPIInklusiService interface {
	Create(req *dtos.PPIInklusiCreateRequest, userID uint) (*dtos.PPIInklusiCreateResponse, error)
	GetAllWithFilter(params repositories.GetPPIInklusiParams) (*dtos.PPIInklusiListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.PPIInklusiResponse, error)
	Update(req *dtos.PPIInklusiUpdateRequest, userID uint) (*dtos.PPIInklusiResponse, error)
	Delete(id uint) error
	GetSummary(req *dtos.PPIInklusiSummaryRequest) (*dtos.PPIInklusiSummaryResponse, error)
}

type ppiInklusiService struct {
	repo      repositories.PPIInklusiRepository
	r2Storage *utils.R2Storage
}

func NewPPIInklusiService(repo repositories.PPIInklusiRepository, r2Storage *utils.R2Storage) PPIInklusiService {
	return &ppiInklusiService{
		repo:      repo,
		r2Storage: r2Storage,
	}
}

func (s *ppiInklusiService) Create(req *dtos.PPIInklusiCreateRequest, userID uint) (*dtos.PPIInklusiCreateResponse, error) {
	var createdPPIList []models.PPIInklusi
	var responses []dtos.PPIInklusiResponse

	// Loop through each program in the list
	for _, program := range req.ProgramList {
		// Create model for each program
		ppiInklusi := &models.PPIInklusi{
			AnakInklusiRombelID:  req.AnakInklusiRombelID,
			TahunPelajaranID:     req.TahunPelajaranID,
			AspekPembelajaran:    program.AspekPembelajaran,
			NamaProgram:          program.NamaProgram,
			TujuanPembelajaran:   program.TujuanPembelajaran,
			StrategiPembelajaran: program.StrategiPembelajaran,
			TargetWaktu:          program.TargetWaktu,
			CapaianSaatIni:       "", // Empty on create
			StatusCapaian:        "Belum Dimulai", // Default status
			GuruPembuatID:        req.GuruPembuatID,
			CreatedByID:          &userID,
			UpdatedByID:          &userID,
		}

		// Save to database
		if err := s.repo.Create(ppiInklusi); err != nil {
			return nil, fmt.Errorf("gagal menyimpan data PPI: %w", err)
		}

		createdPPIList = append(createdPPIList, *ppiInklusi)
	}

	// Get all created data with relations
	for _, ppi := range createdPPIList {
		created, err := s.repo.GetByID(ppi.ID)
		if err != nil {
			continue
		}
		responses = append(responses, *s.toResponseDTO(created))
	}

	return &dtos.PPIInklusiCreateResponse{
		TotalCreated: len(responses),
		Data:         responses,
	}, nil
}

func (s *ppiInklusiService) GetAllWithFilter(params repositories.GetPPIInklusiParams) (*dtos.PPIInklusiListWithPaginationResponse, error) {
	ppiInklusiList, total, err := s.repo.GetAllWithFilter(params)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data PPI: %w", err)
	}

	// Convert to response DTOs
	var responses []dtos.PPIInklusiResponse
	for _, ppiInklusi := range ppiInklusiList {
		responses = append(responses, *s.toResponseDTO(&ppiInklusi))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	page := (params.Offset / params.Limit) + 1

	return &dtos.PPIInklusiListWithPaginationResponse{
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

func (s *ppiInklusiService) GetByID(id uint) (*dtos.PPIInklusiResponse, error) {
	ppiInklusi, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("data PPI tidak ditemukan")
	}

	return s.toResponseDTO(ppiInklusi), nil
}

func (s *ppiInklusiService) Update(req *dtos.PPIInklusiUpdateRequest, userID uint) (*dtos.PPIInklusiResponse, error) {
	// Check if PPI exists
	ppiInklusi, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("data PPI tidak ditemukan")
	}

	// Update fields
	ppiInklusi.AspekPembelajaran = req.AspekPembelajaran
	ppiInklusi.NamaProgram = req.NamaProgram
	ppiInklusi.TujuanPembelajaran = req.TujuanPembelajaran
	ppiInklusi.StrategiPembelajaran = req.StrategiPembelajaran
	ppiInklusi.TargetWaktu = req.TargetWaktu
	ppiInklusi.GuruPembuatID = req.GuruPembuatID
	ppiInklusi.StatusCapaian = req.StatusCapaian
	ppiInklusi.UpdatedByID = &userID

	// Save to database
	if err := s.repo.Update(ppiInklusi); err != nil {
		return nil, fmt.Errorf("gagal memperbarui data PPI: %w", err)
	}

	// Get updated data with relations
	updated, err := s.repo.GetByID(ppiInklusi.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data PPI yang diperbarui")
	}

	return s.toResponseDTO(updated), nil
}

func (s *ppiInklusiService) Delete(id uint) error {
	// Check if PPI exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("data PPI tidak ditemukan")
	}

	// Delete from database
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus data PPI: %w", err)
	}

	return nil
}

func (s *ppiInklusiService) GetSummary(req *dtos.PPIInklusiSummaryRequest) (*dtos.PPIInklusiSummaryResponse, error) {
	// Prepare filter strings
	tahunPelajaranIDStr := ""
	if req.TahunPelajaranID > 0 {
		tahunPelajaranIDStr = fmt.Sprint(req.TahunPelajaranID)
	}

	rombelIDStr := ""
	if req.Search.RombelID > 0 {
		rombelIDStr = fmt.Sprint(req.Search.RombelID)
	}

	guruPembuatIDStr := ""
	if req.Search.GuruPembuatID > 0 {
		guruPembuatIDStr = fmt.Sprint(req.Search.GuruPembuatID)
	}

	anakInklusiRombelIDStr := ""
	if req.Search.AnakInklusiRombelID > 0 {
		anakInklusiRombelIDStr = fmt.Sprint(req.Search.AnakInklusiRombelID)
	}

	aspekPembelajaran := req.Search.AspekPembelajaran

	// Get summary data from repository with filters
	summaryData, err := s.repo.GetSummaryByTahunPelajaran(repositories.GetPPIInklusiFilter{
		TahunPelajaranID:    tahunPelajaranIDStr,
		RombelID:            rombelIDStr,
		GuruPembuatID:       guruPembuatIDStr,
		AnakInklusiRombelID: anakInklusiRombelIDStr,
		AspekPembelajaran:   aspekPembelajaran,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil summary PPI: %w", err)
	}

	// Get tahun pelajaran info
	tahunPelajaran, err := s.repo.GetTahunPelajaranByID(req.TahunPelajaranID)
	if err != nil {
		return nil, fmt.Errorf("tahun pelajaran tidak ditemukan")
	}

	// Build status capaian breakdown with percentage
	var statusBreakdown []dtos.StatusCapaianBreakdown
	for status, total := range summaryData.StatusCapaianBreakdown {
		percentage := 0.0
		if summaryData.TotalPPI > 0 {
			percentage = float64(total) / float64(summaryData.TotalPPI) * 100
		}
		statusBreakdown = append(statusBreakdown, dtos.StatusCapaianBreakdown{
			StatusCapaian: status,
			Total:         total,
			Percentage:    percentage,
		})
	}

	// Build aspek pembelajaran breakdown with percentage
	var aspekPembelajaranBreakdown []dtos.AspekPembelajaranBreakdown
	for _, asp := range summaryData.AspekPembelajaranBreakdown {
		percentage := 0.0
		if summaryData.TotalPPI > 0 {
			percentage = float64(asp.Total) / float64(summaryData.TotalPPI) * 100
		}
		
		aspekPembelajaranBreakdown = append(aspekPembelajaranBreakdown, dtos.AspekPembelajaranBreakdown{
			AspekPembelajaran: asp.AspekPembelajaran,
			Total:             asp.Total,
			Percentage:        percentage,
		})
	}

	return &dtos.PPIInklusiSummaryResponse{
		TahunPelajaran: dtos.TahunPelajaranSimple{
			ID:             tahunPelajaran.ID,
			TahunPelajaran: tahunPelajaran.TahunPelajaran,
			Status:         tahunPelajaran.Status,
		},
		TotalPPI:                       summaryData.TotalPPI,
		TotalAnakInklusi:               summaryData.TotalAnakInklusi,
		StatusCapaianBreakdown:         statusBreakdown,
		AspekPembelajaranBreakdown:     aspekPembelajaranBreakdown,
	}, nil
}

func (s *ppiInklusiService) toResponseDTO(ppiInklusi *models.PPIInklusi) *dtos.PPIInklusiResponse {
	resp := &dtos.PPIInklusiResponse{
		ID:                    ppiInklusi.ID,
		AnakInklusiRombelID:   ppiInklusi.AnakInklusiRombelID,
		TahunPelajaranID:      ppiInklusi.TahunPelajaranID,
		AspekPembelajaran:     ppiInklusi.AspekPembelajaran,
		NamaProgram:           ppiInklusi.NamaProgram,
		TujuanPembelajaran:    ppiInklusi.TujuanPembelajaran,
		StrategiPembelajaran:  ppiInklusi.StrategiPembelajaran,
		TargetWaktu:           ppiInklusi.TargetWaktu,
		CapaianSaatIni:        ppiInklusi.CapaianSaatIni,
		StatusCapaian:         ppiInklusi.StatusCapaian,
		GuruPembuatID:         ppiInklusi.GuruPembuatID,
		CreatedAt:             ppiInklusi.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             ppiInklusi.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Include TahunPelajaran data if preloaded
	if ppiInklusi.TahunPelajaran != nil {
		resp.TahunPelajaran = &dtos.TahunPelajaranSimple{
			ID:             ppiInklusi.TahunPelajaran.ID,
			TahunPelajaran: ppiInklusi.TahunPelajaran.TahunPelajaran,
			Status:         ppiInklusi.TahunPelajaran.Status,
		}
	}

	// Include GuruPembuat data if preloaded
	if ppiInklusi.GuruPembuat != nil {
		resp.GuruPembuat = &dtos.KepegawaianSimple{
			ID:       ppiInklusi.GuruPembuat.ID,
			Nama:     ppiInklusi.GuruPembuat.Nama,
			NIP:      ppiInklusi.GuruPembuat.NIP,
			Foto:     ppiInklusi.GuruPembuat.Foto,
			FotoURL:  s.r2Storage.GetPublicURL(ppiInklusi.GuruPembuat.Foto),
			Kategori: ppiInklusi.GuruPembuat.Kategori,
			Jabatan:  ppiInklusi.GuruPembuat.Jabatan,
		}
	}

	// Include AnakInklusiRombel data if preloaded (simplified version without nested data)
	if ppiInklusi.AnakInklusiRombel != nil {
		resp.AnakInklusiRombel = &dtos.AnakInklusiRombelResponse{
			ID:                     ppiInklusi.AnakInklusiRombel.ID,
			AnakInklusiID:          ppiInklusi.AnakInklusiRombel.AnakInklusiID,
			PesertaDidikRombelID:   ppiInklusi.AnakInklusiRombel.PesertaDidikRombelID,
			GuruPendampingKhususID: ppiInklusi.AnakInklusiRombel.GuruPendampingKhususID,
			GuruKelasID:            ppiInklusi.AnakInklusiRombel.GuruKelasID,
			Catatan:                ppiInklusi.AnakInklusiRombel.Catatan,
			CreatedAt:              ppiInklusi.AnakInklusiRombel.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:              ppiInklusi.AnakInklusiRombel.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Add nested data if available
		if ppiInklusi.AnakInklusiRombel.AnakInklusi != nil && ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik != nil {
			resp.AnakInklusiRombel.AnakInklusi = &dtos.AnakInklusiResponse{
				ID:             ppiInklusi.AnakInklusiRombel.AnakInklusi.ID,
				PesertaDidikID: ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidikID,
				JenisHambatan:  ppiInklusi.AnakInklusiRombel.AnakInklusi.JenisHambatan,
				Status:         ppiInklusi.AnakInklusiRombel.AnakInklusi.Status,
				PesertaDidik: &dtos.PesertaDidikSimple{
					ID:   ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.ID,
					NIS:  ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.NIS,
					NISN: ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.NISN,
					Nama: ppiInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.Nama,
				},
			}
		}

		if ppiInklusi.AnakInklusiRombel.PesertaDidikRombel != nil && ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel != nil {
			resp.AnakInklusiRombel.PesertaDidikRombel = &dtos.PesertaDidikRombelSimple{
				ID:               ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.ID,
				PesertaDidikID:   ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.PesertaDidikID,
				RombelID:         ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.RombelID,
				TahunPelajaranID: ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.TahunPelajaranID,
				Rombel: &dtos.RombelSimple{
					ID:     ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel.ID,
					Name:   ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel.Name,
					Status: ppiInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel.Status,
				},
			}
		}
	}

	return resp
}
