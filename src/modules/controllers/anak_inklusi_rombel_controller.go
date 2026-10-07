package controllers

import (
	"fmt"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

type AnakInklusiRombelController struct {
	service services.AnakInklusiRombelService
}

func NewAnakInklusiRombelController(service services.AnakInklusiRombelService) *AnakInklusiRombelController {
	return &AnakInklusiRombelController{service: service}
}

// SyncByTahunPelajaran syncs anak inklusi rombel data based on tahun pelajaran
// @Summary Sync anak inklusi rombel by tahun pelajaran
// @Description Automatically sync anak_inklusi_rombel based on peserta_didik_rombel and anak_inklusi data for specified tahun pelajaran
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.AnakInklusiRombelSyncRequest true "Request body with tahun_pelajaran_id"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiRombelSyncResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/sync-data-induk-inklusi-rombel [post]
func (c *AnakInklusiRombelController) SyncByTahunPelajaran(ctx *gin.Context) {
	var req dtos.AnakInklusiRombelSyncRequest
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

	data, err := c.service.SyncByTahunPelajaran(req.TahunPelajaranID, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all anak inklusi rombel with filters and pagination
// @Summary Get all anak inklusi rombel
// @Description Retrieve all anak inklusi rombel records with filters and pagination by tahun pelajaran
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.AnakInklusiRombelGetAllRequest true "Request body with filters"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiRombelListWithPaginationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-data-induk-inklusi-rombel [post]
func (c *AnakInklusiRombelController) GetAll(ctx *gin.Context) {
	var req dtos.AnakInklusiRombelGetAllRequest
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

	// Call service with filters
	rombelIDStr := ""
	if req.Search.RombelID > 0 {
		rombelIDStr = fmt.Sprint(req.Search.RombelID)
	}
	
	data, err := c.service.GetAllWithFilter(repositories.GetAnakInklusiRombelParams{
		Filter: repositories.GetAnakInklusiRombelFilter{
			TahunPelajaranID: fmt.Sprint(req.TahunPelajaranID),
			NamaPesertaDidik: req.Search.NamaPesertaDidik,
			NIS:              req.Search.NIS,
			NISN:             req.Search.NISN,
			JenisHambatan:    req.Search.JenisHambatan,
			Status:           req.Search.Status,
			RombelID:         rombelIDStr,
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

// GetByID retrieves anak inklusi rombel by ID
// @Summary Get anak inklusi rombel by ID
// @Description Retrieve anak inklusi rombel details by ID
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiRombelResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-data-induk-inklusi-rombel-by-id [post]
func (c *AnakInklusiRombelController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data anak inklusi rombel tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates anak inklusi rombel
// @Summary Update anak inklusi rombel
// @Description Update anak inklusi rombel (guru_pendamping_khusus_id, guru_kelas_id, catatan)
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.AnakInklusiRombelUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.AnakInklusiRombelResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/update-data-induk-inklusi-rombel [post]
func (c *AnakInklusiRombelController) Update(ctx *gin.Context) {
	var req dtos.AnakInklusiRombelUpdateRequest
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

	data, err := c.service.Update(req.ID, &req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes anak inklusi rombel by ID
// @Summary Delete anak inklusi rombel
// @Description Delete anak inklusi rombel by ID (soft delete)
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/delete-data-induk-inklusi-rombel [post]
func (c *AnakInklusiRombelController) Delete(ctx *gin.Context) {
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
		"message": "Data anak inklusi rombel berhasil dihapus",
	})
}
