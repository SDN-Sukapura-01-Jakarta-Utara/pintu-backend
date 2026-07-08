package services

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
	"crypto/rand"
	"math/big"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type PesertaDidikService interface {
	Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	GetByID(id uint) (*dtos.PesertaDidikResponse, error)
	GetByNIS(nis string) (*dtos.PesertaDidikResponse, error)
	GetAll(limit int, offset int) (*dtos.PesertaDidikListResponse, error)
	GetAllWithFilter(params repositories.GetPesertaDidikParams) (*dtos.PesertaDidikListWithPaginationResponse, error)
	Update(id uint, photo *multipart.FileHeader, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	Delete(id uint) error
	ImportExcel(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error)
	ImportSiswaLulus(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error)
	DownloadTemplate() (*excelize.File, error)
	DownloadTemplateSiswaLulus() (*excelize.File, error)
	ExportDataIndukSiswaExcel(status string) (*excelize.File, error)
	ExportDataIndukSiswaPDF(status string) ([]byte, error)
	ExportPemetaanRombelExcel(rombelID uint, tahunPelajaranID uint) (*excelize.File, error)
	ExportPemetaanRombelPDF(rombelID uint, tahunPelajaranID uint) ([]byte, error)
	DownloadKartuPelajar(pesertaDidikIDs []uint) ([]byte, error)
	GetTotalSiswa() (*dtos.TotalSiswaResponse, error)
	GenerateBarcodeAllPesertaDidik() (*dtos.GenerateBarcodeResponse, error)
	GenerateBarcodePesertaDidikByID(id uint) (*dtos.GenerateBarcodeResponse, error)
}

type PesertaDidikServiceImpl struct {
	repository repositories.PesertaDidikRepository
	r2Storage  *utils.R2Storage
}

// NewPesertaDidikService creates a new PesertaDidik service
func NewPesertaDidikService(repository repositories.PesertaDidikRepository, r2Storage *utils.R2Storage) PesertaDidikService {
	return &PesertaDidikServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new PesertaDidik
func (s *PesertaDidikServiceImpl) Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Check if NIS already exists
	existing, _ := s.repository.GetByNIS(req.NIS)
	if existing != nil {
		return nil, errors.New("NIS sudah ada")
	}

	// Check if NISN already exists
	existingNISN, _ := s.repository.GetByNISN(req.NISN)
	if existingNISN != nil {
		return nil, errors.New("NISN sudah ada")
	}

	// Check if username already exists (if provided)
	if req.Username != "" {
		existing, _ := s.repository.GetByUsername(req.Username)
		if existing != nil {
			return nil, errors.New("username sudah ada")
		}
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal hash password")
		}
		hashedPassword = string(hp)
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Parse tanggal_lahir (YYYY-MM-DD format)
	var tanggalLahir *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalLahir = &t
	}

	// Create peserta didik record
	data := &models.PesertaDidik{
		Nama:               req.Nama,
		NIS:                req.NIS,
		JenisKelamin:       req.JenisKelamin,
		NISN:               req.NISN,
		TempatLahir:        req.TempatLahir,
		TanggalLahir:       tanggalLahir,
		NIK:                req.NIK,
		Agama:              req.Agama,
		Alamat:             req.Alamat,
		RT:                 req.RT,
		RW:                 req.RW,
		Kelurahan:          req.Kelurahan,
		Kecamatan:          req.Kecamatan,
		KodePos:            req.KodePos,
		NamaAyah:           req.NamaAyah,
		NamaIbu:            req.NamaIbu,
		Status:             status,
		Username:           req.Username,
		Password:           hashedPassword,
		CreatedByID:        &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	// Assign roles
	if err := s.repository.AssignRoles(data.ID, req.RoleIDs); err != nil {
		return nil, err
	}

	// Reload data with roles preloaded
	dataWithRoles, err := s.repository.GetByIDWithDetails(data.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(dataWithRoles), nil
}

// GetByID retrieves PesertaDidik by ID with details
func (s *PesertaDidikServiceImpl) GetByID(id uint) (*dtos.PesertaDidikResponse, error) {
	data, err := s.repository.GetByIDWithDetails(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetByNIS retrieves PesertaDidik by NIS
func (s *PesertaDidikServiceImpl) GetByNIS(nis string) (*dtos.PesertaDidikResponse, error) {
	data, err := s.repository.GetByNIS(nis)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all PesertaDidik
func (s *PesertaDidikServiceImpl) GetAll(limit int, offset int) (*dtos.PesertaDidikListResponse, error) {
	// Set default limit and offset
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	data, total, err := s.repository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.PesertaDidikResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.PesertaDidikListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves PesertaDidik with filters and pagination
func (s *PesertaDidikServiceImpl) GetAllWithFilter(params repositories.GetPesertaDidikParams) (*dtos.PesertaDidikListWithPaginationResponse, error) {
	// Validate and set default limit and offset
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.PesertaDidikResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.PesertaDidikListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationInfo{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       (params.Offset / params.Limit) + 1,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// Update updates PesertaDidik
func (s *PesertaDidikServiceImpl) Update(id uint, photo *multipart.FileHeader, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("peserta didik tidak ditemukan")
	}

	// Upload photo if provided
	if photo != nil {
		// Validate photo file
		if photo.Size > 5*1024*1024 { // 5MB
			return nil, errors.New("ukuran photo maksimal 5MB")
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
			"image/webp": true,
		}
		contentType := photo.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, errors.New("hanya file gambar yang diperbolehkan (jpeg, jpg, png, webp)")
		}

		// Delete old photo if exists
		if existing.Photo != "" {
			_ = s.r2Storage.DeleteFile(existing.Photo) // Ignore error if file doesn't exist
		}

		// Upload photo to R2
		fileKey, err := s.r2Storage.UploadFile(photo, "peserta-didik")
		if err != nil {
			return nil, err
		}
		existing.Photo = fileKey
	}

	// Update basic fields if provided
	if req.Nama != "" {
		existing.Nama = req.Nama
	}
	if req.JenisKelamin != "" {
		existing.JenisKelamin = req.JenisKelamin
	}
	if req.NIS != "" && req.NIS != existing.NIS {
		// Check if new NIS already exists
		existingNIS, _ := s.repository.GetByNIS(req.NIS)
		if existingNIS != nil && existingNIS.ID != existing.ID {
			return nil, errors.New("NIS sudah ada")
		}
		existing.NIS = req.NIS
	}
	if req.NISN != "" && req.NISN != existing.NISN {
		// Check if new NISN already exists
		existingNISN, _ := s.repository.GetByNISN(req.NISN)
		if existingNISN != nil && existingNISN.ID != existing.ID {
			return nil, errors.New("NISN sudah ada")
		}
		existing.NISN = req.NISN
	}
	if req.TempatLahir != "" {
		existing.TempatLahir = req.TempatLahir
	}
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
		}
		existing.TanggalLahir = &t
	}
	if req.NIK != "" {
		existing.NIK = req.NIK
	}
	if req.Agama != "" {
		existing.Agama = req.Agama
	}
	if req.Alamat != "" {
		existing.Alamat = req.Alamat
	}
	if req.RT != "" {
		existing.RT = req.RT
	}
	if req.RW != "" {
		existing.RW = req.RW
	}
	if req.Kelurahan != "" {
		existing.Kelurahan = req.Kelurahan
	}
	if req.Kecamatan != "" {
		existing.Kecamatan = req.Kecamatan
	}
	if req.KodePos != "" {
		existing.KodePos = req.KodePos
	}
	if req.NamaAyah != "" {
		existing.NamaAyah = req.NamaAyah
	}
	if req.NamaIbu != "" {
		existing.NamaIbu = req.NamaIbu
	}
	if req.Username != "" && req.Username != existing.Username {
		// Check if new username already exists
		existingUsername, _ := s.repository.GetByUsername(req.Username)
		if existingUsername != nil && existingUsername.ID != existing.ID {
			return nil, errors.New("username sudah ada")
		}
		existing.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal hash password")
		}
		existing.Password = string(hashedPassword)
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	existing.UpdatedByID = &userID

	// Update record
	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	// Update roles only if RoleIDs is provided (even if empty to clear roles)
	// Note: Controller will only set RoleIDs if the field was sent in the request
	if req.RoleIDs != nil {
		if len(*req.RoleIDs) > 0 {
			if err := s.repository.AssignRoles(existing.ID, *req.RoleIDs); err != nil {
				return nil, errors.New("gagal mengupdate roles peserta didik")
			}
		} else {
			// If empty array is sent explicitly, remove all roles
			if err := s.repository.RemoveRoles(existing.ID); err != nil {
				return nil, errors.New("gagal menghapus roles peserta didik")
			}
		}
	}
	// If RoleIDs is nil, don't touch roles at all

	// Reload data with roles preloaded
	dataWithRoles, err := s.repository.GetByIDWithDetails(existing.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(dataWithRoles), nil
}

// Delete deletes PesertaDidik by ID
func (s *PesertaDidikServiceImpl) Delete(id uint) error {
	// Validate peserta didik exists before delete
	data, err := s.repository.GetByID(id)
	if err != nil || data == nil {
		return errors.New("peserta didik tidak ditemukan atau sudah dihapus")
	}

	return s.repository.Delete(id)
}

// ImportExcel imports PesertaDidik data from an Excel file
func (s *PesertaDidikServiceImpl) ImportExcel(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("gagal membuka file excel")
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, errors.New("gagal membaca Sheet1")
	}

	successCount := 0
	failedCount := 0
	var importErrors []dtos.ImportExcelRowError

	for i, row := range rows {
		// Skip header row
		if i == 0 {
			continue
		}

		rowNum := i + 1

		// Skip empty rows
		if len(row) == 0 {
			continue
		}

		// Helper to safely get column value
		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		username := getCol(0)
		password := getCol(1)
		namaLengkap := getCol(2)
		nis := getCol(3)
		nisn := getCol(4)
		jenisKelamin := getCol(5)
		tempatLahir := getCol(6)
		tanggalLahirStr := getCol(7)
		nik := getCol(8)
		agama := getCol(9)
		alamat := getCol(10)
		rt := getCol(11)
		rw := getCol(12)
		kelurahan := getCol(13)
		kecamatan := getCol(14)
		kodePos := getCol(15)
		namaAyah := getCol(16)
		namaIbu := getCol(17)
		roleIDStr := getCol(18)
		status := getCol(19)

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: "gagal hash password"})
			continue
		}

		// Parse tanggal_lahir
		var tanggalLahir *time.Time
		if tanggalLahirStr != "" {
			t, err := time.Parse("2006-01-02", tanggalLahirStr)
			if err != nil {
				failedCount++
				importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: "format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
				continue
			}
			tanggalLahir = &t
		}

		// Check duplicate NIS
		existingNIS, _ := s.repository.GetByNIS(nis)
		if existingNIS != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("NIS '%s' sudah ada", nis)})
			continue
		}

		// Check duplicate NISN
		existingNISN, _ := s.repository.GetByNISN(nisn)
		if existingNISN != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("NISN '%s' sudah ada", nisn)})
			continue
		}

		// Check duplicate username
		if username != "" {
			existingUsername, _ := s.repository.GetByUsername(username)
			if existingUsername != nil {
				failedCount++
				importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("username '%s' sudah ada", username)})
				continue
			}
		}

		// Parse role_id string "[1,2,3]" into []uint
		var roleIDs []uint
		if roleIDStr != "" {
			trimmed := strings.Trim(roleIDStr, "[]")
			if trimmed != "" {
				parts := strings.Split(trimmed, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					id, err := strconv.ParseUint(part, 10, 32)
					if err != nil {
						failedCount++
						importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("role_id '%s' tidak valid", roleIDStr)})
						continue
					}
					roleIDs = append(roleIDs, uint(id))
				}
			}
		}

		// Set default status
		if status == "" {
			status = "active"
		}

		// Create the PesertaDidik record
		data := &models.PesertaDidik{
			Nama:             namaLengkap,
			NIS:              nis,
			NISN:             nisn,
			JenisKelamin:     jenisKelamin,
			TempatLahir:      tempatLahir,
			TanggalLahir:     tanggalLahir,
			NIK:              nik,
			Agama:            agama,
			Alamat:           alamat,
			RT:               rt,
			RW:               rw,
			Kelurahan:        kelurahan,
			Kecamatan:        kecamatan,
			KodePos:          kodePos,
			NamaAyah:         namaAyah,
			NamaIbu:          namaIbu,
			Status:           status,
			Username:         username,
			Password:         string(hashedPassword),
			CreatedByID:      &userID,
		}

		if err := s.repository.Create(data); err != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("gagal menyimpan data: %s", err.Error())})
			continue
		}

		// Assign roles
		if len(roleIDs) > 0 {
			if err := s.repository.AssignRoles(data.ID, roleIDs); err != nil {
				failedCount++
				importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("gagal assign roles: %s", err.Error())})
				continue
			}
		}

		successCount++
	}

	return &dtos.ImportExcelResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Errors:       importErrors,
	}, nil
}

