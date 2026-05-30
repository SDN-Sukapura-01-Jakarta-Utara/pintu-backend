package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
)

type KelulusanService interface {
	CreateKelulusan(req *dtos.KelulusanCreateRequest, file *multipart.FileHeader, userID uint) (*dtos.KelulusanResponse, error)
	DownloadTemplate(mapelList []string) (*excelize.File, error)
	ImportExcel(file multipart.File, userID uint) (*dtos.ImportKelulusanResponse, error)
	GetAllWithFilter(params repositories.GetKelulusanParams) (*dtos.KelulusanListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.KelulusanResponse, error)
	Update(id uint, req *dtos.KelulusanUpdateRequest, file *multipart.FileHeader, userID uint) (*dtos.KelulusanResponse, error)
	Delete(id uint) error
}

type KelulusanServiceImpl struct {
	repository repositories.KelulusanRepository
	r2Storage  *utils.R2Storage
}

// NewKelulusanService creates a new Kelulusan service
func NewKelulusanService(repository repositories.KelulusanRepository) KelulusanService {
	return &KelulusanServiceImpl{
		repository: repository,
		r2Storage:  utils.NewR2Storage(),
	}
}

// CreateKelulusan creates a new kelulusan record with optional SKL file upload
func (s *KelulusanServiceImpl) CreateKelulusan(req *dtos.KelulusanCreateRequest, file *multipart.FileHeader, userID uint) (*dtos.KelulusanResponse, error) {
	// Parse tanggal_lahir (YYYY-MM-DD format)
	tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
	if err != nil {
		return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
	}

	// Check if nomor_peserta already exists
	existing, _ := s.repository.GetByNomorPeserta(req.NomorPeserta)
	if existing != nil {
		return nil, errors.New("nomor peserta sudah terdaftar")
	}

	// Convert nilai map to JSON
	nilaiJSON, err := json.Marshal(req.Nilai)
	if err != nil {
		return nil, errors.New("gagal mengkonversi nilai ke JSON")
	}

	// Handle SKL file upload if provided
	var sklPath string
	if file != nil {
		// Upload to R2 in kelulusan-skl folder
		uploadedPath, err := s.r2Storage.UploadFile(file, "kelulusan-skl")
		if err != nil {
			return nil, fmt.Errorf("gagal upload file SKL: %s", err.Error())
		}
		sklPath = uploadedPath
	}

	// Create kelulusan record
	kelulusan := &models.Kelulusan{
		NomorPeserta: req.NomorPeserta,
		NISN:         req.NISN,
		Nama:         req.Nama,
		TanggalLahir: tanggalLahir,
		Nilai:        datatypes.JSON(nilaiJSON),
		Lulus:        req.Lulus,
		SKL:          sklPath,
		CreatedByID:  &userID,
		UpdatedByID:  &userID,
	}

	if err := s.repository.Create(kelulusan); err != nil {
		// Delete uploaded file from R2 if save to DB failed
		if sklPath != "" {
			s.r2Storage.DeleteFile(sklPath)
		}
		return nil, errors.New("gagal menyimpan data kelulusan")
	}

	// Map to response
	response := s.mapToResponse(kelulusan)

	return response, nil
}

// mapToResponse maps Kelulusan model to KelulusanResponse DTO
func (s *KelulusanServiceImpl) mapToResponse(data *models.Kelulusan) *dtos.KelulusanResponse {
	// Parse nilai JSON to map
	var nilaiMap map[string]interface{}
	if err := json.Unmarshal(data.Nilai, &nilaiMap); err != nil {
		nilaiMap = make(map[string]interface{})
	}

	// Calculate rata-rata nilai (average)
	var totalNilai float64
	var countMapel int
	
	for _, nilai := range nilaiMap {
		// Convert nilai to float64
		switch v := nilai.(type) {
		case float64:
			totalNilai += v
			countMapel++
		case int:
			totalNilai += float64(v)
			countMapel++
		case int64:
			totalNilai += float64(v)
			countMapel++
		}
	}
	
	// Calculate average and round to 2 decimal places
	var rataRata float64
	if countMapel > 0 {
		rataRata = totalNilai / float64(countMapel)
		// Round to 2 decimal places
		rataRata = float64(int(rataRata*100)) / 100
	}

	// Generate full URL for SKL file
	sklURL := s.r2Storage.GetPublicURL(data.SKL)

	response := &dtos.KelulusanResponse{
		ID:            data.ID,
		NomorPeserta:  data.NomorPeserta,
		NISN:          data.NISN,
		Nama:          data.Nama,
		TanggalLahir:  data.TanggalLahir.Format("2006-01-02"),
		Nilai:         nilaiMap,
		RataRataNilai: rataRata,
		Lulus:         data.Lulus,
		SKL:           sklURL,
		CreatedAt:     data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedByID:   data.CreatedByID,
		UpdatedByID:   data.UpdatedByID,
	}

	return response
}


