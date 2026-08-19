package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaKelasRoutes registers SIEKSA kelas routes
func RegisterSieksaKelasRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller (reuse existing implementations)
	kelasRepo := repositories.NewKelasRepository(db)
	kelasService := services.NewKelasService(kelasRepo)
	kelasController := controllers.NewKelasController(kelasService)

	// Protected routes for SIEKSA (auth required with SIEKSA token)
	protected := router.Group("/api/sieksa/kelas")
	protected.Use(middleware.AuthMiddleware()) // Middleware will auto-detect SIEKSA from path
	{
		// Get all kelas
		protected.POST("/get-kelas", kelasController.GetAll)
	}
}