// DownloadTemplate generates an Excel template for PesertaDidik import
func (s *PesertaDidikServiceImpl) DownloadTemplate() (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sheet1"

	headers := []string{
		"username", "password", "nama_lengkap", "nis", "nisn",
		"jenis_kelamin", "tempat_lahir", "tanggal_lahir", "nik", "agama",
		"alamat", "rt", "rw", "kelurahan", "kecamatan", "kode_pos",
		"nama_ayah", "nama_ibu", "role_id", "status",
		"catatan",
	}

	// Create header style with bold font and blue background
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style header")
	}

	// Create note style with yellow background
	noteStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Italic: true,
			Color:  "#666666",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFF2CC"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style note")
	}

	// Set headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set example values in row 2
	examples := []string{
		"siswa01", "password123", "Ahmad Fauzi", "001234", "1234567890",
		"L", "Jakarta", "2015-07-15", "3171234567890001", "Islam",
		"Jl. Merdeka No. 1", "001", "002", "Sukapura", "Cilincing", "14140",
		"Budi Santoso", "Siti Aminah", "[1,2]", "active",
		"Untuk role_id, buka menu Master Data > tab Role untuk melihat ID role",
	}

	for i, val := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Set note style on catatan column row 2 (column U = col 21)
	f.SetCellStyle(sheetName, "U2", "U2", noteStyle)

	// Add dropdown for agama (column J = col 10, rows 2-1000)
	agamaValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       true,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	agamaValidation.Sqref = "J2:J1000"
	agamaValidation.SetDropList([]string{"Islam", "Kristen", "Katolik", "Hindu", "Buddha", "Konghucu"})
	f.AddDataValidation(sheetName, agamaValidation)

	// Add dropdown for jenis_kelamin (column F = col 6, rows 2-1000)
	jkValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       true,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	jkValidation.Sqref = "F2:F1000"
	jkValidation.SetDropList([]string{"L", "P"})
	f.AddDataValidation(sheetName, jkValidation)

	// Add dropdown for status (column T = col 20, rows 2-1000)
	statusValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       true,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	statusValidation.Sqref = "T2:T1000"
	statusValidation.SetDropList([]string{"active", "inactive"})
	f.AddDataValidation(sheetName, statusValidation)

	// Set column widths
	colWidths := map[string]float64{
		"A": 15, "B": 15, "C": 20, "D": 12, "E": 15,
		"F": 15, "G": 15, "H": 15, "I": 20, "J": 12,
		"K": 25, "L": 8, "M": 8, "N": 15, "O": 15, "P": 10,
		"Q": 20, "R": 20, "S": 12, "T": 10, "U": 55,
	}

	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	return f, nil
}

