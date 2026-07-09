package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// MutasiSiswaService handles business logic for Mutasi Siswa
type MutasiSiswaService interface {
	CreatePublic(req *dtos.MutasiSiswaCreateRequest, files map[string]*multipart.FileHeader) (*dtos.MutasiSiswaResponse, error)
	GetAllWithFilter(req *dtos.MutasiSiswaGetAllRequest) (*dtos.MutasiSiswaListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.MutasiSiswaResponse, error)
	Update(req *dtos.MutasiSiswaUpdateRequest, files map[string]*multipart.FileHeader) (*dtos.MutasiSiswaResponse, error)
	ExportFormulirPDF(id uint) ([]byte, error)
	Delete(id uint) error
}

type MutasiSiswaServiceImpl struct {
	repository repositories.MutasiSiswaRepository
	r2Storage  *utils.R2Storage
}

// NewMutasiSiswaService creates a new Mutasi Siswa service
func NewMutasiSiswaService(repository repositories.MutasiSiswaRepository, r2Storage *utils.R2Storage) MutasiSiswaService {
	return &MutasiSiswaServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// CreatePublic creates a new Mutasi Siswa from public form
func (s *MutasiSiswaServiceImpl) CreatePublic(req *dtos.MutasiSiswaCreateRequest, files map[string]*multipart.FileHeader) (*dtos.MutasiSiswaResponse, error) {
	// Generate registration number
	nomorPendaftaran, err := s.generateRegistrationNumber(req.TahunPelajaranID, req.Semester)
	if err != nil {
		return nil, err
	}

	// Parse tanggal lahir
	tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
	}

	// Upload files to R2
	var raporPath, akteKelahiranPath, kartuKeluargaPath, sptjmPath *string

	if files["rapor"] != nil {
		path, err := s.r2Storage.UploadFile(files["rapor"], "mutasi-siswa/rapor")
		if err != nil {
			return nil, fmt.Errorf("gagal upload rapor: %w", err)
		}
		raporPath = &path
	}

	if files["akte_kelahiran"] != nil {
		path, err := s.r2Storage.UploadFile(files["akte_kelahiran"], "mutasi-siswa/akte")
		if err != nil {
			return nil, fmt.Errorf("gagal upload akte kelahiran: %w", err)
		}
		akteKelahiranPath = &path
	}

	if files["kartu_keluarga"] != nil {
		path, err := s.r2Storage.UploadFile(files["kartu_keluarga"], "mutasi-siswa/kk")
		if err != nil {
			return nil, fmt.Errorf("gagal upload kartu keluarga: %w", err)
		}
		kartuKeluargaPath = &path
	}

	if files["sptjm"] != nil {
		path, err := s.r2Storage.UploadFile(files["sptjm"], "mutasi-siswa/sptjm")
		if err != nil {
			return nil, fmt.Errorf("gagal upload SPTJM: %w", err)
		}
		sptjmPath = &path
	}

	// Create mutasi siswa record
	data := &models.MutasiSiswa{
		NomorPendaftaran: nomorPendaftaran,
		TahunPelajaranID: req.TahunPelajaranID,
		Semester:         req.Semester,
		NamaLengkap:      req.NamaLengkap,
		NamaPanggilan:    req.NamaPanggilan,
		NISN:             req.NISN,
		TempatLahir:      req.TempatLahir,
		TanggalLahir:     tanggalLahir,
		JenisKelamin:     req.JenisKelamin,
		Agama:            req.Agama,
		GolonganDarah:    req.GolonganDarah,
		AnakKe:           req.AnakKe,
		JumlahSaudara:    req.JumlahSaudara,
		StatusAnak:       req.StatusAnak,
		Alamat:           req.Alamat,
		RT:               req.RT,
		RW:               req.RW,
		Kelurahan:        req.Kelurahan,
		Kecamatan:        req.Kecamatan,
		Kota:             req.Kota,
		Provinsi:         req.Provinsi,
		NamaAyah:         req.NamaAyah,
		NamaIbu:          req.NamaIbu,
		PendidikanAyah:   req.PendidikanAyah,
		PendidikanIbu:    req.PendidikanIbu,
		PekerjaanAyah:    req.PekerjaanAyah,
		PekerjaanIbu:     req.PekerjaanIbu,
		PenghasilanAyah:  req.PenghasilanAyah,
		PenghasilanIbu:   req.PenghasilanIbu,
		NomorHPOrtu:      req.NomorHPOrtu,
		NamaWali:         req.NamaWali,
		PendidikanWali:   req.PendidikanWali,
		HubunganWali:     req.HubunganWali,
		PekerjaanWali:    req.PekerjaanWali,
		NomorHPWali:      req.NomorHPWali,
		PindahanKelas:    req.PindahanKelas,
		AsalSekolah:      req.AsalSekolah,
		NamaAsalSekolah:  req.NamaAsalSekolah,
		Rapor:            raporPath,
		AkteKelahiran:    akteKelahiranPath,
		KartuKeluarga:    kartuKeluargaPath,
		SPTJM:            sptjmPath,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// generateRegistrationNumber generates a new registration number based on tahun pelajaran and semester
func (s *MutasiSiswaServiceImpl) generateRegistrationNumber(tahunPelajaranID, semester int) (string, error) {
	// Get last registration number
	lastNumber, err := s.repository.GetLastRegistrationNumber(tahunPelajaranID, semester)
	if err != nil {
		return "", err
	}

	var nextNumber int
	if lastNumber == "" {
		// First registration
		nextNumber = 1
	} else {
		// Extract the numeric part from the last registration number
		// Format expected: PREFIX-XXX or just XXX (3 digits)
		parts := strings.Split(lastNumber, "-")
		numericPart := parts[len(parts)-1] // Get last part
		
		lastNum, err := strconv.Atoi(numericPart)
		if err != nil {
			// If parsing fails, start from 1
			nextNumber = 1
		} else {
			nextNumber = lastNum + 1
		}
	}

	// Format with 3 digits (001, 002, ..., 999)
	return fmt.Sprintf("%03d", nextNumber), nil
}

// mapToResponse maps MutasiSiswa model to response DTO
func (s *MutasiSiswaServiceImpl) mapToResponse(data *models.MutasiSiswa) *dtos.MutasiSiswaResponse {
	// Convert file keys to public URLs
	var raporURL, akteKelahiranURL, kartuKeluargaURL, sptjmURL *string
	
	if data.Rapor != nil && *data.Rapor != "" {
		url := s.r2Storage.GetPublicURL(*data.Rapor)
		raporURL = &url
	}
	
	if data.AkteKelahiran != nil && *data.AkteKelahiran != "" {
		url := s.r2Storage.GetPublicURL(*data.AkteKelahiran)
		akteKelahiranURL = &url
	}
	
	if data.KartuKeluarga != nil && *data.KartuKeluarga != "" {
		url := s.r2Storage.GetPublicURL(*data.KartuKeluarga)
		kartuKeluargaURL = &url
	}
	
	if data.SPTJM != nil && *data.SPTJM != "" {
		url := s.r2Storage.GetPublicURL(*data.SPTJM)
		sptjmURL = &url
	}

	// Map TahunPelajaran if loaded
	var tahunPelajaran *struct {
		ID             uint   `json:"id"`
		TahunPelajaran string `json:"tahun_pelajaran"`
		Status         string `json:"status"`
	}
	if data.TahunPelajaran != nil {
		tahunPelajaran = &struct {
			ID             uint   `json:"id"`
			TahunPelajaran string `json:"tahun_pelajaran"`
			Status         string `json:"status"`
		}{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
		}
	}
	
	return &dtos.MutasiSiswaResponse{
		ID:               data.ID,
		NomorPendaftaran: data.NomorPendaftaran,
		TahunPelajaranID: data.TahunPelajaranID,
		TahunPelajaran:   tahunPelajaran,
		Semester:         data.Semester,
		NamaLengkap:      data.NamaLengkap,
		NamaPanggilan:    data.NamaPanggilan,
		NISN:             data.NISN,
		TempatLahir:      data.TempatLahir,
		TanggalLahir:     data.TanggalLahir.Format("2006-01-02"),
		JenisKelamin:     data.JenisKelamin,
		Agama:            data.Agama,
		GolonganDarah:    data.GolonganDarah,
		AnakKe:           data.AnakKe,
		JumlahSaudara:    data.JumlahSaudara,
		StatusAnak:       data.StatusAnak,
		Alamat:           data.Alamat,
		RT:               data.RT,
		RW:               data.RW,
		Kelurahan:        data.Kelurahan,
		Kecamatan:        data.Kecamatan,
		Kota:             data.Kota,
		Provinsi:         data.Provinsi,
		NamaAyah:         data.NamaAyah,
		NamaIbu:          data.NamaIbu,
		PendidikanAyah:   data.PendidikanAyah,
		PendidikanIbu:    data.PendidikanIbu,
		PekerjaanAyah:    data.PekerjaanAyah,
		PekerjaanIbu:     data.PekerjaanIbu,
		PenghasilanAyah:  data.PenghasilanAyah,
		PenghasilanIbu:   data.PenghasilanIbu,
		NomorHPOrtu:      data.NomorHPOrtu,
		NamaWali:         data.NamaWali,
		PendidikanWali:   data.PendidikanWali,
		HubunganWali:     data.HubunganWali,
		PekerjaanWali:    data.PekerjaanWali,
		NomorHPWali:      data.NomorHPWali,
		PindahanKelas:    data.PindahanKelas,
		AsalSekolah:      data.AsalSekolah,
		NamaAsalSekolah:  data.NamaAsalSekolah,
		Rapor:            raporURL,
		AkteKelahiran:    akteKelahiranURL,
		KartuKeluarga:    kartuKeluargaURL,
		SPTJM:            sptjmURL,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}


// GetAllWithFilter retrieves all Mutasi Siswa with filters and pagination
func (s *MutasiSiswaServiceImpl) GetAllWithFilter(req *dtos.MutasiSiswaGetAllRequest) (*dtos.MutasiSiswaListWithPaginationResponse, error) {
	// Set default pagination
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Build filter params
	params := repositories.GetMutasiSiswaParams{
		Filter: repositories.GetMutasiSiswaFilter{
			TahunPelajaranID: req.Search.TahunPelajaranID,
			Semester:         req.Search.Semester,
			StartDate:        req.Search.StartDate,
			EndDate:          req.Search.EndDate,
			NamaSiswa:        req.Search.NamaSiswa,
			NISN:             req.Search.NISN,
			TempatLahir:      req.Search.TempatLahir,
			JenisKelamin:     req.Search.JenisKelamin,
			PindahanKelas:    req.Search.PindahanKelas,
		},
		Limit:  limit,
		Offset: offset,
	}

	// Get data from repository
	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	var responses []dtos.MutasiSiswaResponse
	for _, item := range data {
		responses = append(responses, *s.mapToResponse(&item))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dtos.MutasiSiswaListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationMeta{
			Limit:      limit,
			Offset:     offset,
			Page:       page,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}


// GetByID retrieves Mutasi Siswa by ID
func (s *MutasiSiswaServiceImpl) GetByID(id uint) (*dtos.MutasiSiswaResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("mutasi siswa dengan ID %d tidak ditemukan", id)
	}

	return s.mapToResponse(data), nil
}


// Update updates Mutasi Siswa data
func (s *MutasiSiswaServiceImpl) Update(req *dtos.MutasiSiswaUpdateRequest, files map[string]*multipart.FileHeader) (*dtos.MutasiSiswaResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("mutasi siswa dengan ID %d tidak ditemukan", req.ID)
	}

	// Store old file paths for cleanup
	oldRapor := existing.Rapor
	oldAkte := existing.AkteKelahiran
	oldKK := existing.KartuKeluarga
	oldSPTJM := existing.SPTJM

	// Update basic fields if provided
	if req.TahunPelajaranID != nil {
		existing.TahunPelajaranID = *req.TahunPelajaranID
	}
	if req.Semester != nil {
		existing.Semester = *req.Semester
	}
	if req.NamaLengkap != "" {
		existing.NamaLengkap = req.NamaLengkap
	}
	if req.NamaPanggilan != nil {
		existing.NamaPanggilan = req.NamaPanggilan
	}
	if req.NISN != nil {
		existing.NISN = req.NISN
	}
	if req.TempatLahir != "" {
		existing.TempatLahir = req.TempatLahir
	}
	if req.TanggalLahir != "" {
		tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, fmt.Errorf("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
		}
		existing.TanggalLahir = tanggalLahir
	}
	if req.JenisKelamin != "" {
		existing.JenisKelamin = req.JenisKelamin
	}
	if req.Agama != "" {
		existing.Agama = req.Agama
	}
	if req.GolonganDarah != nil {
		existing.GolonganDarah = req.GolonganDarah
	}
	if req.AnakKe != nil {
		existing.AnakKe = req.AnakKe
	}
	if req.JumlahSaudara != nil {
		existing.JumlahSaudara = req.JumlahSaudara
	}
	if req.StatusAnak != nil {
		existing.StatusAnak = req.StatusAnak
	}
	if req.Alamat != "" {
		existing.Alamat = req.Alamat
	}
	if req.RT != nil {
		existing.RT = req.RT
	}
	if req.RW != nil {
		existing.RW = req.RW
	}
	if req.Kelurahan != nil {
		existing.Kelurahan = req.Kelurahan
	}
	if req.Kecamatan != nil {
		existing.Kecamatan = req.Kecamatan
	}
	if req.Kota != nil {
		existing.Kota = req.Kota
	}
	if req.Provinsi != nil {
		existing.Provinsi = req.Provinsi
	}
	if req.NamaAyah != nil {
		existing.NamaAyah = req.NamaAyah
	}
	if req.NamaIbu != nil {
		existing.NamaIbu = req.NamaIbu
	}
	if req.PendidikanAyah != nil {
		existing.PendidikanAyah = req.PendidikanAyah
	}
	if req.PendidikanIbu != nil {
		existing.PendidikanIbu = req.PendidikanIbu
	}
	if req.PekerjaanAyah != nil {
		existing.PekerjaanAyah = req.PekerjaanAyah
	}
	if req.PekerjaanIbu != nil {
		existing.PekerjaanIbu = req.PekerjaanIbu
	}
	if req.PenghasilanAyah != nil {
		existing.PenghasilanAyah = req.PenghasilanAyah
	}
	if req.PenghasilanIbu != nil {
		existing.PenghasilanIbu = req.PenghasilanIbu
	}
	if req.NomorHPOrtu != nil {
		existing.NomorHPOrtu = req.NomorHPOrtu
	}
	if req.NamaWali != nil {
		existing.NamaWali = req.NamaWali
	}
	if req.PendidikanWali != nil {
		existing.PendidikanWali = req.PendidikanWali
	}
	if req.HubunganWali != nil {
		existing.HubunganWali = req.HubunganWali
	}
	if req.PekerjaanWali != nil {
		existing.PekerjaanWali = req.PekerjaanWali
	}
	if req.NomorHPWali != nil {
		existing.NomorHPWali = req.NomorHPWali
	}
	if req.PindahanKelas != nil {
		existing.PindahanKelas = req.PindahanKelas
	}
	if req.AsalSekolah != nil {
		existing.AsalSekolah = req.AsalSekolah
	}
	if req.NamaAsalSekolah != nil {
		existing.NamaAsalSekolah = req.NamaAsalSekolah
	}

	// Update Rapor if provided
	if files["rapor"] != nil {
		// Upload new rapor
		path, err := s.r2Storage.UploadFile(files["rapor"], "mutasi-siswa/rapor")
		if err != nil {
			return nil, fmt.Errorf("gagal upload rapor: %w", err)
		}

		// Delete old rapor if exists
		if oldRapor != nil && *oldRapor != "" {
			_ = s.r2Storage.DeleteFile(*oldRapor)
		}

		existing.Rapor = &path
	}

	// Update Akte Kelahiran if provided
	if files["akte_kelahiran"] != nil {
		// Upload new akte kelahiran
		path, err := s.r2Storage.UploadFile(files["akte_kelahiran"], "mutasi-siswa/akte")
		if err != nil {
			return nil, fmt.Errorf("gagal upload akte kelahiran: %w", err)
		}

		// Delete old akte if exists
		if oldAkte != nil && *oldAkte != "" {
			_ = s.r2Storage.DeleteFile(*oldAkte)
		}

		existing.AkteKelahiran = &path
	}

	// Update Kartu Keluarga if provided
	if files["kartu_keluarga"] != nil {
		// Upload new kartu keluarga
		path, err := s.r2Storage.UploadFile(files["kartu_keluarga"], "mutasi-siswa/kk")
		if err != nil {
			return nil, fmt.Errorf("gagal upload kartu keluarga: %w", err)
		}

		// Delete old KK if exists
		if oldKK != nil && *oldKK != "" {
			_ = s.r2Storage.DeleteFile(*oldKK)
		}

		existing.KartuKeluarga = &path
	}

	// Update SPTJM if provided
	if files["sptjm"] != nil {
		// Upload new SPTJM
		path, err := s.r2Storage.UploadFile(files["sptjm"], "mutasi-siswa/sptjm")
		if err != nil {
			return nil, fmt.Errorf("gagal upload SPTJM: %w", err)
		}

		// Delete old SPTJM if exists
		if oldSPTJM != nil && *oldSPTJM != "" {
			_ = s.r2Storage.DeleteFile(*oldSPTJM)
		}

		existing.SPTJM = &path
	}

	// Save to database
	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Mutasi Siswa and all associated files from R2 storage
func (s *MutasiSiswaServiceImpl) Delete(id uint) error {
	// Get existing data to retrieve file paths
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("mutasi siswa dengan ID %d tidak ditemukan", id)
	}

	// Delete files from R2 storage
	if existing.Rapor != nil && *existing.Rapor != "" {
		_ = s.r2Storage.DeleteFile(*existing.Rapor)
	}
	if existing.AkteKelahiran != nil && *existing.AkteKelahiran != "" {
		_ = s.r2Storage.DeleteFile(*existing.AkteKelahiran)
	}
	if existing.KartuKeluarga != nil && *existing.KartuKeluarga != "" {
		_ = s.r2Storage.DeleteFile(*existing.KartuKeluarga)
	}
	if existing.SPTJM != nil && *existing.SPTJM != "" {
		_ = s.r2Storage.DeleteFile(*existing.SPTJM)
	}

	// Delete from database
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus mutasi siswa: %w", err)
	}

	return nil
}
