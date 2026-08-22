package services

import (
	"errors"
	"fmt"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"time"

	"gorm.io/gorm"
)

type AbsensiEkskulService struct {
	repo *repositories.AbsensiEkskulRepository
	db   *gorm.DB
}

func NewAbsensiEkskulService(repo *repositories.AbsensiEkskulRepository, db *gorm.DB) *AbsensiEkskulService {
	return &AbsensiEkskulService{
		repo: repo,
		db:   db,
	}
}

// Create creates new ekstrakurikuler attendance (kegiatan + bulk absensi siswa + absensi pelatih)
func (s *AbsensiEkskulService) Create(req *dtos.AbsensiEkskulCreateRequest, createdByID *uint) (*dtos.AbsensiEkskulResponse, error) {
	// Validate ekstrakurikuler exists
	exists, err := s.repo.CheckEkstrakurikulerExists(req.EkstrakurikulerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("ekstrakurikuler not found")
	}

	// Validate tahun pelajaran exists
	exists, err = s.repo.CheckTahunPelajaranExists(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("tahun pelajaran not found")
	}

	// Check if kegiatan already exists for same ekstrakurikuler, tahun pelajaran, and date
	tanggalStr := req.TanggalKegiatan.Format("2006-01-02")
	exists, err = s.repo.CheckKegiatanExists(req.EkstrakurikulerID, req.TahunPelajaranID, tanggalStr)
	if err != nil {
		return nil, err
	}
	if exists {
		// Get ekstrakurikuler and tahun pelajaran names for better error message
		ekskulName, _ := s.repo.GetEkstrakurikulerName(req.EkstrakurikulerID)
		tahunName, _ := s.repo.GetTahunPelajaranName(req.TahunPelajaranID)
		
		if ekskulName != "" && tahunName != "" {
			return nil, fmt.Errorf("tanggal %s sudah digunakan pada ekstrakurikuler %s di tahun pelajaran %s", tanggalStr, ekskulName, tahunName)
		}
		return nil, fmt.Errorf("tanggal %s sudah digunakan pada ekstrakurikuler ini di tahun pelajaran ini", tanggalStr)
	}

	// Validate all peserta didik rombel exist
	for _, siswa := range req.AbsensiSiswa {
		exists, err := s.repo.CheckPesertaDidikRombelExists(siswa.PesertaDidikRombelID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("peserta didik rombel with id %d not found", siswa.PesertaDidikRombelID)
		}
	}

	// Validate all pelatih exist
	for _, pelatih := range req.AbsensiPelatih {
		exists, err := s.repo.CheckPelatihExists(pelatih.PelatihID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("pelatih with id %d not found", pelatih.PelatihID)
		}
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Ensure rollback on any error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw panic after rollback
		}
	}()

	// Create repository with transaction
	txRepo := repositories.NewAbsensiEkskulRepository(tx)

	// 1. Create kegiatan ekskul (foto_kegiatan is NULL for create absensi)
	kegiatan := &models.KegiatanEkskul{
		EkstrakurikulerID: req.EkstrakurikulerID,
		TahunPelajaranID:  req.TahunPelajaranID,
		TanggalKegiatan:   req.TanggalKegiatan,
		WaktuMulai:        req.WaktuMulai,
		WaktuSelesai:      req.WaktuSelesai,
		MateriKegiatan:    req.MateriKegiatan,
		FotoKegiatan:      nil, // NOT saved in create absensi
		CreatedByID:       createdByID,
		UpdatedByID:       createdByID,
	}

	if err := txRepo.CreateKegiatanEkskul(kegiatan); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create kegiatan ekskul: %w", err)
	}

	// 2. Bulk create absensi siswa
	absensiSiswaList := make([]models.AbsensiEkskul, len(req.AbsensiSiswa))
	for i, siswa := range req.AbsensiSiswa {
		// Normalize status: "alpha" → "alpa" for consistency (Indonesian spelling)
		status := siswa.Status
		if status == "alpha" {
			status = "alpa"
		}
		
		absensiSiswaList[i] = models.AbsensiEkskul{
			KegiatanEkskulID:     kegiatan.ID,
			PesertaDidikRombelID: siswa.PesertaDidikRombelID,
			Status:               status,
			Keterangan:           siswa.Keterangan,
		}
	}

	if err := txRepo.BulkCreateAbsensiEkskul(absensiSiswaList); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create absensi siswa: %w", err)
	}

	// 3. Bulk create absensi pelatih
	absensiPelatihList := make([]models.AbsensiPelatihEkskul, len(req.AbsensiPelatih))
	for i, pelatih := range req.AbsensiPelatih {
		absensiPelatihList[i] = models.AbsensiPelatihEkskul{
			KegiatanEkskulID: kegiatan.ID,
			PelatihID:        pelatih.PelatihID,
		}
	}

	if err := txRepo.BulkCreateAbsensiPelatih(absensiPelatihList); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create absensi pelatih: %w", err)
	}

	// Commit transaction - if this fails, everything is rolled back
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Calculate statistics for response
	var totalHadir, totalSakit, totalIzin, totalAlpha int
	for _, siswa := range req.AbsensiSiswa {
		// Normalize status for counting (alpha → alpa)
		status := siswa.Status
		if status == "alpha" {
			status = "alpa"
		}
		
		switch status {
		case "hadir":
			totalHadir++
		case "sakit":
			totalSakit++
		case "izin":
			totalIzin++
		case "alpa":
			totalAlpha++
		}
	}

	// Build response
	response := &dtos.AbsensiEkskulResponse{
		KegiatanEkskulID:  kegiatan.ID,
		EkstrakurikulerID: kegiatan.EkstrakurikulerID,
		TahunPelajaranID:  kegiatan.TahunPelajaranID,
		TanggalKegiatan:   kegiatan.TanggalKegiatan.Format("2006-01-02"),
		MateriKegiatan:    kegiatan.MateriKegiatan,
		TotalSiswaHadir:   totalHadir,
		TotalSiswaSakit:   totalSakit,
		TotalSiswaIzin:    totalIzin,
		TotalSiswaAlpha:   totalAlpha,
		TotalPelatihHadir: len(req.AbsensiPelatih),
	}

	return response, nil
}

