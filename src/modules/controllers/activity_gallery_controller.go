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

// ActivityGalleryController handles HTTP requests for ActivityGallery
type ActivityGalleryController struct {
	service services.ActivityGalleryService
}

// NewActivityGalleryController creates a new ActivityGallery controller
func NewActivityGalleryController(service services.ActivityGalleryService) *ActivityGalleryController {
	return &ActivityGalleryController{service: service}
}

// Create creates a new ActivityGallery with foto uploads
// @Summary Create new Activity Gallery
// @Description Create a new Activity Gallery with multiple foto uploads to Cloudflare R2
// @Tags activity-gallery
// @Accept multipart/form-data
// @Produce json
// @Param judul formData string true "Gallery title"
// @Param tanggal formData string true "Gallery date (YYYY-MM-DD)"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param foto formData file true "Foto files - multiple files allowed (jpeg, png, gif, webp) - max 10MB each"
// @Success 201 {object} gin.H{data=dtos.ActivityGalleryResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/create-gallery [post]
func (c *ActivityGalleryController) Create(ctx *gin.Context) {
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

	// Get optional fields
	statusPublikasi := ctx.PostForm("status_publikasi")
	status := ctx.PostForm("status")

	// Get fotos (required, multiple)
	fotos := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFotos, exists := form.File["foto"]; exists {
			fotos = uploadedFotos
		}
	}

	if len(fotos) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "at least one foto is required"})
		return
	}

	// Create request DTO
	req := &dtos.ActivityGalleryCreateRequest{
		Judul:           judul,
		Tanggal:         tanggal,
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
	data, err := c.service.Create(fotos, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all activity galleries
// @Summary Get all Activity Galleries
// @Description Retrieve all Activity Gallery records with pagination
// @Tags activity-gallery
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} gin.H{data=dtos.ActivityGalleryListResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/get-galleries [post]
func (c *ActivityGalleryController) GetAll(ctx *gin.Context) {
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

// GetByID retrieves ActivityGallery by ID
// @Summary Get Activity Gallery by ID
// @Description Retrieve activity gallery details by ID
// @Tags activity-gallery
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.ActivityGalleryResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/get-gallery-by-id [post]
func (c *ActivityGalleryController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity Gallery not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates ActivityGallery
// @Summary Update Activity Gallery
// @Description Update Activity Gallery details (all fields optional, including foto uploads)
// @Tags activity-gallery
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Activity Gallery ID"
// @Param judul formData string false "Gallery title"
// @Param tanggal formData string false "Gallery date (YYYY-MM-DD)"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param foto formData file false "Foto files - multiple files allowed (jpeg, png, gif, webp) - max 10MB each"
// @Success 200 {object} gin.H{data=dtos.ActivityGalleryResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/update-gallery [post]
func (c *ActivityGalleryController) Update(ctx *gin.Context) {
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
	statusPublikasi := ctx.PostForm("status_publikasi")
	status := ctx.PostForm("status")

	// Get fotos (optional, multiple)
	fotos := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFotos, exists := form.File["foto"]; exists {
			fotos = uploadedFotos
		}
	}

	// Parse foto_to_delete if provided
	var fotoToDelete []string
	if fotoDeleteStr := ctx.PostForm("foto_to_delete"); fotoDeleteStr != "" {
		_ = json.Unmarshal([]byte(fotoDeleteStr), &fotoToDelete)
	}

	// Create request DTO
	req := &dtos.ActivityGalleryUpdateRequest{
		ID:              uint(id),
		Judul:           judul,
		Tanggal:         tanggal,
		StatusPublikasi: statusPublikasi,
		Status:          status,
		FotoToDelete:    fotoToDelete,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(uint(id), fotos, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes ActivityGallery by ID
// @Summary Delete Activity Gallery
// @Description Delete Activity Gallery by ID (also deletes all fotos from R2)
// @Tags activity-gallery
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/delete-gallery [post]
func (c *ActivityGalleryController) Delete(ctx *gin.Context) {
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
		"message": "Activity Gallery deleted successfully",
	})
}
