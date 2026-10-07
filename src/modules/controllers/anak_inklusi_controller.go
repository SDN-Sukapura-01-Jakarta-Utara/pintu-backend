package controllers

import (
	"net/http"
	"strconv"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

type AnakInklusiController struct {
	service services.AnakInklusiService
}

func NewAnakInklusiController(service services.AnakInklusiService) *AnakInklusiController {
	return &AnakInklusiController{service: service}
}

// Create creates a new anak inklusi with file upload
// @Summary Create new anak inklusi
// @Description Create a new anak inklusi with file surat dokter upload to Cloudflare R2
// @Tags monitoring-inklusi
// @Accept multipart/form-data
// @Produce json
// @Param peserta_didik_id formData uint true "Peserta Didik ID"
// @Param jenis_hambatan formData string true "Jenis hambatan"
// @Param tanggal_identifikasi formData string false "Tanggal identifikasi (YYYY-MM-DD)"
// @Param tanggal_diagnosa formData string false "Tanggal diagnosa (YYYY-MM-DD)"
// @Param tanggal_kadaluarsa_surat formData string false "Tanggal kadaluarsa surat (YYYY-MM-DD)"
// @Param status formData string false "Status (identified/diagnosed/lulus)"
// @Param catatan formData string false "Catatan"
// @Param file_surat_dokter formData file false "File surat dokter (PDF) - max 5MB"
// @Success 201 {object} gin.H{data=dtos.AnakInklusiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/create-data-induk-inklusi [post]
func (c *AnakInklusiController) Create(ctx *gin.Context) {
	// Parse multipart form (max 10MB)
	if err := ctx.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get file surat dokter (optional)
	fileSuratDokter, _ := ctx.FormFile("file_surat_dokter")

	// Create request DTO from form data
	req := &dtos.AnakInklusiCreateRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.Create(fileSuratDokter, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all anak inklusi with filters and pagination
// @Summary Get all anak inklusi
// @Description Retrieve all anak inklusi records with filters and pagination
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.AnakInklusiListWithPaginationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-data-induk-inklusi [post]
func (c *AnakInklusiController) GetAll(ctx *gin.Context) {
	var req dtos.AnakInklusiGetAllRequest
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

	// Parse date filters
	var startDate, endDate time.Time
	if req.Search.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.StartDate); err == nil {
			startDate = parsed
		}
	}
	if req.Search.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.EndDate); err == nil {
			// Set to end of day for inclusive range
			endDate = parsed.Add(time.Hour * 24).Add(-time.Nanosecond)
		}
	}

	// Call service with filters (only active peserta_didik)
	data, err := c.service.GetAllWithFilter(repositories.GetAnakInklusiParams{
		Filter: repositories.GetAnakInklusiFilter{
			NamaPesertaDidik:       req.Search.NamaPesertaDidik,
			NIS:                    req.Search.NIS,
			NISN:                   req.Search.NISN,
			JenisHambatan:          req.Search.JenisHambatan,
			StartDate:              startDate,
			EndDate:                endDate,
			Status:                 req.Search.Status,
			PesertaDidikStatusOnly: "active", // Only show active students
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

// GetHistory retrieves all anak inklusi history (non-active students) with filters and pagination
// @Summary Get all anak inklusi history
// @Description Retrieve all anak inklusi records with non-active peserta_didik status, with filters and pagination
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.AnakInklusiListWithPaginationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-history-data-induk-inklusi [post]
func (c *AnakInklusiController) GetHistory(ctx *gin.Context) {
	var req dtos.AnakInklusiGetAllRequest
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

	// Parse date filters
	var startDate, endDate time.Time
	if req.Search.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.StartDate); err == nil {
			startDate = parsed
		}
	}
	if req.Search.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.EndDate); err == nil {
			// Set to end of day for inclusive range
			endDate = parsed.Add(time.Hour * 24).Add(-time.Nanosecond)
		}
	}

	// Call service with filters (only non-active peserta_didik)
	data, err := c.service.GetHistoryWithFilter(repositories.GetAnakInklusiParams{
		Filter: repositories.GetAnakInklusiFilter{
			NamaPesertaDidik:       req.Search.NamaPesertaDidik,
			NIS:                    req.Search.NIS,
			NISN:                   req.Search.NISN,
			JenisHambatan:          req.Search.JenisHambatan,
			StartDate:              startDate,
			EndDate:                endDate,
			Status:                 req.Search.Status,
			PesertaDidikStatusOnly: "non-active", // Only show non-active students
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

// GetByID retrieves anak inklusi by ID
// @Summary Get anak inklusi by ID
// @Description Retrieve anak inklusi details by ID
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-data-induk-inklusi-by-id [post]
func (c *AnakInklusiController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data anak inklusi tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates anak inklusi
// @Summary Update anak inklusi
// @Description Update anak inklusi details (all fields optional, including file upload)
// @Tags monitoring-inklusi
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Anak Inklusi ID"
// @Param jenis_hambatan formData string false "Jenis hambatan"
// @Param tanggal_identifikasi formData string false "Tanggal identifikasi (YYYY-MM-DD)"
// @Param tanggal_diagnosa formData string false "Tanggal diagnosa (YYYY-MM-DD)"
// @Param tanggal_kadaluarsa_surat formData string false "Tanggal kadaluarsa surat (YYYY-MM-DD)"
// @Param status formData string false "Status (identified/diagnosed/lulus)"
// @Param catatan formData string false "Catatan"
// @Param delete_file formData bool false "Set true to delete existing file"
// @Param file_surat_dokter formData file false "File surat dokter (PDF) - max 5MB (replaces existing)"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/update-data-induk-inklusi [post]
func (c *AnakInklusiController) Update(ctx *gin.Context) {
	// Parse multipart form (max 10MB)
	if err := ctx.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID from form
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	// Get optional fields
	jenisHambatan := ctx.PostForm("jenis_hambatan")
	tanggalIdentifikasi := ctx.PostForm("tanggal_identifikasi")
	tanggalDiagnosa := ctx.PostForm("tanggal_diagnosa")
	tanggalKadaluarsaSurat := ctx.PostForm("tanggal_kadaluarsa_surat")
	status := ctx.PostForm("status")
	catatan := ctx.PostForm("catatan")
	deleteFileStr := ctx.PostForm("delete_file")

	// Parse delete_file
	deleteFile := false
	if deleteFileStr == "true" || deleteFileStr == "1" {
		deleteFile = true
	}

	// Get file surat dokter (optional)
	fileSuratDokter, _ := ctx.FormFile("file_surat_dokter")

	// Create request DTO
	req := &dtos.AnakInklusiUpdateRequest{
		ID:                     uint(id),
		JenisHambatan:          jenisHambatan,
		TanggalIdentifikasi:    tanggalIdentifikasi,
		TanggalDiagnosa:        tanggalDiagnosa,
		TanggalKadaluarsaSurat: tanggalKadaluarsaSurat,
		Status:                 status,
		Catatan:                catatan,
		DeleteFile:             deleteFile,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(uint(id), fileSuratDokter, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes anak inklusi by ID
// @Summary Delete anak inklusi
// @Description Delete anak inklusi by ID (also deletes file from R2)
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/delete-data-induk-inklusi [post]
func (c *AnakInklusiController) Delete(ctx *gin.Context) {
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
		"message": "Data anak inklusi berhasil dihapus",
	})
}