// GetAbsensiSiswa gets all attendance data by ekstrakurikuler and tahun pelajaran
func (s *AbsensiEkskulService) GetAbsensiSiswa(req *dtos.AbsensiEkskulGetRequest) (*dtos.AbsensiEkskulGetResponse, error) {
	// Validate ekstrakurikuler exists
	exists, err := s.repo.CheckEkstrakurikulerExists(req.EkstrakurikulerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("ekstrakurikuler not found")
	}

	// Validate tahun pelajaran exists
	exists, err = s.repo.CheckTahunPelajaranExists(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("tahun pelajaran not found")
	}

	// Get ekstrakurikuler and tahun pelajaran names
	ekskulName, err := s.repo.GetEkstrakurikulerName(req.EkstrakurikulerID)
	if err != nil {
		return nil, err
	}

	tahunName, err := s.repo.GetTahunPelajaranName(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}

	// Get all kegiatan with relations
	kegiatanList, err := s.repo.GetKegiatanByEkskulAndTahun(req.EkstrakurikulerID, req.TahunPelajaranID, req.Nama, req.RombelID, req.Bulan, req.Tahun)
	if err != nil {
		return nil, err
	}

	// Build response
	kegiatanDetails := make([]dtos.KegiatanEkskulDetail, 0, len(kegiatanList))
	for _, kegiatan := range kegiatanList {
		// Build absensi siswa details
		absensiSiswaList := make([]dtos.AbsensiSiswaDetail, 0, len(kegiatan.AbsensiEkskul))
		var totalHadir, totalSakit, totalIzin, totalAlpha int
		
		for _, absensi := range kegiatan.AbsensiEkskul {
			// Normalize status for counting (alpha → alpa)
			status := absensi.Status
			if status == "alpha" {
				status = "alpa"
			}
			
			// Count by status
			switch status {
			case "hadir":
				totalHadir++
			case "sakit":
				totalSakit++
			case "izin":
				totalIzin++
			case "alpa":
				totalAlpha++
			}

			// Build detail
			detail := dtos.AbsensiSiswaDetail{
				ID:                   absensi.ID,
				PesertaDidikRombelID: absensi.PesertaDidikRombelID,
				Status:               absensi.Status,
				Keterangan:           absensi.Keterangan,
			}

			// Add student info if available
			if absensi.PesertaDidikRombel != nil {
				if absensi.PesertaDidikRombel.PesertaDidik != nil {
					detail.NamaSiswa = absensi.PesertaDidikRombel.PesertaDidik.Nama
					detail.NIS = absensi.PesertaDidikRombel.PesertaDidik.NIS
					detail.NISN = absensi.PesertaDidikRombel.PesertaDidik.NISN
				}
				if absensi.PesertaDidikRombel.Rombel != nil {
					detail.NamaRombel = absensi.PesertaDidikRombel.Rombel.Name
					if absensi.PesertaDidikRombel.Rombel.Kelas != nil {
						detail.NamaKelas = absensi.PesertaDidikRombel.Rombel.Kelas.Name
					}
				}
			}

			absensiSiswaList = append(absensiSiswaList, detail)
		}

		// Build absensi pelatih details
		absensiPelatihList := make([]dtos.AbsensiPelatihDetail, 0, len(kegiatan.AbsensiPelatih))
		for _, absensi := range kegiatan.AbsensiPelatih {
			detail := dtos.AbsensiPelatihDetail{
				ID:        absensi.ID,
				PelatihID: absensi.PelatihID,
			}

			// Add pelatih info if available
			if absensi.Pelatih != nil {
				detail.NamaPelatih = absensi.Pelatih.Nama
				detail.Telepon = absensi.Pelatih.Telepon
			}

			absensiPelatihList = append(absensiPelatihList, detail)
		}

		// Build kegiatan detail
		kegiatanDetail := dtos.KegiatanEkskulDetail{
			ID:                kegiatan.ID,
			TanggalKegiatan:   kegiatan.TanggalKegiatan.Format("2006-01-02"),
			WaktuMulai:        kegiatan.WaktuMulai,
			WaktuSelesai:      kegiatan.WaktuSelesai,
			MateriKegiatan:    kegiatan.MateriKegiatan,
			FotoKegiatan:      kegiatan.FotoKegiatan,
			AbsensiSiswa:      absensiSiswaList,
			AbsensiPelatih:    absensiPelatihList,
			TotalSiswaHadir:   totalHadir,
			TotalSiswaSakit:   totalSakit,
			TotalSiswaIzin:    totalIzin,
			TotalSiswaAlpha:   totalAlpha,
			TotalPelatihHadir: len(absensiPelatihList),
		}

		kegiatanDetails = append(kegiatanDetails, kegiatanDetail)
	}

	// Build final response
	response := &dtos.AbsensiEkskulGetResponse{
		EkstrakurikulerID:   req.EkstrakurikulerID,
		NamaEkstrakurikuler: ekskulName,
		TahunPelajaranID:    req.TahunPelajaranID,
		TahunPelajaran:      tahunName,
		TotalKegiatan:       len(kegiatanDetails),
		Kegiatan:            kegiatanDetails,
	}

	return response, nil
}


