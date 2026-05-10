package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// PengaduanController handles HTTP requests for Pengaduan
type PengaduanController struct {
	service services.PengaduanService
}

// NewPengaduanController creates a new Pengaduan controller
func NewPengaduanController(service services.PengaduanService) *PengaduanController {
	return &PengaduanController{service: service}
}

// CreatePublic creates a new Pengaduan from public form (no auth required)
// @Summary Create new Pengaduan (Public)
// @Description Create a new pengaduan from public form with file uploads. Nama dan email tidak wajib diisi untuk pengaduan anonim.
// @Tags pengaduan
// @Accept multipart/form-data
// @Produce json
// @Param tipe_pelapor formData string false "Tipe pelapor (anonim/teridentifikasi)" default(anonim)
// @Param nama formData string false "Nama pengirim (opsional untuk pengaduan anonim)"
// @Param email formData string false "Email pengirim (opsional untuk pengaduan anonim)"
// @Param telepon formData string false "Nomor telepon"
// @Param kategori formData string true "Kategori pengaduan"
// @Param prioritas formData string false "Prioritas pengaduan" default(Sedang)
// @Param judul formData string true "Judul pengaduan"
// @Param deskripsi formData string true "Deskripsi pengaduan"
// @Param file_pengaduan formData file false "File pengaduan - multiple files allowed - max 10MB each"
// @Success 201 {object} gin.H{data=dtos.PengaduanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/create-pengaduan [post]
func (c *PengaduanController) CreatePublic(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Bind form data
	var req dtos.PengaduanCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get files (optional, multiple)
	files := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["file_pengaduan"]; exists {
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

// TrackPengaduan tracks pengaduan status by ID Tiket (no auth required)
// @Summary Track Pengaduan by ID Tiket (Public)
// @Description Track pengaduan status using ID Tiket
// @Tags pengaduan
// @Accept json
// @Produce json
// @Param body body dtos.PengaduanTrackRequest true "Request body with ID Tiket"
// @Success 200 {object} gin.H{data=dtos.PengaduanTrackResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/public/track-pengaduan [post]
func (c *PengaduanController) TrackPengaduan(ctx *gin.Context) {
	var req dtos.PengaduanTrackRequest
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

// GetAll retrieves all pengaduan with filters and pagination (auth required)
// @Summary Get all Pengaduan with filters
// @Description Retrieve all pengaduan with filters, sorting by prioritas and tanggal pengajuan
// @Tags pengaduan
// @Accept json
// @Produce json
// @Param body body dtos.PengaduanGetAllRequest true "Request body with filters and pagination"
// @Success 200 {object} dtos.PengaduanListWithPaginationResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/get-pengaduan [post]
func (c *PengaduanController) GetAll(ctx *gin.Context) {
	var req dtos.PengaduanGetAllRequest
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

// GetByID retrieves pengaduan by ID (auth required)
// @Summary Get Pengaduan by ID
// @Description Retrieve pengaduan details by ID
// @Tags pengaduan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PengaduanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/get-pengaduan-by-id [post]
func (c *PengaduanController) GetByID(ctx *gin.Context) {
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

// SendReply sends email reply to pengaduan (auth required)
// @Summary Send email reply to pengaduan
// @Description Send email reply with answer and update pengaduan status. Hanya bisa mengirim email jika pengaduan memiliki alamat email.
// @Tags pengaduan
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Pengaduan ID"
// @Param judul_jawaban formData string true "Judul jawaban"
// @Param deskripsi_jawaban formData string true "Deskripsi jawaban"
// @Param file_jawaban formData file false "File jawaban - multiple files allowed - max 10MB each"
// @Success 200 {object} gin.H{data=dtos.PengaduanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/send-reply [post]
func (c *PengaduanController) SendReply(ctx *gin.Context) {
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
	req := &dtos.PengaduanSendReplyRequest{
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

// SaveTindakLanjut saves tindak lanjut for pengaduan (requires auth)
// @Summary Save tindak lanjut for pengaduan
// @Description Save or update tindak lanjut with optional file uploads and deletion
// @Tags pengaduan
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id formData int true "Pengaduan ID"
// @Param tindak_lanjut formData string true "Tindak lanjut text"
// @Param files_to_delete formData []string false "Array of file IDs to delete (comma-separated)"
// @Param file_tindak_lanjut formData file false "File tindak lanjut - multiple files allowed - max 10MB each"
// @Success 200 {object} gin.H{data=dtos.PengaduanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/save-tindak-lanjut [post]
func (c *PengaduanController) SaveTindakLanjut(ctx *gin.Context) {
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

	// Get tindak_lanjut from form
	tindakLanjut := ctx.PostForm("tindak_lanjut")
	if tindakLanjut == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tindak_lanjut is required"})
		return
	}

	// Parse files_to_delete if provided (sent as JSON string)
	var filesToDelete []string
	if filesDeleteStr := ctx.PostForm("files_to_delete"); filesDeleteStr != "" {
		if err := json.Unmarshal([]byte(filesDeleteStr), &filesToDelete); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid files_to_delete format"})
			return
		}
	}

	// Get files from form
	var files []*multipart.FileHeader
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["file_tindak_lanjut"]; exists {
			files = uploadedFiles
		}
	}

	// Create request DTO
	req := dtos.PengaduanSaveTindakLanjutRequest{
		ID:            uint(id),
		TindakLanjut:  tindakLanjut,
		FilesToDelete: filesToDelete,
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	result, err := c.service.SaveTindakLanjut(files, &req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tindak lanjut berhasil disimpan",
		"data":    result,
	})
}

// ClosePengaduan closes pengaduan by ID (auth required)
// @Summary Close Pengaduan
// @Description Close pengaduan by setting status to closed and tanggal_selesai
// @Tags pengaduan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Pengaduan ID"
// @Success 200 {object} gin.H{message=string,data=dtos.PengaduanResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/close-pengaduan [post]
func (c *PengaduanController) ClosePengaduan(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.ClosePengaduan(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pengaduan berhasil ditutup",
		"data":    data,
	})
}

// DeletePengaduan soft deletes pengaduan by ID (auth required)
// @Summary Delete Pengaduan
// @Description Soft delete pengaduan by setting deleted_at and deleted_by_id
// @Tags pengaduan
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Pengaduan ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/pengaduan/delete-pengaduan [post]
func (c *PengaduanController) DeletePengaduan(ctx *gin.Context) {
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

	err := c.service.DeletePengaduan(req.ID, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pengaduan berhasil dihapus",
	})
}
