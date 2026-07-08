package routes

import (
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAbsensiScanRoutes registers all absensi scan routes
func RegisterAbsensiScanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewAbsensiScanRepository(db)
	service := services.NewAbsensiScanService(repository)
	controller := controllers.NewAbsensiScanController(service)

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/absensi-siswa", controller.ScanAbsensi)
	}
}
