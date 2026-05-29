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
)

type AbsensiService interface {
	CreateAbsensiManual(req *dtos.AbsensiManualCreateRequest, files map[uint][]*multipart.FileHeader, userID uint) (*dtos.AbsensiManualCreateResponse, error)
	GetRekapAbsensi(req *dtos.AbsensiRekapRequest) (*dtos.AbsensiRekapResponse, error)
	UpdateRekapAbsensi(id uint, req *dtos.AbsensiUpdateRequest, file *multipart.FileHeader, userID uint) (*dtos.AbsensiUpdateResponse, error)
}

type AbsensiServiceImpl struct {
	repository repositories.AbsensiRepository
	r2Storage  *utils.R2Storage
}

// NewAbsensiService creates a new Absensi service
func NewAbsensiService(repository repositories.AbsensiRepository) AbsensiService {
	return &AbsensiServiceImpl{
		repository: repository,
		r2Storage:  utils.NewR2Storage(),
	}
}

// CreateAbsensiManual creates multiple absensi records (bulk input) with file upload support
func (s *AbsensiServiceImpl) CreateAbsensiManual(req *dtos.AbsensiManualCreateRequest, files map[uint][]*multipart.FileHeader, userID uint) (*dtos.AbsensiManualCreateResponse, error) {
	// Parse tanggal (YYYY-MM-DD format)
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	// Validasi pertemuan_ke untuk guru mapel (bidang_studi_id NOT NULL)
	if req.BidangStudiID != nil && req.PertemuanKe != nil {
		// Extract bulan dan tahun dari tanggal
		bulan := int(tanggal.Month())
		tahun := tanggal.Year()
		
		// Check if pertemuan already exists in this month
		existingAbsensi, err := s.repository.CheckPertemuanExists(
			req.RombelID,
			*req.BidangStudiID,
			req.TahunPelajaranID,
			req.Semester,
			bulan,
			tahun,
			*req.PertemuanKe,
		)
		
		if err != nil {
			return nil, errors.New("gagal memeriksa data pertemuan")
		}
		
		if existingAbsensi != nil {
			rombelNama := "rombel ini"
			if existingAbsensi.Rombel != nil {
				rombelNama = existingAbsensi.Rombel.Name
			}
			
			mapelNama := "mata pelajaran ini"
			if existingAbsensi.BidangStudi != nil {
				mapelNama = existingAbsensi.BidangStudi.Name
			}
			
			return nil, fmt.Errorf("pertemuan ke-%d untuk %s di %s sudah ada di bulan %d tahun %d", 
				*req.PertemuanKe, mapelNama, rombelNama, bulan, tahun)
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

	// Process each student in the list
	for _, item := range req.AbsensiList {
		// Check if absensi already exists for this student on this date and mapel
		existing, _ := s.repository.GetByPesertaDidikTanggalMapel(item.PesertaDidikID, tanggal, req.BidangStudiID)
		if existing != nil {
			totalFailed++
			var errorMsg string
			if req.BidangStudiID == nil {
				errorMsg = "absensi untuk tanggal ini sudah ada"
			} else {
				errorMsg = "absensi untuk tanggal ini di mata pelajaran ini sudah ada"
			}
			errorItems = append(errorItems, dtos.AbsensiCreateErrorItem{
				PesertaDidikID: item.PesertaDidikID,
				Message:        errorMsg,
			})
			continue
		}

		// Handle file upload for this student (if any)
		var fileSuratPath string
		if fileHeaders, ok := files[item.PesertaDidikID]; ok && len(fileHeaders) > 0 {
			// Only take the first file if multiple files uploaded
			fileHeader := fileHeaders[0]
			
			// Upload to R2 in absensi-siswa folder
			uploadedPath, err := s.r2Storage.UploadFile(fileHeader, "absensi-siswa")
			if err != nil {
				totalFailed++
				errorItems = append(errorItems, dtos.AbsensiCreateErrorItem{
					PesertaDidikID: item.PesertaDidikID,
					Message:        fmt.Sprintf("gagal upload file: %s", err.Error()),
				})
				continue
			}
			fileSuratPath = uploadedPath
		}

		// Create absensi record
		absensi := &models.Absensi{
			PesertaDidikID:   item.PesertaDidikID,
			RombelID:         &req.RombelID,
			TahunPelajaranID: req.TahunPelajaranID,
			Semester:         req.Semester,
			Tanggal:          tanggal,
			BidangStudiID:    req.BidangStudiID, // NULL = guru kelas, NOT NULL = guru mapel
			PertemuanKe:      req.PertemuanKe,   // NULL = guru kelas, NOT NULL = guru mapel
			Status:           item.Status,
			WaktuAbsen:       waktuAbsen,
			MetodeInput:      "manual",
			Keterangan:       item.Keterangan,
			FileSurat:        fileSuratPath,
			DicatatOlehID:    &userID,
		}

		if err := s.repository.Create(absensi); err != nil {
			totalFailed++
			errorItems = append(errorItems, dtos.AbsensiCreateErrorItem{
				PesertaDidikID: item.PesertaDidikID,
				Message:        fmt.Sprintf("gagal menyimpan data: %s", err.Error()),
			})
			
			// Delete uploaded file from R2 if save to DB failed
			if fileSuratPath != "" {
				s.r2Storage.DeleteFile(fileSuratPath)
			}
			continue
		}

		totalSuccess++
	}

	message := fmt.Sprintf("Berhasil menyimpan %d absensi", totalSuccess)
	if totalFailed > 0 {
		message += fmt.Sprintf(", %d gagal", totalFailed)
	}

	return &dtos.AbsensiManualCreateResponse{
		TotalSuccess: totalSuccess,
		TotalFailed:  totalFailed,
		Message:      message,
		Errors:       errorItems,
	}, nil
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

	// Group absensi by peserta_didik_id
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

		// Initialize siswa map if not exists
		if _, exists := siswaMap[absensi.PesertaDidikID]; !exists {
			siswaMap[absensi.PesertaDidikID] = &dtos.AbsensiRekapSiswa{
				PesertaDidikID:   absensi.PesertaDidikID,
				NIS:              absensi.PesertaDidik.NIS,
				Nama:             absensi.PesertaDidik.Nama,
				JenisKelamin:     absensi.PesertaDidik.JenisKelamin,
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

		siswa := siswaMap[absensi.PesertaDidikID]

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
func (s *AbsensiServiceImpl) mapToResponse(data *models.Absensi) *dtos.AbsensiResponse {
	response := &dtos.AbsensiResponse{
		ID:               data.ID,
		PesertaDidikID:   data.PesertaDidikID,
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

	if data.PesertaDidik != nil {
		response.PesertaDidikNama = data.PesertaDidik.Nama
	}

	if data.Rombel != nil {
		response.RombelNama = data.Rombel.Name
	}

	if data.BidangStudi != nil {
		response.BidangStudiNama = data.BidangStudi.Name
	}

	return response
}
