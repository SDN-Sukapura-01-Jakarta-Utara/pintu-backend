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

// RegisterKutipanKepsekRoutes registers all kutipan kepsek routes
func RegisterKutipanKepsekRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	kutipanKepsekRepo := repositories.NewKutipanKepsekRepository(db)
	kutipanKepsekService := services.NewKutipanKepsekService(kutipanKepsekRepo, r2Storage)
	kutipanKepsekController := controllers.NewKutipanKepsekController(kutipanKepsekService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/kutipan-kepsek")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create kutipan kepsek with file upload
		protected.POST("/create-kutipan-kepsek", kutipanKepsekController.Create)

		// Get all kutipan kepsek
		protected.POST("/get-kutipan-kepsek", kutipanKepsekController.GetAll)

		// Get kutipan kepsek by ID
		protected.POST("/get-kutipan-kepsek-by-id", kutipanKepsekController.GetByID)

		// Update kutipan kepsek
		protected.POST("/update-kutipan-kepsek", kutipanKepsekController.Update)

		// Delete kutipan kepsek
		protected.POST("/delete-kutipan-kepsek", kutipanKepsekController.Delete)
	}
}
