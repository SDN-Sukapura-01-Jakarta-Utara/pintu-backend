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

// AnnouncementController handles HTTP requests for Announcement
type AnnouncementController struct {
	service services.AnnouncementService
}

// NewAnnouncementController creates a new Announcement controller
func NewAnnouncementController(service services.AnnouncementService) *AnnouncementController {
	return &AnnouncementController{service: service}
}

// Create creates a new Announcement with file uploads
// @Summary Create new Announcement
// @Description Create a new Announcement with gambar and files upload to Cloudflare R2
// @Tags announcement
// @Accept multipart/form-data
// @Produce json
// @Param judul formData string true "Announcement title"
// @Param tanggal formData string true "Announcement date (YYYY-MM-DD)"
// @Param deskripsi formData string false "Announcement description"
// @Param penulis formData string true "Announcement author"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param gambar formData file false "Featured image (jpeg, png, gif, webp) - max 5MB"
// @Param files formData file false "Attachment files - multiple files allowed (max 10MB each)"
// @Success 201 {object} gin.H{data=dtos.AnnouncementResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/announcements/create-announcement [post]
func (c *AnnouncementController) Create(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get judul
	judul := ctx.PostForm("judul")
	if judul == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judul is required"})
		return
	}

	// Get tanggal
	tanggal := ctx.PostForm("tanggal")
	if tanggal == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tanggal is required"})
		return
	}

	// Get penulis
	penulis := ctx.PostForm("penulis")
	if penulis == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "penulis is required"})
		return
	}

	// Get optional fields
	deskripsi := ctx.PostForm("deskripsi")
	statusPublikasi := ctx.PostForm("status_publikasi")
	status := ctx.PostForm("status")

	// Get gambar (optional)
	gambar, _ := ctx.FormFile("gambar")

	// Get files (optional, multiple)
	files := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["files"]; exists {
			files = uploadedFiles
		}
	}

	// Create request DTO
	req := &dtos.AnnouncementCreateRequest{
		Judul:           judul,
		Tanggal:         tanggal,
		Deskripsi:       deskripsi,
		Penulis:         penulis,
		StatusPublikasi: statusPublikasi,
		Status:          status,
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.Create(gambar, files, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all announcements
// @Summary Get all Announcements
// @Description Retrieve all Announcement records with pagination
// @Tags announcement
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} gin.H{data=dtos.AnnouncementListResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/announcements/get-announcements [post]
func (c *AnnouncementController) GetAll(ctx *gin.Context) {
	// Parse query parameters
	limit := 10
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	data, err := c.service.GetAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   data.Data,
		"limit":  data.Limit,
		"offset": data.Offset,
		"total":  data.Total,
	})
}

// GetByID retrieves Announcement by ID
// @Summary Get Announcement by ID
// @Description Retrieve announcement details by ID
// @Tags announcement
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.AnnouncementResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/announcements/get-announcement-by-id [post]
func (c *AnnouncementController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Announcement
// @Summary Update Announcement
// @Description Update Announcement details (all fields optional, including file uploads)
// @Tags announcement
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Announcement ID"
// @Param judul formData string false "Announcement title"
// @Param tanggal formData string false "Announcement date (YYYY-MM-DD)"
// @Param deskripsi formData string false "Announcement description"
// @Param penulis formData string false "Announcement author"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param gambar formData file false "Featured image (jpeg, png, gif, webp) - max 5MB (replaces existing)"
// @Param files formData file false "Attachment files - multiple files allowed (max 10MB each, replaces existing)"
// @Success 200 {object} gin.H{data=dtos.AnnouncementResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/announcements/update-announcement [post]
func (c *AnnouncementController) Update(ctx *gin.Context) {
	// Parse multipart form
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

	// Get optional fields
	judul := ctx.PostForm("judul")
	tanggal := ctx.PostForm("tanggal")
	deskripsi := ctx.PostForm("deskripsi")
	penulis := ctx.PostForm("penulis")
	statusPublikasi := ctx.PostForm("status_publikasi")
	status := ctx.PostForm("status")

	// Get gambar (optional)
	gambar, _ := ctx.FormFile("gambar")

	// Get files (optional, multiple)
	files := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["files"]; exists {
			files = uploadedFiles
		}
	}

	// Parse files_to_delete if provided
	var filesToDelete []string
	if filesDeleteStr := ctx.PostForm("files_to_delete"); filesDeleteStr != "" {
		_ = json.Unmarshal([]byte(filesDeleteStr), &filesToDelete)
	}

	// Create request DTO
	req := &dtos.AnnouncementUpdateRequest{
		ID:              uint(id),
		Judul:           judul,
		Tanggal:         tanggal,
		Deskripsi:       deskripsi,
		Penulis:         penulis,
		StatusPublikasi: statusPublikasi,
		Status:          status,
		FilesToDelete:   filesToDelete,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(uint(id), gambar, files, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Announcement by ID
// @Summary Delete Announcement
// @Description Delete Announcement by ID (also deletes all files from R2)
// @Tags announcement
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/announcements/delete-announcement [post]
func (c *AnnouncementController) Delete(ctx *gin.Context) {
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
		"message": "Announcement deleted successfully",
	})
}
