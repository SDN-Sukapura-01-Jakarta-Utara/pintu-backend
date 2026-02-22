package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// BidangStudiController handles HTTP requests for BidangStudi
type BidangStudiController struct {
	service services.BidangStudiService
}

// NewBidangStudiController creates a new BidangStudi controller
func NewBidangStudiController(service services.BidangStudiService) *BidangStudiController {
	return &BidangStudiController{service: service}
}

// Create creates a new BidangStudi
// @Summary Create new BidangStudi
// @Description Create a new BidangStudi with name and status
// @Tags bidang_studi
// @Accept json
// @Produce json
// @Param body body dtos.BidangStudiCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.BidangStudiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/bidang-studi/create-bidang-studi [post]
func (c *BidangStudiController) Create(ctx *gin.Context) {
	var req dtos.BidangStudiCreateRequest
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

// GetAll retrieves all bidang_studi
// @Summary Get all BidangStudi
// @Description Retrieve all BidangStudi records
// @Tags bidang_studi
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.BidangStudiResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/bidang-studi/get-bidang-studi [post]
func (c *BidangStudiController) GetAll(ctx *gin.Context) {
	var req dtos.BidangStudiGetAllRequest
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

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetBidangStudiParams{
		Filter: repositories.GetBidangStudiFilter{
			Name:   req.Search.Name,
			Status: req.Search.Status,
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

// GetByID retrieves BidangStudi by ID
// @Summary Get BidangStudi by ID
// @Description Retrieve bidang_studi details by ID
// @Tags bidang_studi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.BidangStudiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/bidang-studi/get-bidang-studi-by-id [post]
func (c *BidangStudiController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Bidang studi not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates BidangStudi
// @Summary Update BidangStudi
// @Description Update BidangStudi details
// @Tags bidang_studi
// @Accept json
// @Produce json
// @Param body body dtos.BidangStudiUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.BidangStudiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/bidang-studi/update-bidang-studi [post]
func (c *BidangStudiController) Update(ctx *gin.Context) {
	var req dtos.BidangStudiUpdateRequest
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

// Delete deletes BidangStudi by ID
// @Summary Delete BidangStudi
// @Description Delete BidangStudi by ID
// @Tags bidang_studi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/bidang-studi/delete-bidang-studi [post]
func (c *BidangStudiController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Bidang studi not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Bidang studi deleted successfully",
	})
}