// GetAbsensiSiswaByID gets single absensi siswa detail by ID
func (s *AbsensiEkskulService) GetAbsensiSiswaByID(req *dtos.AbsensiSiswaGetByIDRequest) (*dtos.AbsensiSiswaDetailResponse, error) {
	// Get absensi siswa by ID
	absensi, err := s.repo.GetAbsensiSiswaByID(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("absensi siswa not found")
		}
		return nil, err
	}

	// Build response
	response := &dtos.AbsensiSiswaDetailResponse{
		ID:                   absensi.ID,
		KegiatanEkskulID:     absensi.KegiatanEkskulID,
		PesertaDidikRombelID: absensi.PesertaDidikRombelID,
		Status:               absensi.Status,
		Keterangan:           absensi.Keterangan,
	}

	// Add student info
	if absensi.PesertaDidikRombel != nil {
		if absensi.PesertaDidikRombel.PesertaDidik != nil {
			response.NamaSiswa = absensi.PesertaDidikRombel.PesertaDidik.Nama
			response.NISN = absensi.PesertaDidikRombel.PesertaDidik.NISN
		}
		if absensi.PesertaDidikRombel.Rombel != nil {
			response.NamaRombel = absensi.PesertaDidikRombel.Rombel.Name
			if absensi.PesertaDidikRombel.Rombel.Kelas != nil {
				response.NamaKelas = absensi.PesertaDidikRombel.Rombel.Kelas.Name
			}
		}
	}

	// Add kegiatan info
	if absensi.KegiatanEkskul != nil {
		response.TanggalKegiatan = absensi.KegiatanEkskul.TanggalKegiatan.Format("2006-01-02")
		response.WaktuMulai = absensi.KegiatanEkskul.WaktuMulai
		response.WaktuSelesai = absensi.KegiatanEkskul.WaktuSelesai
		response.MateriKegiatan = absensi.KegiatanEkskul.MateriKegiatan

		// Add ekstrakurikuler info
		if absensi.KegiatanEkskul.Ekstrakurikuler != nil {
			response.EkstrakurikulerID = absensi.KegiatanEkskul.Ekstrakurikuler.ID
			response.NamaEkstrakurikuler = absensi.KegiatanEkskul.Ekstrakurikuler.Name
		}

		// Add tahun pelajaran info
		if absensi.KegiatanEkskul.TahunPelajaran != nil {
			response.TahunPelajaranID = absensi.KegiatanEkskul.TahunPelajaran.ID
			response.TahunPelajaran = absensi.KegiatanEkskul.TahunPelajaran.TahunPelajaran
		}
	}

	return response, nil
}

