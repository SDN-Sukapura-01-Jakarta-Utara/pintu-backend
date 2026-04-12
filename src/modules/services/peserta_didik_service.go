package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type PesertaDidikService interface {
	Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	GetByID(id uint) (*dtos.PesertaDidikResponse, error)
	GetByNIS(nis string) (*dtos.PesertaDidikResponse, error)
	GetAll(limit int, offset int) (*dtos.PesertaDidikListResponse, error)
	GetAllWithFilter(params repositories.GetPesertaDidikParams) (*dtos.PesertaDidikListWithPaginationResponse, error)
	Update(id uint, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	Delete(id uint) error
	ImportExcel(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error)
	DownloadTemplate() (*excelize.File, error)
}

type PesertaDidikServiceImpl struct {
	repository repositories.PesertaDidikRepository
}

// NewPesertaDidikService creates a new PesertaDidik service
func NewPesertaDidikService(repository repositories.PesertaDidikRepository) PesertaDidikService {
	return &PesertaDidikServiceImpl{
		repository: repository,
	}
}

// Create creates a new PesertaDidik
func (s *PesertaDidikServiceImpl) Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Check if NIS already exists in same tahun_pelajaran
	existing, _ := s.repository.GetByNISAndTahunPelajaran(req.NIS, req.TahunPelajaranID)
	if existing != nil {
		return nil, errors.New("NIS sudah ada di tahun pelajaran ini")
	}

	// Check if username already exists in same tahun_pelajaran (if provided)
	if req.Username != "" {
		existing, _ := s.repository.GetByUsernameAndTahunPelajaran(req.Username, req.TahunPelajaranID)
		if existing != nil {
			return nil, errors.New("username sudah ada di tahun pelajaran ini")
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
		RombelID:           req.RombelID,
		TahunPelajaranID:   req.TahunPelajaranID,
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

	return s.mapToResponse(data), nil
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
func (s *PesertaDidikServiceImpl) Update(id uint, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("peserta didik tidak ditemukan")
	}

	// Update basic fields if provided
	if req.Nama != "" {
		existing.Nama = req.Nama
	}
	if req.JenisKelamin != "" {
		existing.JenisKelamin = req.JenisKelamin
	}
	if req.NIS != "" && req.NIS != existing.NIS {
		// Check if new NIS already exists for other peserta didik in same tahun_pelajaran
		// Only check if NIS actually changed
		tahunPelajaranID := existing.TahunPelajaranID
		if existing.TahunPelajaranID == nil && req.TahunPelajaranID == nil {
			// Skip check if both are null
		} else if existing.TahunPelajaranID != nil && req.TahunPelajaranID == nil {
			tahunPelajaranID = existing.TahunPelajaranID
		} else if existing.TahunPelajaranID == nil && req.TahunPelajaranID != nil {
			tahunPelajaranID = req.TahunPelajaranID
		}
		
		existingNIS, _ := s.repository.GetByNISAndTahunPelajaran(req.NIS, tahunPelajaranID)
		if existingNIS != nil {
			return nil, errors.New("NIS sudah ada di tahun pelajaran ini")
		}
		existing.NIS = req.NIS
	}
	if req.NISN != "" {
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
	if req.RombelID != nil {
		existing.RombelID = req.RombelID
	}
	if req.TahunPelajaranID != nil {
		existing.TahunPelajaranID = req.TahunPelajaranID
	}
	if req.Username != "" && req.Username != existing.Username {
		// Check if new username already exists for other peserta didik in same tahun_pelajaran
		// Only check if username actually changed
		tahunPelajaranID := existing.TahunPelajaranID
		if existing.TahunPelajaranID == nil && req.TahunPelajaranID == nil {
			// Skip check if both are null
		} else if existing.TahunPelajaranID != nil && req.TahunPelajaranID == nil {
			tahunPelajaranID = existing.TahunPelajaranID
		} else if existing.TahunPelajaranID == nil && req.TahunPelajaranID != nil {
			tahunPelajaranID = req.TahunPelajaranID
		}
		
		existingUsername, _ := s.repository.GetByUsernameAndTahunPelajaran(req.Username, tahunPelajaranID)
		if existingUsername != nil {
			return nil, errors.New("username sudah ada di tahun pelajaran ini")
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

	// Update roles
	if len(req.RoleIDs) > 0 {
		if err := s.repository.AssignRoles(existing.ID, req.RoleIDs); err != nil {
			return nil, errors.New("gagal mengupdate roles peserta didik")
		}
	} else {
		// If no roles provided, remove all roles
		if err := s.repository.RemoveRoles(existing.ID); err != nil {
			return nil, errors.New("gagal menghapus roles peserta didik")
		}
	}

	return s.mapToResponse(existing), nil
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

	// Pre-load rombel and tahun_pelajaran data into maps for fast lookup
	rombelMap := make(map[string]uint)
	rombels, _ := s.repository.GetAllRombels()
	for _, r := range rombels {
		rombelMap[strings.ToLower(r.Name)] = r.ID
	}

	tpMap := make(map[string]uint)
	tpList, _ := s.repository.GetAllTahunPelajaran()
	for _, tp := range tpList {
		tpMap[strings.ToLower(tp.TahunPelajaran)] = tp.ID
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
		rombelName := getCol(18)
		tahunPelajaranName := getCol(19)
		roleIDStr := getCol(20)
		status := getCol(21)

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

		// Look up rombel by name from cached map
		var rombelID *uint
		if rombelName != "" {
			id, ok := rombelMap[strings.ToLower(rombelName)]
			if !ok {
				failedCount++
				importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("rombel '%s' tidak ditemukan", rombelName)})
				continue
			}
			rombelID = &id
		}

		// Look up tahun_pelajaran by name from cached map
		var tahunPelajaranID *uint
		if tahunPelajaranName != "" {
			id, ok := tpMap[strings.ToLower(tahunPelajaranName)]
			if !ok {
				failedCount++
				importErrors = append(importErrors, dtos.ImportExcelRowError{Row: rowNum, Message: fmt.Sprintf("tahun pelajaran '%s' tidak ditemukan", tahunPelajaranName)})
				continue
			}
			tahunPelajaranID = &id
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
			RombelID:         rombelID,
			TahunPelajaranID: tahunPelajaranID,
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
		"nama_ayah", "nama_ibu", "rombel", "tahun_pelajaran", "role_id", "status",
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

	// Fetch rombel data from DB for dropdown
	rombels, _ := s.repository.GetAllRombels()
	var rombelNames []string
	rombelExample := "1A"
	for _, r := range rombels {
		rombelNames = append(rombelNames, r.Name)
	}
	if len(rombelNames) > 0 {
		rombelExample = rombelNames[0]
	}

	// Fetch tahun pelajaran data from DB for dropdown
	tahunPelajaranList, _ := s.repository.GetAllTahunPelajaran()
	var tahunPelajaranNames []string
	tahunPelajaranExample := "2024/2025"
	for _, tp := range tahunPelajaranList {
		tahunPelajaranNames = append(tahunPelajaranNames, tp.TahunPelajaran)
	}
	if len(tahunPelajaranNames) > 0 {
		tahunPelajaranExample = tahunPelajaranNames[0]
	}

	// Set example values in row 2
	examples := []string{
		"siswa01", "password123", "Ahmad Fauzi", "001234", "1234567890",
		"Laki-laki", "Jakarta", "2015-07-15", "3171234567890001", "Islam",
		"Jl. Merdeka No. 1", "001", "002", "Sukapura", "Cilincing", "14140",
		"Budi Santoso", "Siti Aminah", rombelExample, tahunPelajaranExample, "[1,2]", "active",
		"Untuk role_id, buka menu Master Data > tab Role untuk melihat ID role",
	}

	for i, val := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Set note style on catatan column row 2
	f.SetCellStyle(sheetName, "W2", "W2", noteStyle)

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
	jkValidation.SetDropList([]string{"Laki-laki", "Perempuan"})
	f.AddDataValidation(sheetName, jkValidation)

	// Add dropdown for status (column V = col 22, rows 2-1000)
	statusValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       true,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	statusValidation.Sqref = "V2:V1000"
	statusValidation.SetDropList([]string{"active", "inactive"})
	f.AddDataValidation(sheetName, statusValidation)

	// Add dropdown for rombel (column S = col 19, rows 2-1000) if data exists
	if len(rombelNames) > 0 {
		rombelValidation := &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       true,
			ShowInputMessage: true,
			ShowErrorMessage: true,
		}
		rombelValidation.Sqref = "S2:S1000"
		rombelValidation.SetDropList(rombelNames)
		f.AddDataValidation(sheetName, rombelValidation)
	}

	// Add dropdown for tahun_pelajaran (column T = col 20, rows 2-1000) if data exists
	if len(tahunPelajaranNames) > 0 {
		tpValidation := &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       true,
			ShowInputMessage: true,
			ShowErrorMessage: true,
		}
		tpValidation.Sqref = "T2:T1000"
		tpValidation.SetDropList(tahunPelajaranNames)
		f.AddDataValidation(sheetName, tpValidation)
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 15, "B": 15, "C": 20, "D": 12, "E": 15,
		"F": 15, "G": 15, "H": 15, "I": 20, "J": 12,
		"K": 25, "L": 8, "M": 8, "N": 15, "O": 15, "P": 10,
		"Q": 20, "R": 20, "S": 12, "T": 20, "U": 12, "V": 10,
		"W": 55,
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

	// Map rombel
	var rombel *dtos.RombelDetailResponse
	if data.Rombel != nil {
		var kelas *dtos.KelasDetailResponse
		if data.Rombel.Kelas != nil {
			kelas = &dtos.KelasDetailResponse{
				ID:        data.Rombel.Kelas.ID,
				Name:      data.Rombel.Kelas.Name,
				Status:    data.Rombel.Kelas.Status,
				CreatedAt: data.Rombel.Kelas.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: data.Rombel.Kelas.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}
		rombel = &dtos.RombelDetailResponse{
			ID:        data.Rombel.ID,
			Name:      data.Rombel.Name,
			Status:    data.Rombel.Status,
			KelasID:   data.Rombel.KelasID,
			Kelas:     kelas,
			CreatedAt: data.Rombel.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: data.Rombel.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// Map tahun pelajaran
	var tahunPelajaran *dtos.TahunPelajaranDetailResponse
	if data.TahunPelajaran != nil {
		tahunPelajaran = &dtos.TahunPelajaranDetailResponse{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
			CreatedAt:      data.TahunPelajaran.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      data.TahunPelajaran.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			CreatedByID:    data.TahunPelajaran.CreatedByID,
			UpdatedByID:    data.TahunPelajaran.UpdatedByID,
		}
	}

	tanggalLahirStr := ""
	if data.TanggalLahir != nil {
		tanggalLahirStr = data.TanggalLahir.Format("2006-01-02")
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
		RombelID:         data.RombelID,
		Rombel:           rombel,
		TahunPelajaranID: data.TahunPelajaranID,
		TahunPelajaran:   tahunPelajaran,
		Status:           data.Status,
		Username:         data.Username,
		Roles:            roles,
		CreatedAt:        data.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedByID:      data.CreatedByID,
		UpdatedByID:      data.UpdatedByID,
	}
}