// DownloadTemplate generates an Excel template for Kelulusan import
func (s *KelulusanServiceImpl) DownloadTemplate(mapelList []string) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sheet1"

	// Base headers (fixed columns)
	baseHeaders := []string{
		"nomor_peserta", "nisn", "nama", "tanggal_lahir", "lulus",
	}

	// Combine base headers with dynamic mapel columns
	headers := append(baseHeaders, mapelList...)
	headers = append(headers, "catatan")

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
	baseExamples := []string{
		"2024-001-001", "0012345678", "Ahmad Fauzi", "2006-05-15", "true",
	}

	examples := append([]string{}, baseExamples...)
	
	// Add example nilai for each mapel
	for range mapelList {
		examples = append(examples, "85")
	}
	
	// Add catatan
	examples = append(examples, "CATATAN: Kolom mata pelajaran (setelah kolom 'lulus') bersifat dinamis. Anda bisa menambah/mengurangi kolom sesuai kebutuhan. Kolom 'lulus' isi dengan: true atau false")

	for i, val := range examples {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Set note style on catatan column row 2
	catatanCol := len(headers)
	catatanCell, _ := excelize.CoordinatesToCellName(catatanCol, 2)
	f.SetCellStyle(sheetName, catatanCell, catatanCell, noteStyle)

	// Add dropdown for lulus (column E = col 5, rows 2-1000)
	lulusValidation := &excelize.DataValidation{
		Type:             "list",
		AllowBlank:       false,
		ShowInputMessage: true,
		ShowErrorMessage: true,
	}
	lulusValidation.Sqref = "E2:E1000"
	lulusValidation.SetDropList([]string{"true", "false"})
	f.AddDataValidation(sheetName, lulusValidation)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 15) // nomor_peserta
	f.SetColWidth(sheetName, "B", "B", 12) // nisn
	f.SetColWidth(sheetName, "C", "C", 25) // nama
	f.SetColWidth(sheetName, "D", "D", 15) // tanggal_lahir
	f.SetColWidth(sheetName, "E", "E", 10) // lulus
	
	// Set width for mapel columns
	for i := range mapelList {
		col, _ := excelize.ColumnNumberToName(6 + i)
		f.SetColWidth(sheetName, col, col, 12)
	}
	
	// Set width for catatan column
	catatanColName, _ := excelize.ColumnNumberToName(catatanCol)
	f.SetColWidth(sheetName, catatanColName, catatanColName, 80)

	return f, nil
}

