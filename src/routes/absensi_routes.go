package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAbsensiRoutes registers all Absensi routes
func RegisterAbsensiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewAbsensiRepository(db)
	service := services.NewAbsensiService(repository)
	controller := controllers.NewAbsensiController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1/absensi-siswa")
	api.Use(middleware.AuthMiddleware())
	{
		// Create absensi manual (bulk input with file upload)
		api.POST("/create-absensi-manual", controller.CreateAbsensiManual)
		
		// Get rekap absensi
		api.POST("/get-rekap-absensi", controller.GetRekapAbsensi)
		
		// Update rekap absensi
		api.POST("/update-rekap-absensi", controller.UpdateRekapAbsensi)
		
		// Dashboard monitoring
		api.POST("/dashboard-summary", controller.GetDashboardSummary)
		api.POST("/grafik-kehadiran", controller.GetGrafikKehadiran)
		api.POST("/statistik-per-hari", controller.GetStatistikPerHari)
		api.POST("/perbandingan-rombel", controller.GetPerbandinganRombel)
		api.POST("/siswa-terendah", controller.GetSiswaTerendah)
		api.POST("/dashboard-siswa", controller.GetDashboardSiswa)
	}
}
