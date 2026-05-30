package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPengumumanKelulusanRoutes registers all PengumumanKelulusan routes
func RegisterPengumumanKelulusanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewPengumumanKelulusanRepository(db)
	service := services.NewPengumumanKelulusanService(repository)
	controller := controllers.NewPengumumanKelulusanController(service)

	// Public routes (no authentication required)
	publicAPI := router.Group("/api/v1/public")
	{
		// Get setting pengumuman kelulusan (ID 1)
		publicAPI.POST("/get-setting-pengumuman-kelulusan", controller.GetSettingPengumumanPublic)
	}

	// Protected routes (require authentication)
	api := router.Group("/api/v1/kelulusan")
	api.Use(middleware.AuthMiddleware())
	{
		// Configure pengumuman (create or update)
		api.POST("/konfigurasi-pengumuman", controller.ConfigurePengumuman)
		
		// Get pengumuman configuration
		api.POST("/get-konfigurasi-pengumuman", controller.GetPengumuman)
	}
}
