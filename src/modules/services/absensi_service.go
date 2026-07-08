package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"sort"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type AbsensiService interface {
	CreateAbsensiManual(req *dtos.AbsensiManualCreateRequest, files map[uint][]*multipart.FileHeader, userID uint) (*dtos.AbsensiManualCreateResponse, error)
	CreateAbsensiManualByID(req *dtos.AbsensiManualCreateByIDRequest, file *multipart.FileHeader, userID uint) (*dtos.AbsensiResponse, error)
	GetRekapAbsensi(req *dtos.AbsensiRekapRequest) (*dtos.AbsensiRekapResponse, error)
	UpdateRekapAbsensi(id uint, req *dtos.AbsensiUpdateRequest, file *multipart.FileHeader, userID uint) (*dtos.AbsensiUpdateResponse, error)
	GetDashboardSummary(req *dtos.DashboardSummaryRequest) (*dtos.DashboardSummaryResponse, error)
	GetGrafikKehadiran(req *dtos.GrafikKehadiranRequest) (*dtos.GrafikKehadiranResponse, error)
	GetStatistikPerHari(req *dtos.StatistikPerHariRequest) (*dtos.StatistikPerHariResponse, error)
	GetPerbandinganRombel(req *dtos.PerbandinganRombelRequest) (*dtos.PerbandinganRombelResponse, error)
	GetSiswaTerendah(req *dtos.SiswaTerendahRequest) (*dtos.SiswaTerendahResponse, error)
	GetDashboardSiswa(req *dtos.DashboardSiswaRequest) (*dtos.DashboardSiswaResponse, error)
	SynchronizeAbsensi(req *dtos.AbsensiSyncRequest, userID uint) (*dtos.AbsensiSyncResponse, error)
	ExportAbsensiExcel(req *dtos.ExportAbsensiExcelRequest) (*excelize.File, error)
	ExportAbsensiPDF(req *dtos.ExportAbsensiExcelRequest) ([]byte, error)
}

type AbsensiServiceImpl struct {
	repository repositories.AbsensiRepository
	r2Storage  *utils.R2Storage
	db         *gorm.DB
}

// NewAbsensiService creates a new Absensi service
func NewAbsensiService(repository repositories.AbsensiRepository, db *gorm.DB) AbsensiService {
	return &AbsensiServiceImpl{
		repository: repository,
		r2Storage:  utils.NewR2Storage(),
		db:         db,
	}
}

// CreateAbsensiManual creates multiple absensi records (bulk input) with file upload support
func (s *AbsensiServiceImpl) CreateAbsensiManual(req *dtos.AbsensiManualCreateRequest, files map[uint][]*multipart.FileHeader, userID uint) (*dtos.AbsensiManualCreateResponse, error) {
	// Parse tanggal (YYYY-MM-DD format)
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	// Validasi duplicate untuk guru kelas (bidang_studi_id = NULL)
	if req.BidangStudiID == nil {
		// Check if already exists for this rombel, tahun_pelajaran, semester, tanggal
		existing, _ := s.repository.CheckDuplicateGuruKelas(req.RombelID, req.TahunPelajaranID, req.Semester, tanggal)
		if existing != nil {
			return nil, fmt.Errorf("Absensi untuk rombel ini di tanggal %s sudah ada", req.Tanggal)
		}
	}

	// Validasi duplicate untuk guru mapel (bidang_studi_id NOT NULL)
	if req.BidangStudiID != nil && req.PertemuanKe != nil {
		// Extract bulan dan tahun dari tanggal
		bulan := int(tanggal.Month())
		tahun := tanggal.Year()
		
		// Check if already exists for this rombel, tahun_pelajaran, semester, bidang_studi, pertemuan_ke in the same month
		existing, _ := s.repository.CheckDuplicateGuruMapel(
			req.RombelID,
			req.TahunPelajaranID,
			*req.BidangStudiID,
			req.Semester,
			*req.PertemuanKe,
			bulan,
			tahun,
		)
		if existing != nil {
			return nil, fmt.Errorf("Pertemuan ke-%d untuk mata pelajaran ini di bulan %d tahun %d sudah ada", *req.PertemuanKe, bulan, tahun)
		}
	}

	// Parse waktu_absen if provided (YYYY-MM-DD HH:MM:SS format)
	var waktuAbsen *time.Time
	if req.WaktuAbsen != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.WaktuAbsen)
		if err != nil {
			return nil, errors.New("format waktu_absen tidak valid, gunakan YYYY-MM-DD HH:MM:SS")
		}
		waktuAbsen = &t
	} else {
		// Default to current time if not provided
		now := time.Now()
		waktuAbsen = &now
	}

	totalSuccess := 0
	totalFailed := 0
	var errorItems []dtos.AbsensiCreateErrorItem
	var uploadedFiles []string // Track uploaded files for cleanup

	// Start database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, errors.New("gagal memulai transaksi database")
	}

	// Defer rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Clean up uploaded files
			for _, filePath := range uploadedFiles {
				s.r2Storage.DeleteFile(filePath)
			}
		}
	}()

	// Process each student in the list
	for _, item := range req.AbsensiList {
		// Validate peserta_didik_rombel exists
		_, err := s.repository.GetPesertaDidikRombelByID(item.PesertaDidikRombelID)
		if err != nil {
			tx.Rollback()
			// Clean up uploaded files
			for _, filePath := range uploadedFiles {
				s.r2Storage.DeleteFile(filePath)
			}
			return nil, fmt.Errorf("data peserta didik rombel ID %d tidak ditemukan", item.PesertaDidikRombelID)
		}

		// Check if absensi already exists for this student on this date and mapel
		existing, _ := s.repository.GetByPesertaDidikTanggalMapel(item.PesertaDidikRombelID, tanggal, req.BidangStudiID)
		if existing != nil {
			tx.Rollback()
			// Clean up uploaded files
			for _, filePath := range uploadedFiles {
				s.r2Storage.DeleteFile(filePath)
			}
			var errorMsg string
			if req.BidangStudiID == nil {
				errorMsg = "absensi untuk tanggal ini sudah ada"
			} else {
				errorMsg = "absensi untuk tanggal ini di mata pelajaran ini sudah ada"
			}
			return nil, fmt.Errorf("peserta didik rombel ID %d: %s", item.PesertaDidikRombelID, errorMsg)
		}

		// Handle file upload for this student (if any)
		var fileSuratPath string
		if fileHeaders, ok := files[item.PesertaDidikRombelID]; ok && len(fileHeaders) > 0 {
			// Only take the first file if multiple files uploaded
			fileHeader := fileHeaders[0]
			
			// Upload to R2 in absensi-siswa folder
			uploadedPath, err := s.r2Storage.UploadFile(fileHeader, "absensi-siswa")
			if err != nil {
				tx.Rollback()
				// Clean up uploaded files
				for _, filePath := range uploadedFiles {
					s.r2Storage.DeleteFile(filePath)
				}
				return nil, fmt.Errorf("gagal upload file untuk peserta didik rombel ID %d: %s", item.PesertaDidikRombelID, err.Error())
			}
			fileSuratPath = uploadedPath
			uploadedFiles = append(uploadedFiles, uploadedPath) // Track for cleanup
		}

		// Create absensi record with peserta_didik_rombel_id using transaction
		absensi := &models.RekapitulasiAbsensi{
			PesertaDidikRombelID: item.PesertaDidikRombelID,
			RombelID:             &req.RombelID,
			TahunPelajaranID:     req.TahunPelajaranID,
			Semester:             req.Semester,
			Tanggal:              tanggal,
			BidangStudiID:        req.BidangStudiID, // NULL = guru kelas, NOT NULL = guru mapel
			PertemuanKe:          req.PertemuanKe,   // NULL = guru kelas, NOT NULL = guru mapel
			Status:               item.Status,
			WaktuAbsen:           waktuAbsen,
			MetodeInput:          "manual",
			Keterangan:           item.Keterangan,
			FileSurat:            fileSuratPath,
			DicatatOlehID:        &userID,
		}

		if err := tx.Create(absensi).Error; err != nil {
			tx.Rollback()
			// Clean up uploaded files
			for _, filePath := range uploadedFiles {
				s.r2Storage.DeleteFile(filePath)
			}
			return nil, fmt.Errorf("gagal menyimpan data untuk peserta didik rombel ID %d: %s", item.PesertaDidikRombelID, err.Error())
		}

		totalSuccess++
	}

	// Commit transaction if all succeeded
	if err := tx.Commit().Error; err != nil {
		// Clean up uploaded files if commit fails
		for _, filePath := range uploadedFiles {
			s.r2Storage.DeleteFile(filePath)
		}
		return nil, errors.New("gagal menyimpan transaksi ke database")
	}

	message := fmt.Sprintf("Berhasil menyimpan %d absensi", totalSuccess)

	return &dtos.AbsensiManualCreateResponse{
		TotalSuccess: totalSuccess,
		TotalFailed:  totalFailed,
		Message:      message,
		Errors:       errorItems,
	}, nil
}

