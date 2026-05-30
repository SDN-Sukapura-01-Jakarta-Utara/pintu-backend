package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// KelulusanController handles HTTP requests for Kelulusan
type KelulusanController struct {
	service services.KelulusanService
}

// NewKelulusanController creates a new Kelulusan controller
func NewKelulusanController(service services.KelulusanService) *KelulusanController {
	return &KelulusanController{service: service}
}

// CreateKelulusan creates a new kelulusan record with optional SKL file upload
func (c *KelulusanController) CreateKelulusan(ctx *gin.Context) {
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
	var req dtos.KelulusanCreateRequest
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

	// Get SKL file if uploaded
	var file *multipart.FileHeader
	fileHeader, err := ctx.FormFile("skl")
	if err == nil {
		file = fileHeader
	}

	// Call service
	result, err := c.service.CreateKelulusan(&req, file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Data kelulusan berhasil ditambahkan",
		"data":    result,
	})
}

// DownloadTemplate downloads the Excel template for kelulusan import
func (c *KelulusanController) DownloadTemplate(ctx *gin.Context) {
	var req dtos.KelulusanDownloadTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	f, err := c.service.DownloadTemplate(req.MapelList)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat template"})
		return
	}
	defer f.Close()

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_kelulusan.xlsx")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ImportExcel imports kelulusan data from Excel file
func (c *KelulusanController) ImportExcel(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file excel wajib diunggah"})
		return
	}
	defer file.Close()

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	result, err := c.service.ImportExcel(file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}


// GetAll retrieves all Kelulusan with pagination and filters
func (c *KelulusanController) GetAll(ctx *gin.Context) {
	var req dtos.KelulusanGetAllRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Build filter
	filter := repositories.GetKelulusanFilter{
		Nama:         req.Search.Nama,
		NomorPeserta: req.Search.NomorPeserta,
		NISN:         req.Search.NISN,
		Lulus:        req.Search.Lulus,
	}

	// Set default pagination
	limit := req.Pagination.Limit
	page := req.Pagination.Page

	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// Call service with filter
	params := repositories.GetKelulusanParams{
		Filter: filter,
		Limit:  limit,
		Offset: offset,
	}

	result, err := c.service.GetAllWithFilter(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}


// GetByID retrieves a Kelulusan by ID
func (c *KelulusanController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "data kelulusan tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Update updates a Kelulusan record
func (c *KelulusanController) Update(ctx *gin.Context) {
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
	var req dtos.KelulusanUpdateRequest
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

	// Get SKL file if uploaded
	var file *multipart.FileHeader
	fileHeader, err := ctx.FormFile("skl")
	if err == nil {
		file = fileHeader
	}

	// Call service
	result, err := c.service.Update(req.ID, &req, file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data kelulusan berhasil diupdate",
		"data":    result,
	})
}


// Delete deletes a Kelulusan record
func (c *KelulusanController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kelulusan berhasil dihapus",
	})
}
