package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// EkstrakurikulerController handles HTTP requests for Ekstrakurikuler
type EkstrakurikulerController struct {
	service services.EkstrakurikulerService
}

// NewEkstrakurikulerController creates a new Ekstrakurikuler controller
func NewEkstrakurikulerController(service services.EkstrakurikulerService) *EkstrakurikulerController {
	return &EkstrakurikulerController{service: service}
}

// Create creates a new Ekstrakurikuler
// @Summary Create new Ekstrakurikuler
// @Description Create a new Ekstrakurikuler with name, kelas_ids, kategori, and status
// @Tags ekstrakurikuler
// @Accept json
// @Produce json
// @Param body body dtos.EkstrakurikulerCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.EkstrakurikulerResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/create-ekstrakurikuler [post]
func (c *EkstrakurikulerController) Create(ctx *gin.Context) {
	var req dtos.EkstrakurikulerCreateRequest
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

// GetAll retrieves all ekstrakurikuler
// @Summary Get all Ekstrakurikuler
// @Description Retrieve all Ekstrakurikuler records
// @Tags ekstrakurikuler
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.EkstrakurikulerResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/get-ekstrakurikuler [post]
func (c *EkstrakurikulerController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll(10, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data})
}

// GetByID retrieves Ekstrakurikuler by ID
// @Summary Get Ekstrakurikuler by ID
// @Description Retrieve ekstrakurikuler details by ID
// @Tags ekstrakurikuler
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.EkstrakurikulerResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/get-ekstrakurikuler-by-id [post]
func (c *EkstrakurikulerController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Ekstrakurikuler not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Ekstrakurikuler
// @Summary Update Ekstrakurikuler
// @Description Update Ekstrakurikuler details
// @Tags ekstrakurikuler
// @Accept json
// @Produce json
// @Param body body dtos.EkstrakurikulerUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.EkstrakurikulerResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/update-ekstrakurikuler [post]
func (c *EkstrakurikulerController) Update(ctx *gin.Context) {
	var req dtos.EkstrakurikulerUpdateRequest
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

// Delete deletes Ekstrakurikuler by ID
// @Summary Delete Ekstrakurikuler
// @Description Delete Ekstrakurikuler by ID
// @Tags ekstrakurikuler
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/delete-ekstrakurikuler [post]
func (c *EkstrakurikulerController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Ekstrakurikuler not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ekstrakurikuler deleted successfully",
	})
}
