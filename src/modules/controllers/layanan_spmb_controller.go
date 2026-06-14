package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// LayananSPMBController handles HTTP requests for Layanan SPMB
type LayananSPMBController struct {
	service services.LayananSPMBService
}

// NewLayananSPMBController creates a new Layanan SPMB controller
func NewLayananSPMBController(service services.LayananSPMBService) *LayananSPMBController {
	return &LayananSPMBController{service: service}
}

// CreatePublic creates a new Layanan SPMB from public form (no auth required)
// @Summary Create new Layanan SPMB (Public)
// @Description Create a new layanan SPMB dari form public untuk orang tua murid
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.LayananSPMBCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.LayananSPMBResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/layanan-spmb [post]
func (c *LayananSPMBController) CreatePublic(ctx *gin.Context) {
	var req dtos.LayananSPMBCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.CreatePublic(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all layanan SPMB with filters and pagination (auth required)
// @Summary Get all Layanan SPMB with filters
// @Description Retrieve all layanan SPMB dengan filters, sorting by status (pending first) then tanggal laporan (DESC)
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.LayananSPMBGetAllRequest true "Request body with filters and pagination"
// @Success 200 {object} dtos.LayananSPMBListWithPaginationResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb/get-layanan-spmb [post]
func (c *LayananSPMBController) GetAll(ctx *gin.Context) {
	var req dtos.LayananSPMBGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetAllWithFilter(&req)
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

// GetByID retrieves layanan SPMB by ID (auth required)
// @Summary Get Layanan SPMB by ID
// @Description Retrieve layanan SPMB details by ID
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.LayananSPMBResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb/get-layanan-spmb-by-id [post]
func (c *LayananSPMBController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// UpdateStatus updates status layanan SPMB (auth required)
// @Summary Update status Layanan SPMB
// @Description Update status layanan SPMB (pending/selesai)
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.LayananSPMBUpdateStatusRequest true "Request body with ID and status"
// @Success 200 {object} gin.H{message=string,data=dtos.LayananSPMBResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb/set-status-selesai [post]
func (c *LayananSPMBController) UpdateStatus(ctx *gin.Context) {
	var req dtos.LayananSPMBUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.UpdateStatus(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Status berhasil diupdate",
		"data":    data,
	})
}

// DeleteLayananSPMB soft deletes layanan SPMB (auth required)
// @Summary Delete Layanan SPMB
// @Description Soft delete layanan SPMB by setting deleted_at
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Layanan SPMB ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb/delete-layanan-spmb [post]
func (c *LayananSPMBController) DeleteLayananSPMB(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err := c.service.DeleteLayananSPMB(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Layanan SPMB berhasil dihapus",
	})
}
