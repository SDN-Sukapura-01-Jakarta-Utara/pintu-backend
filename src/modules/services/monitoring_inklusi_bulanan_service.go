package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type MonitoringInklusiBulananService interface {
	Create(files []*multipart.FileHeader, req *dtos.MonitoringInklusiBulananCreateRequest, userID uint) (*dtos.MonitoringInklusiBulananResponse, error)
	GetByBulanTahunRombel(req *dtos.MonitoringInklusiBulananGetRequest) (*dtos.MonitoringInklusiBulananGetAllResponse, error)
	GetByID(id uint) (*dtos.MonitoringInklusiBulananResponse, error)
	Update(files []*multipart.FileHeader, req *dtos.MonitoringInklusiBulananUpdateRequest, userID uint) (*dtos.MonitoringInklusiBulananResponse, error)
	GetSummary(req *dtos.MonitoringInklusiBulananSummaryRequest) (*dtos.MonitoringInklusiBulananSummaryResponse, error)
}

type monitoringInklusiBulananService struct {
	repo      repositories.MonitoringInklusiBulananRepository
	r2Storage *utils.R2Storage
}

func NewMonitoringInklusiBulananService(repo repositories.MonitoringInklusiBulananRepository, r2Storage *utils.R2Storage) MonitoringInklusiBulananService {
	return &monitoringInklusiBulananService{
		repo:      repo,
		r2Storage: r2Storage,
	}
}

func (s *monitoringInklusiBulananService) Create(files []*multipart.FileHeader, req *dtos.MonitoringInklusiBulananCreateRequest, userID uint) (*dtos.MonitoringInklusiBulananResponse, error) {
	// Validate PPI Inklusi exists
	ppi, err := s.repo.GetPPIByID(req.PPIInklusiID)
	if err != nil {
		return nil, errors.New("PPI tidak ditemukan")
	}

	// If PPI status is not "Sedang Berlangsung", update it
	if ppi.StatusCapaian != "Sedang Berlangsung" {
		if err := s.repo.UpdatePPIStatus(req.PPIInklusiID, "Sedang Berlangsung"); err != nil {
			return nil, fmt.Errorf("gagal memperbarui status PPI: %w", err)
		}
	}

	// Handle file uploads
	var filePaths []string
	if len(files) > 0 {
		// Validate and upload each file
		for _, file := range files {
			// Validate file type (pdf, jpg, jpeg, png)
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				// Clean up already uploaded files
				for _, path := range filePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("file %s harus berformat PDF, JPG, JPEG, atau PNG", file.Filename)
			}

			// Validate file size (max 5MB per file)
			if file.Size > 5*1024*1024 {
				// Clean up already uploaded files
				for _, path := range filePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("ukuran file %s maksimal 5MB", file.Filename)
			}

			// Upload to R2
			uploadedPath, err := s.r2Storage.UploadFile(file, "inklusi/monitoring-bulanan")
			if err != nil {
				// Clean up already uploaded files
				for _, path := range filePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("gagal mengunggah file: %w", err)
			}
			filePaths = append(filePaths, uploadedPath)
		}
	}

	// Create monitoring record
	monitoring := &models.MonitoringInklusiBulanan{
		PPIInklusiID:          req.PPIInklusiID,
		Bulan:                 req.Bulan,
		Tahun:                 req.Tahun,
		DeskripsiPerkembangan: req.DeskripsiPerkembangan,
		KendalaDitemui:        req.KendalaDitemui,
		TindakLanjut:          req.TindakLanjut,
		FilePendukung:         filePaths,
		GuruPengisiID:         req.GuruPengisiID,
		CreatedByID:           &userID,
		UpdatedByID:           &userID,
	}

	if err := s.repo.Create(monitoring); err != nil {
		// Clean up uploaded files on database error
		for _, path := range filePaths {
			s.r2Storage.DeleteFile(path)
		}
		return nil, fmt.Errorf("gagal menyimpan data monitoring: %w", err)
	}

	// Get created data with relations
	created, err := s.repo.GetByID(monitoring.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring yang dibuat")
	}

	return s.toResponseDTO(created), nil
}

