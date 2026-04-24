package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// StrukturOrganisasiController handles HTTP requests for StrukturOrganisasi
type StrukturOrganisasiController struct {
	service services.StrukturOrganisasiService
}

// NewStrukturOrganisasiController creates a new StrukturOrganisasi controller
func NewStrukturOrganisasiController(service services.StrukturOrganisasiService) *StrukturOrganisasiController {
	return &StrukturOrganisasiController{service: service}
}

// Create creates a new StrukturOrganisasi
func (c *StrukturOrganisasiController) Create(ctx *gin.Context) {
	var req dtos.StrukturOrganisasiCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.Create(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all StrukturOrganisasi with filters and pagination
func (c *StrukturOrganisasiController) GetAll(ctx *gin.Context) {
	var req dtos.StrukturOrganisasiGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetStrukturOrganisasiParams{
		Filter: repositories.GetStrukturOrganisasiFilter{
			Nama:    req.Search.Nama,
			Urutan:  req.Search.Urutan,
			Relasi:  req.Search.Relasi,
			Jabatan: req.Search.Jabatan,
			Status:  req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data.Data,
		"pagination": gin.H{
			"limit":       data.Pagination.Limit,
			"offset":      data.Pagination.Offset,
			"page":        data.Pagination.Page,
			"total":       data.Pagination.Total,
			"total_pages": data.Pagination.TotalPages,
		},
	})
}

// GetByID retrieves StrukturOrganisasi by ID
func (c *StrukturOrganisasiController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates StrukturOrganisasi
func (c *StrukturOrganisasiController) Update(ctx *gin.Context) {
	// Read and check raw body FIRST to track which fields are present
	var rawBody map[string]interface{}
	if err := ctx.ShouldBindJSON(&rawBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Now convert to typed struct
	var req dtos.StrukturOrganisasiUpdateRequest
	req.ID = uint(rawBody["id"].(float64))
	
	// Check if pegawai_id is explicitly present
	if pegawaiID, exists := rawBody["pegawai_id"]; exists {
		req.PegawaiIDSet = true
		if pegawaiID != nil {
			pegawaiIDUint := uint(pegawaiID.(float64))
			req.PegawaiID = &pegawaiIDUint
		} else {
			// Explicitly set to nil when pegawai_id is null in request
			req.PegawaiID = nil
		}
	}

	// Map other fields
	if nama, ok := rawBody["nama_non_pegawai"].(string); ok && nama != "" {
		req.NamaNonPegawai = &nama
	}
	if jabatan, ok := rawBody["jabatan_non_pegawai"].(string); ok && jabatan != "" {
		req.JabatanNonPegawai = &jabatan
	}
	if urutan, ok := rawBody["urutan"].(float64); ok {
		urutanInt := int(urutan)
		req.Urutan = &urutanInt
	}
	if relasi, ok := rawBody["relasi"].(string); ok && relasi != "" {
		req.Relasi = &relasi
	}
	if status, ok := rawBody["status"].(string); ok && status != "" {
		req.Status = &status
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes StrukturOrganisasi by ID
func (c *StrukturOrganisasiController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Struktur organisasi deleted successfully",
	})
}

// GetPublic retrieves all active StrukturOrganisasi for public display (no auth required)
func (c *StrukturOrganisasiController) GetPublic(ctx *gin.Context) {
	data, err := c.service.GetPublic()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
