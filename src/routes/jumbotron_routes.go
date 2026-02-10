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

// RegisterJumbotronRoutes registers all jumbotron routes
func RegisterJumbotronRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	jumbotronRepo := repositories.NewJumbotronRepository(db)
	jumbotronService := services.NewJumbotronService(jumbotronRepo, r2Storage)
	jumbotronController := controllers.NewJumbotronController(jumbotronService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/jumbotron")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create jumbotron with file upload
		protected.POST("/create-jumbotron", jumbotronController.Create)

		// Get all jumbotron
		protected.POST("/get-jumbotron", jumbotronController.GetAll)

		// Get jumbotron by ID
		protected.POST("/get-jumbotron-by-id", jumbotronController.GetByID)

		// Update jumbotron
		protected.POST("/update-jumbotron", jumbotronController.Update)

		// Delete jumbotron
		protected.POST("/delete-jumbotron", jumbotronController.Delete)
	}
}
