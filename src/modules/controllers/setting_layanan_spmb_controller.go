package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// SettingLayananSPMBController handles HTTP requests for Setting Layanan SPMB
type SettingLayananSPMBController struct {
	service services.SettingLayananSPMBService
}

// NewSettingLayananSPMBController creates a new Setting Layanan SPMB controller
func NewSettingLayananSPMBController(service services.SettingLayananSPMBService) *SettingLayananSPMBController {
	return &SettingLayananSPMBController{service: service}
}

// UpsertSetting creates or updates setting layanan SPMB (auth required)
// @Summary Create or Update Setting Layanan SPMB
// @Description Create or update setting layanan SPMB dengan ID = 1
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Param body body dtos.SettingLayananSPMBRequest true "Request body"
// @Success 200 {object} gin.H{message=string,data=dtos.SettingLayananSPMBResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb/setting-layanan-spmb [post]
func (c *SettingLayananSPMBController) UpsertSetting(ctx *gin.Context) {
	var req dtos.SettingLayananSPMBRequest
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

	data, err := c.service.UpsertSetting(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Setting layanan SPMB berhasil disimpan",
		"data":    data,
	})
}

// GetGrupWAPublic gets grup WA link for public (no auth required)
// @Summary Get Grup WA SPMB (Public)
// @Description Get grup WA link untuk layanan SPMB
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.GrupWASPMBResponse}
// @Router /api/v1/public/get-grup-wa-spmb [post]
func (c *SettingLayananSPMBController) GetGrupWAPublic(ctx *gin.Context) {
	data, err := c.service.GetGrupWAPublic()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetSetting gets setting layanan SPMB (auth required)
// @Summary Get Setting Layanan SPMB
// @Description Get setting layanan SPMB (ID = 1)
// @Tags layanan-spmb
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.SettingLayananSPMBResponse}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb/get-setting-layanan-spmb [post]
func (c *SettingLayananSPMBController) GetSetting(ctx *gin.Context) {
	data, err := c.service.GetSetting()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