// mapToResponse maps model to DTO response
func (s *PesertaDidikServiceImpl) mapToResponse(data *models.PesertaDidik) *dtos.PesertaDidikResponse {
	// Map roles
	roles := make([]dtos.RoleResponse, len(data.Roles))
	for i, role := range data.Roles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			SystemID:    role.SystemID,
			System:      system,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			CreatedByID: role.CreatedByID,
			UpdatedByID: role.UpdatedByID,
		}
	}

	tanggalLahirStr := ""
	if data.TanggalLahir != nil {
		tanggalLahirStr = data.TanggalLahir.Format("2006-01-02")
	}

	barcodeGeneratedAtStr := ""
	if data.BarcodeGeneratedAt != nil {
		barcodeGeneratedAtStr = data.BarcodeGeneratedAt.Format("2006-01-02T15:04:05Z")
	}

	// Generate public URL for photo if exists
	photoURL := ""
	if data.Photo != "" {
		photoURL = s.r2Storage.GetPublicURL(data.Photo)
	}

	return &dtos.PesertaDidikResponse{
		ID:               data.ID,
		Nama:             data.Nama,
		NIS:              data.NIS,
		JenisKelamin:     data.JenisKelamin,
		NISN:             data.NISN,
		TempatLahir:      data.TempatLahir,
		TanggalLahir:     tanggalLahirStr,
		NIK:              data.NIK,
		Agama:            data.Agama,
		Alamat:           data.Alamat,
		RT:               data.RT,
		RW:               data.RW,
		Kelurahan:        data.Kelurahan,
		Kecamatan:        data.Kecamatan,
		KodePos:          data.KodePos,
		NamaAyah:         data.NamaAyah,
		NamaIbu:          data.NamaIbu,
		Status:           data.Status,
		Username:         data.Username,
		Photo:            photoURL,
		Barcode:          data.Barcode,
		BarcodeGeneratedAt: barcodeGeneratedAtStr,
		Roles:            roles,
		CreatedAt:        data.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedByID:      data.CreatedByID,
		UpdatedByID:      data.UpdatedByID,
	}
}

