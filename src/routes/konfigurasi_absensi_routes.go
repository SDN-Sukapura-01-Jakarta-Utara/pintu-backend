package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterKonfigurasiAbsensiRoutes registers all konfigurasi absensi routes
func RegisterKonfigurasiAbsensiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewKonfigurasiAbsensiRepository(db)
	service := services.NewKonfigurasiAbsensiService(repository)
	controller := controllers.NewKonfigurasiAbsensiController(service)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/absensi-siswa")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/setting-konfigurasi-absensi", controller.UpsertKonfigurasi)
		protected.POST("/get-konfigurasi-absensi", controller.GetKonfigurasi)
	}
}
