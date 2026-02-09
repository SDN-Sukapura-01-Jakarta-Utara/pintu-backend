package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRombelRoutes registers all rombel routes
func RegisterRombelRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	rombelRepo := repositories.NewRombelRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	rombelService := services.NewRombelService(rombelRepo, kelasRepo)
	rombelController := controllers.NewRombelController(rombelService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/rombel")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create rombel
		protected.POST("/create-rombel", rombelController.Create)

		// Get all rombel
		protected.POST("/get-rombel", rombelController.GetAll)

		// Get rombel by ID
		protected.POST("/get-rombel-by-id", rombelController.GetByID)

		// Update rombel
		protected.POST("/update-rombel", rombelController.Update)

		// Delete rombel
		protected.POST("/delete-rombel", rombelController.Delete)
	}
}
