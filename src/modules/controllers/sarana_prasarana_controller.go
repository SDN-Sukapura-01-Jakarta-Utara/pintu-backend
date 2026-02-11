package controllers

import (
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// SaranaPrasaranaController handles HTTP requests for SaranaPrasarana
type SaranaPrasaranaController struct {
	service services.SaranaPrasaranaService
}

// NewSaranaPrasaranaController creates a new SaranaPrasarana controller
func NewSaranaPrasaranaController(service services.SaranaPrasaranaService) *SaranaPrasaranaController {
	return &SaranaPrasaranaController{service: service}
}

// Create creates a new SaranaPrasarana with file upload
func (c *SaranaPrasaranaController) Create(ctx *gin.Context) {
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
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// Create request DTO
	req := &dtos.SaranaPrasaranaCreateRequest{
		Name:   name,
		Status: status,
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

// GetAll retrieves all SaranaPrasarana
func (c *SaranaPrasaranaController) GetAll(ctx *gin.Context) {
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

// GetByID retrieves SaranaPrasarana by ID
func (c *SaranaPrasaranaController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Sarana prasarana not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates SaranaPrasarana
func (c *SaranaPrasaranaController) Update(ctx *gin.Context) {
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
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")

	// Create request DTO
	req := &dtos.SaranaPrasaranaUpdateRequest{
		ID: id,
	}

	if name != "" {
		req.Name = &name
	}
	if status != "" {
		req.Status = &status
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

// Delete deletes SaranaPrasarana by ID
func (c *SaranaPrasaranaController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Sarana prasarana not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sarana prasarana deleted successfully",
	})
}
