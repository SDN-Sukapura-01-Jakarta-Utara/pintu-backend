package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// KonfigurasiAbsensiController handles HTTP requests for Konfigurasi Absensi
type KonfigurasiAbsensiController struct {
	service services.KonfigurasiAbsensiService
}

// NewKonfigurasiAbsensiController creates a new Konfigurasi Absensi controller
func NewKonfigurasiAbsensiController(service services.KonfigurasiAbsensiService) *KonfigurasiAbsensiController {
	return &KonfigurasiAbsensiController{service: service}
}

// UpsertKonfigurasi creates or updates konfigurasi absensi (auth required)
// @Summary Create or Update Konfigurasi Absensi
// @Description Create or update konfigurasi absensi dengan ID = 1
// @Tags absensi-siswa
// @Accept json
// @Produce json
// @Param body body dtos.KonfigurasiAbsensiRequest true "Request body"
// @Success 200 {object} gin.H{message=string,data=dtos.KonfigurasiAbsensiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/absensi-siswa/setting-konfigurasi-absensi [post]
func (c *KonfigurasiAbsensiController) UpsertKonfigurasi(ctx *gin.Context) {
	var req dtos.KonfigurasiAbsensiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.UpsertKonfigurasi(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Konfigurasi absensi berhasil disimpan",
		"data":    data,
	})
}

// GetKonfigurasi gets konfigurasi absensi (auth required)
// @Summary Get Konfigurasi Absensi
// @Description Get konfigurasi absensi dengan ID = 1
// @Tags absensi-siswa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.KonfigurasiAbsensiResponse}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/absensi-siswa/get-konfigurasi-absensi [post]
func (c *KonfigurasiAbsensiController) GetKonfigurasi(ctx *gin.Context) {
	data, err := c.service.GetKonfigurasi()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
