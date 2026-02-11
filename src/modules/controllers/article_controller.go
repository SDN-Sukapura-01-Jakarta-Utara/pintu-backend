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

// GetAll retrieves all articles
// @Summary Get all Articles
// @Description Retrieve all Article records with pagination
// @Tags article
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} gin.H{data=dtos.ArticleListResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/articles/get-articles [post]
func (c *ArticleController) GetAll(ctx *gin.Context) {
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
