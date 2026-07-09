package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// KonfigurasiMutasiSiswaController handles HTTP requests for Konfigurasi Mutasi Siswa
type KonfigurasiMutasiSiswaController struct {
	service services.KonfigurasiMutasiSiswaService
}

// NewKonfigurasiMutasiSiswaController creates a new Konfigurasi Mutasi Siswa controller
func NewKonfigurasiMutasiSiswaController(service services.KonfigurasiMutasiSiswaService) *KonfigurasiMutasiSiswaController {
	return &KonfigurasiMutasiSiswaController{service: service}
}

// UpsertSetting creates or updates konfigurasi mutasi siswa (auth required)
// @Summary Create or Update Konfigurasi Mutasi Siswa
// @Description Create or update konfigurasi mutasi siswa dengan ID = 1
// @Tags mutasi-siswa
// @Accept multipart/form-data
// @Produce json
// @Param tanggal_buka_pendaftaran formData string true "Tanggal Buka Pendaftaran (YYYY-MM-DD)"
// @Param tanggal_tutup_pendaftaran formData string true "Tanggal Tutup Pendaftaran (YYYY-MM-DD)"
// @Param nama_kepala_sekolah formData string true "Nama Kepala Sekolah"
// @Param nip_kepala_sekolah formData string true "NIP Kepala Sekolah"
// @Param nama_ketua_panitia formData string true "Nama Ketua Panitia"
// @Param nip_ketua_panitia formData string true "NIP Ketua Panitia"
// @Param grup_wa formData string false "Link Grup WhatsApp"
// @Param template_sptjm formData file false "File Template SPTJM (PDF)"
// @Success 200 {object} gin.H{message=string,data=dtos.KonfigurasiMutasiSiswaResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/setting-konfigurasi-mutasi-siswa [post]
func (c *KonfigurasiMutasiSiswaController) UpsertSetting(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil { // 10MB
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var req dtos.KonfigurasiMutasiSiswaRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get template SPTJM file if provided
	file, _ := ctx.FormFile("template_sptjm")

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.UpsertSetting(&req, file, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Konfigurasi mutasi siswa berhasil disimpan",
		"data":    data,
	})
}

// GetSetting gets konfigurasi mutasi siswa (auth required)
// @Summary Get Konfigurasi Mutasi Siswa
// @Description Get konfigurasi mutasi siswa (ID = 1), returns null if not found
// @Tags mutasi-siswa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.KonfigurasiMutasiSiswaResponse}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/get-konfigurasi-mutasi-siswa [post]
func (c *KonfigurasiMutasiSiswaController) GetSetting(ctx *gin.Context) {
	data, err := c.service.GetSetting()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If data is nil (not found), return null
	if data == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}


// GetSettingPublic gets konfigurasi mutasi siswa for public (no auth required)
// @Summary Get Konfigurasi Mutasi Siswa (Public)
// @Description Get konfigurasi mutasi siswa (ID = 1) for public, returns null if not found
// @Tags mutasi-siswa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.KonfigurasiMutasiSiswaResponse}
// @Router /api/v1/public/get-konfigurasi-mutasi-siswa [post]
func (c *KonfigurasiMutasiSiswaController) GetSettingPublic(ctx *gin.Context) {
	data, err := c.service.GetSetting()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If data is nil (not found), return null
	if data == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
