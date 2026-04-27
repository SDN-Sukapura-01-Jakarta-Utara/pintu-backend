package controllers

import (
	"net/http"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// KritikSaranController handles HTTP requests for KritikSaran
type KritikSaranController struct {
	service services.KritikSaranService
}

// NewKritikSaranController creates a new KritikSaran controller
func NewKritikSaranController(service services.KritikSaranService) *KritikSaranController {
	return &KritikSaranController{service: service}
}

// CreatePublic creates a new KritikSaran (public endpoint, no auth required)
// @Summary Create new KritikSaran (public)
// @Description Create a new kritik saran from public (no authentication required)
// @Tags kritik_saran
// @Accept json
// @Produce json
// @Param body body dtos.KritikSaranCreateRequest true "Request body"
// @Success 201 {object} gin.H{message=string,data=dtos.KritikSaranResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/create-kritik-saran [post]
func (c *KritikSaranController) CreatePublic(ctx *gin.Context) {
	var req dtos.KritikSaranCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	data, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Kritik dan saran berhasil dikirim",
		"data":    data,
	})
}

// GetByID retrieves KritikSaran by ID (protected endpoint)
// @Summary Get KritikSaran by ID
// @Description Retrieve kritik saran details by ID
// @Tags kritik_saran
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.KritikSaranResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/kritik-saran/get-kritik-saran-by-id [post]
func (c *KritikSaranController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Kritik saran not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all kritik saran (protected endpoint)
// @Summary Get all KritikSaran
// @Description Retrieve all kritik saran records with date filters and pagination
// @Tags kritik_saran
// @Accept json
// @Produce json
// @Param body body dtos.KritikSaranGetAllRequest true "Request body with filters"
// @Success 200 {object} gin.H{data=dtos.KritikSaranListResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/kritik-saran/get-kritik-saran [post]
func (c *KritikSaranController) GetAll(ctx *gin.Context) {
	var req dtos.KritikSaranGetAllRequest
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
	data, err := c.service.GetAllWithFilter(repositories.GetKritikSaranParams{
		Filter: repositories.GetKritikSaranFilter{
			StartDate: startDate,
			EndDate:   endDate,
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes KritikSaran by ID (protected endpoint)
// @Summary Delete KritikSaran
// @Description Delete kritik saran by ID
// @Tags kritik_saran
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/kritik-saran/delete-kritik-saran [post]
func (c *KritikSaranController) Delete(ctx *gin.Context) {
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
		"message": "Kritik saran deleted successfully",
	})
}