// UpdateAbsensiSiswa updates or creates absensi siswa (upsert)
func (s *AbsensiEkskulService) UpdateAbsensiSiswa(req *dtos.AbsensiSiswaUpdateRequest) (*dtos.AbsensiSiswaDetailResponse, error) {
	// Normalize status: "alpha" → "alpa" for consistency (Indonesian spelling)
	status := req.Status
	if status == "alpha" {
		status = "alpa"
	}

	var absensi *models.AbsensiEkskul
	var err error

	// If ID is provided, update existing record
	if req.ID != nil && *req.ID > 0 {
		// Get existing absensi siswa
		absensi, err = s.repo.GetAbsensiSiswaByID(*req.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New("absensi siswa not found")
			}
			return nil, err
		}

		// Update fields
		absensi.Status = status
		absensi.Keterangan = req.Keterangan

		// Save to database
		if err := s.repo.UpdateAbsensiSiswa(absensi); err != nil {
			return nil, fmt.Errorf("failed to update absensi siswa: %w", err)
		}
	} else {
		// Create new record
		// Validate kegiatan exists
		var kegiatan models.KegiatanEkskul
		if err := s.db.First(&kegiatan, req.KegiatanEkskulID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New("kegiatan ekstrakurikuler not found")
			}
			return nil, err
		}

		// Validate peserta didik rombel exists
		var pesertaDidikRombel models.PesertaDidikRombel
		if err := s.db.First(&pesertaDidikRombel, req.PesertaDidikRombelID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New("peserta didik rombel not found")
			}
			return nil, err
		}

		// Check if absensi already exists for this kegiatan and peserta didik
		var existingAbsensi models.AbsensiEkskul
		err := s.db.Where("kegiatan_ekskul_id = ? AND peserta_didik_rombel_id = ?", 
			req.KegiatanEkskulID, req.PesertaDidikRombelID).First(&existingAbsensi).Error
		if err == nil {
			// Already exists, return error
			return nil, errors.New("absensi siswa already exists for this kegiatan, use update with ID")
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		// Create new absensi
		absensi = &models.AbsensiEkskul{
			KegiatanEkskulID:     req.KegiatanEkskulID,
			PesertaDidikRombelID: req.PesertaDidikRombelID,
			Status:               status,
			Keterangan:           req.Keterangan,
		}

		if err := s.db.Create(absensi).Error; err != nil {
			return nil, fmt.Errorf("failed to create absensi siswa: %w", err)
		}
	}

	// Get updated/created data with relations
	finalAbsensi, err := s.repo.GetAbsensiSiswaByID(absensi.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dtos.AbsensiSiswaDetailResponse{
		ID:                   finalAbsensi.ID,
		KegiatanEkskulID:     finalAbsensi.KegiatanEkskulID,
		PesertaDidikRombelID: finalAbsensi.PesertaDidikRombelID,
		Status:               finalAbsensi.Status,
		Keterangan:           finalAbsensi.Keterangan,
	}

	// Add student info
	if finalAbsensi.PesertaDidikRombel != nil {
		if finalAbsensi.PesertaDidikRombel.PesertaDidik != nil {
			response.NamaSiswa = finalAbsensi.PesertaDidikRombel.PesertaDidik.Nama
			response.NIS = finalAbsensi.PesertaDidikRombel.PesertaDidik.NIS
			response.NISN = finalAbsensi.PesertaDidikRombel.PesertaDidik.NISN
		}
		if finalAbsensi.PesertaDidikRombel.Rombel != nil {
			response.NamaRombel = finalAbsensi.PesertaDidikRombel.Rombel.Name
			if finalAbsensi.PesertaDidikRombel.Rombel.Kelas != nil {
				response.NamaKelas = finalAbsensi.PesertaDidikRombel.Rombel.Kelas.Name
			}
		}
	}

	// Add kegiatan info
	if finalAbsensi.KegiatanEkskul != nil {
		response.TanggalKegiatan = finalAbsensi.KegiatanEkskul.TanggalKegiatan.Format("2006-01-02")
		response.WaktuMulai = finalAbsensi.KegiatanEkskul.WaktuMulai
		response.WaktuSelesai = finalAbsensi.KegiatanEkskul.WaktuSelesai
		response.MateriKegiatan = finalAbsensi.KegiatanEkskul.MateriKegiatan

		// Add ekstrakurikuler info
		if finalAbsensi.KegiatanEkskul.Ekstrakurikuler != nil {
			response.EkstrakurikulerID = finalAbsensi.KegiatanEkskul.Ekstrakurikuler.ID
			response.NamaEkstrakurikuler = finalAbsensi.KegiatanEkskul.Ekstrakurikuler.Name
		}

		// Add tahun pelajaran info
		if finalAbsensi.KegiatanEkskul.TahunPelajaran != nil {
			response.TahunPelajaranID = finalAbsensi.KegiatanEkskul.TahunPelajaran.ID
			response.TahunPelajaran = finalAbsensi.KegiatanEkskul.TahunPelajaran.TahunPelajaran
		}
	}

	return response, nil
}


