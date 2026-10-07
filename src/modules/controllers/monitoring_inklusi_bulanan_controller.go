package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

type MonitoringInklusiBulananController struct {
	service services.MonitoringInklusiBulananService
}

func NewMonitoringInklusiBulananController(service services.MonitoringInklusiBulananService) *MonitoringInklusiBulananController {
	return &MonitoringInklusiBulananController{service: service}
}

// Create creates new monitoring inklusi bulanan
// @Summary Create new monitoring inklusi bulanan
// @Description Create monthly monitoring for inclusive students
// @Tags monitoring-inklusi
// @Accept multipart/form-data
// @Produce json
// @Param ppi_inklusi_id formData integer true "PPI Inklusi ID"
// @Param bulan formData integer true "Bulan (1-12)"
// @Param tahun formData integer true "Tahun"
// @Param deskripsi_perkembangan formData string true "Deskripsi Perkembangan"
// @Param kendala_ditemui formData string false "Kendala Ditemui"
// @Param tindak_lanjut formData string false "Tindak Lanjut"
// @Param file_pendukung formData file false "File Pendukung (multiple files)"
// @Param guru_pengisi_id formData integer false "Guru Pengisi ID"
// @Success 201 {object} gin.H{data=dtos.MonitoringInklusiBulananResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/create-monitoring-bulanan [post]
func (c *MonitoringInklusiBulananController) Create(ctx *gin.Context) {
	// Parse form data
	var req dtos.MonitoringInklusiBulananCreateRequest

	// Parse ppi_inklusi_id
	ppiInklusiID, err := strconv.ParseUint(ctx.PostForm("ppi_inklusi_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ppi_inklusi_id harus berupa angka"})
		return
	}
	req.PPIInklusiID = uint(ppiInklusiID)

	// Parse bulan
	bulan, err := strconv.ParseInt(ctx.PostForm("bulan"), 10, 16)
	if err != nil || bulan < 1 || bulan > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bulan harus berupa angka 1-12"})
		return
	}
	req.Bulan = int16(bulan)

	// Parse tahun
	tahun, err := strconv.ParseInt(ctx.PostForm("tahun"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tahun harus berupa angka"})
		return
	}
	req.Tahun = int16(tahun)

	// Parse deskripsi_perkembangan
	req.DeskripsiPerkembangan = ctx.PostForm("deskripsi_perkembangan")
	if req.DeskripsiPerkembangan == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deskripsi_perkembangan harus diisi"})
		return
	}

	// Parse optional fields
	req.KendalaDitemui = ctx.PostForm("kendala_ditemui")
	req.TindakLanjut = ctx.PostForm("tindak_lanjut")

	// Parse guru_pengisi_id (optional)
	if guruPengisiIDStr := ctx.PostForm("guru_pengisi_id"); guruPengisiIDStr != "" {
		guruPengisiID, err := strconv.ParseUint(guruPengisiIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "guru_pengisi_id harus berupa angka"})
			return
		}
		guruID := uint(guruPengisiID)
		req.GuruPengisiID = &guruID
	}

	// Get file pendukung (multiple files)
	form, err := ctx.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil && form != nil && form.File["file_pendukung"] != nil {
		files = form.File["file_pendukung"]
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Create(files, &req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetByBulanTahunRombel retrieves all monitoring for specific month, year, and optional rombel
// @Summary Get monitoring inklusi bulanan by bulan, tahun, optional rombel
// @Description Retrieve all monthly monitoring records for specific month, year with optional rombel and PPI filter
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.MonitoringInklusiBulananGetRequest true "Request body with bulan, tahun, optional rombel_id and ppi_inklusi_id"
// @Success 200 {object} gin.H{data=dtos.MonitoringInklusiBulananGetAllResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-monitoring-bulanan [post]
func (c *MonitoringInklusiBulananController) GetByBulanTahunRombel(ctx *gin.Context) {
	var req dtos.MonitoringInklusiBulananGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Expected JSON with fields: bulan (1-12), tahun (year), optional rombel_id (number), optional ppi_inklusi_id (number). Detail: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Bulan < 1 || req.Bulan > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bulan harus antara 1-12"})
		return
	}
	if req.Tahun == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tahun harus diisi"})
		return
	}

	data, err := c.service.GetByBulanTahunRombel(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetByID retrieves monitoring inklusi bulanan by ID
// @Summary Get monitoring inklusi bulanan by ID
// @Description Retrieve monitoring inklusi bulanan details by ID with complete data
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.MonitoringInklusiBulananResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-monitoring-bulanan-by-id [post]
func (c *MonitoringInklusiBulananController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format for get-monitoring-bulanan-by-id. Expected JSON with field: id (number). Detail: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id harus diisi"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates monitoring inklusi bulanan
// @Summary Update monitoring inklusi bulanan
// @Description Update monthly monitoring for inclusive students including file management
// @Tags monitoring-inklusi
// @Accept multipart/form-data
// @Produce json
// @Param id formData integer true "Monitoring ID"
// @Param deskripsi_perkembangan formData string true "Deskripsi Perkembangan"
// @Param kendala_ditemui formData string false "Kendala Ditemui"
// @Param tindak_lanjut formData string false "Tindak Lanjut"
// @Param files_to_delete formData string false "Files to delete (comma-separated paths)"
// @Param file_pendukung formData file false "New file pendukung (multiple files)"
// @Param guru_pengisi_id formData integer false "Guru Pengisi ID"
// @Success 200 {object} gin.H{data=dtos.MonitoringInklusiBulananResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/update-monitoring-bulanan [post]
func (c *MonitoringInklusiBulananController) Update(ctx *gin.Context) {
	// Parse form data
	var req dtos.MonitoringInklusiBulananUpdateRequest

	// Parse id
	id, err := strconv.ParseUint(ctx.PostForm("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id harus berupa angka"})
		return
	}
	req.ID = uint(id)

	// Parse deskripsi_perkembangan
	req.DeskripsiPerkembangan = ctx.PostForm("deskripsi_perkembangan")
	if req.DeskripsiPerkembangan == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deskripsi_perkembangan harus diisi"})
		return
	}

	// Parse optional fields
	req.KendalaDitemui = ctx.PostForm("kendala_ditemui")
	req.TindakLanjut = ctx.PostForm("tindak_lanjut")

	// Parse files_to_delete (comma-separated)
	if filesToDeleteStr := ctx.PostForm("files_to_delete"); filesToDeleteStr != "" {
		req.FilesToDelete = strings.Split(filesToDeleteStr, ",")
		// Trim whitespace from each path
		for i, path := range req.FilesToDelete {
			req.FilesToDelete[i] = strings.TrimSpace(path)
		}
	}

	// Parse guru_pengisi_id (optional)
	if guruPengisiIDStr := ctx.PostForm("guru_pengisi_id"); guruPengisiIDStr != "" {
		guruPengisiID, err := strconv.ParseUint(guruPengisiIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "guru_pengisi_id harus berupa angka"})
			return
		}
		guruID := uint(guruPengisiID)
		req.GuruPengisiID = &guruID
	}

	// Get new file pendukung (multiple files)
	form, err := ctx.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil && form != nil && form.File["file_pendukung"] != nil {
		files = form.File["file_pendukung"]
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(files, &req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetSummary retrieves summary and statistics for monitoring inklusi bulanan
// @Summary Get summary monitoring inklusi bulanan
// @Description Retrieve comprehensive summary, statistics, and charts data for monitoring inklusi bulanan by year with optional month and rombel filter
// @Tags monitoring-inklusi
// @Accept json
// @Produce json
// @Param body body dtos.MonitoringInklusiBulananSummaryRequest true "Request body with tahun, optional bulan (1-12), and optional rombel_id"
// @Success 200 {object} gin.H{data=dtos.MonitoringInklusiBulananSummaryResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/monitoring-inklusi/get-summary-monitoring-bulanan [post]
func (c *MonitoringInklusiBulananController) GetSummary(ctx *gin.Context) {
	var req dtos.MonitoringInklusiBulananSummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Expected JSON with fields: tahun (year), optional bulan (1-12), optional rombel_id (number). Detail: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Tahun == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tahun harus diisi"})
		return
	}

	// Validate optional bulan
	if req.Bulan != nil && (*req.Bulan < 1 || *req.Bulan > 12) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bulan harus antara 1-12"})
		return
	}

	data, err := c.service.GetSummary(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