// CreateAbsensiManualByID creates a single absensi record by peserta didik rombel ID with auto semester detection
func (s *AbsensiServiceImpl) CreateAbsensiManualByID(req *dtos.AbsensiManualCreateByIDRequest, file *multipart.FileHeader, userID uint) (*dtos.AbsensiResponse, error) {
	// Parse tanggal (YYYY-MM-DD format)
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	// Auto detect semester based on month
	// July-December (7-12) = Semester 1
	// January-June (1-6) = Semester 2
	month := int(tanggal.Month())
	semester := 2
	if month >= 7 && month <= 12 {
		semester = 1
	}

	// Parse waktu_absen if provided (YYYY-MM-DD HH:MM:SS format)
	var waktuAbsen *time.Time
	if req.WaktuAbsen != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.WaktuAbsen)
		if err != nil {
			return nil, errors.New("format waktu_absen tidak valid, gunakan YYYY-MM-DD HH:MM:SS")
		}
		waktuAbsen = &t
	} else {
		// Default to current time if not provided
		now := time.Now()
		waktuAbsen = &now
	}

	// Validate peserta_didik_rombel exists
	_, err = s.repository.GetPesertaDidikRombelByID(req.PesertaDidikRombelID)
	if err != nil {
		return nil, fmt.Errorf("data peserta didik rombel ID %d tidak ditemukan", req.PesertaDidikRombelID)
	}

	// Check if absensi already exists for this student on this date and mapel
	existing, _ := s.repository.GetByPesertaDidikTanggalMapel(req.PesertaDidikRombelID, tanggal, req.BidangStudiID)
	if existing != nil {
		var errorMsg string
		if req.BidangStudiID == nil {
			errorMsg = "absensi untuk tanggal ini sudah ada"
		} else {
			errorMsg = "absensi untuk tanggal ini di mata pelajaran ini sudah ada"
		}
		return nil, fmt.Errorf("peserta didik rombel ID %d: %s", req.PesertaDidikRombelID, errorMsg)
	}

	// Handle file upload (if any)
	var fileSuratPath string
	if file != nil {
		// Upload to R2 in absensi-siswa folder
		uploadedPath, err := s.r2Storage.UploadFile(file, "absensi-siswa")
		if err != nil {
			return nil, fmt.Errorf("gagal upload file: %s", err.Error())
		}
		fileSuratPath = uploadedPath
	}

	// Create absensi record
	absensi := &models.RekapitulasiAbsensi{
		PesertaDidikRombelID: req.PesertaDidikRombelID,
		RombelID:             &req.RombelID,
		TahunPelajaranID:     req.TahunPelajaranID,
		Semester:             semester, // Auto detected
		Tanggal:              tanggal,
		BidangStudiID:        req.BidangStudiID, // NULL = guru kelas, NOT NULL = guru mapel
		PertemuanKe:          req.PertemuanKe,   // NULL = guru kelas, NOT NULL = guru mapel
		Status:               req.Status,
		WaktuAbsen:           waktuAbsen,
		MetodeInput:          "manual",
		Keterangan:           req.Keterangan,
		FileSurat:            fileSuratPath,
		DicatatOlehID:        &userID,
	}

	if err := s.db.Create(absensi).Error; err != nil {
		// Clean up uploaded file if save fails
		if fileSuratPath != "" {
			s.r2Storage.DeleteFile(fileSuratPath)
		}
		return nil, fmt.Errorf("gagal menyimpan data absensi: %s", err.Error())
	}

	// Load relationships for response
	s.db.Preload("PesertaDidikRombel.PesertaDidik").
		Preload("Rombel").
		Preload("BidangStudi").
		Preload("DicatatOleh").
		First(absensi, absensi.ID)

	// Map to response
	response := s.mapToResponse(absensi)

	return response, nil
}

// GetRekapAbsensi retrieves attendance recap with summary per student
func (s *AbsensiServiceImpl) GetRekapAbsensi(req *dtos.AbsensiRekapRequest) (*dtos.AbsensiRekapResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai if provided
	var tanggalMulai, tanggalSelesai *time.Time
	
	if req.TanggalMulai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMulai = &t
	}
	
	if req.TanggalSelesai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	// Handle bidang_studi_id filter:
	// - If req.BidangStudiID is nil (null in JSON): filter for NULL (guru kelas only)
	// - If req.BidangStudiID is not nil: filter for specific bidang_studi_id (guru mapel)
	var bidangStudiFilter *uint
	if req.BidangStudiID == nil {
		// Explicitly set to nil to trigger NULL filter in repository
		bidangStudiFilter = nil
	} else {
		// Use the provided bidang_studi_id value
		bidangStudiFilter = req.BidangStudiID
	}
	
	// Get all absensi records based on filters
	absensiList, err := s.repository.GetRekapAbsensi(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.Bulan,
		req.Tahun,
		tanggalMulai,
		tanggalSelesai,
		bidangStudiFilter,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Group absensi by peserta_didik_rombel_id
	siswaMap := make(map[uint]*dtos.AbsensiRekapSiswa)
	var rombelNama string
	var bidangStudiNama string

	for _, absensi := range absensiList {
		// Get rombel name from first record
		if rombelNama == "" && absensi.Rombel != nil {
			rombelNama = absensi.Rombel.Name
		}

		// Get bidang studi name from first record
		if bidangStudiNama == "" && absensi.BidangStudi != nil {
			bidangStudiNama = absensi.BidangStudi.Name
		}

		// Get peserta_didik_id from PesertaDidikRombel
		pesertaDidikID := uint(0)
		var nis, nama, jenisKelamin string
		if absensi.PesertaDidikRombel != nil && absensi.PesertaDidikRombel.PesertaDidik != nil {
			pesertaDidikID = absensi.PesertaDidikRombel.PesertaDidik.ID
			nis = absensi.PesertaDidikRombel.PesertaDidik.NIS
			nama = absensi.PesertaDidikRombel.PesertaDidik.Nama
			jenisKelamin = absensi.PesertaDidikRombel.PesertaDidik.JenisKelamin
		}

		// Initialize siswa map if not exists
		if _, exists := siswaMap[pesertaDidikID]; !exists {
			siswaMap[pesertaDidikID] = &dtos.AbsensiRekapSiswa{
				PesertaDidikID:   pesertaDidikID,
				NIS:              nis,
				Nama:             nama,
				JenisKelamin:     jenisKelamin,
				TotalHadir:       0,
				TotalSakit:       0,
				TotalIzin:        0,
				TotalAlpa:        0,
				TotalAbsen:       0,
				TotalPertemuan:   0,
				PersentaseHadir:  0,
				DetailPerTanggal: []dtos.AbsensiDetailTanggal{},
			}
		}

		siswa := siswaMap[pesertaDidikID]

		// Count status
		switch absensi.Status {
		case "hadir":
			siswa.TotalHadir++
		case "sakit":
			siswa.TotalSakit++
		case "izin":
			siswa.TotalIzin++
		case "alpa":
			siswa.TotalAlpa++
		}

		siswa.TotalPertemuan++

		// Add detail per tanggal
		waktuAbsen := ""
		if absensi.WaktuAbsen != nil {
			waktuAbsen = absensi.WaktuAbsen.Format("2006-01-02 15:04:05")
		}

		dicatatOleh := ""
		if absensi.DicatatOleh != nil {
			dicatatOleh = absensi.DicatatOleh.Nama
		}

		// Generate full URL for file_surat
		fileSuratURL := s.r2Storage.GetPublicURL(absensi.FileSurat)

		siswa.DetailPerTanggal = append(siswa.DetailPerTanggal, dtos.AbsensiDetailTanggal{
			ID:            absensi.ID,
			Tanggal:       absensi.Tanggal.Format("2006-01-02"),
			PertemuanKe:   absensi.PertemuanKe, // Will be nil for guru kelas, has value for guru mapel
			Status:        absensi.Status,
			WaktuAbsen:    waktuAbsen,
			MetodeInput:   absensi.MetodeInput,
			Keterangan:    absensi.Keterangan,
			FileSurat:     fileSuratURL,
			DicatatOleh:   dicatatOleh,
			DicatatOlehID: absensi.DicatatOlehID,
		})
	}

	// Calculate total_absen, persentase hadir and convert map to slice
	var dataSiswa []dtos.AbsensiRekapSiswa
	for _, siswa := range siswaMap {
		// Calculate total_absen (sakit + izin + alpa)
		siswa.TotalAbsen = siswa.TotalSakit + siswa.TotalIzin + siswa.TotalAlpa
		
		// Calculate persentase hadir
		if siswa.TotalPertemuan > 0 {
			siswa.PersentaseHadir = float64(siswa.TotalHadir) / float64(siswa.TotalPertemuan) * 100
		}
		dataSiswa = append(dataSiswa, *siswa)
	}

	// Sort by nama (A-Z)
	sort.Slice(dataSiswa, func(i, j int) bool {
		return dataSiswa[i].Nama < dataSiswa[j].Nama
	})

	return &dtos.AbsensiRekapResponse{
		TahunPelajaranID: req.TahunPelajaranID,
		RombelID:         req.RombelID,
		RombelNama:       rombelNama,
		Semester:         req.Semester,
		Bulan:            req.Bulan,
		Tahun:            req.Tahun,
		BidangStudiID:    req.BidangStudiID,
		BidangStudiNama:  bidangStudiNama,
		TotalSiswa:       len(dataSiswa),
		DataSiswa:        dataSiswa,
	}, nil
}

