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
	service := services.NewAbsensiService(repository, db)
	controller := controllers.NewAbsensiController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1/absensi-siswa")
	api.Use(middleware.AuthMiddleware())
	{
		// Create absensi manual (bulk input with file upload)
		api.POST("/create-absensi-manual", controller.CreateAbsensiManual)
		
		// Create absensi manual by ID (single student with auto semester detection)
		api.POST("/create-absensi-manual-by-id", controller.CreateAbsensiManualByID)
		
		// Synchronize absensi from scanner to rekapitulasi
		api.POST("/synchronize-absensi-siswa", controller.SynchronizeAbsensi)
		
		// Get rekap absensi
		api.POST("/get-rekap-absensi", controller.GetRekapAbsensi)
		
		// Update rekap absensi
		api.POST("/update-rekap-absensi", controller.UpdateRekapAbsensi)
		
		// Export absensi to Excel
		api.POST("/export-excel-absensi-siswa", controller.ExportAbsensiExcel)
		
		// Export absensi to PDF
		api.POST("/export-pdf-absensi-siswa", controller.ExportAbsensiPDF)
		
		// Dashboard monitoring
		api.POST("/dashboard-summary", controller.GetDashboardSummary)
		api.POST("/grafik-kehadiran", controller.GetGrafikKehadiran)
		api.POST("/statistik-per-hari", controller.GetStatistikPerHari)
		api.POST("/perbandingan-rombel", controller.GetPerbandinganRombel)
		api.POST("/siswa-terendah", controller.GetSiswaTerendah)
		api.POST("/dashboard-siswa", controller.GetDashboardSiswa)
	}
}
