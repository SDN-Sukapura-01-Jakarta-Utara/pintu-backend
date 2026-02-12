package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterContactRoutes registers all contact routes
func RegisterContactRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	contactRepo := repositories.NewContactRepository(db)
	contactService := services.NewContactService(contactRepo)
	contactController := controllers.NewContactController(contactService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/contacts")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create contact
		protected.POST("/create-contact", contactController.Create)

		// Get all contacts
		protected.POST("/get-contacts", contactController.GetAll)

		// Get contact by ID
		protected.POST("/get-contact-by-id", contactController.GetByID)

		// Update contact
		protected.POST("/update-contact", contactController.Update)

		// Delete contact
		protected.POST("/delete-contact", contactController.Delete)
	}
}
