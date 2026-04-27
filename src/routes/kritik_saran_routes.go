package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterKritikSaranRoutes registers all kritik saran routes
func RegisterKritikSaranRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	kritikSaranRepo := repositories.NewKritikSaranRepository(db)
	kritikSaranService := services.NewKritikSaranService(kritikSaranRepo)
	kritikSaranController := controllers.NewKritikSaranController(kritikSaranService)

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/create-kritik-saran", kritikSaranController.CreatePublic)
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/kritik-saran")
	protected.Use(middleware.AuthMiddleware())
	{
		// Get all kritik saran
		protected.POST("/get-kritik-saran", kritikSaranController.GetAll)

		// Get kritik saran by ID
		protected.POST("/get-kritik-saran-by-id", kritikSaranController.GetByID)

		// Delete kritik saran
		protected.POST("/delete-kritik-saran", kritikSaranController.Delete)
	}
}