// GetTotalSiswa retrieves total count of peserta didik with active tahun pelajaran
func (s *PesertaDidikServiceImpl) GetTotalSiswa() (*dtos.TotalSiswaResponse, error) {
	total, err := s.repository.GetTotalSiswaByActiveTahunPelajaran()
	if err != nil {
		return nil, err
	}

	return &dtos.TotalSiswaResponse{
		TotalSiswa: total,
	}, nil
}

// generateRandomString generates a random alphanumeric string of specified length
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// generateBarcode generates a unique barcode for a peserta didik
// Format: {NIS}-{10_RANDOM}
// Example: 001234-A7B9X2K4M1
func generateBarcode(nis string) string {
	random := generateRandomString(10)
	return fmt.Sprintf("%s-%s", nis, random)
}

// GenerateBarcodeAllPesertaDidik generates barcodes for all peserta didik
func (s *PesertaDidikServiceImpl) GenerateBarcodeAllPesertaDidik() (*dtos.GenerateBarcodeResponse, error) {
	// Get all peserta didik with status active
	pesertaDidikList, err := s.repository.GetAllPesertaDidikActive()
	if err != nil {
		return nil, errors.New("gagal mengambil data peserta didik")
	}

	if len(pesertaDidikList) == 0 {
		return nil, errors.New("tidak ada peserta didik aktif")
	}

	totalGenerated := 0
	var errorMessages []string

	// Generate barcode for each peserta didik
	for _, pd := range pesertaDidikList {
		// Skip jika sudah punya barcode
		if pd.Barcode != "" {
			continue
		}
		
		barcode := generateBarcode(pd.NIS)
		
		// Update barcode
		if err := s.repository.UpdateBarcode(pd.ID, barcode); err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Gagal generate barcode untuk NIS %s: %s", pd.NIS, err.Error()))
			continue
		}
		
		totalGenerated++
	}

	message := fmt.Sprintf("Berhasil generate barcode untuk %d peserta didik", totalGenerated)
	if len(errorMessages) > 0 {
		message += fmt.Sprintf(", %d gagal", len(errorMessages))
	}

	return &dtos.GenerateBarcodeResponse{
		TotalGenerated: totalGenerated,
		Message:        message,
		Errors:         errorMessages,
	}, nil
}

// GenerateBarcodePesertaDidikByID generates or regenerates barcode for a specific peserta didik by ID
func (s *PesertaDidikServiceImpl) GenerateBarcodePesertaDidikByID(id uint) (*dtos.GenerateBarcodeResponse, error) {
	// Get peserta didik by ID
	pesertaDidik, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("peserta didik tidak ditemukan")
	}

	// Generate barcode
	barcode := generateBarcode(pesertaDidik.NIS)
	
	// Update barcode (will overwrite if exists)
	if err := s.repository.UpdateBarcode(pesertaDidik.ID, barcode); err != nil {
		return nil, fmt.Errorf("gagal generate barcode untuk NIS %s: %s", pesertaDidik.NIS, err.Error())
	}

	message := fmt.Sprintf("Berhasil generate barcode untuk peserta didik %s (NIS: %s)", pesertaDidik.Nama, pesertaDidik.NIS)
	
	return &dtos.GenerateBarcodeResponse{
		TotalGenerated: 1,
		Message:        message,
		Errors:         []string{},
	}, nil
}