// GetAbsensiPelatih gets pelatih attendance by pelatih and tahun pelajaran
func (s *AbsensiEkskulService) GetAbsensiPelatih(req *dtos.AbsensiPelatihGetRequest) (*dtos.AbsensiPelatihGetResponse, error) {
	// Validate pelatih exists
	pelatih, err := s.repo.GetPelatihByID(req.PelatihID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("pelatih not found")
		}
		return nil, err
	}

	// Validate tahun pelajaran exists
	exists, err := s.repo.CheckTahunPelajaranExists(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("tahun pelajaran not found")
	}

	// Get tahun pelajaran name
	tahunName, err := s.repo.GetTahunPelajaranName(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}

	// Get all kegiatan for this pelatih
	kegiatanList, err := s.repo.GetKegiatanByPelatihAndTahun(req.PelatihID, req.TahunPelajaranID, req.EkstrakurikulerID, req.Bulan, req.Tahun)
	if err != nil {
		return nil, err
	}

	// Build response
	kegiatanDetails := make([]dtos.KegiatanPelatihDetail, 0, len(kegiatanList))
	totalHadir := 0
	totalTidakHadir := 0

	for _, kegiatan := range kegiatanList {
		// Check if pelatih hadir (has record in absensi_pelatih)
		isHadir := len(kegiatan.AbsensiPelatih) > 0
		if isHadir {
			totalHadir++
		} else {
			totalTidakHadir++
		}

		// Count student attendance statistics
		var totalSiswaHadir, totalSiswaSakit, totalSiswaIzin, totalSiswaAlpa int
		for _, absensi := range kegiatan.AbsensiEkskul {
			// Normalize status for counting
			status := absensi.Status
			if status == "alpha" {
				status = "alpa"
			}

			switch status {
			case "hadir":
				totalSiswaHadir++
			case "sakit":
				totalSiswaSakit++
			case "izin":
				totalSiswaIzin++
			case "alpa":
				totalSiswaAlpa++
			}
		}

		// Build kegiatan detail
		kegiatanDetail := dtos.KegiatanPelatihDetail{
			ID:                kegiatan.ID,
			TanggalKegiatan:   kegiatan.TanggalKegiatan.Format("2006-01-02"),
			WaktuMulai:        kegiatan.WaktuMulai,
			WaktuSelesai:      kegiatan.WaktuSelesai,
			MateriKegiatan:    kegiatan.MateriKegiatan,
			FotoKegiatan:      kegiatan.FotoKegiatan,
			TotalSiswaHadir:   totalSiswaHadir,
			TotalSiswaSakit:   totalSiswaSakit,
			TotalSiswaIzin:    totalSiswaIzin,
			TotalSiswaAlpa:    totalSiswaAlpa,
			IsHadir:           isHadir,
		}

		// Add ekstrakurikuler info
		if kegiatan.Ekstrakurikuler != nil {
			kegiatanDetail.EkstrakurikulerID = kegiatan.Ekstrakurikuler.ID
			kegiatanDetail.NamaEkstrakurikuler = kegiatan.Ekstrakurikuler.Name
		}

		kegiatanDetails = append(kegiatanDetails, kegiatanDetail)
	}

	// Build final response
	response := &dtos.AbsensiPelatihGetResponse{
		PelatihID:        pelatih.ID,
		NamaPelatih:      pelatih.Nama,
		Telepon:          pelatih.Telepon,
		TahunPelajaranID: req.TahunPelajaranID,
		TahunPelajaran:   tahunName,
		TotalKegiatan:    len(kegiatanDetails),
		TotalHadir:       totalHadir,
		TotalTidakHadir:  totalTidakHadir,
		Kegiatan:         kegiatanDetails,
	}

	return response, nil
}


