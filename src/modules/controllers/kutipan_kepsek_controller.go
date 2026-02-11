package controllers

import (
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// KutipanKepsekController handles HTTP requests for KutipanKepsek
type KutipanKepsekController struct {
	service services.KutipanKepsekService
}

// NewKutipanKepsekController creates a new KutipanKepsek controller
func NewKutipanKepsekController(service services.KutipanKepsekService) *KutipanKepsekController {
	return &KutipanKepsekController{service: service}
}

// Create creates a new KutipanKepsek with file upload
func (c *KutipanKepsekController) Create(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Get form data
	namaKepsek := ctx.PostForm("nama_kepsek")
	kutipanKepsek := ctx.PostForm("kutipan_kepsek")

	if namaKepsek == "" || kutipanKepsek == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nama_kepsek and kutipan_kepsek are required"})
		return
	}

	// Create request DTO
	req := &dtos.KutipanKepsekCreateRequest{
		NamaKepsek:    namaKepsek,
		KutipanKepsek: kutipanKepsek,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.Create(file, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all KutipanKepsek
func (c *KutipanKepsekController) GetAll(ctx *gin.Context) {
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

// GetByID retrieves KutipanKepsek by ID
func (c *KutipanKepsekController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Kutipan kepsek not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates KutipanKepsek
func (c *KutipanKepsekController) Update(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID from form
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var id uint
	if _, err := strconv.Atoi(idStr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}
	id = uint(toInt(idStr))

	// Get file (optional)
	file, _ := ctx.FormFile("file")

	// Get form data (optional)
	namaKepsek := ctx.PostForm("nama_kepsek")
	kutipanKepsek := ctx.PostForm("kutipan_kepsek")

	// Create request DTO
	req := &dtos.KutipanKepsekUpdateRequest{
		ID: id,
	}

	if namaKepsek != "" {
		req.NamaKepsek = &namaKepsek
	}
	if kutipanKepsek != "" {
		req.KutipanKepsek = &kutipanKepsek
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.UpdateWithFile(id, file, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes KutipanKepsek by ID
func (c *KutipanKepsekController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Kutipan kepsek not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kutipan kepsek deleted successfully",
	})
}
