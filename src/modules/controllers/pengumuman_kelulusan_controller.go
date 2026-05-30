package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PengumumanKelulusanController handles HTTP requests for PengumumanKelulusan
type PengumumanKelulusanController struct {
	service services.PengumumanKelulusanService
}

// NewPengumumanKelulusanController creates a new PengumumanKelulusan controller
func NewPengumumanKelulusanController(service services.PengumumanKelulusanService) *PengumumanKelulusanController {
	return &PengumumanKelulusanController{service: service}
}

// ConfigurePengumuman creates or updates pengumuman kelulusan configuration
func (c *PengumumanKelulusanController) ConfigurePengumuman(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.PengumumanKelulusanConfigRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Manual validation using validator
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Get foto_kepsek file if uploaded
	var fotoKepsek *multipart.FileHeader
	fotoKepsekHeader, err := ctx.FormFile("foto_kepsek")
	if err == nil {
		fotoKepsek = fotoKepsekHeader
	}

	// Get ttd_kepsek file if uploaded
	var ttdKepsek *multipart.FileHeader
	ttdKepsekHeader, err := ctx.FormFile("ttd_kepsek")
	if err == nil {
		ttdKepsek = ttdKepsekHeader
	}

	// Call service
	result, err := c.service.ConfigurePengumuman(&req, fotoKepsek, ttdKepsek, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := "Konfigurasi pengumuman berhasil disimpan"
	if req.ID != nil && *req.ID > 0 {
		message = "Konfigurasi pengumuman berhasil diupdate"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    result,
	})
}

// GetPengumuman retrieves the pengumuman kelulusan configuration
func (c *PengumumanKelulusanController) GetPengumuman(ctx *gin.Context) {
	result, err := c.service.GetPengumuman()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetSettingPengumumanPublic retrieves the pengumuman kelulusan configuration (public API)
func (c *PengumumanKelulusanController) GetSettingPengumumanPublic(ctx *gin.Context) {
	result, err := c.service.GetSettingPengumumanPublic()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
