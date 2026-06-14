package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterLayananSPMBRoutes registers all layanan SPMB routes
func RegisterLayananSPMBRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize Layanan SPMB repository, service, and controller
	layananSPMBRepo := repositories.NewLayananSPMBRepository(db)
	layananSPMBService := services.NewLayananSPMBService(layananSPMBRepo)
	layananSPMBController := controllers.NewLayananSPMBController(layananSPMBService)

	// Initialize Setting Layanan SPMB repository, service, and controller
	settingLayananSPMBRepo := repositories.NewSettingLayananSPMBRepository(db)
	settingLayananSPMBService := services.NewSettingLayananSPMBService(settingLayananSPMBRepo)
	settingLayananSPMBController := controllers.NewSettingLayananSPMBController(settingLayananSPMBService)

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/layanan-spmb", layananSPMBController.CreatePublic)
		public.POST("/get-grup-wa-spmb", settingLayananSPMBController.GetGrupWAPublic)
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/spmb")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/get-layanan-spmb", layananSPMBController.GetAll)
		protected.POST("/get-layanan-spmb-by-id", layananSPMBController.GetByID)
		protected.POST("/set-status-selesai", layananSPMBController.UpdateStatus)
		protected.POST("/delete-layanan-spmb", layananSPMBController.DeleteLayananSPMB)
		protected.POST("/setting-layanan-spmb", settingLayananSPMBController.UpsertSetting)
		protected.POST("/get-setting-layanan-spmb", settingLayananSPMBController.GetSetting)
	}
}