// GetKegiatanEkskul gets kegiatan list by ekstrakurikuler and tahun pelajaran
func (s *AbsensiEkskulService) GetKegiatanEkskul(req *dtos.KegiatanEkskulGetRequest) (*dtos.KegiatanEkskulGetResponse, error) {
	// Validate ekstrakurikuler exists
	exists, err := s.repo.CheckEkstrakurikulerExists(req.EkstrakurikulerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("ekstrakurikuler not found")
	}

	// Validate tahun pelajaran exists
	exists, err = s.repo.CheckTahunPelajaranExists(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("tahun pelajaran not found")
	}

	// Get ekstrakurikuler and tahun pelajaran names
	ekskulName, err := s.repo.GetEkstrakurikulerName(req.EkstrakurikulerID)
	if err != nil {
		return nil, err
	}

	tahunName, err := s.repo.GetTahunPelajaranName(req.TahunPelajaranID)
	if err != nil {
		return nil, err
	}

	// Get kegiatan list
	kegiatanList, err := s.repo.GetKegiatanListByEkskulAndTahun(req.EkstrakurikulerID, req.TahunPelajaranID, req.Bulan, req.Tahun)
	if err != nil {
		return nil, err
	}

	// Build response
	kegiatanItems := make([]dtos.KegiatanEkskulItem, 0, len(kegiatanList))
	
	for _, kegiatan := range kegiatanList {
		// Count student attendance statistics
		var totalHadir, totalSakit, totalIzin, totalAlpa int
		for _, absensi := range kegiatan.AbsensiEkskul {
			// Normalize status for counting
			status := absensi.Status
			if status == "alpha" {
				status = "alpa"
			}

			switch status {
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

		// Collect pelatih names who attended
		pelatihNames := make([]string, 0) // Initialize empty slice instead of nil
		for _, absensiPelatih := range kegiatan.AbsensiPelatih {
			if absensiPelatih.Pelatih != nil {
				pelatihNames = append(pelatihNames, absensiPelatih.Pelatih.Nama)
			}
		}

		// Build kegiatan item
		item := dtos.KegiatanEkskulItem{
			ID:                kegiatan.ID,
			TanggalKegiatan:   kegiatan.TanggalKegiatan.Format("2006-01-02"),
			WaktuMulai:        kegiatan.WaktuMulai,
			WaktuSelesai:      kegiatan.WaktuSelesai,
			MateriKegiatan:    kegiatan.MateriKegiatan,
			FotoKegiatan:      kegiatan.FotoKegiatan,
			TotalSiswa:        len(kegiatan.AbsensiEkskul),
			TotalSiswaHadir:   totalHadir,
			TotalSiswaSakit:   totalSakit,
			TotalSiswaIzin:    totalIzin,
			TotalSiswaAlpa:    totalAlpa,
			TotalPelatihHadir: len(kegiatan.AbsensiPelatih),
			PelatihHadir:      pelatihNames,
		}

		kegiatanItems = append(kegiatanItems, item)
	}

	// Build final response
	response := &dtos.KegiatanEkskulGetResponse{
		EkstrakurikulerID:   req.EkstrakurikulerID,
		NamaEkstrakurikuler: ekskulName,
		TahunPelajaranID:    req.TahunPelajaranID,
		TahunPelajaran:      tahunName,
		TotalKegiatan:       len(kegiatanItems),
		Kegiatan:            kegiatanItems,
	}

	return response, nil
}


// UpdateAbsensiPelatih updates pelatih attendance based on status (create if true, delete if false)
func (s *AbsensiEkskulService) UpdateAbsensiPelatih(req *dtos.AbsensiPelatihUpdateRequest) (*dtos.AbsensiPelatihUpdateResponse, error) {
	// Validate kegiatan ekskul exists
	kegiatan, err := s.repo.GetKegiatanEkskulByID(req.KegiatanEkskulID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("kegiatan ekstrakurikuler not found")
		}
		return nil, err
	}

	// Validate pelatih exists
	pelatih, err := s.repo.GetPelatihByID(req.PelatihID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("pelatih not found")
		}
		return nil, err
	}

	// Check if absensi already exists
	existingAbsensi, err := s.repo.CheckAbsensiPelatihExists(req.KegiatanEkskulID, req.PelatihID)
	if err != nil {
		return nil, err
	}

	var isHadir bool
	var message string

	if req.Status {
		// Status true = hadir (create if not exists)
		if existingAbsensi == nil {
			// Record doesn't exist, create it
			newAbsensi := &models.AbsensiPelatihEkskul{
				KegiatanEkskulID: req.KegiatanEkskulID,
				PelatihID:        req.PelatihID,
			}
			if err := s.repo.CreateAbsensiPelatih(newAbsensi); err != nil {
				return nil, fmt.Errorf("failed to create absensi pelatih: %w", err)
			}
			message = "Pelatih marked as hadir"
		} else {
			// Record already exists, no need to create again
			message = "Pelatih already marked as hadir"
		}
		isHadir = true
	} else {
		// Status false = tidak hadir (delete if exists)
		if existingAbsensi != nil {
			// Record exists, delete it
			if err := s.repo.DeleteAbsensiPelatih(existingAbsensi.ID); err != nil {
				return nil, fmt.Errorf("failed to delete absensi pelatih: %w", err)
			}
			message = "Pelatih marked as tidak hadir"
		} else {
			// Record doesn't exist, already tidak hadir
			message = "Pelatih already marked as tidak hadir"
		}
		isHadir = false
	}

	// Build response
	response := &dtos.AbsensiPelatihUpdateResponse{
		KegiatanEkskulID: kegiatan.ID,
		PelatihID:        pelatih.ID,
		NamaPelatih:      pelatih.Nama,
		IsHadir:          isHadir,
		Message:          message,
	}

	return response, nil
}


// UploadFotoKegiatan uploads multiple photos for kegiatan ekstrakurikuler
func (s *AbsensiEkskulService) UploadFotoKegiatan(kegiatanID uint, fotoUrls []string, fotoToDelete []string) (*dtos.FotoKegiatanUploadResponse, error) {
	// Validate kegiatan exists
	kegiatan, err := s.repo.GetKegiatanEkskulByID(kegiatanID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("kegiatan ekstrakurikuler not found")
		}
		return nil, err
	}

	// Get existing foto URLs from database
	var existingFotos []string
	if kegiatan.FotoKegiatan != nil && *kegiatan.FotoKegiatan != "" && *kegiatan.FotoKegiatan != "[]" {
		// Parse existing JSON array
		fotoJSON := *kegiatan.FotoKegiatan
		// Remove brackets and quotes, split by comma
		fotoJSON = fotoJSON[1 : len(fotoJSON)-1] // Remove [ and ]
		if fotoJSON != "" {
			// Simple JSON array parsing (assuming clean data)
			var currentURL string
			inQuote := false
			for _, char := range fotoJSON {
				if char == '"' {
					if inQuote {
						// End of URL
						existingFotos = append(existingFotos, currentURL)
						currentURL = ""
						inQuote = false
					} else {
						// Start of URL
						inQuote = true
					}
				} else if inQuote {
					currentURL += string(char)
				}
			}
		}
	}

	// Filter out URLs to delete
	var finalFotos []string
	for _, existingURL := range existingFotos {
		shouldDelete := false
		for _, deleteURL := range fotoToDelete {
			if existingURL == deleteURL {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			finalFotos = append(finalFotos, existingURL)
		}
	}

	// Add new uploaded URLs
	finalFotos = append(finalFotos, fotoUrls...)

	// Convert foto URLs array to JSON string
	fotoJSON := "[]"
	if len(finalFotos) > 0 {
		// Build JSON array string manually
		fotoJSON = "["
		for i, url := range finalFotos {
			if i > 0 {
				fotoJSON += ","
			}
			// Escape double quotes in URL if any
			escapedURL := fmt.Sprintf("\"%s\"", url)
			fotoJSON += escapedURL
		}
		fotoJSON += "]"
	}

	// Update foto_kegiatan in database
	if err := s.repo.UpdateFotoKegiatan(kegiatanID, fotoJSON); err != nil {
		return nil, fmt.Errorf("failed to update foto kegiatan: %w", err)
	}

	// Build response
	response := &dtos.FotoKegiatanUploadResponse{
		KegiatanEkskulID: kegiatan.ID,
		FotoUrls:         finalFotos,
		TotalFoto:        len(finalFotos),
		DeletedFoto:      len(fotoToDelete),
		UploadedFoto:     len(fotoUrls),
		Message:          fmt.Sprintf("Successfully uploaded %d photo(s) and deleted %d photo(s)", len(fotoUrls), len(fotoToDelete)),
	}

	return response, nil
}

