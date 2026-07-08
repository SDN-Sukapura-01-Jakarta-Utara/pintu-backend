package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPesertaDidikRoutes registers all PesertaDidik routes
func RegisterPesertaDidikRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	repository := repositories.NewPesertaDidikRepository(db)
	service := services.NewPesertaDidikService(repository, r2Storage)
	controller := controllers.NewPesertaDidikController(service)

	// Public routes (no authentication required)
	public := router.Group("/api/v1/public")
	{
		// Get total siswa with active tahun pelajaran and active status
		public.POST("/get-total-siswa", controller.GetTotalSiswa)
	}

	// Protected routes (require authentication)
	api := router.Group("/api/v1/peserta-didik")
	api.Use(middleware.AuthMiddleware())
	{
		// Create
		api.POST("/create-peserta-didik", controller.Create)

		// Read
		api.POST("/get-peserta-didik", controller.GetAll)
		api.POST("/get-peserta-didik-by-id", controller.GetByID)
		api.POST("/get-peserta-didik-by-nis", controller.GetByNIS)

		// Update
		api.POST("/update-peserta-didik", controller.Update)

		// Delete
		api.POST("/delete-peserta-didik", controller.Delete)

		// Import Excel
		api.POST("/import-excel", controller.ImportExcel)
		api.POST("/import-siswa-lulus", controller.ImportSiswaLulus)

		// Download Template
		api.POST("/download-template", controller.DownloadTemplate)
		api.POST("/download-template-siswa-lulus", controller.DownloadTemplateSiswaLulus)
		api.POST("/export-data-induk-siswa-excel", controller.ExportDataIndukSiswaExcel)
		api.POST("/export-data-induk-siswa-pdf", controller.ExportDataIndukSiswaPDF)
		api.POST("/export-pemetaan-rombel-excel", controller.ExportPemetaanRombelExcel)
		api.POST("/export-pemetaan-rombel-pdf", controller.ExportPemetaanRombelPDF)
		api.POST("/download-kartu-pelajar", controller.DownloadKartuPelajar)

		// Generate Barcode
		api.POST("/generate-barcode-all-peserta-didik", controller.GenerateBarcodeAllPesertaDidik)
		api.POST("/generate-barcode-peserta-didik-by-id", controller.GenerateBarcodePesertaDidikByID)
	}
}
