package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterVisiMisiRoutes registers all visi misi routes
func RegisterVisiMisiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	visiMisiRepo := repositories.NewVisiMisiRepository(db)
	visiMisiService := services.NewVisiMisiService(visiMisiRepo)
	visiMisiController := controllers.NewVisiMisiController(visiMisiService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/visi-misi")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create visi misi
		protected.POST("/create-visi-misi", visiMisiController.Create)

		// Get all visi misi
		protected.POST("/get-visi-misi", visiMisiController.GetAll)

		// Get visi misi by ID
		protected.POST("/get-visi-misi-by-id", visiMisiController.GetByID)

		// Update visi misi
		protected.POST("/update-visi-misi", visiMisiController.Update)

		// Delete visi misi
		protected.POST("/delete-visi-misi", visiMisiController.Delete)
	}
}