// UpdateKegiatanEkskul updates kegiatan ekstrakurikuler with all fields including photos
func (s *AbsensiEkskulService) UpdateKegiatanEkskul(kegiatanID uint, fotoUrls []string, fotoToDelete []string, tanggalKegiatan, waktuMulai, waktuSelesai, materiKegiatan *string) (*dtos.KegiatanEkskulUpdateResponse, error) {
	// Validate kegiatan exists
	kegiatan, err := s.repo.GetKegiatanEkskulByID(kegiatanID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("kegiatan ekstrakurikuler not found")
		}
		return nil, err
	}

	// Update basic fields if provided
	if tanggalKegiatan != nil && *tanggalKegiatan != "" {
		// Parse date string to time.Time
		parsedDate, err := time.Parse("2006-01-02", *tanggalKegiatan)
		if err != nil {
			return nil, fmt.Errorf("invalid tanggal_kegiatan format, use YYYY-MM-DD: %w", err)
		}
		kegiatan.TanggalKegiatan = parsedDate
	}

	if waktuMulai != nil {
		kegiatan.WaktuMulai = waktuMulai
	}

	if waktuSelesai != nil {
		kegiatan.WaktuSelesai = waktuSelesai
	}

	if materiKegiatan != nil && *materiKegiatan != "" {
		kegiatan.MateriKegiatan = *materiKegiatan
	}

	// Handle foto updates
	var existingFotos []string
	if kegiatan.FotoKegiatan != nil && *kegiatan.FotoKegiatan != "" && *kegiatan.FotoKegiatan != "[]" {
		// Parse existing JSON array
		fotoJSON := *kegiatan.FotoKegiatan
		fotoJSON = fotoJSON[1 : len(fotoJSON)-1] // Remove [ and ]
		if fotoJSON != "" {
			var currentURL string
			inQuote := false
			for _, char := range fotoJSON {
				if char == '"' {
					if inQuote {
						existingFotos = append(existingFotos, currentURL)
						currentURL = ""
						inQuote = false
					} else {
						inQuote = true
					}
				} else if inQuote {
					currentURL += string(char)
				}
			}
		}
	}

	// Filter out URLs to delete
	var finalFotos []string
	for _, existingURL := range existingFotos {
		shouldDelete := false
		for _, deleteURL := range fotoToDelete {
			if existingURL == deleteURL {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			finalFotos = append(finalFotos, existingURL)
		}
	}

	// Add new uploaded URLs
	finalFotos = append(finalFotos, fotoUrls...)

	// Convert foto URLs array to JSON string
	fotoJSON := "[]"
	if len(finalFotos) > 0 {
		fotoJSON = "["
		for i, url := range finalFotos {
			if i > 0 {
				fotoJSON += ","
			}
			escapedURL := fmt.Sprintf("\"%s\"", url)
			fotoJSON += escapedURL
		}
		fotoJSON += "]"
	}

	// Update foto_kegiatan field
	kegiatan.FotoKegiatan = &fotoJSON

	// Save to database
	if err := s.repo.UpdateKegiatanEkskul(kegiatan); err != nil {
		return nil, fmt.Errorf("failed to update kegiatan ekstrakurikuler: %w", err)
	}

	// Build response
	response := &dtos.KegiatanEkskulUpdateResponse{
		ID:              kegiatan.ID,
		TanggalKegiatan: kegiatan.TanggalKegiatan.Format("2006-01-02"),
		WaktuMulai:      kegiatan.WaktuMulai,
		WaktuSelesai:    kegiatan.WaktuSelesai,
		MateriKegiatan:  kegiatan.MateriKegiatan,
		FotoUrls:        finalFotos,
		TotalFoto:       len(finalFotos),
		UploadedFoto:    len(fotoUrls),
		DeletedFoto:     len(fotoToDelete),
		Message:         fmt.Sprintf("Successfully updated kegiatan, uploaded %d photo(s) and deleted %d photo(s)", len(fotoUrls), len(fotoToDelete)),
	}

	return response, nil
}


