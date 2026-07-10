package services

import (
	"fmt"
	"mime/multipart"
	"strings"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
	"github.com/xuri/excelize/v2"
)

type PesertaDidikRombelService interface {
	BulkCreate(req *dtos.PesertaDidikRombelCreateRequest, userID uint) (*dtos.PesertaDidikRombelBulkCreateResponse, error)
	GetAllWithFilter(params repositories.GetPesertaDidikRombelParams) (*dtos.PesertaDidikRombelListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.PesertaDidikRombelResponse, error)
	Update(id uint, req *dtos.PesertaDidikRombelUpdateRequest, userID uint) (*dtos.PesertaDidikRombelResponse, error)
	Delete(id uint) error
	DownloadTemplate() (*excelize.File, error)
	ImportExcel(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error)
	Reset(req *dtos.PesertaDidikRombelResetRequest) (*dtos.PesertaDidikRombelResetResponse, error)
}

type PesertaDidikRombelServiceImpl struct {
	repository               repositories.PesertaDidikRombelRepository
	pesertaDidikRepository   repositories.PesertaDidikRepository
	r2Storage                *utils.R2Storage
}

// NewPesertaDidikRombelService creates a new PesertaDidikRombel service
func NewPesertaDidikRombelService(
	repository repositories.PesertaDidikRombelRepository,
	pesertaDidikRepository repositories.PesertaDidikRepository,
	r2Storage *utils.R2Storage,
) PesertaDidikRombelService {
	return &PesertaDidikRombelServiceImpl{
		repository:             repository,
		pesertaDidikRepository: pesertaDidikRepository,
		r2Storage:              r2Storage,
	}
}

