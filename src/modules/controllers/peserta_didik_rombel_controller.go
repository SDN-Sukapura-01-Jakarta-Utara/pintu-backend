package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
)

// PesertaDidikRombelController handles HTTP requests for PesertaDidikRombel
type PesertaDidikRombelController struct {
	service services.PesertaDidikRombelService
}

// NewPesertaDidikRombelController creates a new PesertaDidikRombel controller
func NewPesertaDidikRombelController(service services.PesertaDidikRombelService) *PesertaDidikRombelController {
	return &PesertaDidikRombelController{service: service}
}

// BulkCreate creates multiple PesertaDidikRombel mappings
func (c *PesertaDidikRombelController) BulkCreate(ctx *gin.Context) {
	var req dtos.PesertaDidikRombelCreateRequest

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
	result, err := c.service.BulkCreate(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if there are validation errors
	if result.FailedCount > 0 && result.SuccessCount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"data": result})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// GetAll retrieves all PesertaDidikRombel with pagination and filters
func (c *PesertaDidikRombelController) GetAll(ctx *gin.Context) {
	var req dtos.PesertaDidikRombelGetAllRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Build filter
	filter := repositories.GetPesertaDidikRombelFilter{
		Nama:             req.Search.Nama,
		RombelID:         req.Search.RombelID,
		TahunPelajaranID: req.Search.TahunPelajaranID,
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
	params := repositories.GetPesertaDidikRombelParams{
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

// GetByID retrieves a PesertaDidikRombel by ID
func (c *PesertaDidikRombelController) GetByID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "pemetaan rombel tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Update updates a PesertaDidikRombel
func (c *PesertaDidikRombelController) Update(ctx *gin.Context) {
	var req dtos.PesertaDidikRombelUpdateRequest

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

// Delete deletes a PesertaDidikRombel
func (c *PesertaDidikRombelController) Delete(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "pemetaan rombel berhasil dihapus"})
}

// DownloadTemplate downloads the Excel template for pemetaan rombel import
func (c *PesertaDidikRombelController) DownloadTemplate(ctx *gin.Context) {
	f, err := c.service.DownloadTemplate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat template"})
		return
	}
	defer f.Close()

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_pemetaan_rombel.xlsx")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ImportExcel imports pemetaan rombel data from Excel file
func (c *PesertaDidikRombelController) ImportExcel(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "data": result})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Reset deletes pemetaan rombel data by rombel_id or tahun_pelajaran_id or both
func (c *PesertaDidikRombelController) Reset(ctx *gin.Context) {
	var req dtos.PesertaDidikRombelResetRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.Reset(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