// ExportDataIndukSiswaExcel exports data induk siswa to Excel file
func (s *PesertaDidikServiceImpl) ExportDataIndukSiswaExcel(status string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Data Induk Siswa"
	
	// Rename default sheet
	f.SetSheetName("Sheet1", sheetName)

	// Headers
	headers := []string{
		"NAMA", "NIS", "JENIS KELAMIN", "NISN", "TEMPAT LAHIR", "TANGGAL LAHIR",
		"NIK", "AGAMA", "ALAMAT", "RT", "RW", "KELURAHAN", "KODE POS",
		"NAMA AYAH", "NAMA IBU",
	}

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style header")
	}

	// Set headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, // NAMA
		"B": 12, // NIS
		"C": 15, // JENIS KELAMIN
		"D": 15, // NISN
		"E": 20, // TEMPAT LAHIR
		"F": 15, // TANGGAL LAHIR
		"G": 20, // NIK
		"H": 12, // AGAMA
		"I": 35, // ALAMAT
		"J": 8,  // RT
		"K": 8,  // RW
		"L": 20, // KELURAHAN
		"M": 12, // KODE POS
		"N": 25, // NAMA AYAH
		"O": 25, // NAMA IBU
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Get data from repository
	var siswaList []models.PesertaDidik
	var err2 error
	
	if status != "" {
		// Filter by status
		siswaList, _, err2 = s.repository.GetAll(10000, 0) // Get all with high limit
		if err2 != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err2.Error())
		}
		// Filter by status manually
		filtered := []models.PesertaDidik{}
		for _, siswa := range siswaList {
			if siswa.Status == status {
				filtered = append(filtered, siswa)
			}
		}
		siswaList = filtered
	} else {
		// Get all
		siswaList, _, err2 = s.repository.GetAll(10000, 0)
		if err2 != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err2.Error())
		}
	}

	// Create data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style data")
	}

	// Fill data starting from row 2
	for idx, siswa := range siswaList {
		row := idx + 2
		
		// Format tanggal lahir
		tanggalLahir := ""
		if siswa.TanggalLahir != nil {
			tanggalLahir = siswa.TanggalLahir.Format("02-01-2006")
		}
		
		// Jenis kelamin
		jenisKelamin := ""
		if siswa.JenisKelamin == "L" {
			jenisKelamin = "Laki-laki"
		} else if siswa.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}

		data := []interface{}{
			siswa.Nama,
			siswa.NIS,
			jenisKelamin,
			siswa.NISN,
			siswa.TempatLahir,
			tanggalLahir,
			siswa.NIK,
			siswa.Agama,
			siswa.Alamat,
			siswa.RT,
			siswa.RW,
			siswa.Kelurahan,
			siswa.KodePos,
			siswa.NamaAyah,
			siswa.NamaIbu,
		}

		for colIdx, value := range data {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Freeze first row (header)
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	return f, nil
}

