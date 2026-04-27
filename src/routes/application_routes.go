package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterApplicationRoutes registers all application routes
func RegisterApplicationRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	applicationRepo := repositories.NewApplicationRepository(db)
	applicationService := services.NewApplicationService(applicationRepo)
	applicationController := controllers.NewApplicationController(applicationService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/application")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create application
		protected.POST("/create-application", applicationController.Create)

		// Get all applications
		protected.POST("/get-application", applicationController.GetAll)

		// Get application by ID
		protected.POST("/get-application-by-id", applicationController.GetByID)

		// Update application
		protected.POST("/update-application", applicationController.Update)

		// Delete application
		protected.POST("/delete-application", applicationController.Delete)
	}
}
