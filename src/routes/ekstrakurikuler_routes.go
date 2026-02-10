package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterEkstrakurikulerRoutes registers all ekstrakurikuler routes
func RegisterEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	ekstrakurikulerService := services.NewEkstrakurikulerService(ekstrakurikulerRepo, kelasRepo)
	ekstrakurikulerController := controllers.NewEkstrakurikulerController(ekstrakurikulerService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/ekstrakurikuler")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create ekstrakurikuler
		protected.POST("/create-ekstrakurikuler", ekstrakurikulerController.Create)

		// Get all ekstrakurikuler
		protected.POST("/get-ekstrakurikuler", ekstrakurikulerController.GetAll)

		// Get ekstrakurikuler by ID
		protected.POST("/get-ekstrakurikuler-by-id", ekstrakurikulerController.GetByID)

		// Update ekstrakurikuler
		protected.POST("/update-ekstrakurikuler", ekstrakurikulerController.Update)

		// Delete ekstrakurikuler
		protected.POST("/delete-ekstrakurikuler", ekstrakurikulerController.Delete)
	}
}
