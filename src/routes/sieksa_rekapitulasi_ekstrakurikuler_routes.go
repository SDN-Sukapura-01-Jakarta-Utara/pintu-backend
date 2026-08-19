package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaRekapitulasiEkstrakurikulerRoutes registers SIEKSA rekapitulasi ekstrakurikuler routes
func RegisterSieksaRekapitulasiEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services
	pesertaDidikEkskulRepo := repositories.NewPesertaDidikEkstrakurikulerRepository(db)
	pesertaDidikRombelRepo := repositories.NewPesertaDidikRombelRepository(db)
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	pesertaDidikEkskulService := services.NewPesertaDidikEkstrakurikulerService(pesertaDidikEkskulRepo, pesertaDidikRombelRepo, ekstrakurikulerRepo)
	pesertaDidikEkskulController := controllers.NewPesertaDidikEkstrakurikulerController(pesertaDidikEkskulService)

	// Protected routes for SIEKSA rekapitulasi ekstrakurikuler (auth required with SIEKSA token)
	rekapGroup := router.Group("/api/sieksa/rekapitulasi-ekstrakurikuler")
	rekapGroup.Use(middleware.AuthMiddleware())
	{
		// Get rekapitulasi data per ekstrakurikuler
		rekapGroup.POST("/rekapitulasi-data-per-ekskul", pesertaDidikEkskulController.GetRekapitulasiPerEkskul)

		// Get rekapitulasi data per rombel
		rekapGroup.POST("/rekapitulasi-data-per-rombel", pesertaDidikEkskulController.GetRekapitulasiPerRombel)

		// Export Excel per ekstrakurikuler
		rekapGroup.POST("/download-excel-data-per-ekskul", pesertaDidikEkskulController.ExportExcelPerEkskul)

		// Export Excel per rombel
		rekapGroup.POST("/download-excel-data-per-rombel", pesertaDidikEkskulController.ExportExcelPerRombel)
	}
}
