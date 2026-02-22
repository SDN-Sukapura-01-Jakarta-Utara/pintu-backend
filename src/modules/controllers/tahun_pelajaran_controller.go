package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// TahunPelajaranController handles HTTP requests for TahunPelajaran
type TahunPelajaranController struct {
	service services.TahunPelajaranService
}

// NewTahunPelajaranController creates a new TahunPelajaran controller
func NewTahunPelajaranController(service services.TahunPelajaranService) *TahunPelajaranController {
	return &TahunPelajaranController{service: service}
}

// Create creates a new TahunPelajaran
// @Summary Create new TahunPelajaran
// @Description Create a new TahunPelajaran with tahun_pelajaran and status
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Param body body dtos.TahunPelajaranCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.TahunPelajaranResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/tahun-pelajaran/create-tahun-pelajaran [post]
func (c *TahunPelajaranController) Create(ctx *gin.Context) {
	var req dtos.TahunPelajaranCreateRequest
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

// GetAll retrieves all tahun_pelajaran
// @Summary Get all TahunPelajaran
// @Description Retrieve all TahunPelajaran records
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.TahunPelajaranResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/tahun-pelajaran/get-tahun-pelajaran [post]
func (c *TahunPelajaranController) GetAll(ctx *gin.Context) {
	var req dtos.TahunPelajaranGetAllRequest
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
	data, err := c.service.GetAllWithFilter(repositories.GetTahunPelajaranParams{
		Filter: repositories.GetTahunPelajaranFilter{
			TahunPelajaran: req.Search.TahunPelajaran,
			Status:         req.Search.Status,
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

// GetByID retrieves TahunPelajaran by ID
// @Summary Get TahunPelajaran by ID
// @Description Retrieve tahun_pelajaran details by ID
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.TahunPelajaranResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/tahun-pelajaran/get-tahun-pelajaran-by-id [post]
func (c *TahunPelajaranController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tahun pelajaran not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates TahunPelajaran
// @Summary Update TahunPelajaran
// @Description Update TahunPelajaran details
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Param body body dtos.TahunPelajaranUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.TahunPelajaranResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/tahun-pelajaran/update-tahun-pelajaran [post]
func (c *TahunPelajaranController) Update(ctx *gin.Context) {
	var req dtos.TahunPelajaranUpdateRequest
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

// Delete deletes TahunPelajaran by ID
// @Summary Delete TahunPelajaran
// @Description Delete TahunPelajaran by ID
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/tahun-pelajaran/delete-tahun-pelajaran [post]
func (c *TahunPelajaranController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tahun pelajaran not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tahun pelajaran deleted successfully",
	})
}
