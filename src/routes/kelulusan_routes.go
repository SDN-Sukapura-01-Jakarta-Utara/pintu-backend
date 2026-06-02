package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterKelulusanRoutes registers all Kelulusan routes
func RegisterKelulusanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories, service, and controller
	kelulusanRepository := repositories.NewKelulusanRepository(db)
	tahunPelajaranRepository := repositories.NewTahunPelajaranRepository(db)
	pengumumanKelulusanRepository := repositories.NewPengumumanKelulusanRepository(db)
	service := services.NewKelulusanService(kelulusanRepository, tahunPelajaranRepository, pengumumanKelulusanRepository)
	controller := controllers.NewKelulusanController(service)

	// Public routes (no authentication required)
	publicAPI := router.Group("/api/v1/public")
	{
		// Cek nilai kelulusan by NISN and tanggal lahir (without lulus info)
		publicAPI.POST("/cek-nilai-kelulusan", controller.CekNilaiKelulusan)
		
		// Cek kelulusan by NISN and tanggal lahir (with lulus info and SKL)
		publicAPI.POST("/cek-kelulusan", controller.CekKelulusan)
		
		// Download laporan nilai kelulusan PDF
		publicAPI.POST("/download-laporan-nilai-kelulusan", controller.DownloadLaporanNilaiKelulusan)
	}

	// Protected routes (require authentication)
	api := router.Group("/api/v1/kelulusan")
	api.Use(middleware.AuthMiddleware())
	{
		// Create data kelulusan
		api.POST("/create-data-kelulusan", controller.CreateKelulusan)
		
		// Get all data kelulusan with filters
		api.POST("/get-data-kelulusan", controller.GetAll)
		
		// Get by ID
		api.POST("/get-data-kelulusan-by-id", controller.GetByID)
		
		// Update data kelulusan
		api.POST("/update-data-kelulusan", controller.Update)
		
		// Delete data kelulusan
		api.POST("/delete-data-kelulusan", controller.Delete)
		
		// Download template
		api.POST("/download-template", controller.DownloadTemplate)
		
		// Import Excel
		api.POST("/import-excel", controller.ImportExcel)
	}
}