// ExportDataIndukSiswaPDF exports data induk siswa to PDF file
func (s *PesertaDidikServiceImpl) ExportDataIndukSiswaPDF(status string) ([]byte, error) {
	// Get data from repository
	var siswaList []models.PesertaDidik
	var err error
	
	if status != "" {
		// Filter by status
		siswaList, _, err = s.repository.GetAll(10000, 0) // Get all with high limit
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
		// Filter by status manually
		filtered := []models.PesertaDidik{}
		for _, siswa := range siswaList {
			if siswa.Status == status {
				filtered = append(filtered, siswa)
			}
		}
		siswaList = filtered
	} else {
		// Get all
		siswaList, _, err = s.repository.GetAll(10000, 0)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
	}

	// Create PDF landscape A4 with custom margins
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5) // Left, Top, Right margins dikecilkan
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()
	
	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "DATA INDUK SISWA", "", 1, "C", false, 0, "")
	pdf.Ln(3)
	
	// Table headers
	headers := []string{
		"NO", "NAMA", "NIS", "JENIS KELAMIN", "NISN", "TEMPAT LAHIR", 
		"TANGGAL LAHIR", "NIK", "AGAMA", "ALAMAT", "RT", "RW", 
		"KELURAHAN", "KODE POS", "NAMA AYAH", "NAMA IBU",
	}
	
	// Column widths (total should be ~277mm for landscape A4)
	colWidths := []float64{
		7, 30, 13, 20, 16, 20, 20, 25, 11, 28, 7, 7, 22, 14, 23, 23,
	}
	// NO(7), NAMA(30↑), NIS(13), JENIS KELAMIN(20↑), NISN(16), TEMPAT LAHIR(20), 
	// TANGGAL LAHIR(20↑), NIK(25↑), AGAMA(11), ALAMAT(28), RT(7), RW(7), 
	// KELURAHAN(22↑), KODE POS(14↑), NAMA AYAH(23), NAMA IBU(23)
	
	// Header background color (abu tua)
	pdf.SetFillColor(80, 80, 80) // Dark gray
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Arial", "B", 7)
	
	// Draw header row
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	
	// Reset text color for data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 6)
	
	// Draw data rows
	for idx, siswa := range siswaList {
		// Check if we need a new page
		if pdf.GetY() > 180 {
			pdf.AddPage()
			// Redraw header on new page
			pdf.SetFillColor(80, 80, 80)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "B", 7)
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Arial", "", 6)
		}
		
		// Format data
		tanggalLahir := ""
		if siswa.TanggalLahir != nil {
			tanggalLahir = siswa.TanggalLahir.Format("02-01-2006")
		}
		
		jenisKelamin := ""
		if siswa.JenisKelamin == "L" {
			jenisKelamin = "Laki-laki"
		} else if siswa.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}
		
		// Alternate row colors (light gray for even rows)
		if idx%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		
		rowData := []string{
			fmt.Sprintf("%d", idx+1),
			siswa.Nama,
			siswa.NIS,
			jenisKelamin,
			siswa.NISN,
			siswa.TempatLahir,
			tanggalLahir,
			siswa.NIK,
			siswa.Agama,
			siswa.Alamat,
			siswa.RT,
			siswa.RW,
			siswa.Kelurahan,
			siswa.KodePos,
			siswa.NamaAyah,
			siswa.NamaIbu,
		}
		
		// Calculate max height needed for this row (with wrap text support)
		lineHeight := 4.0
		maxLines := 1
		for i, data := range rowData {
			// Estimate how many lines needed based on string length and column width
			strWidth := pdf.GetStringWidth(data)
			if strWidth > colWidths[i]-2 { // -2 for padding
				lines := int(strWidth/(colWidths[i]-2)) + 1
				if lines > maxLines {
					maxLines = lines
				}
			}
		}
		rowHeight := float64(maxLines) * lineHeight
		
		// Save current Y position
		startY := pdf.GetY()
		currentX := pdf.GetX()
		
		// Draw cells with MultiCell for wrap text
		for i, data := range rowData {
			x := currentX
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}
			
			// Set position for this cell
			pdf.SetXY(x, startY)
			
			// Draw cell border first
			pdf.Rect(x, startY, colWidths[i], rowHeight, "D")
			
			// Fill background
			pdf.SetXY(x, startY)
			pdf.CellFormat(colWidths[i], rowHeight, "", "", 0, "", true, 0, "")
			
			// Draw text with MultiCell (supports wrap)
			pdf.SetXY(x+0.5, startY+0.5) // Small padding
			pdf.MultiCell(colWidths[i]-1, lineHeight, data, "", "L", false)
		}
		
		// Move to next row
		pdf.SetXY(currentX, startY+rowHeight)
	}
	
	// Output PDF to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %s", err.Error())
	}
	
	return buf.Bytes(), nil
}

// ExportPemetaanRombelExcel exports pemetaan rombel to Excel file
func (s *PesertaDidikServiceImpl) ExportPemetaanRombelExcel(rombelID uint, tahunPelajaranID uint) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Pemetaan Rombel"
	
	// Rename default sheet
	f.SetSheetName("Sheet1", sheetName)

	// Headers - sama dengan data induk siswa
	headers := []string{
		"NAMA", "NIS", "JENIS KELAMIN", "NISN", "TEMPAT LAHIR", "TANGGAL LAHIR",
		"NIK", "AGAMA", "ALAMAT", "RT", "RW", "KELURAHAN", "KODE POS",
		"NAMA AYAH", "NAMA IBU",
	}

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style header")
	}

	// Set headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, // NAMA
		"B": 12, // NIS
		"C": 15, // JENIS KELAMIN
		"D": 15, // NISN
		"E": 20, // TEMPAT LAHIR
		"F": 15, // TANGGAL LAHIR
		"G": 20, // NIK
		"H": 12, // AGAMA
		"I": 35, // ALAMAT
		"J": 8,  // RT
		"K": 8,  // RW
		"L": 20, // KELURAHAN
		"M": 12, // KODE POS
		"N": 25, // NAMA AYAH
		"O": 25, // NAMA IBU
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Get data from repository by rombel_id and/or tahun_pelajaran_id
	var siswaList []models.PesertaDidik
	var err2 error
	
	if rombelID > 0 || tahunPelajaranID > 0 {
		// Filter by rombel_id and/or tahun_pelajaran_id
		siswaList, err2 = s.repository.GetPesertaDidikByRombelAndTahunPelajaran(rombelID, tahunPelajaranID)
		if err2 != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err2.Error())
		}
	} else {
		// Get all
		siswaList, _, err2 = s.repository.GetAll(10000, 0)
		if err2 != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err2.Error())
		}
	}

	// Create data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		return nil, errors.New("gagal membuat style data")
	}

	// Fill data starting from row 2
	for idx, siswa := range siswaList {
		row := idx + 2
		
		// Format tanggal lahir
		tanggalLahir := ""
		if siswa.TanggalLahir != nil {
			tanggalLahir = siswa.TanggalLahir.Format("02-01-2006")
		}
		
		// Jenis kelamin
		jenisKelamin := ""
		if siswa.JenisKelamin == "L" {
			jenisKelamin = "Laki-laki"
		} else if siswa.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}

		data := []interface{}{
			siswa.Nama,
			siswa.NIS,
			jenisKelamin,
			siswa.NISN,
			siswa.TempatLahir,
			tanggalLahir,
			siswa.NIK,
			siswa.Agama,
			siswa.Alamat,
			siswa.RT,
			siswa.RW,
			siswa.Kelurahan,
			siswa.KodePos,
			siswa.NamaAyah,
			siswa.NamaIbu,
		}

		for colIdx, value := range data {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Freeze first row (header)
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	return f, nil
}

