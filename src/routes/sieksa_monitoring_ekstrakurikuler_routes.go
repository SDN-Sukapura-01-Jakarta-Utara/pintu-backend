package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaMonitoringEkstrakurikulerRoutes registers SIEKSA monitoring ekstrakurikuler routes
func RegisterSieksaMonitoringEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services
	pesertaDidikEkskulRepo := repositories.NewPesertaDidikEkstrakurikulerRepository(db)
	pesertaDidikRombelRepo := repositories.NewPesertaDidikRombelRepository(db)
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	pesertaDidikEkskulService := services.NewPesertaDidikEkstrakurikulerService(pesertaDidikEkskulRepo, pesertaDidikRombelRepo, ekstrakurikulerRepo)
	pesertaDidikEkskulController := controllers.NewPesertaDidikEkstrakurikulerController(pesertaDidikEkskulService)

	// Protected routes for SIEKSA monitoring ekstrakurikuler (auth required with SIEKSA token)
	monitoringGroup := router.Group("/api/sieksa/monitoring-ekstrakurikuler")
	monitoringGroup.Use(middleware.AuthMiddleware())
	{
		// Get comprehensive statistics for monitoring
		monitoringGroup.POST("/get-all-statistic-ekstrakurikuler-siswa", pesertaDidikEkskulController.GetStatistikEkstrakurikuler)
	}
}
