package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// RombelController handles HTTP requests for Rombel
type RombelController struct {
	service services.RombelService
}

// NewRombelController creates a new Rombel controller
func NewRombelController(service services.RombelService) *RombelController {
	return &RombelController{service: service}
}

// Create creates a new Rombel
// @Summary Create new Rombel
// @Description Create a new Rombel with name, status, and optional kelas_id
// @Tags rombel
// @Accept json
// @Produce json
// @Param body body dtos.RombelCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.RombelResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/rombel/create-rombel [post]
func (c *RombelController) Create(ctx *gin.Context) {
	var req dtos.RombelCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Create(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all rombel
// @Summary Get all Rombel
// @Description Retrieve all Rombel records
// @Tags rombel
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.RombelResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/rombel/get-rombel [post]
func (c *RombelController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll(10, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data})
}

// GetByID retrieves Rombel by ID
// @Summary Get Rombel by ID
// @Description Retrieve rombel details by ID
// @Tags rombel
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.RombelResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/rombel/get-rombel-by-id [post]
func (c *RombelController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rombel not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Rombel
// @Summary Update Rombel
// @Description Update Rombel details
// @Tags rombel
// @Accept json
// @Produce json
// @Param body body dtos.RombelUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.RombelResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/rombel/update-rombel [post]
func (c *RombelController) Update(ctx *gin.Context) {
	var req dtos.RombelUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Rombel by ID
// @Summary Delete Rombel
// @Description Delete Rombel by ID
// @Tags rombel
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/rombel/delete-rombel [post]
func (c *RombelController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rombel not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Rombel deleted successfully",
	})
}
