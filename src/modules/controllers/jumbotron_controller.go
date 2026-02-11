package controllers

import (
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// JumbotronController handles HTTP requests for Jumbotron
type JumbotronController struct {
	service services.JumbotronService
}

// NewJumbotronController creates a new Jumbotron controller
func NewJumbotronController(service services.JumbotronService) *JumbotronController {
	return &JumbotronController{service: service}
}

// Create creates a new Jumbotron with file upload
// @Summary Create new Jumbotron
// @Description Create a new Jumbotron with file upload to Cloudflare R2
// @Tags jumbotron
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file (jpeg, png, gif, webp)"
// @Param status formData string false "Status (active/inactive)"
// @Success 201 {object} gin.H{data=dtos.JumbotronResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/jumbotron/create-jumbotron [post]
func (c *JumbotronController) Create(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(5 * 1024 * 1024); err != nil { // 5MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Get status from form (optional)
	status := ctx.PostForm("status")

	// Create request DTO
	req := &dtos.JumbotronCreateRequest{
		Status: status,
	}

	// Get user ID from context (set by auth middleware)
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

// GetAll retrieves all jumbotron
// @Summary Get all Jumbotron
// @Description Retrieve all Jumbotron records with pagination
// @Tags jumbotron
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} gin.H{data=dtos.JumbotronListResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/jumbotron/get-jumbotron [post]
func (c *JumbotronController) GetAll(ctx *gin.Context) {
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

// GetByID retrieves Jumbotron by ID
// @Summary Get Jumbotron by ID
// @Description Retrieve jumbotron details by ID
// @Tags jumbotron
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.JumbotronResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/jumbotron/get-jumbotron-by-id [post]
func (c *JumbotronController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Jumbotron not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Jumbotron
// @Summary Update Jumbotron
// @Description Update Jumbotron details (status and/or file)
// @Tags jumbotron
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Jumbotron ID"
// @Param file formData file false "Image file (jpeg, png, gif, webp) - optional"
// @Param status formData string false "Status (active/inactive) - optional"
// @Success 200 {object} gin.H{data=dtos.JumbotronResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/jumbotron/update-jumbotron [post]
func (c *JumbotronController) Update(ctx *gin.Context) {
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

	// Get status (optional)
	status := ctx.PostForm("status")

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.UpdateWithFile(id, file, status, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Helper function to convert string to int
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Delete deletes Jumbotron by ID
// @Summary Delete Jumbotron
// @Description Delete Jumbotron by ID (also deletes file from R2)
// @Tags jumbotron
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/jumbotron/delete-jumbotron [post]
func (c *JumbotronController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Jumbotron not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Jumbotron deleted successfully",
	})
}