// ExportPemetaanRombelPDF exports pemetaan rombel to PDF file
func (s *PesertaDidikServiceImpl) ExportPemetaanRombelPDF(rombelID uint, tahunPelajaranID uint) ([]byte, error) {
	// Get rombel name if rombelID is provided
	rombelName := ""
	if rombelID > 0 {
		rombel, err := s.repository.GetRombelByID(rombelID)
		if err == nil && rombel != nil {
			rombelName = rombel.Name
		}
	}
	
	// Get data from repository by rombel_id and/or tahun_pelajaran_id
	var siswaList []models.PesertaDidik
	var err error
	
	if rombelID > 0 || tahunPelajaranID > 0 {
		// Filter by rombel_id and/or tahun_pelajaran_id
		siswaList, err = s.repository.GetPesertaDidikByRombelAndTahunPelajaran(rombelID, tahunPelajaranID)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
	} else {
		// Get all
		siswaList, _, err = s.repository.GetAll(10000, 0)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
	}

	// Create PDF landscape A4 with custom margins
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()
	
	// Title dengan nama rombel
	pdf.SetFont("Arial", "B", 16)
	title := "DAFTAR SISWA"
	if rombelName != "" {
		title = fmt.Sprintf("DAFTAR SISWA KELAS %s", rombelName)
	}
	pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")
	pdf.Ln(3)
	
	// Table headers
	headers := []string{
		"NO", "NAMA", "NIS", "JENIS KELAMIN", "NISN", "TEMPAT LAHIR", 
		"TANGGAL LAHIR", "NIK", "AGAMA", "ALAMAT", "RT", "RW", 
		"KELURAHAN", "KODE POS", "NAMA AYAH", "NAMA IBU",
	}
	
	// Column widths
	colWidths := []float64{
		7, 30, 13, 20, 16, 20, 20, 25, 11, 28, 7, 7, 22, 14, 23, 23,
	}
	
	// Header background color (abu tua)
	pdf.SetFillColor(80, 80, 80)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 7)
	
	// Draw header row
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	
	// Reset text color for data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 6)
	
	// Draw data rows
	for idx, siswa := range siswaList {
		// Check if we need a new page
		if pdf.GetY() > 180 {
			pdf.AddPage()
			// Redraw header on new page
			pdf.SetFillColor(80, 80, 80)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "B", 7)
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Arial", "", 6)
		}
		
		// Format data
		tanggalLahir := ""
		if siswa.TanggalLahir != nil {
			tanggalLahir = siswa.TanggalLahir.Format("02-01-2006")
		}
		
		jenisKelamin := ""
		if siswa.JenisKelamin == "L" {
			jenisKelamin = "Laki-laki"
		} else if siswa.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}
		
		// Alternate row colors
		if idx%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		
		rowData := []string{
			fmt.Sprintf("%d", idx+1),
			siswa.Nama,
			siswa.NIS,
			jenisKelamin,
			siswa.NISN,
			siswa.TempatLahir,
			tanggalLahir,
			siswa.NIK,
			siswa.Agama,
			siswa.Alamat,
			siswa.RT,
			siswa.RW,
			siswa.Kelurahan,
			siswa.KodePos,
			siswa.NamaAyah,
			siswa.NamaIbu,
		}
		
		// Calculate max height needed for this row (with wrap text support)
		lineHeight := 4.0
		maxLines := 1
		for i, data := range rowData {
			strWidth := pdf.GetStringWidth(data)
			if strWidth > colWidths[i]-2 {
				lines := int(strWidth/(colWidths[i]-2)) + 1
				if lines > maxLines {
					maxLines = lines
				}
			}
		}
		rowHeight := float64(maxLines) * lineHeight
		
		// Save current Y position
		startY := pdf.GetY()
		currentX := pdf.GetX()
		
		// Draw cells with MultiCell for wrap text
		for i, data := range rowData {
			x := currentX
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}
			
			pdf.SetXY(x, startY)
			pdf.Rect(x, startY, colWidths[i], rowHeight, "D")
			
			pdf.SetXY(x, startY)
			pdf.CellFormat(colWidths[i], rowHeight, "", "", 0, "", true, 0, "")
			
			pdf.SetXY(x+0.5, startY+0.5)
			pdf.MultiCell(colWidths[i]-1, lineHeight, data, "", "L", false)
		}
		
		pdf.SetXY(currentX, startY+rowHeight)
	}
	
	// Output PDF to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %s", err.Error())
	}
	
	return buf.Bytes(), nil
}

// DownloadTemplateSiswaLulus generates an Excel template for siswa lulus with nama and nis columns only
func (s *PesertaDidikServiceImpl) DownloadTemplateSiswaLulus() (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sheet1"

	headers := []string{
		"nama", "nis",
	}

	// Create header style with bold font and blue background
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat style header: %s", err.Error())
	}

	// Set headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set example values in row 2
	examples := []string{
		"Ahmad Fauzi",
		"001234",
	}

	for i, val := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 30, // nama
		"B": 15, // nis
	}

	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	return f, nil
}

