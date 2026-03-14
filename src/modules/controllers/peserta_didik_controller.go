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
