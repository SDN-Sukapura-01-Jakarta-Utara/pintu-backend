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

// RegisterSaranaPrasaranaRoutes registers all sarana prasarana routes
func RegisterSaranaPrasaranaRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	saranaPrasaranaRepo := repositories.NewSaranaPrasaranaRepository(db)
	saranaPrasaranaService := services.NewSaranaPrasaranaService(saranaPrasaranaRepo, r2Storage)
	saranaPrasaranaController := controllers.NewSaranaPrasaranaController(saranaPrasaranaService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/sarana-prasarana")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create sarana prasarana with file upload
		protected.POST("/create-sarana-prasarana", saranaPrasaranaController.Create)

		// Get all sarana prasarana
		protected.POST("/get-sarana-prasarana", saranaPrasaranaController.GetAll)

		// Get sarana prasarana by ID
		protected.POST("/get-sarana-prasarana-by-id", saranaPrasaranaController.GetByID)

		// Update sarana prasarana
		protected.POST("/update-sarana-prasarana", saranaPrasaranaController.Update)

		// Delete sarana prasarana
		protected.POST("/delete-sarana-prasarana", saranaPrasaranaController.Delete)
	}
}
