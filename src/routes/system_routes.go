package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSystemRoutes registers all system routes
func RegisterSystemRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	systemRepo := repositories.NewSystemRepository(db)
	systemService := services.NewSystemService(systemRepo)
	systemController := controllers.NewSystemController(systemService)

	// Group routes under /api/v1/systems with auth middleware
	api := router.Group("/api/v1/systems")
	api.Use(middleware.AuthMiddleware()) // Require authentication
	{
		api.POST("/create-system", systemController.Create)      // Create system
		api.POST("/get-systems", systemController.GetAll)        // Get all systems
		api.POST("/get-system-by-id", systemController.GetByID)  // Get system by ID
		api.POST("/update-system", systemController.Update)      // Update system
		api.POST("/delete-system", systemController.Delete)      // Delete system
	}
}
