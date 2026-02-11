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

// RegisterAnnouncementRoutes registers all announcement routes
func RegisterAnnouncementRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	announcementRepo := repositories.NewAnnouncementRepository(db)
	announcementService := services.NewAnnouncementService(announcementRepo, r2Storage)
	announcementController := controllers.NewAnnouncementController(announcementService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/announcements")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create announcement with gambar and files upload
		protected.POST("/create-announcement", announcementController.Create)

		// Get all announcements
		protected.POST("/get-announcements", announcementController.GetAll)

		// Get announcement by ID
		protected.POST("/get-announcement-by-id", announcementController.GetByID)

		// Update announcement (handle update fields, add files, delete files)
		protected.POST("/update-announcement", announcementController.Update)

		// Delete announcement
		protected.POST("/delete-announcement", announcementController.Delete)
	}
}
