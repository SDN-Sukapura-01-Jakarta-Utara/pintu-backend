package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

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
	// Files are expected with field name pattern: file_surat_{peserta_didik_id}
	// Example: file_surat_1, file_surat_2, etc.
	files := make(map[uint][]*multipart.FileHeader)
	
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Check if field name starts with "file_surat_"
			if len(fieldName) > 11 && fieldName[:11] == "file_surat_" {
				// Extract peserta_didik_id from field name
				pesertaDidikIDStr := fieldName[11:]
				pesertaDidikID, err := strconv.ParseUint(pesertaDidikIDStr, 10, 32)
				if err != nil {
					continue // Skip invalid field names
				}
				files[uint(pesertaDidikID)] = fileHeaders
			}
		}
	}

	// Call service
	result, err := c.service.CreateAbsensiManual(&req, files, userIDUint)
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
