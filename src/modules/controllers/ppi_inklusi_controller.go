package controllers

import (
	"fmt"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

type PPIInklusiController struct {
	service services.PPIInklusiService
}

func NewPPIInklusiController(service services.PPIInklusiService) *PPIInklusiController {
	return &PPIInklusiController{service: service}
}

// Create creates new PPI (Program Pembelajaran Individual) in bulk
// @Summary Create new PPI in bulk
// @Description Create multiple Program Pembelajaran Individual for anak inklusi in one request
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.PPIInklusiCreateRequest true "PPI data with multiple programs"
// @Success 201 {object} gin.H{data=dtos.PPIInklusiCreateResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/create-ppi [post]
func (c *PPIInklusiController) Create(ctx *gin.Context) {
	var req dtos.PPIInklusiCreateRequest
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

	data, err := c.service.Create(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all PPI with filters and pagination
// @Summary Get all PPI
// @Description Retrieve all PPI records with filters and pagination
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.PPIInklusiGetAllRequest true "Request body with filters"
// @Success 200 {object} gin.H{data=dtos.PPIInklusiListWithPaginationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-ppi [post]
func (c *PPIInklusiController) GetAll(ctx *gin.Context) {
	var req dtos.PPIInklusiGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values for pagination
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Prepare filter strings
	tahunPelajaranIDStr := ""
	if req.Search.TahunPelajaranID > 0 {
		tahunPelajaranIDStr = fmt.Sprint(req.Search.TahunPelajaranID)
	}

	rombelIDStr := ""
	if req.Search.RombelID > 0 {
		rombelIDStr = fmt.Sprint(req.Search.RombelID)
	}

	guruPembuatIDStr := ""
	if req.Search.GuruPembuatID > 0 {
		guruPembuatIDStr = fmt.Sprint(req.Search.GuruPembuatID)
	}

	anakInklusiRombelIDStr := ""
	if req.Search.AnakInklusiRombelID > 0 {
		anakInklusiRombelIDStr = fmt.Sprint(req.Search.AnakInklusiRombelID)
	}

	aspekPembelajaran := req.Search.AspekPembelajaran

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetPPIInklusiParams{
		Filter: repositories.GetPPIInklusiFilter{
			TahunPelajaranID:    tahunPelajaranIDStr,
			RombelID:            rombelIDStr,
			GuruPembuatID:       guruPembuatIDStr,
			AnakInklusiRombelID: anakInklusiRombelIDStr,
			AspekPembelajaran:   aspekPembelajaran,
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

// GetByID retrieves PPI by ID
// @Summary Get PPI by ID
// @Description Retrieve PPI details by ID with complete data
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PPIInklusiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-ppi-by-id [post]
func (c *PPIInklusiController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data PPI tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates PPI data
// @Summary Update PPI
// @Description Update PPI (Program Pembelajaran Individual) data
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.PPIInklusiUpdateRequest true "Request body with PPI data to update"
// @Success 200 {object} gin.H{data=dtos.PPIInklusiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/update-ppi [post]
func (c *PPIInklusiController) Update(ctx *gin.Context) {
	var req dtos.PPIInklusiUpdateRequest
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

// Delete deletes PPI by ID
// @Summary Delete PPI
// @Description Delete PPI (Program Pembelajaran Individual) by ID (soft delete)
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/delete-ppi [post]
func (c *PPIInklusiController) Delete(ctx *gin.Context) {
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
		"message": "Data PPI berhasil dihapus",
	})
}

// GetSummary retrieves PPI summary by tahun pelajaran
// @Summary Get PPI summary
// @Description Retrieve PPI summary statistics by tahun pelajaran including total, status breakdown, and bidang studi breakdown
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.PPIInklusiSummaryRequest true "Request body with tahun_pelajaran_id"
// @Success 200 {object} gin.H{data=dtos.PPIInklusiSummaryResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/summary-ppi [post]
func (c *PPIInklusiController) GetSummary(ctx *gin.Context) {
	var req dtos.PPIInklusiSummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetSummary(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
