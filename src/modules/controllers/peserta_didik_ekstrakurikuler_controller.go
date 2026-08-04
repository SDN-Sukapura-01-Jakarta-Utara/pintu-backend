package controllers

import (
	"fmt"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// PesertaDidikEkstrakurikulerController handles HTTP requests
type PesertaDidikEkstrakurikulerController struct {
	service services.PesertaDidikEkstrakurikulerService
}

// NewPesertaDidikEkstrakurikulerController creates a new controller
func NewPesertaDidikEkstrakurikulerController(service services.PesertaDidikEkstrakurikulerService) *PesertaDidikEkstrakurikulerController {
	return &PesertaDidikEkstrakurikulerController{service: service}
}

// RegisterOrUpdateEkstrakurikuler handles registration/update of multiple ekstrakurikuler for a student
// @Summary Register or update student ekstrakurikuler
// @Description Register a student to ekstrakurikuler or update existing registration (add/remove)
// @Tags ekstrakurikuler-registration
// @Accept json
// @Produce json
// @Param body body dtos.RegisterOrUpdateEkstrakurikulerRequest true "Request body"
// @Success 200 {object} dtos.RegisterOrUpdateEkstrakurikulerResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/register-ekskul-peserta-didik [post]
func (c *PesertaDidikEkstrakurikulerController) RegisterOrUpdateEkstrakurikuler(ctx *gin.Context) {
	var req dtos.RegisterOrUpdateEkstrakurikulerRequest
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

	data, err := c.service.RegisterOrUpdateEkstrakurikuler(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// GetEkstrakurikulerByPesertaDidik handles getting all ekstrakurikuler for a student
// @Summary Get student's ekstrakurikuler
// @Description Get all ekstrakurikuler registered by a student
// @Tags ekstrakurikuler-registration
// @Accept json
// @Produce json
// @Param body body dtos.GetEkstrakurikulerPesertaDidikRequest true "Request body"
// @Success 200 {object} dtos.GetEkstrakurikulerPesertaDidikResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/get-ekskul-peserta-didik [post]
func (c *PesertaDidikEkstrakurikulerController) GetEkstrakurikulerByPesertaDidik(ctx *gin.Context) {
	var req dtos.GetEkstrakurikulerPesertaDidikRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetEkstrakurikulerByPesertaDidik(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAllEkstrakurikulerSiswa handles getting all students with their ekstrakurikuler by rombel
// @Summary Get all students' ekstrakurikuler by rombel
// @Description Get all students in a rombel with their registered ekstrakurikuler
// @Tags ekstrakurikuler-registration
// @Accept json
// @Produce json
// @Param body body dtos.GetAllEkstrakurikulerSiswaRequest true "Request body"
// @Success 200 {object} dtos.GetAllEkstrakurikulerSiswaResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/get-all-ekstrakurikuler-siswa [post]
func (c *PesertaDidikEkstrakurikulerController) GetAllEkstrakurikulerSiswa(ctx *gin.Context) {
	var req dtos.GetAllEkstrakurikulerSiswaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetAllEkstrakurikulerSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}


// RegisterAllEkstrakurikulerSiswa handles bulk registration/update of ekstrakurikuler for multiple students
// @Summary Bulk register/update ekstrakurikuler for multiple students
// @Description Register or update ekstrakurikuler for multiple students at once (for admin/teacher)
// @Tags ekstrakurikuler-registration
// @Accept json
// @Produce json
// @Param body body dtos.RegisterAllEkstrakurikulerSiswaRequest true "Request body"
// @Success 200 {object} dtos.RegisterAllEkstrakurikulerSiswaResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/register-all-ekstrakurikuler-siswa [post]
func (c *PesertaDidikEkstrakurikulerController) RegisterAllEkstrakurikulerSiswa(ctx *gin.Context) {
	var req dtos.RegisterAllEkstrakurikulerSiswaRequest
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

	data, err := c.service.RegisterAllEkstrakurikulerSiswa(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}


// GetStatistikEkstrakurikuler handles getting comprehensive statistics for ekstrakurikuler monitoring
// @Summary Get ekstrakurikuler statistics
// @Description Get comprehensive statistics for monitoring ekstrakurikuler registration (per ekskul, per rombel, students without ekskul)
// @Tags ekstrakurikuler-statistics
// @Accept json
// @Produce json
// @Param body body dtos.GetStatistikEkstrakurikulerRequest true "Request body"
// @Success 200 {object} dtos.GetStatistikEkstrakurikulerResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/get-all-statistic-ekstrakurikuler-siswa [post]
func (c *PesertaDidikEkstrakurikulerController) GetStatistikEkstrakurikuler(ctx *gin.Context) {
	var req dtos.GetStatistikEkstrakurikulerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetStatistikEkstrakurikuler(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}


// GetRekapitulasiPerEkskul handles getting rekapitulasi data per ekstrakurikuler
// @Summary Get rekapitulasi per ekstrakurikuler
// @Description Get list of students per ekstrakurikuler with filters and pagination
// @Tags ekstrakurikuler-rekapitulasi
// @Accept json
// @Produce json
// @Param body body dtos.RekapitulasiPerEkskulRequest true "Request body"
// @Success 200 {object} dtos.RekapitulasiPerEkskulResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/rekapitulasi-data-per-ekskul [post]
func (c *PesertaDidikEkstrakurikulerController) GetRekapitulasiPerEkskul(ctx *gin.Context) {
	var req dtos.RekapitulasiPerEkskulRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetRekapitulasiPerEkskul(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetRekapitulasiPerRombel handles getting rekapitulasi data per rombel
// @Summary Get rekapitulasi per rombel
// @Description Get list of students in a rombel with their ekstrakurikuler
// @Tags ekstrakurikuler-rekapitulasi
// @Accept json
// @Produce json
// @Param body body dtos.RekapitulasiPerRombelRequest true "Request body"
// @Success 200 {object} dtos.RekapitulasiPerRombelResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/rekapitulasi-data-per-rombel [post]
func (c *PesertaDidikEkstrakurikulerController) GetRekapitulasiPerRombel(ctx *gin.Context) {
	var req dtos.RekapitulasiPerRombelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetRekapitulasiPerRombel(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// ExportExcelPerEkskul handles exporting ekstrakurikuler data per ekskul to Excel
// @Summary Export Excel per ekstrakurikuler
// @Description Export list of students per ekstrakurikuler to Excel file
// @Tags ekstrakurikuler-export
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param body body dtos.ExportExcelPerEkskulRequest true "Request body"
// @Success 200 {file} file
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/download-excel-data-per-ekskul [post]
func (c *PesertaDidikEkstrakurikulerController) ExportExcelPerEkskul(ctx *gin.Context) {
	var req dtos.ExportExcelPerEkskulRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, filename, err := c.service.ExportExcelPerEkskul(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// ExportExcelPerRombel handles exporting ekstrakurikuler data per rombel to Excel
// @Summary Export Excel per rombel
// @Description Export list of students in a rombel with their ekstrakurikuler to Excel file
// @Tags ekstrakurikuler-export
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param body body dtos.ExportExcelPerRombelRequest true "Request body"
// @Success 200 {file} file
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/ekstrakurikuler/download-excel-data-per-rombel [post]
func (c *PesertaDidikEkstrakurikulerController) ExportExcelPerRombel(ctx *gin.Context) {
	var req dtos.ExportExcelPerRombelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, filename, err := c.service.ExportExcelPerRombel(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
