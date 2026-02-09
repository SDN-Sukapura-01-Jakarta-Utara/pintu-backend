package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterKelasRoutes registers all kelas routes
func RegisterKelasRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	kelasRepo := repositories.NewKelasRepository(db)
	kelasService := services.NewKelasService(kelasRepo)
	kelasController := controllers.NewKelasController(kelasService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/kelas")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create kelas
		protected.POST("/create-kelas", kelasController.Create)

		// Get all kelas
		protected.POST("/get-kelas", kelasController.GetAll)

		// Get kelas by ID
		protected.POST("/get-kelas-by-id", kelasController.GetByID)

		// Update kelas
		protected.POST("/update-kelas", kelasController.Update)

		// Delete kelas
		protected.POST("/delete-kelas", kelasController.Delete)
	}
}
