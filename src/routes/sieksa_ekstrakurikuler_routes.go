package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaEkstrakurikulerRoutes registers all SIEKSA ekstrakurikuler CRUD routes
func RegisterSieksaEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller (reuse existing implementations)
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	ekstrakurikulerService := services.NewEkstrakurikulerService(ekstrakurikulerRepo, kelasRepo)
	ekstrakurikulerController := controllers.NewEkstrakurikulerController(ekstrakurikulerService)

	// Public routes for SIEKSA (no authentication required)
	public := router.Group("/api/sieksa/public")
	{
		// Get total ekstrakurikuler with status "active"
		public.POST("/get-total-ekskul", ekstrakurikulerController.GetTotalEkskul)
	}

	// Protected routes for SIEKSA ekstrakurikuler CRUD (auth required with SIEKSA token)
	protected := router.Group("/api/sieksa/ekstrakurikuler")
	protected.Use(middleware.AuthMiddleware()) // Middleware will auto-detect SIEKSA from path
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