// UpdateRekapAbsensi updates a single absensi record
func (s *AbsensiServiceImpl) UpdateRekapAbsensi(id uint, req *dtos.AbsensiUpdateRequest, file *multipart.FileHeader, userID uint) (*dtos.AbsensiUpdateResponse, error) {
	// Get existing absensi record
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("data absensi tidak ditemukan")
	}

	oldFileSurat := existing.FileSurat

	// Update status
	existing.Status = req.Status

	// Update keterangan
	existing.Keterangan = req.Keterangan

	// Update dicatat_oleh_id
	existing.DicatatOlehID = &userID

	// Handle file deletion if requested
	if req.DeleteFileSurat {
		// Delete old file from R2 if exists
		if oldFileSurat != "" {
			_ = s.r2Storage.DeleteFile(oldFileSurat)
		}
		existing.FileSurat = ""
	}

	// Handle file upload if provided (this will override delete_file_surat if both are sent)
	if file != nil {
		// Upload new file to R2
		uploadedPath, err := s.r2Storage.UploadFile(file, "absensi-siswa")
		if err != nil {
			return nil, fmt.Errorf("gagal upload file: %s", err.Error())
		}

		// Delete old file from R2 if exists (only if different from new file)
		if oldFileSurat != "" && oldFileSurat != uploadedPath {
			_ = s.r2Storage.DeleteFile(oldFileSurat)
		}

		existing.FileSurat = uploadedPath
	}

	// Save to database
	if err := s.repository.Update(existing); err != nil {
		// If update failed and new file was uploaded, delete the new file
		if file != nil && existing.FileSurat != "" && existing.FileSurat != oldFileSurat {
			s.r2Storage.DeleteFile(existing.FileSurat)
		}
		return nil, errors.New("gagal mengupdate data absensi")
	}

	// Map to response
	response := s.mapToResponse(existing)

	return &dtos.AbsensiUpdateResponse{
		Message: "Data absensi berhasil diupdate",
		Data:    response,
	}, nil
}

// mapToResponse maps Absensi model to AbsensiResponse DTO
func (s *AbsensiServiceImpl) mapToResponse(data *models.RekapitulasiAbsensi) *dtos.AbsensiResponse {
	pesertaDidikID := uint(0)
	pesertaDidikNama := ""
	if data.PesertaDidikRombel != nil && data.PesertaDidikRombel.PesertaDidik != nil {
		pesertaDidikID = data.PesertaDidikRombel.PesertaDidik.ID
		pesertaDidikNama = data.PesertaDidikRombel.PesertaDidik.Nama
	}

	response := &dtos.AbsensiResponse{
		ID:               data.ID,
		PesertaDidikID:   pesertaDidikID,
		PesertaDidikNama: pesertaDidikNama,
		RombelID:         data.RombelID,
		TahunPelajaranID: data.TahunPelajaranID,
		Semester:         data.Semester,
		Tanggal:          data.Tanggal.Format("2006-01-02"),
		BidangStudiID:    data.BidangStudiID,
		PertemuanKe:      data.PertemuanKe,
		Status:           data.Status,
		MetodeInput:      data.MetodeInput,
		Keterangan:       data.Keterangan,
		FileSurat:        s.r2Storage.GetPublicURL(data.FileSurat),
		DicatatOlehID:    data.DicatatOlehID,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if data.WaktuAbsen != nil {
		response.WaktuAbsen = data.WaktuAbsen.Format("2006-01-02 15:04:05")
	}

	if data.Rombel != nil {
		response.RombelNama = data.Rombel.Name
	}

	if data.BidangStudi != nil {
		response.BidangStudiNama = data.BidangStudi.Name
	}

	return response
}

// GetDashboardSummary retrieves dashboard summary statistics
func (s *AbsensiServiceImpl) GetDashboardSummary(req *dtos.DashboardSummaryRequest) (*dtos.DashboardSummaryResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai if provided
	var tanggalMulai, tanggalSelesai *time.Time
	
	if req.TanggalMulai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMulai = &t
	}
	
	if req.TanggalSelesai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	// Get all absensi records based on filters
	absensiList, err := s.repository.GetDashboardSummary(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
		tanggalMulai,
		tanggalSelesai,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	// Count unique students
	totalSiswa, err := s.repository.CountUniqueSiswa(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
	)
	if err != nil {
		return nil, errors.New("gagal menghitung jumlah siswa")
	}
	
	// Calculate summary statistics
	totalHadir := 0
	totalSakit := 0
	totalIzin := 0
	totalAlpa := 0
	
	// Count unique dates (pertemuan)
	dateMap := make(map[string]bool)
	
	for _, absensi := range absensiList {
		// Count by status
		switch absensi.Status {
		case "hadir":
			totalHadir++
		case "sakit":
			totalSakit++
		case "izin":
			totalIzin++
		case "alpa":
			totalAlpa++
		}
		
		// Track unique dates
		dateKey := absensi.Tanggal.Format("2006-01-02")
		dateMap[dateKey] = true
	}
	
	totalPertemuan := len(dateMap)
	
	// Calculate persentase kehadiran (rounded to 2 decimal places)
	totalKehadiran := totalHadir + totalSakit + totalIzin + totalAlpa
	persentaseKehadiran := 0.0
	if totalKehadiran > 0 {
		persentaseKehadiran = float64(int(float64(totalHadir) / float64(totalKehadiran) * 10000)) / 100
	}
	
	// Build response
	response := &dtos.DashboardSummaryResponse{
		TotalSiswa:     totalSiswa,
		TotalPertemuan: totalPertemuan,
		Summary: dtos.SummaryKehadiran{
			TotalHadir:          totalHadir,
			TotalSakit:          totalSakit,
			TotalIzin:           totalIzin,
			TotalAlpa:           totalAlpa,
			PersentaseKehadiran: persentaseKehadiran,
		},
	}
	
	// Calculate trend (always show, even if only 1 date)
	trend := s.calculateTrendFromData(absensiList, dateMap)
	response.Trend = trend
	
	return response, nil
}

// calculateTrendFromData calculates attendance trend from the last 2 dates in data
func (s *AbsensiServiceImpl) calculateTrendFromData(absensiList []models.RekapitulasiAbsensi, dateMap map[string]bool) *dtos.TrendKehadiran {
	// Get all unique dates and sort them
	var dates []time.Time
	for dateStr := range dateMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}
	
	// Sort dates descending (newest first)
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	
	// If no dates, return zero trend
	if len(dates) == 0 {
		return &dtos.TrendKehadiran{
			HadirKemarin: "0",
			HadirHariIni: "0",
			Perubahan:    "+0.0%",
		}
	}
	
	// If only 1 date, show it as "hari ini" with 0 for "kemarin"
	if len(dates) == 1 {
		hariIni := dates[0]
		hadirHariIni := 0
		
		for _, absensi := range absensiList {
			if absensi.Status == "hadir" && absensi.Tanggal.Format("2006-01-02") == hariIni.Format("2006-01-02") {
				hadirHariIni++
			}
		}
		
		return &dtos.TrendKehadiran{
			HadirKemarin: "0",
			HadirHariIni: fmt.Sprintf("%d", hadirHariIni),
			Perubahan:    "+0.0%",
		}
	}
	
	// Get the last 2 dates
	hariIni := dates[0]
	kemarin := dates[1]
	
	hadirKemarin := 0
	hadirHariIni := 0
	
	for _, absensi := range absensiList {
		if absensi.Status == "hadir" {
			if absensi.Tanggal.Format("2006-01-02") == kemarin.Format("2006-01-02") {
				hadirKemarin++
			} else if absensi.Tanggal.Format("2006-01-02") == hariIni.Format("2006-01-02") {
				hadirHariIni++
			}
		}
	}
	
	// Calculate percentage change (rounded to 2 decimal places)
	perubahan := 0.0
	if hadirKemarin > 0 {
		perubahan = float64(hadirHariIni-hadirKemarin) / float64(hadirKemarin) * 100
		perubahan = float64(int(perubahan*100)) / 100 // Round to 2 decimal places
	} else if hadirHariIni > 0 && hadirKemarin == 0 {
		perubahan = 100.0 // If no data yesterday but has data today
	}
	
	perubahanStr := fmt.Sprintf("%+.2f%%", perubahan)
	
	return &dtos.TrendKehadiran{
		HadirKemarin: fmt.Sprintf("%d", hadirKemarin),
		HadirHariIni: fmt.Sprintf("%d", hadirHariIni),
		Perubahan:    perubahanStr,
	}
}

