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

// ArticleController handles HTTP requests for Article
type ArticleController struct {
	service services.ArticleService
}

// NewArticleController creates a new Article controller
func NewArticleController(service services.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

// Create creates a new Article with file uploads
// @Summary Create new Article
// @Description Create a new Article with gambar and files upload to Cloudflare R2
// @Tags article
// @Accept multipart/form-data
// @Produce json
// @Param judul formData string true "Article title"
// @Param tanggal formData string true "Article date (YYYY-MM-DD)"
// @Param kategori formData string true "Article category"
// @Param deskripsi formData string false "Article description"
// @Param penulis formData string true "Article author"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param gambar formData file false "Featured image (jpeg, png, gif, webp) - max 5MB"
// @Param files formData file false "Attachment files - multiple files allowed (max 10MB each)"
// @Success 201 {object} gin.H{data=dtos.ArticleResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/articles/create-article [post]
func (c *ArticleController) Create(ctx *gin.Context) {
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

	// Get kategori
	kategori := ctx.PostForm("kategori")
	if kategori == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "kategori is required"})
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
	req := &dtos.ArticleCreateRequest{
		Judul:           judul,
		Tanggal:         tanggal,
		Kategori:        kategori,
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

// GetAll retrieves all articles with filters and pagination
// @Summary Get all Articles
// @Description Retrieve all Article records with filters and pagination
// @Tags article
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.ArticleListWithPaginationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/articles/get-articles [post]
func (c *ArticleController) GetAll(ctx *gin.Context) {
	var req dtos.ArticleGetAllRequest
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
	data, err := c.service.GetAllWithFilter(repositories.GetArticleParams{
		Filter: repositories.GetArticleFilter{
			Judul:           req.Search.Judul,
			StartDate:       startDate,
			EndDate:         endDate,
			Kategori:        req.Search.Kategori,
			Penulis:         req.Search.Penulis,
			StatusPublikasi: req.Search.StatusPublikasi,
			Status:          req.Search.Status,
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

// GetByID retrieves Article by ID
// @Summary Get Article by ID
// @Description Retrieve article details by ID
// @Tags article
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.ArticleResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/articles/get-article-by-id [post]
func (c *ArticleController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Article
// @Summary Update Article
// @Description Update Article details (all fields optional, including file uploads)
// @Tags article
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Article ID"
// @Param judul formData string false "Article title"
// @Param tanggal formData string false "Article date (YYYY-MM-DD)"
// @Param kategori formData string false "Article category"
// @Param deskripsi formData string false "Article description"
// @Param penulis formData string false "Article author"
// @Param status_publikasi formData string false "Publication status (draft/published/archived)"
// @Param status formData string false "Status (active/inactive)"
// @Param gambar formData file false "Featured image (jpeg, png, gif, webp) - max 5MB (replaces existing)"
// @Param files formData file false "Attachment files - multiple files allowed (max 10MB each, replaces existing)"
// @Success 200 {object} gin.H{data=dtos.ArticleResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/articles/update-article [post]
func (c *ArticleController) Update(ctx *gin.Context) {
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
	kategori := ctx.PostForm("kategori")
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
	req := &dtos.ArticleUpdateRequest{
		ID:              uint(id),
		Judul:           judul,
		Tanggal:         tanggal,
		Kategori:        kategori,
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

// Delete deletes Article by ID
// @Summary Delete Article
// @Description Delete Article by ID (also deletes all files from R2)
// @Tags article
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/articles/delete-article [post]
func (c *ArticleController) Delete(ctx *gin.Context) {
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
		"message": "Article deleted successfully",
	})
}

// GetPublicLatest retrieves 10 latest published and active articles for public display (no auth required)
// @Summary Get latest Articles for public
// @Description Retrieve 10 latest published and active articles ordered by tanggal DESC (no authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {object} dtos.ArticlePublicListResponse
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/public/get-data-artikel [post]
func (c *ArticleController) GetPublicLatest(ctx *gin.Context) {
	data, err := c.service.GetPublicLatest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
