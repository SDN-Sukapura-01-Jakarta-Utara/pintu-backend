package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// PertanyaanController handles HTTP requests for Pertanyaan
type PertanyaanController struct {
	service services.PertanyaanService
}

// NewPertanyaanController creates a new Pertanyaan controller
func NewPertanyaanController(service services.PertanyaanService) *PertanyaanController {
	return &PertanyaanController{service: service}
}

// CreatePublic creates a new Pertanyaan from public form (no auth required)
// @Summary Create new Pertanyaan (Public)
// @Description Create a new pertanyaan/pengaduan from public form with file uploads
// @Tags pertanyaan
// @Accept multipart/form-data
// @Produce json
// @Param nama formData string true "Nama pengirim"
// @Param email formData string true "Email pengirim"
// @Param telepon formData string false "Nomor telepon"
// @Param kategori formData string true "Kategori pertanyaan"
// @Param judul formData string true "Judul pertanyaan"
// @Param deskripsi formData string true "Deskripsi pertanyaan"
// @Param file_pertanyaan formData file false "File pertanyaan - multiple files allowed - max 10MB each"
// @Success 201 {object} gin.H{data=dtos.PertanyaanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/create-pertanyaan [post]
func (c *PertanyaanController) CreatePublic(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Bind form data
	var req dtos.PertanyaanCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get files (optional, multiple)
	files := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["file_pertanyaan"]; exists {
			files = uploadedFiles
		}
	}

	// Call service
	data, err := c.service.CreatePublic(files, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}


// TrackPertanyaan tracks pertanyaan status by ID Tiket (no auth required)
// @Summary Track Pertanyaan by ID Tiket (Public)
// @Description Track pertanyaan/pengaduan status using ID Tiket
// @Tags pertanyaan
// @Accept json
// @Produce json
// @Param body body dtos.PertanyaanTrackRequest true "Request body with ID Tiket"
// @Success 200 {object} gin.H{data=dtos.PertanyaanTrackResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/public/track-pertanyaan [post]
func (c *PertanyaanController) TrackPertanyaan(ctx *gin.Context) {
	var req dtos.PertanyaanTrackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.TrackByIDTiket(req.IDTiket)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}


// GetAll retrieves all pertanyaan with filters and pagination (auth required)
// @Summary Get all Pertanyaan with filters
// @Description Retrieve all pertanyaan with filters, sorting by prioritas and tanggal pengajuan
// @Tags pertanyaan
// @Accept json
// @Produce json
// @Param body body dtos.PertanyaanGetAllRequest true "Request body with filters and pagination"
// @Success 200 {object} dtos.PertanyaanListWithPaginationResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/pertanyaan/get-pertanyaan [post]
func (c *PertanyaanController) GetAll(ctx *gin.Context) {
	var req dtos.PertanyaanGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetAllWithFilter(&req)
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


// GetByID retrieves pertanyaan by ID (auth required)
// @Summary Get Pertanyaan by ID
// @Description Retrieve pertanyaan details by ID
// @Tags pertanyaan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PertanyaanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pertanyaan/get-pertanyaan-by-id [post]
func (c *PertanyaanController) GetByID(ctx *gin.Context) {
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


// SendReply sends email reply to pertanyaan (auth required)
// @Summary Send email reply to pertanyaan
// @Description Send email reply with answer and update pertanyaan status
// @Tags pertanyaan
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Pertanyaan ID"
// @Param judul_jawaban formData string true "Judul jawaban"
// @Param deskripsi_jawaban formData string true "Deskripsi jawaban"
// @Param file_jawaban formData file false "File jawaban - multiple files allowed - max 10MB each"
// @Success 200 {object} gin.H{data=dtos.PertanyaanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/pertanyaan/send-reply [post]
func (c *PertanyaanController) SendReply(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
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

	// Get required fields
	judulJawaban := ctx.PostForm("judul_jawaban")
	if judulJawaban == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judul_jawaban is required"})
		return
	}

	deskripsiJawaban := ctx.PostForm("deskripsi_jawaban")
	if deskripsiJawaban == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deskripsi_jawaban is required"})
		return
	}

	// Get files (optional, multiple)
	files := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["file_jawaban"]; exists {
			files = uploadedFiles
		}
	}

	// Create request DTO
	req := &dtos.PertanyaanSendReplyRequest{
		ID:               uint(id),
		JudulJawaban:     judulJawaban,
		DeskripsiJawaban: deskripsiJawaban,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.SendReply(files, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email berhasil dikirim dan data tersimpan",
		"data":    data,
	})
}

// ClosePertanyaan closes pertanyaan by ID (auth required)
// @Summary Close Pertanyaan
// @Description Close pertanyaan and set tanggal_selesai to current time
// @Tags pertanyaan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PertanyaanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pertanyaan/close-pertanyaan [post]
func (c *PertanyaanController) ClosePertanyaan(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.ClosePertanyaan(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pertanyaan berhasil ditutup",
		"data":    data,
	})
}

// DeletePertanyaan soft deletes pertanyaan by ID (auth required)
// @Summary Delete Pertanyaan
// @Description Soft delete pertanyaan by setting deleted_at and deleted_by_id
// @Tags pertanyaan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pertanyaan/delete-pertanyaan [post]
func (c *PertanyaanController) DeletePertanyaan(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := c.service.DeletePertanyaan(req.ID, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pertanyaan berhasil dihapus",
	})
}