// ImportExcel imports Kelulusan data from an Excel file
func (s *KelulusanServiceImpl) ImportExcel(file multipart.File, userID uint) (*dtos.ImportKelulusanResponse, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("gagal membuka file excel")
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, errors.New("gagal membaca Sheet1")
	}

	if len(rows) < 2 {
		return nil, errors.New("file excel kosong atau tidak ada data")
	}

	// Get headers from first row
	headers := rows[0]
	
	// Find base column indices
	nomorPesertaIdx := -1
	nisnIdx := -1
	namaIdx := -1
	tanggalLahirIdx := -1
	lulusIdx := -1
	
	for i, header := range headers {
		header = strings.ToLower(strings.TrimSpace(header))
		switch header {
		case "nomor_peserta":
			nomorPesertaIdx = i
		case "nisn":
			nisnIdx = i
		case "nama":
			namaIdx = i
		case "tanggal_lahir":
			tanggalLahirIdx = i
		case "lulus":
			lulusIdx = i
		}
	}

	// Validate required columns
	if nomorPesertaIdx == -1 || nisnIdx == -1 || namaIdx == -1 || tanggalLahirIdx == -1 || lulusIdx == -1 {
		return nil, errors.New("kolom wajib tidak lengkap: nomor_peserta, nisn, nama, tanggal_lahir, lulus")
	}

	// Get mapel columns (columns after lulus, before catatan)
	mapelStartIdx := lulusIdx + 1
	mapelColumns := []string{}
	for i := mapelStartIdx; i < len(headers); i++ {
		header := strings.TrimSpace(headers[i])
		if strings.ToLower(header) == "catatan" {
			break
		}
		if header != "" {
			mapelColumns = append(mapelColumns, header)
		}
	}

	// Track nomor_peserta to detect duplicates
	excelNomorPeserta := make(map[string]int) // key: nomor_peserta -> first row number

	successCount := 0
	failedCount := 0
	var importErrors []dtos.ImportKelulusanRowError

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

		nomorPeserta := getCol(nomorPesertaIdx)
		nisn := getCol(nisnIdx)
		nama := getCol(namaIdx)
		tanggalLahirStr := getCol(tanggalLahirIdx)
		lulusStr := getCol(lulusIdx)

		// Validate required fields
		if nomorPeserta == "" {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: "nomor_peserta wajib diisi"})
			continue
		}

		if nisn == "" {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: "nisn wajib diisi"})
			continue
		}

		if nama == "" {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: "nama wajib diisi"})
			continue
		}

		// Check duplicate nomor_peserta in Excel
		if firstRow, exists := excelNomorPeserta[strings.ToLower(nomorPeserta)]; exists {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: fmt.Sprintf("nomor_peserta '%s' duplikat dengan baris %d", nomorPeserta, firstRow)})
			continue
		}
		excelNomorPeserta[strings.ToLower(nomorPeserta)] = rowNum

		// Check duplicate nomor_peserta in DB
		existing, _ := s.repository.GetByNomorPeserta(nomorPeserta)
		if existing != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: fmt.Sprintf("nomor_peserta '%s' sudah terdaftar", nomorPeserta)})
			continue
		}

		// Parse tanggal_lahir - support multiple formats including Excel date serial
		var tanggalLahir time.Time
		var parseErr error
		
		// First, try to parse as Excel date serial number
		if excelDate, err := strconv.ParseFloat(tanggalLahirStr, 64); err == nil && excelDate > 0 {
			// Excel stores dates as serial numbers (days since 1900-01-01)
			// Use excelize to convert Excel date to time.Time
			tanggalLahir, parseErr = excelize.ExcelDateToTime(excelDate, false)
			if parseErr == nil {
				// Successfully parsed as Excel date
				goto dateParseSuccess
			}
		}
		
		// Try format: MM-DD-YY (e.g., 06-16-13 = June 16, 2013)
		tanggalLahir, parseErr = time.Parse("01-02-06", tanggalLahirStr)
		if parseErr == nil {
			// Adjust year: if year < 50, assume 2000s, else 1900s
			if tanggalLahir.Year() < 50 {
				tanggalLahir = tanggalLahir.AddDate(2000, 0, 0)
			} else if tanggalLahir.Year() < 100 {
				tanggalLahir = tanggalLahir.AddDate(1900, 0, 0)
			}
			goto dateParseSuccess
		}
		
		// Try format: YYYY-MM-DD
		tanggalLahir, parseErr = time.Parse("2006-01-02", tanggalLahirStr)
		if parseErr == nil {
			goto dateParseSuccess
		}
		
		// Try format: DD/MM/YYYY
		tanggalLahir, parseErr = time.Parse("02/01/2006", tanggalLahirStr)
		if parseErr == nil {
			goto dateParseSuccess
		}
		
		// Try format: DD-MM-YYYY
		tanggalLahir, parseErr = time.Parse("02-01-2006", tanggalLahirStr)
		if parseErr == nil {
			goto dateParseSuccess
		}
		
		// Try format: D/M/YYYY (single digit day/month)
		tanggalLahir, parseErr = time.Parse("2/1/2006", tanggalLahirStr)
		if parseErr == nil {
			goto dateParseSuccess
		}
		
		// Try format: MM/DD/YY (e.g., 06/16/13)
		tanggalLahir, parseErr = time.Parse("01/02/06", tanggalLahirStr)
		if parseErr == nil {
			// Adjust year: if year < 50, assume 2000s, else 1900s
			if tanggalLahir.Year() < 50 {
				tanggalLahir = tanggalLahir.AddDate(2000, 0, 0)
			} else if tanggalLahir.Year() < 100 {
				tanggalLahir = tanggalLahir.AddDate(1900, 0, 0)
			}
			goto dateParseSuccess
		}
		
		// All formats failed
		failedCount++
		importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: fmt.Sprintf("format tanggal_lahir tidak valid: '%s', gunakan YYYY-MM-DD atau DD/MM/YYYY", tanggalLahirStr)})
		continue
		
		dateParseSuccess:

		// Parse lulus
		lulus := false
		lulusLower := strings.ToLower(lulusStr)
		if lulusLower == "true" || lulusLower == "1" || lulusLower == "ya" || lulusLower == "yes" {
			lulus = true
		}

		// Parse nilai (mapel columns)
		nilaiMap := make(map[string]interface{})
		hasInvalidNilai := false
		var invalidNilaiMsg string
		
		for idx, mapelName := range mapelColumns {
			colIdx := mapelStartIdx + idx
			nilaiStr := getCol(colIdx)
			
			if nilaiStr != "" {
				// Replace comma with dot for decimal separator
				nilaiStr = strings.Replace(nilaiStr, ",", ".", -1)
				
				// Try to parse as number
				nilai, err := strconv.ParseFloat(nilaiStr, 64)
				if err != nil {
					hasInvalidNilai = true
					invalidNilaiMsg = fmt.Sprintf("nilai '%s' untuk mapel '%s' tidak valid", getCol(colIdx), mapelName)
					break
				}
				nilaiMap[mapelName] = nilai
			}
		}
		
		if hasInvalidNilai {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: invalidNilaiMsg})
			continue
		}

		// Convert nilai map to JSON
		nilaiJSON, err := json.Marshal(nilaiMap)
		if err != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: "gagal mengkonversi nilai ke JSON"})
			continue
		}

		// Create kelulusan record
		kelulusan := &models.Kelulusan{
			NomorPeserta: nomorPeserta,
			NISN:         nisn,
			Nama:         nama,
			TanggalLahir: tanggalLahir,
			Nilai:        datatypes.JSON(nilaiJSON),
			Lulus:        lulus,
			CreatedByID:  &userID,
			UpdatedByID:  &userID,
		}

		if err := s.repository.Create(kelulusan); err != nil {
			failedCount++
			importErrors = append(importErrors, dtos.ImportKelulusanRowError{Row: rowNum, Message: fmt.Sprintf("gagal menyimpan data: %s", err.Error())})
			continue
		}

		successCount++
	}

	return &dtos.ImportKelulusanResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Errors:       importErrors,
	}, nil
}