// GetGrafikKehadiran retrieves attendance chart data
func (s *AbsensiServiceImpl) GetGrafikKehadiran(req *dtos.GrafikKehadiranRequest) (*dtos.GrafikKehadiranResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
	}
	
	tanggalSelesai, err := time.Parse("2006-01-02", req.TanggalSelesai)
	if err != nil {
		return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
	}
	
	// Validate date range
	if tanggalSelesai.Before(tanggalMulai) {
		return nil, errors.New("tanggal_selesai harus setelah atau sama dengan tanggal_mulai")
	}
	
	// Get all absensi records based on filters
	absensiList, err := s.repository.GetDashboardSummary(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
		&tanggalMulai,
		&tanggalSelesai,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	// Group data by periode
	var labels []string
	var dataHadir []int
	var dataSakit []int
	var dataIzin []int
	var dataAlpa []int
	
	switch req.Periode {
	case "harian":
		labels, dataHadir, dataSakit, dataIzin, dataAlpa = s.groupByHarian(absensiList, tanggalMulai, tanggalSelesai)
	case "mingguan":
		labels, dataHadir, dataSakit, dataIzin, dataAlpa = s.groupByMingguan(absensiList, tanggalMulai, tanggalSelesai)
	case "bulanan":
		labels, dataHadir, dataSakit, dataIzin, dataAlpa = s.groupByBulanan(absensiList, tanggalMulai, tanggalSelesai)
	default:
		return nil, errors.New("periode tidak valid, gunakan: harian, mingguan, atau bulanan")
	}
	
	// Build response
	response := &dtos.GrafikKehadiranResponse{
		Labels: labels,
		Datasets: []dtos.DatasetKehadiran{
			{
				Label: "Hadir",
				Data:  dataHadir,
			},
			{
				Label: "Sakit",
				Data:  dataSakit,
			},
			{
				Label: "Izin",
				Data:  dataIzin,
			},
			{
				Label: "Alpa",
				Data:  dataAlpa,
			},
		},
	}
	
	return response, nil
}

