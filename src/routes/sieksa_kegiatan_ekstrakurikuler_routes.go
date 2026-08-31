package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaKegiatanEkstrakurikulerRoutes registers SIEKSA kegiatan ekstrakurikuler routes
func RegisterSieksaKegiatanEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services (reuse absensi ekskul components)
	absensiEkskulRepo := repositories.NewAbsensiEkskulRepository(db)
	absensiEkskulService := services.NewAbsensiEkskulService(absensiEkskulRepo, db)
	absensiEkskulController := controllers.NewAbsensiEkskulController(absensiEkskulService)

	// SIEKSA Kegiatan Ekstrakurikuler routes - uses SIEKSA auth
	protected := router.Group("/api/sieksa/kegiatan-ekstrakurikuler")
	protected.Use(middleware.AuthMiddleware())
	{
		// Get kegiatan list by ekstrakurikuler and tahun pelajaran
		protected.POST("/get-kegiatan-ekskul", absensiEkskulController.GetKegiatanEkskul)
		
		// Get kegiatan detail by ID
		protected.POST("/get-kegiatan-ekskul-by-id", absensiEkskulController.GetKegiatanEkskulByID)
		
		// Update kegiatan (edit tanggal, waktu, materi, upload/delete foto)
		protected.POST("/update-kegiatan-ekskul", absensiEkskulController.UpdateKegiatan)
		
		// Download Word documentation
		protected.POST("/download-word-dokumentasi-ekskul", absensiEkskulController.DownloadWordDokumentasiEkskul)
		
		// Download PDF documentation
		protected.POST("/download-pdf-dokumentasi-ekskul", absensiEkskulController.DownloadPDFDokumentasiEkskul)
	}
}