// GetAllWithFilter retrieves Kelulusan with filters and pagination
func (s *KelulusanServiceImpl) GetAllWithFilter(params repositories.GetKelulusanParams) (*dtos.KelulusanListWithPaginationResponse, error) {
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
	responses := make([]dtos.KelulusanResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.KelulusanListWithPaginationResponse{
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


// GetByID retrieves Kelulusan by ID
func (s *KelulusanServiceImpl) GetByID(id uint) (*dtos.KelulusanResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("data kelulusan tidak ditemukan")
	}

	return s.mapToResponse(data), nil
}

// Update updates Kelulusan record
func (s *KelulusanServiceImpl) Update(id uint, req *dtos.KelulusanUpdateRequest, file *multipart.FileHeader, userID uint) (*dtos.KelulusanResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("data kelulusan tidak ditemukan")
	}

	oldSKL := existing.SKL

	// Update fields if provided
	if req.NomorPeserta != "" && req.NomorPeserta != existing.NomorPeserta {
		// Check if new nomor_peserta already exists
		existingNomor, _ := s.repository.GetByNomorPeserta(req.NomorPeserta)
		if existingNomor != nil && existingNomor.ID != id {
			return nil, errors.New("nomor peserta sudah terdaftar")
		}
		existing.NomorPeserta = req.NomorPeserta
	}

	if req.NISN != "" {
		existing.NISN = req.NISN
	}

	if req.Nama != "" {
		existing.Nama = req.Nama
	}

	if req.TanggalLahir != "" {
		tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
		}
		existing.TanggalLahir = tanggalLahir
	}

	if req.Nilai != nil && len(req.Nilai) > 0 {
		nilaiJSON, err := json.Marshal(req.Nilai)
		if err != nil {
			return nil, errors.New("gagal mengkonversi nilai ke JSON")
		}
		existing.Nilai = datatypes.JSON(nilaiJSON)
	}

	if req.Lulus != nil {
		existing.Lulus = *req.Lulus
	}

	// Handle file deletion if requested
	if req.DeleteSKL {
		// Delete old file from R2 if exists
		if oldSKL != "" {
			_ = s.r2Storage.DeleteFile(oldSKL)
		}
		existing.SKL = ""
	}

	// Handle file upload if provided (this will override delete_skl if both are sent)
	if file != nil {
		// Upload new file to R2
		uploadedPath, err := s.r2Storage.UploadFile(file, "kelulusan-skl")
		if err != nil {
			return nil, fmt.Errorf("gagal upload file SKL: %s", err.Error())
		}

		// Delete old file from R2 if exists (only if different from new file)
		if oldSKL != "" && oldSKL != uploadedPath {
			_ = s.r2Storage.DeleteFile(oldSKL)
		}

		existing.SKL = uploadedPath
	}

	// Update metadata
	existing.UpdatedByID = &userID

	// Save to database
	if err := s.repository.Update(existing); err != nil {
		// If update failed and new file was uploaded, delete the new file
		if file != nil && existing.SKL != "" && existing.SKL != oldSKL {
			s.r2Storage.DeleteFile(existing.SKL)
		}
		return nil, errors.New("gagal mengupdate data kelulusan")
	}

	// Map to response
	response := s.mapToResponse(existing)

	return response, nil
}


// Delete deletes Kelulusan record (soft delete)
func (s *KelulusanServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("data kelulusan tidak ditemukan")
	}

	// Delete file SKL from R2 if exists
	if existing.SKL != "" {
		_ = s.r2Storage.DeleteFile(existing.SKL)
	}

	// Soft delete from database
	if err := s.repository.Delete(id); err != nil {
		return errors.New("gagal menghapus data kelulusan")
	}

	return nil
}