func (s *monitoringInklusiBulananService) toResponseDTO(monitoring *models.MonitoringInklusiBulanan) *dtos.MonitoringInklusiBulananResponse {
	resp := &dtos.MonitoringInklusiBulananResponse{
		ID:                    monitoring.ID,
		PPIInklusiID:          monitoring.PPIInklusiID,
		Bulan:                 monitoring.Bulan,
		Tahun:                 monitoring.Tahun,
		DeskripsiPerkembangan: monitoring.DeskripsiPerkembangan,
		KendalaDitemui:        monitoring.KendalaDitemui,
		TindakLanjut:          monitoring.TindakLanjut,
		FilePendukung:         monitoring.FilePendukung,
		GuruPengisiID:         monitoring.GuruPengisiID,
		CreatedAt:             monitoring.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             monitoring.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Build file URLs
	var fileURLs []string
	for _, filePath := range monitoring.FilePendukung {
		fileURLs = append(fileURLs, s.r2Storage.GetPublicURL(filePath))
	}
	resp.FilePendukungURL = fileURLs

	// Include PPIInklusi data if preloaded
	if monitoring.PPIInklusi != nil {
		resp.PPIInklusi = &dtos.PPIInklusiDetail{
			ID:                monitoring.PPIInklusi.ID,
			AspekPembelajaran: monitoring.PPIInklusi.AspekPembelajaran,
			NamaProgram:       monitoring.PPIInklusi.NamaProgram,
		}
	}

	// Include GuruPengisi data if preloaded
	if monitoring.GuruPengisi != nil {
		resp.GuruPengisi = &dtos.KepegawaianSimple{
			ID:       monitoring.GuruPengisi.ID,
			Nama:     monitoring.GuruPengisi.Nama,
			NIP:      monitoring.GuruPengisi.NIP,
			Foto:     monitoring.GuruPengisi.Foto,
			FotoURL:  s.r2Storage.GetPublicURL(monitoring.GuruPengisi.Foto),
			Kategori: monitoring.GuruPengisi.Kategori,
			Jabatan:  monitoring.GuruPengisi.Jabatan,
		}
	}

	return resp
}

func (s *monitoringInklusiBulananService) GetByBulanTahunRombel(req *dtos.MonitoringInklusiBulananGetRequest) (*dtos.MonitoringInklusiBulananGetAllResponse, error) {
	// Get rombel info if specified
	var rombelName string
	if req.RombelID != nil && *req.RombelID > 0 {
		rombel, err := s.repo.GetRombelByID(*req.RombelID)
		if err != nil {
			return nil, errors.New("rombel tidak ditemukan")
		}
		rombelName = rombel.Name
	}

	// Get all monitoring for this bulan, tahun, and optional rombel/PPI filter
	monitoringList, err := s.repo.GetByBulanTahunRombel(req.Bulan, req.Tahun, req.RombelID, req.PPIInklusiID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring: %w", err)
	}

	// Convert to response DTOs
	var responseList []dtos.MonitoringInklusiBulananResponse
	for _, monitoring := range monitoringList {
		responseList = append(responseList, *s.toResponseDTO(&monitoring))
	}

	return &dtos.MonitoringInklusiBulananGetAllResponse{
		Bulan:           req.Bulan,
		Tahun:           req.Tahun,
		RombelID:        req.RombelID,
		RombelName:      rombelName,
		PPIInklusiID:    req.PPIInklusiID,
		TotalMonitoring: len(responseList),
		Data:            responseList,
	}, nil
}

func (s *monitoringInklusiBulananService) GetByID(id uint) (*dtos.MonitoringInklusiBulananResponse, error) {
	monitoring, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("data monitoring tidak ditemukan")
	}

	return s.toResponseDTO(monitoring), nil
}

func (s *monitoringInklusiBulananService) Update(files []*multipart.FileHeader, req *dtos.MonitoringInklusiBulananUpdateRequest, userID uint) (*dtos.MonitoringInklusiBulananResponse, error) {
	// Get existing monitoring
	monitoring, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("data monitoring tidak ditemukan")
	}

	// Handle file deletions
	existingFiles := monitoring.FilePendukung
	if len(req.FilesToDelete) > 0 {
		// Filter out files to delete
		var remainingFiles []string
		for _, existingFile := range existingFiles {
			shouldDelete := false
			for _, fileToDelete := range req.FilesToDelete {
				if existingFile == fileToDelete {
					shouldDelete = true
					// Delete from R2 storage
					s.r2Storage.DeleteFile(fileToDelete)
					break
				}
			}
			if !shouldDelete {
				remainingFiles = append(remainingFiles, existingFile)
			}
		}
		existingFiles = remainingFiles
	}

	// Handle new file uploads
	var newFilePaths []string
	if len(files) > 0 {
		// Validate and upload each file
		for _, file := range files {
			// Validate file type
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				// Clean up already uploaded files
				for _, path := range newFilePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("file %s harus berformat PDF, JPG, JPEG, atau PNG", file.Filename)
			}

			// Validate file size
			if file.Size > 5*1024*1024 {
				// Clean up already uploaded files
				for _, path := range newFilePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("ukuran file %s maksimal 5MB", file.Filename)
			}

			// Upload to R2
			uploadedPath, err := s.r2Storage.UploadFile(file, "inklusi/monitoring-bulanan")
			if err != nil {
				// Clean up already uploaded files
				for _, path := range newFilePaths {
					s.r2Storage.DeleteFile(path)
				}
				return nil, fmt.Errorf("gagal mengunggah file: %w", err)
			}
			newFilePaths = append(newFilePaths, uploadedPath)
		}
	}

	// Combine existing files and new files
	allFiles := append(existingFiles, newFilePaths...)

	// Update monitoring fields
	monitoring.DeskripsiPerkembangan = req.DeskripsiPerkembangan
	monitoring.KendalaDitemui = req.KendalaDitemui
	monitoring.TindakLanjut = req.TindakLanjut
	monitoring.FilePendukung = allFiles
	monitoring.GuruPengisiID = req.GuruPengisiID
	monitoring.UpdatedByID = &userID

	// Save to database
	if err := s.repo.Update(monitoring); err != nil {
		// Clean up newly uploaded files on database error
		for _, path := range newFilePaths {
			s.r2Storage.DeleteFile(path)
		}
		return nil, fmt.Errorf("gagal memperbarui data monitoring: %w", err)
	}

	// Get updated data with relations
	updated, err := s.repo.GetByID(monitoring.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring yang diperbarui")
	}

	return s.toResponseDTO(updated), nil
}