// GetKegiatanEkskulByID gets kegiatan ekstrakurikuler detail by ID
func (s *AbsensiEkskulService) GetKegiatanEkskulByID(req *dtos.KegiatanEkskulGetByIDRequest) (*dtos.KegiatanEkskulDetailResponse, error) {
	// Get kegiatan with full relations
	kegiatan, err := s.repo.GetKegiatanEkskulByIDWithDetails(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("kegiatan ekstrakurikuler not found")
		}
		return nil, err
	}

	// Parse foto_kegiatan JSON to array
	var fotoUrls []string
	if kegiatan.FotoKegiatan != nil && *kegiatan.FotoKegiatan != "" && *kegiatan.FotoKegiatan != "[]" {
		fotoJSON := *kegiatan.FotoKegiatan
		fotoJSON = fotoJSON[1 : len(fotoJSON)-1] // Remove [ and ]
		if fotoJSON != "" {
			var currentURL string
			inQuote := false
			for _, char := range fotoJSON {
				if char == '"' {
					if inQuote {
						fotoUrls = append(fotoUrls, currentURL)
						currentURL = ""
						inQuote = false
					} else {
						inQuote = true
					}
				} else if inQuote {
					currentURL += string(char)
				}
			}
		}
	}

	// Build absensi siswa details
	absensiSiswaList := make([]dtos.AbsensiSiswaDetail, 0, len(kegiatan.AbsensiEkskul))
	var totalHadir, totalSakit, totalIzin, totalAlpa int

	for _, absensi := range kegiatan.AbsensiEkskul {
		// Normalize status for counting
		status := absensi.Status
		if status == "alpha" {
			status = "alpa"
		}

		// Count by status
		switch status {
		case "hadir":
			totalHadir++
		case "sakit":
			totalSakit++
		case "izin":
			totalIzin++
		case "alpa":
			totalAlpa++
		}

		// Build detail
		detail := dtos.AbsensiSiswaDetail{
			ID:                   absensi.ID,
			PesertaDidikRombelID: absensi.PesertaDidikRombelID,
			Status:               absensi.Status,
			Keterangan:           absensi.Keterangan,
		}

		// Add student info if available
		if absensi.PesertaDidikRombel != nil {
			if absensi.PesertaDidikRombel.PesertaDidik != nil {
				detail.NamaSiswa = absensi.PesertaDidikRombel.PesertaDidik.Nama
				detail.NIS = absensi.PesertaDidikRombel.PesertaDidik.NIS
				detail.NISN = absensi.PesertaDidikRombel.PesertaDidik.NISN
			}
			if absensi.PesertaDidikRombel.Rombel != nil {
				detail.NamaRombel = absensi.PesertaDidikRombel.Rombel.Name
				if absensi.PesertaDidikRombel.Rombel.Kelas != nil {
					detail.NamaKelas = absensi.PesertaDidikRombel.Rombel.Kelas.Name
				}
			}
		}

		absensiSiswaList = append(absensiSiswaList, detail)
	}

	// Build absensi pelatih details
	absensiPelatihList := make([]dtos.AbsensiPelatihDetail, 0, len(kegiatan.AbsensiPelatih))
	for _, absensi := range kegiatan.AbsensiPelatih {
		detail := dtos.AbsensiPelatihDetail{
			ID:        absensi.ID,
			PelatihID: absensi.PelatihID,
		}

		// Add pelatih info if available
		if absensi.Pelatih != nil {
			detail.NamaPelatih = absensi.Pelatih.Nama
			detail.Telepon = absensi.Pelatih.Telepon
		}

		absensiPelatihList = append(absensiPelatihList, detail)
	}

	// Build response
	response := &dtos.KegiatanEkskulDetailResponse{
		ID:                kegiatan.ID,
		TanggalKegiatan:   kegiatan.TanggalKegiatan.Format("2006-01-02"),
		WaktuMulai:        kegiatan.WaktuMulai,
		WaktuSelesai:      kegiatan.WaktuSelesai,
		MateriKegiatan:    kegiatan.MateriKegiatan,
		FotoKegiatan:      fotoUrls,
		AbsensiSiswa:      absensiSiswaList,
		AbsensiPelatih:    absensiPelatihList,
		TotalSiswa:        len(absensiSiswaList),
		TotalSiswaHadir:   totalHadir,
		TotalSiswaSakit:   totalSakit,
		TotalSiswaIzin:    totalIzin,
		TotalSiswaAlpa:    totalAlpa,
		TotalPelatihHadir: len(absensiPelatihList),
	}

	// Add ekstrakurikuler info
	if kegiatan.Ekstrakurikuler != nil {
		response.EkstrakurikulerID = kegiatan.Ekstrakurikuler.ID
		response.NamaEkstrakurikuler = kegiatan.Ekstrakurikuler.Name
	}

	// Add tahun pelajaran info
	if kegiatan.TahunPelajaran != nil {
		response.TahunPelajaranID = kegiatan.TahunPelajaran.ID
		response.TahunPelajaran = kegiatan.TahunPelajaran.TahunPelajaran
	}

	return response, nil
}
