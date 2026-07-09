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

// RegisterMutasiSiswaRoutes registers all mutasi siswa routes
func RegisterMutasiSiswaRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller for Mutasi Siswa
	mutasiSiswaRepo := repositories.NewMutasiSiswaRepository(db)
	mutasiSiswaService := services.NewMutasiSiswaService(mutasiSiswaRepo, r2Storage)
	mutasiSiswaController := controllers.NewMutasiSiswaController(mutasiSiswaService)

	// Initialize repository, service, and controller for Konfigurasi Mutasi Siswa
	konfigurasiRepo := repositories.NewKonfigurasiMutasiSiswaRepository(db)
	konfigurasiService := services.NewKonfigurasiMutasiSiswaService(konfigurasiRepo, r2Storage)
	konfigurasiController := controllers.NewKonfigurasiMutasiSiswaController(konfigurasiService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/spmb-mutasi")
	protected.Use(middleware.AuthMiddleware())
	{
		// Get all mutasi siswa with filters
		protected.POST("/get-mutasi-siswa", mutasiSiswaController.GetAll)
		
		// Get mutasi siswa by ID
		protected.POST("/get-mutasi-siswa-by-id", mutasiSiswaController.GetByID)
		
		// Update mutasi siswa
		protected.POST("/edit-mutasi-siswa", mutasiSiswaController.Update)

		// Delete mutasi siswa
		protected.POST("/delete-mutasi-siswa", mutasiSiswaController.Delete)

		// Export formulir pendaftaran PDF (admin)
		protected.POST("/export-pdf-formulir-mutasi-siswa", mutasiSiswaController.ExportFormulirPDFAuth)

		// Export Excel data mutasi siswa
		protected.POST("/export-excel-mutasi-siswa", mutasiSiswaController.ExportExcel)

		// Export PDF list mutasi siswa
		protected.POST("/export-pdf-mutasi-siswa", mutasiSiswaController.ExportListPDF)

		// Setting konfigurasi mutasi siswa (upsert)
		protected.POST("/setting-konfigurasi-mutasi-siswa", konfigurasiController.UpsertSetting)

		// Get konfigurasi mutasi siswa
		protected.POST("/get-konfigurasi-mutasi-siswa", konfigurasiController.GetSetting)
	}

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		// Create mutasi siswa from public form
		public.POST("/create-mutasi-siswa", mutasiSiswaController.CreatePublic)

		// Get konfigurasi mutasi siswa (public)
		public.POST("/get-konfigurasi-mutasi-siswa", konfigurasiController.GetSettingPublic)

		// Export formulir pendaftaran PDF
		public.POST("/export-pdf-formulir-mutasi-siswa", mutasiSiswaController.ExportFormulirPDF)
	}
}