func (s *monitoringInklusiBulananService) GetSummary(req *dtos.MonitoringInklusiBulananSummaryRequest) (*dtos.MonitoringInklusiBulananSummaryResponse, error) {
	// Get rombel info if specified
	var rombelName string
	if req.RombelID != nil && *req.RombelID > 0 {
		rombel, err := s.repo.GetRombelByID(*req.RombelID)
		if err != nil {
			return nil, errors.New("rombel tidak ditemukan")
		}
		rombelName = rombel.Name
	}

	// Get all monitoring data for the year (and optional month)
	monitoringList, err := s.repo.GetSummaryByTahun(req.Tahun, req.Bulan, req.RombelID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring: %w", err)
	}

	// Get all PPIs for the year
	ppiList, err := s.repo.GetPPIsByTahun(req.Tahun, req.RombelID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data PPI: %w", err)
	}

	// Get all anak inklusi for the year
	anakInklusiList, err := s.repo.GetAnakInklusiByTahun(req.Tahun, req.RombelID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data anak inklusi: %w", err)
	}

	// Calculate statistics
	totalMonitoring := len(monitoringList)
	totalPPI := len(ppiList)
	totalAnakInklusi := len(anakInklusiList)

	// Calculate rata-rata monitoring per PPI
	rataRataMonitoringPerPPI := 0.0
	if totalPPI > 0 {
		rataRataMonitoringPerPPI = float64(totalMonitoring) / float64(totalPPI)
	}

	// Calculate monitoring per bulan (Januari - Desember)
	monitoringPerBulan := make([]dtos.MonitoringBulananPerBulan, 12)
	bulanNames := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	
	// Initialize all months
	for i := 0; i < 12; i++ {
		monitoringPerBulan[i] = dtos.MonitoringBulananPerBulan{
			Bulan:           int16(i + 1),
			NamaBulan:       bulanNames[i],
			TotalMonitoring: 0,
			TotalPPIAktif:   0,
		}
	}

	// Count monitoring per month
	ppiAktifPerBulan := make(map[int16]map[uint]bool) // bulan -> ppi_id -> true
	for _, monitoring := range monitoringList {
		bulanIdx := monitoring.Bulan - 1
		if bulanIdx >= 0 && bulanIdx < 12 {
			monitoringPerBulan[bulanIdx].TotalMonitoring++
			
			// Track unique PPIs per month
			if ppiAktifPerBulan[monitoring.Bulan] == nil {
				ppiAktifPerBulan[monitoring.Bulan] = make(map[uint]bool)
			}
			ppiAktifPerBulan[monitoring.Bulan][monitoring.PPIInklusiID] = true
		}
	}

	// Count unique PPIs per month
	for bulan, ppiMap := range ppiAktifPerBulan {
		bulanIdx := bulan - 1
		if bulanIdx >= 0 && bulanIdx < 12 {
			monitoringPerBulan[bulanIdx].TotalPPIAktif = len(ppiMap)
		}
	}

	// Calculate monitoring per aspek pembelajaran
	aspekMap := make(map[string]*dtos.MonitoringBulananPerAspek)
	ppiPerAspek := make(map[string]map[uint]bool) // aspek -> ppi_id -> true

	for _, monitoring := range monitoringList {
		if monitoring.PPIInklusi != nil {
			aspek := monitoring.PPIInklusi.AspekPembelajaran
			
			if aspekMap[aspek] == nil {
				aspekMap[aspek] = &dtos.MonitoringBulananPerAspek{
					AspekPembelajaran: aspek,
					TotalMonitoring:   0,
					TotalPPI:          0,
				}
				ppiPerAspek[aspek] = make(map[uint]bool)
			}
			
			aspekMap[aspek].TotalMonitoring++
			ppiPerAspek[aspek][monitoring.PPIInklusiID] = true
		}
	}

	// Convert aspek map to slice and calculate averages
	var monitoringPerAspek []dtos.MonitoringBulananPerAspek
	for aspek, data := range aspekMap {
		data.TotalPPI = len(ppiPerAspek[aspek])
		if data.TotalPPI > 0 {
			data.RataRataPerPPI = float64(data.TotalMonitoring) / float64(data.TotalPPI)
		}
		monitoringPerAspek = append(monitoringPerAspek, *data)
	}

	// Calculate monitoring per siswa
	siswaMap := make(map[uint]*dtos.MonitoringBulananPerSiswa)
	ppiDetailMap := make(map[uint]map[uint]*dtos.MonitoringDetail) // anak_inklusi_id -> ppi_id -> detail

	for _, monitoring := range monitoringList {
		if monitoring.PPIInklusi != nil && monitoring.PPIInklusi.AnakInklusiRombel != nil && monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi != nil {
			anakInklusiID := monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.ID
			ppiID := monitoring.PPIInklusiID

			// Initialize siswa data
			if siswaMap[anakInklusiID] == nil {
				siswaMap[anakInklusiID] = &dtos.MonitoringBulananPerSiswa{
					AnakInklusiID:     anakInklusiID,
					JenisHambatan:     monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.JenisHambatan,
					TotalPPI:          0,
					TotalMonitoring:   0,
					MonitoringDetails: []dtos.MonitoringDetail{},
				}

				if monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik != nil {
					siswaMap[anakInklusiID].PesertaDidik = &dtos.PesertaDidikSimple{
						ID:   monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.ID,
						NIS:  monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.NIS,
						NISN: monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.NISN,
						Nama: monitoring.PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik.Nama,
					}
				}

				ppiDetailMap[anakInklusiID] = make(map[uint]*dtos.MonitoringDetail)
			}

			siswaMap[anakInklusiID].TotalMonitoring++

			// Track PPI details
			if ppiDetailMap[anakInklusiID][ppiID] == nil {
				ppiDetailMap[anakInklusiID][ppiID] = &dtos.MonitoringDetail{
					PPIInklusiID:      ppiID,
					AspekPembelajaran: monitoring.PPIInklusi.AspekPembelajaran,
					NamaProgram:       monitoring.PPIInklusi.NamaProgram,
					TotalMonitoring:   0,
					BulanTerakhir:     monitoring.Bulan,
				}
			}

			ppiDetailMap[anakInklusiID][ppiID].TotalMonitoring++
			if monitoring.Bulan > ppiDetailMap[anakInklusiID][ppiID].BulanTerakhir {
				ppiDetailMap[anakInklusiID][ppiID].BulanTerakhir = monitoring.Bulan
			}
		}
	}

	// Convert siswa map to slice
	var monitoringPerSiswa []dtos.MonitoringBulananPerSiswa
	for anakInklusiID, siswaData := range siswaMap {
		// Convert PPI details to slice
		for _, detail := range ppiDetailMap[anakInklusiID] {
			siswaData.MonitoringDetails = append(siswaData.MonitoringDetails, *detail)
		}
		siswaData.TotalPPI = len(siswaData.MonitoringDetails)
		monitoringPerSiswa = append(monitoringPerSiswa, *siswaData)
	}

	// Calculate PPI status breakdown
	ppiStatusBreakdown := dtos.PPIStatusBreakdown{
		BelumDimulai:      0,
		SedangBerlangsung: 0,
		Selesai:           0,
		Dihentikan:        0,
	}

	for _, ppi := range ppiList {
		switch ppi.StatusCapaian {
		case "Belum Dimulai":
			ppiStatusBreakdown.BelumDimulai++
		case "Sedang Berlangsung":
			ppiStatusBreakdown.SedangBerlangsung++
		case "Selesai":
			ppiStatusBreakdown.Selesai++
		case "Dihentikan":
			ppiStatusBreakdown.Dihentikan++
		}
	}

	return &dtos.MonitoringInklusiBulananSummaryResponse{
		Tahun:                    req.Tahun,
		Bulan:                    req.Bulan,
		RombelID:                 req.RombelID,
		RombelName:               rombelName,
		TotalAnakInklusi:         totalAnakInklusi,
		TotalPPI:                 totalPPI,
		TotalMonitoring:          totalMonitoring,
		RataRataMonitoringPerPPI: rataRataMonitoringPerPPI,
		MonitoringPerBulan:       monitoringPerBulan,
		MonitoringPerAspek:       monitoringPerAspek,
		MonitoringPerSiswa:       monitoringPerSiswa,
		PPIAktifPerStatus:        ppiStatusBreakdown,
	}, nil
}
