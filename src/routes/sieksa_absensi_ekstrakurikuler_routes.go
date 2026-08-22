package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaAbsensiEkstrakurikulerRoutes registers SIEKSA absensi ekstrakurikuler routes
func RegisterSieksaAbsensiEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services
	absensiEkskulRepo := repositories.NewAbsensiEkskulRepository(db)
	absensiEkskulService := services.NewAbsensiEkskulService(absensiEkskulRepo, db)
	absensiEkskulController := controllers.NewAbsensiEkskulController(absensiEkskulService)

	// SIEKSA Absensi Ekstrakurikuler routes - uses SIEKSA auth
	protected := router.Group("/api/sieksa/absensi-ekstrakurikuler")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create absensi ekstrakurikuler (kegiatan + bulk absensi siswa + absensi pelatih)
		protected.POST("/create-absensi", absensiEkskulController.Create)
		
		// Get absensi siswa by ekstrakurikuler and tahun pelajaran
		protected.POST("/get-absensi-siswa", absensiEkskulController.GetAbsensiSiswa)
		
		// Get single absensi siswa by ID
		protected.POST("/get-absensi-siswa-by-id", absensiEkskulController.GetAbsensiSiswaByID)
		
		// Update absensi siswa
		protected.POST("/update-absensi-siswa", absensiEkskulController.UpdateAbsensiSiswa)
		
		// Get absensi pelatih by pelatih and tahun pelajaran
		protected.POST("/get-absensi-pelatih", absensiEkskulController.GetAbsensiPelatih)
		
		// Get kegiatan list by ekstrakurikuler and tahun pelajaran
		protected.POST("/get-kegiatan-ekskul", absensiEkskulController.GetKegiatanEkskul)
		
		// Update/toggle pelatih attendance
		protected.POST("/update-absensi-pelatih", absensiEkskulController.UpdateAbsensiPelatih)
		
		// Download Excel absensi siswa
		protected.POST("/download-excel-absensi-siswa", absensiEkskulController.DownloadExcelAbsensiSiswa)
		
		// Download PDF absensi siswa
		protected.POST("/download-pdf-absensi-siswa", absensiEkskulController.DownloadPDFAbsensiSiswa)
		
		// Download Excel absensi pelatih
		protected.POST("/download-excel-absensi-pelatih", absensiEkskulController.DownloadExcelAbsensiPelatih)
		
		// Download PDF absensi pelatih
		protected.POST("/download-pdf-absensi-pelatih", absensiEkskulController.DownloadPDFAbsensiPelatih)
	}
}
