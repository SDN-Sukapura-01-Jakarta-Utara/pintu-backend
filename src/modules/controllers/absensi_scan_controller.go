package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// AbsensiScanController handles HTTP requests for Absensi Scan
type AbsensiScanController struct {
	service services.AbsensiScanService
}

// NewAbsensiScanController creates a new Absensi Scan controller
func NewAbsensiScanController(service services.AbsensiScanService) *AbsensiScanController {
	return &AbsensiScanController{service: service}
}

// ScanAbsensi handles barcode scanning for attendance (public, no auth)
// @Summary Scan Absensi Siswa
// @Description Scan barcode untuk absensi siswa (datang/pulang)
// @Tags absensi-siswa
// @Accept json
// @Produce json
// @Param body body dtos.AbsensiScanRequest true "Request body"
// @Success 200 {object} dtos.AbsensiScanResponse
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/absensi-siswa [post]
func (c *AbsensiScanController) ScanAbsensi(ctx *gin.Context) {
	var req dtos.AbsensiScanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.ScanAbsensi(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan sistem"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
