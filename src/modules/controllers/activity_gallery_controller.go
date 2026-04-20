package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
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

	// Get foto thumbnail info (optional)
	var fotoThumbnails []string
	if fotoThumbnailStr := ctx.PostForm("foto_thumbnails"); fotoThumbnailStr != "" {
		if err := json.Unmarshal([]byte(fotoThumbnailStr), &fotoThumbnails); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid foto_thumbnails JSON format"})
			return
		}
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
	data, err := c.service.Create(fotos, fotoThumbnails, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all activity galleries with filters and pagination
// @Summary Get all Activity Galleries
// @Description Retrieve all Activity Gallery records with filters and pagination
// @Tags activity-gallery
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.ActivityGalleryListWithPaginationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/activity-galleries/get-galleries [post]
func (c *ActivityGalleryController) GetAll(ctx *gin.Context) {
	var req dtos.ActivityGalleryGetAllRequest
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

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetActivityGalleryParams{
		Filter: repositories.GetActivityGalleryFilter{
			Judul:            req.Search.Judul,
			StartDate:        startDate,
			EndDate:          endDate,
			StatusPublikasi:  req.Search.StatusPublikasi,
			Status:           req.Search.Status,
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

	// Get existing foto IDs and their thumbnail updates
	var existingFotoIds []string
	if existingFotoIdsStr := ctx.PostForm("existing_foto_ids"); existingFotoIdsStr != "" {
		if err := json.Unmarshal([]byte(existingFotoIdsStr), &existingFotoIds); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid existing_foto_ids JSON format"})
			return
		}
	}

	// Get existing foto thumbnails (for existing fotos)
	var existingFotoThumbnails []string
	if existingFotoThumbnailsStr := ctx.PostForm("existing_foto_thumbnails"); existingFotoThumbnailsStr != "" {
		if err := json.Unmarshal([]byte(existingFotoThumbnailsStr), &existingFotoThumbnails); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid existing_foto_thumbnails JSON format"})
			return
		}
	}

	// Get new foto thumbnails (for new uploaded fotos)
	var newFotoThumbnails []string
	if newFotoThumbnailsStr := ctx.PostForm("new_foto_thumbnails"); newFotoThumbnailsStr != "" {
		if err := json.Unmarshal([]byte(newFotoThumbnailsStr), &newFotoThumbnails); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_foto_thumbnails JSON format"})
			return
		}
	}

	// Build foto_thumbnail_updates map from existing_foto_ids and existing_foto_thumbnails
	fotoThumbnailUpdates := make(map[string]string)
	if len(existingFotoIds) > 0 && len(existingFotoThumbnails) > 0 {
		for i, fotoID := range existingFotoIds {
			if i < len(existingFotoThumbnails) {
				fotoThumbnailUpdates[fotoID] = existingFotoThumbnails[i]
			}
		}
	}

	// Also support direct foto_thumbnail_updates map format (backward compatibility)
	if fotoThumbnailUpdatesStr := ctx.PostForm("foto_thumbnail_updates"); fotoThumbnailUpdatesStr != "" {
		var directUpdates map[string]string
		if err := json.Unmarshal([]byte(fotoThumbnailUpdatesStr), &directUpdates); err == nil {
			// Merge with existing updates
			for k, v := range directUpdates {
				fotoThumbnailUpdates[k] = v
			}
		}
	}

	// Parse foto_to_delete if provided
	var fotoToDelete []string
	if fotoDeleteStr := ctx.PostForm("foto_to_delete"); fotoDeleteStr != "" {
		_ = json.Unmarshal([]byte(fotoDeleteStr), &fotoToDelete)
	}

	// Create request DTO
	req := &dtos.ActivityGalleryUpdateRequest{
		ID:                   uint(id),
		Judul:                judul,
		Tanggal:              tanggal,
		StatusPublikasi:      statusPublikasi,
		Status:               status,
		FotoToDelete:         fotoToDelete,
		FotoThumbnailUpdates: fotoThumbnailUpdates,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(uint(id), fotos, newFotoThumbnails, req, userID.(uint))
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


// GetPublicLatest retrieves 10 latest published and active activity galleries for public display (no auth required)
// @Summary Get latest Activity Galleries for public
// @Description Retrieve 10 latest published and active activity galleries ordered by tanggal DESC (no authentication required)
// @Tags activity-gallery
// @Accept json
// @Produce json
// @Success 200 {object} dtos.ActivityGalleryPublicListResponse
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/public/get-data-galeri-kegiatan [post]
func (c *ActivityGalleryController) GetPublicLatest(ctx *gin.Context) {
	data, err := c.service.GetPublicLatest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
