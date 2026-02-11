package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterActivityGalleryRoutes registers all activity gallery routes
func RegisterActivityGalleryRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	galleryRepo := repositories.NewActivityGalleryRepository(db)
	galleryService := services.NewActivityGalleryService(galleryRepo, r2Storage)
	galleryController := controllers.NewActivityGalleryController(galleryService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/activity-galleries")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create activity gallery with fotos upload
		protected.POST("/create-gallery", galleryController.Create)

		// Get all activity galleries
		protected.POST("/get-galleries", galleryController.GetAll)

		// Get activity gallery by ID
		protected.POST("/get-gallery-by-id", galleryController.GetByID)

		// Update activity gallery (handle update fields, add fotos, delete fotos)
		protected.POST("/update-gallery", galleryController.Update)

		// Delete activity gallery
		protected.POST("/delete-gallery", galleryController.Delete)
	}
}