// groupByHarian groups attendance data by daily
func (s *AbsensiServiceImpl) groupByHarian(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) ([]string, []int, []int, []int, []int) {
	// Create map to store data per date
	dataMap := make(map[string]map[string]int)
	
	// Initialize all dates in range
	for d := tanggalMulai; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		dataMap[dateKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
	}
	
	// Count attendance by date
	for _, absensi := range absensiList {
		dateKey := absensi.Tanggal.Format("2006-01-02")
		if data, exists := dataMap[dateKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays
	var labels []string
	var dataHadir, dataSakit, dataIzin, dataAlpa []int
	
	for d := tanggalMulai; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		labels = append(labels, d.Format("02 Jan"))
		
		data := dataMap[dateKey]
		dataHadir = append(dataHadir, data["hadir"])
		dataSakit = append(dataSakit, data["sakit"])
		dataIzin = append(dataIzin, data["izin"])
		dataAlpa = append(dataAlpa, data["alpa"])
	}
	
	return labels, dataHadir, dataSakit, dataIzin, dataAlpa
}

// groupByMingguan groups attendance data by weekly
func (s *AbsensiServiceImpl) groupByMingguan(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) ([]string, []int, []int, []int, []int) {
	// Create map to store data per week
	dataMap := make(map[string]map[string]int)
	weekLabels := make(map[string]string)
	
	// Get start of first week (Monday)
	startWeek := tanggalMulai
	for startWeek.Weekday() != time.Monday {
		startWeek = startWeek.AddDate(0, 0, -1)
	}
	
	// Initialize all weeks in range
	for d := startWeek; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 7) {
		weekKey := d.Format("2006-01-02")
		endWeek := d.AddDate(0, 0, 6)
		
		dataMap[weekKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
		
		// Label format: "02-08 Jan" or "30 Dec - 05 Jan"
		if d.Month() == endWeek.Month() {
			weekLabels[weekKey] = fmt.Sprintf("%02d-%02d %s", d.Day(), endWeek.Day(), d.Format("Jan"))
		} else {
			weekLabels[weekKey] = fmt.Sprintf("%02d %s - %02d %s", d.Day(), d.Format("Jan"), endWeek.Day(), endWeek.Format("Jan"))
		}
	}
	
	// Count attendance by week
	for _, absensi := range absensiList {
		// Find which week this date belongs to
		weekStart := absensi.Tanggal
		for weekStart.Weekday() != time.Monday {
			weekStart = weekStart.AddDate(0, 0, -1)
		}
		
		weekKey := weekStart.Format("2006-01-02")
		if data, exists := dataMap[weekKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays (sorted by week)
	var weeks []time.Time
	for d := startWeek; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 7) {
		weeks = append(weeks, d)
	}
	
	var labels []string
	var dataHadir, dataSakit, dataIzin, dataAlpa []int
	
	for _, week := range weeks {
		weekKey := week.Format("2006-01-02")
		labels = append(labels, weekLabels[weekKey])
		
		data := dataMap[weekKey]
		dataHadir = append(dataHadir, data["hadir"])
		dataSakit = append(dataSakit, data["sakit"])
		dataIzin = append(dataIzin, data["izin"])
		dataAlpa = append(dataAlpa, data["alpa"])
	}
	
	return labels, dataHadir, dataSakit, dataIzin, dataAlpa
}

// groupByBulanan groups attendance data by monthly
func (s *AbsensiServiceImpl) groupByBulanan(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) ([]string, []int, []int, []int, []int) {
	// Create map to store data per month
	dataMap := make(map[string]map[string]int)
	
	// Initialize all months in range
	currentMonth := time.Date(tanggalMulai.Year(), tanggalMulai.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(tanggalSelesai.Year(), tanggalSelesai.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	for !currentMonth.After(endMonth) {
		monthKey := currentMonth.Format("2006-01")
		dataMap[monthKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	
	// Count attendance by month
	for _, absensi := range absensiList {
		monthKey := absensi.Tanggal.Format("2006-01")
		if data, exists := dataMap[monthKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays
	var labels []string
	var dataHadir, dataSakit, dataIzin, dataAlpa []int
	
	currentMonth = time.Date(tanggalMulai.Year(), tanggalMulai.Month(), 1, 0, 0, 0, 0, time.UTC)
	for !currentMonth.After(endMonth) {
		monthKey := currentMonth.Format("2006-01")
		labels = append(labels, currentMonth.Format("Jan 2006"))
		
		data := dataMap[monthKey]
		dataHadir = append(dataHadir, data["hadir"])
		dataSakit = append(dataSakit, data["sakit"])
		dataIzin = append(dataIzin, data["izin"])
		dataAlpa = append(dataAlpa, data["alpa"])
		
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	
	return labels, dataHadir, dataSakit, dataIzin, dataAlpa
}

// GetStatistikPerHari retrieves daily attendance statistics (pattern by day of week)
func (s *AbsensiServiceImpl) GetStatistikPerHari(req *dtos.StatistikPerHariRequest) (*dtos.StatistikPerHariResponse, error) {
	// Create date range for the specified month
	startDate := time.Date(req.Tahun, time.Month(req.Bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month
	
	// Get all absensi records for the month
	absensiList, err := s.repository.GetDashboardSummary(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
		&startDate,
		&endDate,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	// Initialize map for each day of week (Monday to Sunday)
	dayNames := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
	dayStats := make(map[string]*struct {
		TotalHadir int
		TotalAbsen int
		Count      int // Number of occurrences of this day in the month
	})
	
	for _, day := range dayNames {
		dayStats[day] = &struct {
			TotalHadir int
			TotalAbsen int
			Count      int
		}{}
	}
	
	// Count occurrences of each day in the month
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dayName := s.getDayNameInIndonesian(d.Weekday())
		dayStats[dayName].Count++
	}
	
	// Group attendance by day of week
	for _, absensi := range absensiList {
		dayName := s.getDayNameInIndonesian(absensi.Tanggal.Weekday())
		
		if absensi.Status == "hadir" {
			dayStats[dayName].TotalHadir++
		} else {
			dayStats[dayName].TotalAbsen++
		}
	}
	
	// Build response
	var data []dtos.HariKehadiran
	
	for _, dayName := range dayNames {
		stats := dayStats[dayName]
		
		// Calculate averages (rounded up to integer)
		rataRataHadir := 0
		rataRataAbsen := 0
		persentaseHadir := 0.0
		
		if stats.Count > 0 {
			// Use math.Ceil to round up
			rataRataHadir = int(float64(stats.TotalHadir) / float64(stats.Count))
			if float64(stats.TotalHadir)/float64(stats.Count) > float64(rataRataHadir) {
				rataRataHadir++ // Round up
			}
			
			rataRataAbsen = int(float64(stats.TotalAbsen) / float64(stats.Count))
			if float64(stats.TotalAbsen)/float64(stats.Count) > float64(rataRataAbsen) {
				rataRataAbsen++ // Round up
			}
		}
		
		// Calculate percentage
		totalKehadiran := stats.TotalHadir + stats.TotalAbsen
		if totalKehadiran > 0 {
			persentaseHadir = float64(stats.TotalHadir) / float64(totalKehadiran) * 100
			persentaseHadir = float64(int(persentaseHadir*100)) / 100 // Round to 2 decimal places
		}
		
		data = append(data, dtos.HariKehadiran{
			Hari:            dayName,
			RataRataHadir:   rataRataHadir,
			RataRataAbsen:   rataRataAbsen,
			PersentaseHadir: persentaseHadir,
		})
	}
	
	return &dtos.StatistikPerHariResponse{
		Data: data,
	}, nil
}

// getDayNameInIndonesian converts time.Weekday to Indonesian day name
func (s *AbsensiServiceImpl) getDayNameInIndonesian(weekday time.Weekday) string {
	dayNames := map[time.Weekday]string{
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
		time.Sunday:    "Minggu",
	}
	return dayNames[weekday]
}

// GetPerbandinganRombel retrieves attendance comparison between rombel (classes)
func (s *AbsensiServiceImpl) GetPerbandinganRombel(req *dtos.PerbandinganRombelRequest) (*dtos.PerbandinganRombelResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai if provided
	var tanggalMulai, tanggalSelesai *time.Time
	
	if req.TanggalMulai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMulai = &t
	}
	
	if req.TanggalSelesai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	// Get all absensi records
	absensiList, err := s.repository.GetPerbandinganRombel(
		req.TahunPelajaranID,
		req.Semester,
		req.BidangStudiID,
		tanggalMulai,
		tanggalSelesai,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	// Group by rombel_id
	rombelMap := make(map[uint]*dtos.RombelKehadiran)
	siswaPerRombel := make(map[uint]map[uint]bool) // Track unique students per rombel
	
	for _, absensi := range absensiList {
		if absensi.RombelID == nil {
			continue
		}
		
		rombelID := *absensi.RombelID
		
		// Initialize rombel data if not exists
		if _, exists := rombelMap[rombelID]; !exists {
			rombelNama := ""
			if absensi.Rombel != nil {
				rombelNama = absensi.Rombel.Name
			}
			
			rombelMap[rombelID] = &dtos.RombelKehadiran{
				RombelID:        rombelID,
				RombelNama:      rombelNama,
				TotalSiswa:      0,
				PersentaseHadir: 0,
				TotalHadir:      0,
				TotalSakit:      0,
				TotalIzin:       0,
				TotalAlpa:       0,
			}
			siswaPerRombel[rombelID] = make(map[uint]bool)
		}
		
		// Track unique students
		if absensi.PesertaDidikRombel != nil && absensi.PesertaDidikRombel.PesertaDidik != nil {
			siswaPerRombel[rombelID][absensi.PesertaDidikRombel.PesertaDidik.ID] = true
		}
		
		// Count by status
		rombel := rombelMap[rombelID]
		switch absensi.Status {
		case "hadir":
			rombel.TotalHadir++
		case "sakit":
			rombel.TotalSakit++
		case "izin":
			rombel.TotalIzin++
		case "alpa":
			rombel.TotalAlpa++
		}
	}
	
	// Calculate total_siswa and persentase_hadir for each rombel
	var data []dtos.RombelKehadiran
	for rombelID, rombel := range rombelMap {
		// Count unique students
		rombel.TotalSiswa = len(siswaPerRombel[rombelID])
		
		// Calculate persentase hadir
		totalKehadiran := rombel.TotalHadir + rombel.TotalSakit + rombel.TotalIzin + rombel.TotalAlpa
		if totalKehadiran > 0 {
			rombel.PersentaseHadir = float64(rombel.TotalHadir) / float64(totalKehadiran) * 100
			rombel.PersentaseHadir = float64(int(rombel.PersentaseHadir*100)) / 100 // Round to 2 decimal places
		}
		
		data = append(data, *rombel)
	}
	
	// Sort by persentase_hadir descending (highest first)
	sort.Slice(data, func(i, j int) bool {
		return data[i].PersentaseHadir > data[j].PersentaseHadir
	})
	
	return &dtos.PerbandinganRombelResponse{
		Data: data,
	}, nil
}

// GetSiswaTerendah retrieves students with lowest attendance
func (s *AbsensiServiceImpl) GetSiswaTerendah(req *dtos.SiswaTerendahRequest) (*dtos.SiswaTerendahResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai if provided
	var tanggalMulai, tanggalSelesai *time.Time
	
	if req.TanggalMulai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMulai = &t
	}
	
	if req.TanggalSelesai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	// Get all absensi records
	absensiList, err := s.repository.GetSiswaTerendah(
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
		tanggalMulai,
		tanggalSelesai,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	// Group by peserta_didik_id
	siswaMap := make(map[uint]*dtos.SiswaKehadiran)
	
	for _, absensi := range absensiList {
		// Get peserta_didik_id from PesertaDidikRombel
		pesertaDidikID := uint(0)
		var nis, nama string
		if absensi.PesertaDidikRombel != nil && absensi.PesertaDidikRombel.PesertaDidik != nil {
			pesertaDidikID = absensi.PesertaDidikRombel.PesertaDidik.ID
			nis = absensi.PesertaDidikRombel.PesertaDidik.NIS
			nama = absensi.PesertaDidikRombel.PesertaDidik.Nama
		}

		// Initialize siswa data if not exists
		if _, exists := siswaMap[pesertaDidikID]; !exists {
			siswaMap[pesertaDidikID] = &dtos.SiswaKehadiran{
				PesertaDidikID:  pesertaDidikID,
				NIS:             nis,
				Nama:            nama,
				TotalHadir:      0,
				TotalAbsen:      0,
				TotalPertemuan:  0,
				PersentaseHadir: 0,
			}
		}
		
		siswa := siswaMap[pesertaDidikID]
		siswa.TotalPertemuan++
		
		// Count by status
		if absensi.Status == "hadir" {
			siswa.TotalHadir++
		} else {
			siswa.TotalAbsen++
		}
	}
	
	// Calculate persentase hadir and convert to slice
	var data []dtos.SiswaKehadiran
	for _, siswa := range siswaMap {
		if siswa.TotalPertemuan > 0 {
			siswa.PersentaseHadir = float64(siswa.TotalHadir) / float64(siswa.TotalPertemuan) * 100
			siswa.PersentaseHadir = float64(int(siswa.PersentaseHadir*100)) / 100 // Round to 2 decimal places
		}
		data = append(data, *siswa)
	}
	
	// Sort by persentase_hadir ascending (lowest first)
	sort.Slice(data, func(i, j int) bool {
		return data[i].PersentaseHadir < data[j].PersentaseHadir
	})
	
	// Apply limit (default 10)
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}
	
	if len(data) > limit {
		data = data[:limit]
	}
	
	return &dtos.SiswaTerendahResponse{
		Data: data,
	}, nil
}

// GetDashboardSiswa retrieves dashboard data for a specific student
func (s *AbsensiServiceImpl) GetDashboardSiswa(req *dtos.DashboardSiswaRequest) (*dtos.DashboardSiswaResponse, error) {
	// Parse tanggal_mulai and tanggal_selesai if provided
	var tanggalMulai, tanggalSelesai *time.Time
	
	if req.TanggalMulai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMulai = &t
	}
	
	if req.TanggalSelesai != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	// Get all absensi records for this student
	absensiList, err := s.repository.GetDashboardSiswa(
		req.PesertaDidikRombelID,
		req.TahunPelajaranID,
		req.RombelID,
		req.Semester,
		req.BidangStudiID,
		tanggalMulai,
		tanggalSelesai,
	)
	if err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}
	
	if len(absensiList) == 0 {
		return nil, errors.New("data absensi tidak ditemukan")
	}
	
	// Get student info from first record
	firstRecord := absensiList[0]
	
	pesertaDidikID := uint(0)
	var nis, nama, jenisKelamin string
	if firstRecord.PesertaDidikRombel != nil && firstRecord.PesertaDidikRombel.PesertaDidik != nil {
		pesertaDidikID = firstRecord.PesertaDidikRombel.PesertaDidik.ID
		nis = firstRecord.PesertaDidikRombel.PesertaDidik.NIS
		nama = firstRecord.PesertaDidikRombel.PesertaDidik.Nama
		jenisKelamin = firstRecord.PesertaDidikRombel.PesertaDidik.JenisKelamin
	}
	
	siswa := dtos.InfoSiswa{
		PesertaDidikID: pesertaDidikID,
		NIS:            nis,
		Nama:           nama,
		JenisKelamin:   jenisKelamin,
		RombelNama:     "",
		Foto:           "", // PesertaDidik model doesn't have Foto field yet
	}
	
	if firstRecord.Rombel != nil {
		siswa.RombelNama = firstRecord.Rombel.Name
	}
	
	// Calculate summary
	totalPertemuan := len(absensiList)
	totalHadir := 0
	totalSakit := 0
	totalIzin := 0
	totalAlpa := 0
	
	for _, absensi := range absensiList {
		switch absensi.Status {
		case "hadir":
			totalHadir++
		case "sakit":
			totalSakit++
		case "izin":
			totalIzin++
		case "alpa":
			totalAlpa++
		}
	}
	
	persentaseHadir := 0.0
	if totalPertemuan > 0 {
		persentaseHadir = float64(totalHadir) / float64(totalPertemuan) * 100
		persentaseHadir = float64(int(persentaseHadir*100)) / 100 // Round to 2 decimal places
	}
	
	// Determine status kehadiran
	statusKehadiran := "Rendah"
	if persentaseHadir >= 90 {
		statusKehadiran = "Sangat Baik"
	} else if persentaseHadir >= 80 {
		statusKehadiran = "Baik"
	} else if persentaseHadir >= 70 {
		statusKehadiran = "Cukup"
	}
	
	summary := dtos.SummarySiswa{
		TotalPertemuan:  totalPertemuan,
		TotalHadir:      totalHadir,
		TotalSakit:      totalSakit,
		TotalIzin:       totalIzin,
		TotalAlpa:       totalAlpa,
		PersentaseHadir: persentaseHadir,
		StatusKehadiran: statusKehadiran,
	}
	
	// Build grafik based on periode
	// Parse tanggal range for grafik
	var grafikTanggalMulai, grafikTanggalSelesai time.Time
	if tanggalMulai != nil {
		grafikTanggalMulai = *tanggalMulai
	} else if len(absensiList) > 0 {
		// Use first record date as start
		grafikTanggalMulai = absensiList[len(absensiList)-1].Tanggal // Last in list (oldest)
	} else {
		grafikTanggalMulai = time.Now()
	}
	
	if tanggalSelesai != nil {
		grafikTanggalSelesai = *tanggalSelesai
	} else if len(absensiList) > 0 {
		// Use last record date as end
		grafikTanggalSelesai = absensiList[0].Tanggal // First in list (newest)
	} else {
		grafikTanggalSelesai = time.Now()
	}
	
	grafik := s.buildGrafikSiswa(absensiList, req.Periode, grafikTanggalMulai, grafikTanggalSelesai)
	
	// Apply limit to riwayat (default 10, max 100)
	limit := 10
	if req.LimitRiwayat > 0 {
		limit = req.LimitRiwayat
	}
	if limit > 100 {
		limit = 100
	}
	
	// Build riwayat absensi (limited)
	var riwayatAbsensi []dtos.RiwayatAbsensiSiswa
	maxItems := len(absensiList)
	if limit < maxItems {
		maxItems = limit
	}
	
	for i := 0; i < maxItems; i++ {
		absensi := absensiList[i]
		waktuAbsen := ""
		if absensi.WaktuAbsen != nil {
			waktuAbsen = absensi.WaktuAbsen.Format("15:04:05")
		}
		
		// Get day name in Indonesian
		hari := s.getDayNameInIndonesian(absensi.Tanggal.Weekday())
		
		riwayat := dtos.RiwayatAbsensiSiswa{
			Tanggal:     absensi.Tanggal.Format("2006-01-02"),
			Hari:        hari,
			Status:      absensi.Status,
			WaktuAbsen:  waktuAbsen,
			MetodeInput: absensi.MetodeInput,
			Keterangan:  absensi.Keterangan,
			FileSurat:   s.r2Storage.GetPublicURL(absensi.FileSurat),
			PertemuanKe: absensi.PertemuanKe,
		}
		
		riwayatAbsensi = append(riwayatAbsensi, riwayat)
	}
	
	return &dtos.DashboardSiswaResponse{
		Siswa:          siswa,
		Summary:        summary,
		Grafik:         grafik,
		RiwayatAbsensi: riwayatAbsensi,
	}, nil
}

// buildGrafikSiswa builds chart data for student dashboard based on periode
func (s *AbsensiServiceImpl) buildGrafikSiswa(absensiList []models.RekapitulasiAbsensi, periode string, tanggalMulai, tanggalSelesai time.Time) dtos.GrafikBulananSiswa {
	switch periode {
	case "harian":
		return s.buildGrafikHarianSiswa(absensiList, tanggalMulai, tanggalSelesai)
	case "mingguan":
		return s.buildGrafikMingguanSiswa(absensiList, tanggalMulai, tanggalSelesai)
	case "bulanan":
		return s.buildGrafikBulananSiswa(absensiList, tanggalMulai, tanggalSelesai)
	default:
		return s.buildGrafikBulananSiswa(absensiList, tanggalMulai, tanggalSelesai)
	}
}

// buildGrafikHarianSiswa builds daily chart data with all dates in range
func (s *AbsensiServiceImpl) buildGrafikHarianSiswa(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) dtos.GrafikBulananSiswa {
	// Create map to store data per date
	dateMap := make(map[string]map[string]int)
	
	// Initialize all dates in range
	for d := tanggalMulai; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		dateMap[dateKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
	}
	
	// Fill in actual data
	for _, absensi := range absensiList {
		dateKey := absensi.Tanggal.Format("2006-01-02")
		if data, exists := dateMap[dateKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays
	var labels []string
	var hadir, sakit, izin, alpa []int
	
	for d := tanggalMulai; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		labels = append(labels, d.Format("02 Jan"))
		
		data := dateMap[dateKey]
		hadir = append(hadir, data["hadir"])
		sakit = append(sakit, data["sakit"])
		izin = append(izin, data["izin"])
		alpa = append(alpa, data["alpa"])
	}
	
	return dtos.GrafikBulananSiswa{
		Labels: labels,
		Hadir:  hadir,
		Sakit:  sakit,
		Izin:   izin,
		Alpa:   alpa,
	}
}

// buildGrafikMingguanSiswa builds weekly chart data with all weeks in range
func (s *AbsensiServiceImpl) buildGrafikMingguanSiswa(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) dtos.GrafikBulananSiswa {
	// Create map to store data per week
	weekMap := make(map[string]map[string]int)
	weekLabels := make(map[string]string)
	
	// Get start of first week (Monday)
	startWeek := tanggalMulai
	for startWeek.Weekday() != time.Monday {
		startWeek = startWeek.AddDate(0, 0, -1)
	}
	
	// Initialize all weeks in range
	for d := startWeek; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 7) {
		weekKey := d.Format("2006-01-02")
		endWeek := d.AddDate(0, 0, 6)
		
		weekMap[weekKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
		
		// Label format
		if d.Month() == endWeek.Month() {
			weekLabels[weekKey] = fmt.Sprintf("%02d-%02d %s", d.Day(), endWeek.Day(), d.Format("Jan"))
		} else {
			weekLabels[weekKey] = fmt.Sprintf("%02d %s - %02d %s", d.Day(), d.Format("Jan"), endWeek.Day(), endWeek.Format("Jan"))
		}
	}
	
	// Fill in actual data
	for _, absensi := range absensiList {
		// Find which week this date belongs to
		weekStart := absensi.Tanggal
		for weekStart.Weekday() != time.Monday {
			weekStart = weekStart.AddDate(0, 0, -1)
		}
		
		weekKey := weekStart.Format("2006-01-02")
		if data, exists := weekMap[weekKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays (sorted by week)
	var weeks []time.Time
	for d := startWeek; !d.After(tanggalSelesai); d = d.AddDate(0, 0, 7) {
		weeks = append(weeks, d)
	}
	
	var labels []string
	var hadir, sakit, izin, alpa []int
	
	for _, week := range weeks {
		weekKey := week.Format("2006-01-02")
		labels = append(labels, weekLabels[weekKey])
		
		data := weekMap[weekKey]
		hadir = append(hadir, data["hadir"])
		sakit = append(sakit, data["sakit"])
		izin = append(izin, data["izin"])
		alpa = append(alpa, data["alpa"])
	}
	
	return dtos.GrafikBulananSiswa{
		Labels: labels,
		Hadir:  hadir,
		Sakit:  sakit,
		Izin:   izin,
		Alpa:   alpa,
	}
}

// buildGrafikBulananSiswa builds monthly chart data with all months in range
func (s *AbsensiServiceImpl) buildGrafikBulananSiswa(absensiList []models.RekapitulasiAbsensi, tanggalMulai, tanggalSelesai time.Time) dtos.GrafikBulananSiswa {
	// Create map to store data per month
	monthMap := make(map[string]map[string]int)
	
	// Initialize all months in range
	currentMonth := time.Date(tanggalMulai.Year(), tanggalMulai.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(tanggalSelesai.Year(), tanggalSelesai.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	for !currentMonth.After(endMonth) {
		monthKey := currentMonth.Format("2006-01")
		monthMap[monthKey] = map[string]int{
			"hadir": 0,
			"sakit": 0,
			"izin":  0,
			"alpa":  0,
		}
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	
	// Fill in actual data
	for _, absensi := range absensiList {
		monthKey := absensi.Tanggal.Format("2006-01")
		if data, exists := monthMap[monthKey]; exists {
			data[absensi.Status]++
		}
	}
	
	// Build arrays
	var labels []string
	var hadir, sakit, izin, alpa []int
	
	currentMonth = time.Date(tanggalMulai.Year(), tanggalMulai.Month(), 1, 0, 0, 0, 0, time.UTC)
	for !currentMonth.After(endMonth) {
		monthKey := currentMonth.Format("2006-01")
		labels = append(labels, currentMonth.Format("Jan 2006"))
		
		data := monthMap[monthKey]
		hadir = append(hadir, data["hadir"])
		sakit = append(sakit, data["sakit"])
		izin = append(izin, data["izin"])
		alpa = append(alpa, data["alpa"])
		
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	
	return dtos.GrafikBulananSiswa{
		Labels: labels,
		Hadir:  hadir,
		Sakit:  sakit,
		Izin:   izin,
		Alpa:   alpa,
	}
}

// SynchronizeAbsensi synchronizes data from absensi (scan) to rekapitulasi_absensi
func (s *AbsensiServiceImpl) SynchronizeAbsensi(req *dtos.AbsensiSyncRequest, userID uint) (*dtos.AbsensiSyncResponse, error) {
	var absensiScanList []models.Absensi
	var err error
	
	// Validate request based on tipe_sync and bidang_studi_id
	if req.BidangStudiID != nil {
		// Guru Bidang Studi: hanya support tipe_sync = tanggal
		if req.TipeSync != "tanggal" {
			return nil, errors.New("untuk guru bidang studi, hanya tipe_sync 'tanggal' yang diperbolehkan")
		}
		if req.PertemuanKe == nil {
			return nil, errors.New("pertemuan_ke wajib diisi untuk guru bidang studi")
		}
		
		// Parse tanggal untuk validasi
		tanggal, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			return nil, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
		}
		
		// Validasi: Cek apakah pertemuan ini sudah ada di bulan yang sama
		bulan := int(tanggal.Month())
		tahun := tanggal.Year()
		
		existingTanggal, err := s.repository.GetPertemuanTanggal(
			req.RombelID,
			*req.BidangStudiID,
			req.TahunPelajaranID,
			bulan,
			tahun,
			*req.PertemuanKe,
		)
		
		if err != nil {
			return nil, fmt.Errorf("gagal memeriksa pertemuan: %s", err.Error())
		}
		
		// Jika pertemuan sudah ada, validasi tanggalnya harus sama
		if existingTanggal != nil {
			existingTanggalStr := existingTanggal.Format("2006-01-02")
			requestTanggalStr := tanggal.Format("2006-01-02")
			
			if existingTanggalStr != requestTanggalStr {
				return nil, fmt.Errorf(
					"pertemuan ke-%d untuk bulan %d tahun %d sudah ada dengan tanggal %s. Gunakan tanggal yang sama untuk sinkronisasi atau ubah nomor pertemuan",
					*req.PertemuanKe,
					bulan,
					tahun,
					existingTanggalStr,
				)
			}
		}
	} else {
		// Guru Kelas: support tipe_sync = tanggal atau bulan
		if req.TipeSync != "tanggal" && req.TipeSync != "bulan" {
			return nil, errors.New("tipe_sync harus 'tanggal' atau 'bulan'")
		}
	}
	
	// Get absensi scan data based on tipe_sync
	if req.TipeSync == "tanggal" {
		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			return nil, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
		}
		
		absensiScanList, err = s.repository.GetAbsensiScanByDate(tanggal)
		if err != nil {
			return nil, errors.New("gagal mengambil data absensi scan")
		}
	} else {
		// By bulan (guru kelas only)
		absensiScanList, err = s.repository.GetAbsensiScanByMonth(req.Bulan, req.Tahun)
		if err != nil {
			return nil, errors.New("gagal mengambil data absensi scan")
		}
	}
	
	if len(absensiScanList) == 0 {
		return nil, errors.New("tidak ada data absensi scan yang ditemukan")
	}
	
	totalProcessed := 0
	totalInserted := 0
	totalUpdated := 0
	totalSkipped := 0
	var details []dtos.AbsensiSyncDetailItem
	
	// Process each absensi scan record
	for _, absensiScan := range absensiScanList {
		totalProcessed++
		
		// Skip if jam_datang is null (student didn't scan in)
		if absensiScan.JamDatang == nil {
			totalSkipped++
			details = append(details, dtos.AbsensiSyncDetailItem{
				PesertaDidikID: absensiScan.PesertaDidikID,
				NIS:            absensiScan.PesertaDidik.NIS,
				Nama:           absensiScan.PesertaDidik.Nama,
				Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
				Action:         "skipped",
				Reason:         "tidak ada jam datang",
			})
			continue
		}
		
		// Find peserta_didik_rombel_id based on peserta_didik_id, rombel_id, and tahun_pelajaran_id
		pesertaDidikRombelID, err := s.repository.GetPesertaDidikRombelID(absensiScan.PesertaDidikID, req.RombelID)
		if err != nil {
			totalSkipped++
			details = append(details, dtos.AbsensiSyncDetailItem{
				PesertaDidikID: absensiScan.PesertaDidikID,
				NIS:            absensiScan.PesertaDidik.NIS,
				Nama:           absensiScan.PesertaDidik.Nama,
				Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
				Action:         "skipped",
				Reason:         "tidak ditemukan di rombel",
			})
			continue
		}
		
		// Determine semester based on month (Juli-Desember = 1, Januari-Juni = 2)
		semester := 1
		if absensiScan.Tanggal.Month() >= 1 && absensiScan.Tanggal.Month() <= 6 {
			semester = 2
		}
		
		// Combine tanggal and jam_datang to create waktu_absen (timestamp)
		waktuAbsenStr := fmt.Sprintf("%s %s", absensiScan.Tanggal.Format("2006-01-02"), *absensiScan.JamDatang)
		waktuAbsen, err := time.Parse("2006-01-02 15:04:05", waktuAbsenStr)
		if err != nil {
			totalSkipped++
			details = append(details, dtos.AbsensiSyncDetailItem{
				PesertaDidikID: absensiScan.PesertaDidikID,
				NIS:            absensiScan.PesertaDidik.NIS,
				Nama:           absensiScan.PesertaDidik.Nama,
				Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
				Action:         "skipped",
				Reason:         "format waktu tidak valid",
			})
			continue
		}
		
		// Check if record already exists in rekapitulasi_absensi
		existing, _ := s.repository.GetByPesertaDidikTanggalMapel(pesertaDidikRombelID, absensiScan.Tanggal, req.BidangStudiID)
		
		if existing != nil {
			// UPDATE existing record
			existing.Status = "hadir"
			existing.WaktuAbsen = &waktuAbsen
			existing.MetodeInput = "auto"
			existing.DicatatOlehID = &userID
			
			if err := s.repository.Update(existing); err != nil {
				totalSkipped++
				details = append(details, dtos.AbsensiSyncDetailItem{
					PesertaDidikID: absensiScan.PesertaDidikID,
					NIS:            absensiScan.PesertaDidik.NIS,
					Nama:           absensiScan.PesertaDidik.Nama,
					Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
					Action:         "skipped",
					Reason:         "gagal update: " + err.Error(),
				})
				continue
			}
			
			totalUpdated++
			details = append(details, dtos.AbsensiSyncDetailItem{
				PesertaDidikID: absensiScan.PesertaDidikID,
				NIS:            absensiScan.PesertaDidik.NIS,
				Nama:           absensiScan.PesertaDidik.Nama,
				Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
				Action:         "updated",
			})
		} else {
			// INSERT new record
			newRekap := &models.RekapitulasiAbsensi{
				PesertaDidikRombelID: pesertaDidikRombelID,
				RombelID:             &req.RombelID,
				TahunPelajaranID:     req.TahunPelajaranID,
				Semester:             semester,
				Tanggal:              absensiScan.Tanggal,
				BidangStudiID:        req.BidangStudiID,
				PertemuanKe:          req.PertemuanKe,
				Status:               "hadir",
				WaktuAbsen:           &waktuAbsen,
				MetodeInput:          "auto",
				Keterangan:           "",
				FileSurat:            "",
				DicatatOlehID:        &userID,
			}
			
			if err := s.repository.Create(newRekap); err != nil {
				totalSkipped++
				details = append(details, dtos.AbsensiSyncDetailItem{
					PesertaDidikID: absensiScan.PesertaDidikID,
					NIS:            absensiScan.PesertaDidik.NIS,
					Nama:           absensiScan.PesertaDidik.Nama,
					Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
					Action:         "skipped",
					Reason:         "gagal insert: " + err.Error(),
				})
				continue
			}
			
			totalInserted++
			details = append(details, dtos.AbsensiSyncDetailItem{
				PesertaDidikID: absensiScan.PesertaDidikID,
				NIS:            absensiScan.PesertaDidik.NIS,
				Nama:           absensiScan.PesertaDidik.Nama,
				Tanggal:        absensiScan.Tanggal.Format("2006-01-02"),
				Action:         "inserted",
			})
		}
	}
	
	message := fmt.Sprintf("Sinkronisasi selesai: %d diproses, %d ditambahkan, %d diupdate, %d dilewati", 
		totalProcessed, totalInserted, totalUpdated, totalSkipped)
	
	return &dtos.AbsensiSyncResponse{
		TotalProcessed: totalProcessed,
		TotalInserted:  totalInserted,
		TotalUpdated:   totalUpdated,
		TotalSkipped:   totalSkipped,
		Message:        message,
		Details:        details,
	}, nil
}
