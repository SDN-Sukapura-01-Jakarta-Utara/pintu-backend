package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
)

// PesertaDidikController handles HTTP requests for PesertaDidik
type PesertaDidikController struct {
	service services.PesertaDidikService
}

// NewPesertaDidikController creates a new PesertaDidik controller
func NewPesertaDidikController(service services.PesertaDidikService) *PesertaDidikController {
	return &PesertaDidikController{service: service}
}

// Create creates a new PesertaDidik
func (c *PesertaDidikController) Create(ctx *gin.Context) {
	var req dtos.PesertaDidikCreateRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service
	result, err := c.service.Create(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// GetByID retrieves a PesertaDidik by ID
func (c *PesertaDidikController) GetByID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "peserta didik tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetByNIS retrieves a PesertaDidik by NIS
func (c *PesertaDidikController) GetByNIS(ctx *gin.Context) {
	var req struct {
		NIS string `json:"nis" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GetByNIS(req.NIS)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "peserta didik dengan NIS tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAll retrieves all PesertaDidik with pagination and filters
func (c *PesertaDidikController) GetAll(ctx *gin.Context) {
	var req dtos.PesertaDidikGetAllRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Build filter
	filter := repositories.GetPesertaDidikFilter{
		TahunPelajaranID: req.Search.TahunPelajaranID,
		RombelID:         req.Search.RombelID,
		Nama:             req.Search.Nama,
		NIS:              req.Search.NIS,
		JenisKelamin:     req.Search.JenisKelamin,
		NISN:             req.Search.NISN,
		TempatLahir:      req.Search.TempatLahir,
		NIK:              req.Search.NIK,
		Agama:            req.Search.Agama,
		Status:           req.Search.Status,
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
	params := repositories.GetPesertaDidikParams{
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

// Update updates a PesertaDidik
func (c *PesertaDidikController) Update(ctx *gin.Context) {
	var req dtos.PesertaDidikUpdateRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service
	result, err := c.service.Update(req.ID, &req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Delete deletes a PesertaDidik
func (c *PesertaDidikController) Delete(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "peserta didik berhasil dihapus"})
}

// ImportExcel imports peserta didik data from Excel file
func (c *PesertaDidikController) ImportExcel(ctx *gin.Context) {
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

// DownloadTemplate downloads the Excel template for peserta didik import
func (c *PesertaDidikController) DownloadTemplate(ctx *gin.Context) {
	f, err := c.service.DownloadTemplate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat template"})
		return
	}
	defer f.Close()

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_peserta_didik.xlsx")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// GetTotalSiswa retrieves total count of peserta didik with active tahun pelajaran (public endpoint)
func (c *PesertaDidikController) GetTotalSiswa(ctx *gin.Context) {
	result, err := c.service.GetTotalSiswa()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
