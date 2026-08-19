package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaRombelRoutes registers all SIEKSA rombel routes
func RegisterSieksaRombelRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller (reuse existing implementations)
	rombelRepo := repositories.NewRombelRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	rombelService := services.NewRombelService(rombelRepo, kelasRepo)
	rombelController := controllers.NewRombelController(rombelService)

	// Public routes for SIEKSA (no authentication required)
	public := router.Group("/api/sieksa/public")
	{
		// Get total rombel with status "active"
		public.POST("/get-total-rombel", rombelController.GetTotalRombel)
	}

	// Protected routes for SIEKSA (auth required with SIEKSA token)
	protected := router.Group("/api/sieksa/rombel")
	protected.Use(middleware.AuthMiddleware()) // Middleware will auto-detect SIEKSA from path
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
