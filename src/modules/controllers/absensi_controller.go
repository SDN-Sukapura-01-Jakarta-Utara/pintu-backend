package controllers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AbsensiController handles HTTP requests for Absensi
type AbsensiController struct {
	service services.AbsensiService
}

// NewAbsensiController creates a new Absensi controller
func NewAbsensiController(service services.AbsensiService) *AbsensiController {
	return &AbsensiController{service: service}
}

// CreateAbsensiManual creates multiple absensi records (bulk input) with file upload support
func (c *AbsensiController) CreateAbsensiManual(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.AbsensiManualCreateRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Manual validation using validator
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Parse uploaded files
	// Files are expected with field name pattern: file_surat_{index}
	// Example: file_surat_1, file_surat_2, file_surat_3 (matching index in absensi_list array, starting from 1)
	filesMap := make(map[uint][]*multipart.FileHeader)
	
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Check if field name starts with "file_surat_"
			if len(fieldName) > 11 && fieldName[:11] == "file_surat_" {
				// Extract index from field name (1-based index)
				indexStr := fieldName[11:]
				index, err := strconv.ParseUint(indexStr, 10, 32)
				if err != nil {
					continue // Skip invalid field names
				}
				
				// Convert 1-based index to 0-based array index
				arrayIndex := int(index) - 1
				
				// Map index to peserta_didik_rombel_id from absensi_list
				if arrayIndex >= 0 && arrayIndex < len(req.AbsensiList) {
					pesertaDidikRombelID := req.AbsensiList[arrayIndex].PesertaDidikRombelID
					filesMap[pesertaDidikRombelID] = fileHeaders
				}
			}
		}
	}

	// Call service
	result, err := c.service.CreateAbsensiManual(&req, filesMap, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// CreateAbsensiManualByID creates a single absensi record by peserta didik rombel ID with auto semester detection
func (c *AbsensiController) CreateAbsensiManualByID(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.AbsensiManualCreateByIDRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Manual validation using validator
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Get file if uploaded
	var file *multipart.FileHeader
	fileHeader, err := ctx.FormFile("file_surat")
	if err == nil {
		file = fileHeader
	}

	// Call service
	result, err := c.service.CreateAbsensiManualByID(&req, file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// GetRekapAbsensi retrieves attendance recap with summary per student
func (c *AbsensiController) GetRekapAbsensi(ctx *gin.Context) {
	var req dtos.AbsensiRekapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetRekapAbsensi(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateRekapAbsensi updates a single absensi record
func (c *AbsensiController) UpdateRekapAbsensi(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.AbsensiUpdateRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Manual validation using validator
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Get file if uploaded
	var file *multipart.FileHeader
	fileHeader, err := ctx.FormFile("file_surat")
	if err == nil {
		file = fileHeader
	}

	// Call service
	result, err := c.service.UpdateRekapAbsensi(req.ID, &req, file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetDashboardSummary retrieves dashboard summary statistics
func (c *AbsensiController) GetDashboardSummary(ctx *gin.Context) {
	var req dtos.DashboardSummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetDashboardSummary(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetGrafikKehadiran retrieves attendance chart data
func (c *AbsensiController) GetGrafikKehadiran(ctx *gin.Context) {
	var req dtos.GrafikKehadiranRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetGrafikKehadiran(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetStatistikPerHari retrieves daily attendance statistics
func (c *AbsensiController) GetStatistikPerHari(ctx *gin.Context) {
	var req dtos.StatistikPerHariRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetStatistikPerHari(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetPerbandinganRombel retrieves attendance comparison between rombel
func (c *AbsensiController) GetPerbandinganRombel(ctx *gin.Context) {
	var req dtos.PerbandinganRombelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetPerbandinganRombel(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetSiswaTerendah retrieves students with lowest attendance
func (c *AbsensiController) GetSiswaTerendah(ctx *gin.Context) {
	var req dtos.SiswaTerendahRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetSiswaTerendah(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetDashboardSiswa retrieves dashboard data for a specific student
func (c *AbsensiController) GetDashboardSiswa(ctx *gin.Context) {
	var req dtos.DashboardSiswaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	result, err := c.service.GetDashboardSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// ExportAbsensiExcel exports absensi data to Excel file
func (c *AbsensiController) ExportAbsensiExcel(ctx *gin.Context) {
	var req dtos.ExportAbsensiExcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	file, err := c.service.ExportAbsensiExcel(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate filename
	filename := fmt.Sprintf("Daftar_Kehadiran_%d.xlsx", time.Now().Unix())

	// Set headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write file to response
	if err := file.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menulis file excel"})
		return
	}
}

// ExportAbsensiPDF exports absensi data to PDF file
func (c *AbsensiController) ExportAbsensiPDF(ctx *gin.Context) {
	var req dtos.ExportAbsensiExcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Call service
	pdfBytes, err := c.service.ExportAbsensiPDF(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate filename
	filename := fmt.Sprintf("Daftar_Kehadiran_%d.pdf", time.Now().Unix())

	// Set headers for file download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write PDF to response
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// SynchronizeAbsensi synchronizes data from absensi scan to rekapitulasi
func (c *AbsensiController) SynchronizeAbsensi(ctx *gin.Context) {
	var req dtos.AbsensiSyncRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service
	result, err := c.service.SynchronizeAbsensi(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
