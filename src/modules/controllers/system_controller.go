package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// SystemController handles HTTP requests for System
type SystemController struct {
	service services.SystemService
}

// NewSystemController creates a new System controller
func NewSystemController(service services.SystemService) *SystemController {
	return &SystemController{service: service}
}

// Create creates a new System
func (c *SystemController) Create(ctx *gin.Context) {
	var req dtos.SystemCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	createdByID := userID.(uint)

	system := &models.System{
		Nama:        req.Nama,
		Description: req.Description,
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	if err := c.service.Create(system); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": system})
}

// GetByID retrieves System by ID
func (c *SystemController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all Systems with filters and pagination
func (c *SystemController) GetAll(ctx *gin.Context) {
	var req dtos.SystemGetAllRequest
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

	// Call service
	systems, total, err := c.service.GetAllWithFilter(repositories.GetSystemsParams{
		Filter: repositories.GetSystemsFilter{
			Nama:   req.Search.Nama,
			Status: req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to response
	var responseData []dtos.SystemResponse
	for _, system := range systems {
		responseData = append(responseData, dtos.SystemResponse{
			ID:          system.ID,
			Nama:        system.Nama,
			Description: system.Description,
			Status:      system.Status,
			CreatedAt:   system.CreatedAt,
			UpdatedAt:   system.UpdatedAt,
			CreatedByID: system.CreatedByID,
			UpdatedByID: system.UpdatedByID,
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	ctx.JSON(http.StatusOK, gin.H{
		"data": responseData,
		"pagination": gin.H{
			"limit":       limit,
			"offset":      offset,
			"page":        page,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// Update updates System
func (c *SystemController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID          uint   `json:"id" binding:"required"`
		Nama        string `json:"nama"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	if req.Nama != "" {
		data.Nama = req.Nama
	}
	if req.Description != "" {
		data.Description = req.Description
	}
	if req.Status != "" {
		data.Status = req.Status
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	updatedByID := userID.(uint)
	data.UpdatedByID = &updatedByID

	if err := c.service.Update(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes System by ID
func (c *SystemController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "System deleted successfully"})
}