// ImportSiswaLulus imports siswa lulus data from Excel and updates status to "lulus" with transaction (all-or-nothing)
func (s *PesertaDidikServiceImpl) ImportSiswaLulus(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file excel: %s", err.Error())
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca Sheet1: %s", err.Error())
	}

	// OPTIMIZATION: Load all peserta didik and create NIS -> peserta didik map
	allPesertaDidik, _, err := s.repository.GetAll(10000, 0) // max 10000 siswa
	if err != nil {
		return nil, fmt.Errorf("gagal load data peserta didik: %s", err.Error())
	}
	pesertaDidikMap := make(map[string]*models.PesertaDidik)
	for i := range allPesertaDidik {
		pesertaDidikMap[allPesertaDidik[i].NIS] = &allPesertaDidik[i]
	}

	// VALIDATION PHASE - Check all data before updating (outside transaction)
	type UpdateRow struct {
		RowNum         int
		PesertaDidikID uint
		Nama           string
		NIS            string
	}

	var validRows []UpdateRow
	var importErrors []dtos.ImportExcelRowError

	for i, row := range rows {
		// Skip header row
		if i == 0 {
			continue
		}

		rowNum := i + 1

		// Skip empty rows
		if len(row) == 0 {
			continue
		}

		// Helper to safely get column value
		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		nama := getCol(0)
		nis := getCol(1)

		// Validate required fields
		if nama == "" || nis == "" {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: "Kolom nama dan nis wajib diisi",
			})
			continue
		}

		// Get peserta didik from map (optimized - no DB query)
		pesertaDidik, exists := pesertaDidikMap[nis]
		if !exists {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Peserta didik dengan NIS '%s' tidak ditemukan", nis),
			})
			continue
		}

		// Optional: Validate nama matches (case-insensitive)
		if strings.ToLower(pesertaDidik.Nama) != strings.ToLower(nama) {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Nama tidak cocok. NIS '%s' terdaftar atas nama '%s', bukan '%s'", nis, pesertaDidik.Nama, nama),
			})
			continue
		}

		// Add to valid rows
		validRows = append(validRows, UpdateRow{
			RowNum:         rowNum,
			PesertaDidikID: pesertaDidik.ID,
			Nama:           pesertaDidik.Nama,
			NIS:            pesertaDidik.NIS,
		})
	}

	// If there are any validation errors, return error immediately (don't update anything)
	if len(importErrors) > 0 {
		return &dtos.ImportExcelResponse{
			SuccessCount: 0,
			FailedCount:  len(importErrors),
			Errors:       importErrors,
		}, fmt.Errorf("validasi gagal, tidak ada data yang diupdate. Total error: %d", len(importErrors))
	}

	// TRANSACTION PHASE - Update all peserta didik status to "lulus" in a single transaction (all or nothing)
	successCount := 0
	err = s.repository.UpdateWithTransaction(func(tx interface{}) error {
		for _, validRow := range validRows {
			// Get peserta didik by ID
			pesertaDidik, err := s.repository.GetByID(validRow.PesertaDidikID)
			if err != nil {
				return fmt.Errorf("gagal mengambil data peserta didik baris %d: %s", validRow.RowNum, err.Error())
			}

			// Update status to "lulus"
			pesertaDidik.Status = "lulus"
			pesertaDidik.UpdatedByID = &userID

			if err := s.repository.UpdateInTransaction(tx, pesertaDidik); err != nil {
				return fmt.Errorf("gagal mengupdate status peserta didik %s (NIS: %s) baris %d: %s", validRow.Nama, validRow.NIS, validRow.RowNum, err.Error())
			}

			successCount++
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan data: %s", err.Error())
	}

	return &dtos.ImportExcelResponse{
		SuccessCount: successCount,
		FailedCount:  0,
		Errors:       []dtos.ImportExcelRowError{},
	}, nil
}

// DownloadKartuPelajar generates PDF for student cards (status active only)
func (s *PesertaDidikServiceImpl) DownloadKartuPelajar(pesertaDidikIDs []uint) ([]byte, error) {
	var siswa []models.PesertaDidik
	var err error
	
	// Jika IDs kosong, ambil semua siswa aktif
	if len(pesertaDidikIDs) == 0 {
		siswa, err = s.repository.GetAllActive()
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
	} else {
		// Jika IDs ada, ambil siswa berdasarkan IDs
		siswa, err = s.repository.GetByIDs(pesertaDidikIDs)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data siswa: %s", err.Error())
		}
	}
	
	if len(siswa) == 0 {
		return nil, fmt.Errorf("tidak ada siswa ditemukan")
	}
	
	// Convert Photo field to full R2 URL for each student
	for i := range siswa {
		if siswa[i].Photo != "" {
			siswa[i].Photo = s.r2Storage.GetPublicURL(siswa[i].Photo)
		}
	}
	
	// Get Kepala Sekolah
	kepalaSekolah, err := s.repository.GetKepalaSekolah()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data kepala sekolah: %s", err.Error())
	}
	
	// Get Visi Misi
	visiMisi, err := s.repository.GetVisiMisi()
	if err != nil {
		// Non-fatal, continue without visi misi
		visiMisi = nil
	}
	
	// Prepare data
	data := &KartuPelajarData{
		Siswa:         siswa,
		KepalaSekolah: kepalaSekolah,
		VisiMisi:      visiMisi,
	}
	
	// Generate PDF
	pdfBytes, err := GenerateKartuPelajarPDF(data)
	if err != nil {
		return nil, err
	}
	
	return pdfBytes, nil
}