// BulkCreate creates multiple PesertaDidikRombel records for multiple students
// Uses transaction to ensure all records are created or none (all or nothing)
func (s *PesertaDidikRombelServiceImpl) BulkCreate(req *dtos.PesertaDidikRombelCreateRequest, userID uint) (*dtos.PesertaDidikRombelBulkCreateResponse, error) {
	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// VALIDATION PHASE - Check all data before saving
	var validationErrors []dtos.BulkCreateError

	for _, pesertaDidikID := range req.PesertaDidikIDs {
		// Check if peserta didik exists
		pesertaDidik, err := s.pesertaDidikRepository.GetByID(pesertaDidikID)
		if err != nil {
			validationErrors = append(validationErrors, dtos.BulkCreateError{
				PesertaDidikID: pesertaDidikID,
				Message:        fmt.Sprintf("Peserta didik dengan ID %d tidak ditemukan", pesertaDidikID),
			})
			continue
		}

		// Check if mapping already exists
		exists, err := s.repository.CheckDuplicateMapping(pesertaDidikID, req.RombelID, req.TahunPelajaranID)
		if err != nil {
			validationErrors = append(validationErrors, dtos.BulkCreateError{
				PesertaDidikID: pesertaDidikID,
				Message:        fmt.Sprintf("Gagal memeriksa duplikasi: %s", err.Error()),
			})
			continue
		}

		if exists {
			// Get rombel and tahun pelajaran data for detailed error message
			rombel, errRombel := s.repository.GetRombelByID(req.RombelID)
			tahunPelajaran, errTP := s.repository.GetTahunPelajaranByID(req.TahunPelajaranID)
			
			rombelName := "rombel ini"
			if errRombel == nil && rombel != nil {
				rombelName = fmt.Sprintf("rombel %s", rombel.Name)
			}
			
			tahunPelajaranName := "tahun pelajaran ini"
			if errTP == nil && tahunPelajaran != nil {
				tahunPelajaranName = fmt.Sprintf("tahun pelajaran %s", tahunPelajaran.TahunPelajaran)
			}
			
			validationErrors = append(validationErrors, dtos.BulkCreateError{
				PesertaDidikID: pesertaDidikID,
				Message:        fmt.Sprintf("Siswa %s (NIS: %s) sudah terdaftar di %s dan %s", pesertaDidik.Nama, pesertaDidik.NIS, rombelName, tahunPelajaranName),
			})
			continue
		}
	}

	// If there are any validation errors, return the response with errors (don't save anything)
	if len(validationErrors) > 0 {
		return &dtos.PesertaDidikRombelBulkCreateResponse{
			SuccessCount: 0,
			FailedCount:  len(validationErrors),
			Data:         []dtos.PesertaDidikRombelResponse{},
			Errors:       validationErrors,
		}, nil
	}

	// TRANSACTION PHASE - Save all data
	// All validations passed, now create all records in a transaction
	var createdIDs []uint
	err := s.repository.CreateWithTransaction(func(tx interface{}) error {
		for _, pesertaDidikID := range req.PesertaDidikIDs {
			data := &models.PesertaDidikRombel{
				PesertaDidikID:   pesertaDidikID,
				RombelID:         req.RombelID,
				TahunPelajaranID: req.TahunPelajaranID,
				Status:           status,
				CreatedByID:      &userID,
			}

			if err := s.repository.CreateInTransaction(tx, data); err != nil {
				return fmt.Errorf("gagal menyimpan data peserta didik ID %d: %s", pesertaDidikID, err.Error())
			}

			createdIDs = append(createdIDs, data.ID)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan data: %s", err.Error())
	}

	// Get all created data with relations
	var createdData []dtos.PesertaDidikRombelResponse
	for _, id := range createdIDs {
		createdRecord, err := s.repository.GetByID(id)
		if err == nil {
			createdData = append(createdData, *s.mapToResponse(createdRecord))
		}
	}

	return &dtos.PesertaDidikRombelBulkCreateResponse{
		SuccessCount: len(createdIDs),
		FailedCount:  0,
		Data:         createdData,
		Errors:       []dtos.BulkCreateError{},
	}, nil
}

// mapToResponse maps model to DTO response
func (s *PesertaDidikRombelServiceImpl) mapToResponse(data *models.PesertaDidikRombel) *dtos.PesertaDidikRombelResponse {
	response := &dtos.PesertaDidikRombelResponse{
		ID:               data.ID,
		PesertaDidikID:   data.PesertaDidikID,
		RombelID:         data.RombelID,
		TahunPelajaranID: data.TahunPelajaranID,
		Status:           data.Status,
		CreatedAt:        data.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedByID:      data.CreatedByID,
		UpdatedByID:      data.UpdatedByID,
	}

	// Map PesertaDidik (complete data)
	if data.PesertaDidik != nil {
		var tanggalLahir string
		if data.PesertaDidik.TanggalLahir != nil {
			tanggalLahir = data.PesertaDidik.TanggalLahir.Format("2006-01-02")
		}

		var barcodeGeneratedAt string
		if data.PesertaDidik.BarcodeGeneratedAt != nil {
			barcodeGeneratedAt = data.PesertaDidik.BarcodeGeneratedAt.Format("2006-01-02T15:04:05Z")
		}

		// Generate public URL for photo if exists
		photoURL := ""
		if data.PesertaDidik.Photo != "" {
			photoURL = s.r2Storage.GetPublicURL(data.PesertaDidik.Photo)
		}

		response.PesertaDidik = &dtos.PesertaDidikResponse{
			ID:                 data.PesertaDidik.ID,
			Nama:               data.PesertaDidik.Nama,
			NIS:                data.PesertaDidik.NIS,
			NISN:               data.PesertaDidik.NISN,
			JenisKelamin:       data.PesertaDidik.JenisKelamin,
			TempatLahir:        data.PesertaDidik.TempatLahir,
			TanggalLahir:       tanggalLahir,
			NIK:                data.PesertaDidik.NIK,
			Agama:              data.PesertaDidik.Agama,
			Alamat:             data.PesertaDidik.Alamat,
			RT:                 data.PesertaDidik.RT,
			RW:                 data.PesertaDidik.RW,
			Kelurahan:          data.PesertaDidik.Kelurahan,
			Kecamatan:          data.PesertaDidik.Kecamatan,
			KodePos:            data.PesertaDidik.KodePos,
			NamaAyah:           data.PesertaDidik.NamaAyah,
			NamaIbu:            data.PesertaDidik.NamaIbu,
			Status:             data.PesertaDidik.Status,
			Username:           data.PesertaDidik.Username,
			Photo:              photoURL,
			Barcode:            data.PesertaDidik.Barcode,
			BarcodeGeneratedAt: barcodeGeneratedAt,
			Roles:              []dtos.RoleResponse{}, // Empty roles for performance
			CreatedAt:          data.PesertaDidik.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:          data.PesertaDidik.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			CreatedByID:        data.PesertaDidik.CreatedByID,
			UpdatedByID:        data.PesertaDidik.UpdatedByID,
		}
	}

	// Map Rombel
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

		response.Rombel = &dtos.RombelDetailResponse{
			ID:        data.Rombel.ID,
			Name:      data.Rombel.Name,
			Status:    data.Rombel.Status,
			KelasID:   data.Rombel.KelasID,
			Kelas:     kelas,
			CreatedAt: data.Rombel.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: data.Rombel.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// Map TahunPelajaran
	if data.TahunPelajaran != nil {
		response.TahunPelajaran = &dtos.TahunPelajaranDetailResponse{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
			CreatedAt:      data.TahunPelajaran.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      data.TahunPelajaran.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			CreatedByID:    data.TahunPelajaran.CreatedByID,
			UpdatedByID:    data.TahunPelajaran.UpdatedByID,
		}
	}

	return response
}

// GetAllWithFilter retrieves PesertaDidikRombel with filters and pagination
func (s *PesertaDidikRombelServiceImpl) GetAllWithFilter(params repositories.GetPesertaDidikRombelParams) (*dtos.PesertaDidikRombelListWithPaginationResponse, error) {
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
	responses := make([]dtos.PesertaDidikRombelResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.PesertaDidikRombelListWithPaginationResponse{
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

// GetByID retrieves PesertaDidikRombel by ID with complete details
func (s *PesertaDidikRombelServiceImpl) GetByID(id uint) (*dtos.PesertaDidikRombelResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// Update updates PesertaDidikRombel
func (s *PesertaDidikRombelServiceImpl) Update(id uint, req *dtos.PesertaDidikRombelUpdateRequest, userID uint) (*dtos.PesertaDidikRombelResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pemetaan rombel tidak ditemukan")
	}

	// Check if peserta didik exists
	pesertaDidik, err := s.pesertaDidikRepository.GetByID(existing.PesertaDidikID)
	if err != nil {
		return nil, fmt.Errorf("peserta didik tidak ditemukan")
	}

	// Check duplicate mapping (excluding current ID)
	isDuplicate, err := s.repository.CheckDuplicateMappingExcludingID(id, existing.PesertaDidikID, req.RombelID, req.TahunPelajaranID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa duplikasi: %s", err.Error())
	}

	if isDuplicate {
		// Get rombel and tahun pelajaran data for detailed error message
		rombel, errRombel := s.repository.GetRombelByID(req.RombelID)
		tahunPelajaran, errTP := s.repository.GetTahunPelajaranByID(req.TahunPelajaranID)
		
		rombelName := "rombel ini"
		if errRombel == nil && rombel != nil {
			rombelName = fmt.Sprintf("rombel %s", rombel.Name)
		}
		
		tahunPelajaranName := "tahun pelajaran ini"
		if errTP == nil && tahunPelajaran != nil {
			tahunPelajaranName = fmt.Sprintf("tahun pelajaran %s", tahunPelajaran.TahunPelajaran)
		}
		
		return nil, fmt.Errorf("siswa %s (NIS: %s) sudah terdaftar di %s dan %s", pesertaDidik.Nama, pesertaDidik.NIS, rombelName, tahunPelajaranName)
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Update fields
	existing.RombelID = req.RombelID
	existing.TahunPelajaranID = req.TahunPelajaranID
	existing.Status = status
	existing.UpdatedByID = &userID

	// Save to database
	if err := s.repository.Update(existing); err != nil {
		return nil, fmt.Errorf("gagal mengupdate pemetaan rombel: %s", err.Error())
	}

	// Reload data with relations
	dataWithRelations, err := s.repository.GetByID(existing.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(dataWithRelations), nil
}

// Delete deletes PesertaDidikRombel by ID
func (s *PesertaDidikRombelServiceImpl) Delete(id uint) error {
	// Validate pemetaan rombel exists before delete
	data, err := s.repository.GetByID(id)
	if err != nil || data == nil {
		return fmt.Errorf("pemetaan rombel tidak ditemukan atau sudah dihapus")
	}

	return s.repository.Delete(id)
}

// DownloadTemplate generates an Excel template for PesertaDidikRombel import
func (s *PesertaDidikRombelServiceImpl) DownloadTemplate() (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sheet1"

	headers := []string{
		"nama", "nis", "rombel", "tahun_pelajaran", "status", "catatan",
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
		return nil, fmt.Errorf("gagal membuat style note: %s", err.Error())
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
		"1A",
		"2025/2026",
		"active",
		"Pastikan data nama dan nis sesuai dengan data peserta didik. Pilih rombel, tahun pelajaran, dan status dari dropdown yang tersedia.",
	}

	for i, val := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Set note style on catatan column row 2
	f.SetCellStyle(sheetName, "F2", "F2", noteStyle)

	// Get all rombels for dropdown
	rombels, err := s.repository.GetAllRombels()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data rombel: %s", err.Error())
	}

	rombelNames := make([]string, len(rombels))
	for i, rombel := range rombels {
		rombelNames[i] = rombel.Name
	}

	// Get all tahun pelajaran for dropdown
	tahunPelajaranList, err := s.repository.GetAllTahunPelajaran()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data tahun pelajaran: %s", err.Error())
	}

	tahunPelajaranNames := make([]string, len(tahunPelajaranList))
	for i, tp := range tahunPelajaranList {
		tahunPelajaranNames[i] = tp.TahunPelajaran
	}

	// Add dropdown for rombel (column C = col 3, rows 2-1000)
	if len(rombelNames) > 0 {
		rombelValidation := &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       false,
			ShowInputMessage: true,
			ShowErrorMessage: true,
		}
		rombelValidation.Sqref = "C2:C1000"
		rombelValidation.SetDropList(rombelNames)
		f.AddDataValidation(sheetName, rombelValidation)
	}

	// Add dropdown for tahun_pelajaran (column D = col 4, rows 2-1000)
	if len(tahunPelajaranNames) > 0 {
		tpValidation := &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       false,
			ShowInputMessage: true,
			ShowErrorMessage: true,
		}
		tpValidation.Sqref = "D2:D1000"
		tpValidation.SetDropList(tahunPelajaranNames)
		f.AddDataValidation(sheetName, tpValidation)
	}

	// Add dropdown for status (column E = col 5, rows 2-1000)
	statusValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       false,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	statusValidation.Sqref = "E2:E1000"
	statusValidation.SetDropList([]string{"active", "inactive"})
	f.AddDataValidation(sheetName, statusValidation)

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, // nama
		"B": 12, // nis
		"C": 12, // rombel
		"D": 18, // tahun_pelajaran
		"E": 12, // status
		"F": 80, // catatan
	}

	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	return f, nil
}

// ImportExcel imports PesertaDidikRombel data from an Excel file with transaction (optimized with caching)
func (s *PesertaDidikRombelServiceImpl) ImportExcel(file multipart.File, userID uint) (*dtos.ImportExcelResponse, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file excel: %s", err.Error())
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca Sheet1: %s", err.Error())
	}

	// OPTIMIZATION: Load all reference data once (cache)
	// Load all peserta didik and create NIS -> ID map
	allPesertaDidik, _, err := s.pesertaDidikRepository.GetAll(10000, 0) // max 10000 siswa
	if err != nil {
		return nil, fmt.Errorf("gagal load data peserta didik: %s", err.Error())
	}
	pesertaDidikMap := make(map[string]*models.PesertaDidik)
	for i := range allPesertaDidik {
		pesertaDidikMap[allPesertaDidik[i].NIS] = &allPesertaDidik[i]
	}

	// Load all rombels and create name -> ID map
	allRombels, err := s.repository.GetAllRombels()
	if err != nil {
		return nil, fmt.Errorf("gagal load data rombel: %s", err.Error())
	}
	rombelMap := make(map[string]*models.Rombel)
	for i := range allRombels {
		rombelMap[strings.ToLower(allRombels[i].Name)] = &allRombels[i]
	}

	// Load all tahun pelajaran and create name -> ID map
	allTahunPelajaran, err := s.repository.GetAllTahunPelajaran()
	if err != nil {
		return nil, fmt.Errorf("gagal load data tahun pelajaran: %s", err.Error())
	}
	tahunPelajaranMap := make(map[string]*models.TahunPelajaran)
	for i := range allTahunPelajaran {
		tahunPelajaranMap[allTahunPelajaran[i].TahunPelajaran] = &allTahunPelajaran[i]
	}

	// VALIDATION PHASE - Check all data before saving (outside transaction)
	type ImportRow struct {
		RowNum           int
		PesertaDidikID   uint
		RombelID         uint
		TahunPelajaranID uint
		Status           string
	}

	var validRows []ImportRow
	var importErrors []dtos.ImportExcelRowError
	// Track combinations to detect duplicates in Excel AND database
	existingCombinations := make(map[string]bool)

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
		rombelName := getCol(2)
		tahunPelajaranName := getCol(3)
		status := getCol(4)

		// Validate required fields
		if nama == "" || nis == "" || rombelName == "" || tahunPelajaranName == "" {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: "Kolom nama, nis, rombel, dan tahun_pelajaran wajib diisi",
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

		// Get rombel from map (optimized - no DB query)
		rombel, exists := rombelMap[strings.ToLower(rombelName)]
		if !exists {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Rombel '%s' tidak ditemukan", rombelName),
			})
			continue
		}

		// Get tahun pelajaran from map (optimized - no DB query)
		tahunPelajaran, exists := tahunPelajaranMap[tahunPelajaranName]
		if !exists {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Tahun pelajaran '%s' tidak ditemukan", tahunPelajaranName),
			})
			continue
		}

		// Create combination key for duplicate check
		combinationKey := fmt.Sprintf("%d-%d-%d", pesertaDidik.ID, rombel.ID, tahunPelajaran.ID)

		// Check if combination already exists in database
		if !existingCombinations[combinationKey] {
			isDuplicate, err := s.repository.CheckDuplicateMapping(pesertaDidik.ID, rombel.ID, tahunPelajaran.ID)
			if err != nil {
				importErrors = append(importErrors, dtos.ImportExcelRowError{
					Row:     rowNum,
					Message: fmt.Sprintf("Gagal memeriksa duplikasi: %s", err.Error()),
				})
				continue
			}

			if isDuplicate {
				importErrors = append(importErrors, dtos.ImportExcelRowError{
					Row:     rowNum,
					Message: fmt.Sprintf("Siswa %s (NIS: %s) sudah terdaftar di rombel %s dan tahun pelajaran %s", pesertaDidik.Nama, pesertaDidik.NIS, rombel.Name, tahunPelajaran.TahunPelajaran),
				})
				existingCombinations[combinationKey] = true // mark as checked
				continue
			}
		} else {
			// Already checked and found duplicate in previous row
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Duplikat dalam file Excel: Siswa %s (NIS: %s) di rombel %s dan tahun pelajaran %s sudah ada di baris sebelumnya", pesertaDidik.Nama, pesertaDidik.NIS, rombel.Name, tahunPelajaran.TahunPelajaran),
			})
			continue
		}

		// Mark this combination as used (prevent duplicate within Excel file)
		existingCombinations[combinationKey] = true

		// Set default status
		if status == "" {
			status = "active"
		}

		// Validate status value
		if status != "active" && status != "inactive" {
			importErrors = append(importErrors, dtos.ImportExcelRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("Status '%s' tidak valid, harus 'active' atau 'inactive'", status),
			})
			continue
		}

		// Add to valid rows
		validRows = append(validRows, ImportRow{
			RowNum:           rowNum,
			PesertaDidikID:   pesertaDidik.ID,
			RombelID:         rombel.ID,
			TahunPelajaranID: tahunPelajaran.ID,
			Status:           status,
		})
	}

	// If there are any validation errors, return error immediately (don't save anything)
	if len(importErrors) > 0 {
		return &dtos.ImportExcelResponse{
			SuccessCount: 0,
			FailedCount:  len(importErrors),
			Errors:       importErrors,
		}, fmt.Errorf("validasi gagal, tidak ada data yang disimpan. Total error: %d", len(importErrors))
	}

	// TRANSACTION PHASE - Save all data in a single transaction (all or nothing)
	successCount := 0
	err = s.repository.CreateWithTransaction(func(tx interface{}) error {
		for _, validRow := range validRows {
			data := &models.PesertaDidikRombel{
				PesertaDidikID:   validRow.PesertaDidikID,
				RombelID:         validRow.RombelID,
				TahunPelajaranID: validRow.TahunPelajaranID,
				Status:           validRow.Status,
				CreatedByID:      &userID,
			}

			if err := s.repository.CreateInTransaction(tx, data); err != nil {
				return fmt.Errorf("gagal menyimpan data baris %d: %s", validRow.RowNum, err.Error())
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

// Reset deletes pemetaan rombel data by tahun_pelajaran_id or both rombel_id and tahun_pelajaran_id
func (s *PesertaDidikRombelServiceImpl) Reset(req *dtos.PesertaDidikRombelResetRequest) (*dtos.PesertaDidikRombelResetResponse, error) {
	// Validate tahun_pelajaran_id is required
	if req.TahunPelajaranID == 0 {
		return nil, fmt.Errorf("tahun_pelajaran_id wajib diisi")
	}

	// Validate tahun pelajaran exists
	tahunPelajaran, errTP := s.repository.GetTahunPelajaranByID(req.TahunPelajaranID)
	if errTP != nil {
		return nil, fmt.Errorf("tahun pelajaran dengan ID %d tidak ditemukan", req.TahunPelajaranID)
	}

	var deletedCount int64
	var err error
	var message string

	// Determine which delete operation to perform based on provided filters
	if req.RombelID != 0 {
		// Both rombel_id and tahun_pelajaran_id provided - delete by both
		// Validate rombel exists
		rombel, errRombel := s.repository.GetRombelByID(req.RombelID)
		if errRombel != nil {
			return nil, fmt.Errorf("rombel dengan ID %d tidak ditemukan", req.RombelID)
		}

		deletedCount, err = s.repository.DeleteByRombelAndTahunPelajaran(req.RombelID, req.TahunPelajaranID)
		if err != nil {
			return nil, fmt.Errorf("gagal menghapus pemetaan rombel: %s", err.Error())
		}

		message = fmt.Sprintf("Berhasil menghapus %d pemetaan rombel untuk rombel %s dan tahun pelajaran %s", deletedCount, rombel.Name, tahunPelajaran.TahunPelajaran)
	} else {
		// Only tahun_pelajaran_id provided - delete all by tahun_pelajaran_id
		deletedCount, err = s.repository.DeleteByTahunPelajaranID(req.TahunPelajaranID)
		if err != nil {
			return nil, fmt.Errorf("gagal menghapus pemetaan rombel: %s", err.Error())
		}

		message = fmt.Sprintf("Berhasil menghapus %d pemetaan rombel untuk tahun pelajaran %s", deletedCount, tahunPelajaran.TahunPelajaran)
	}

	return &dtos.PesertaDidikRombelResetResponse{
		DeletedCount: int(deletedCount),
		Message:      message,
	}, nil
}
