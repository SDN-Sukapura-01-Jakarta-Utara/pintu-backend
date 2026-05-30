package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
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
	var req dtos.PengumumanKelulusanConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service
	result, err := c.service.ConfigurePengumuman(&req, userIDUint)
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
