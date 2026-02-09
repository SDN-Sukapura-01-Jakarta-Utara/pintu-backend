package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterBidangStudiRoutes registers all bidang studi routes
func RegisterBidangStudiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	bidangStudiRepo := repositories.NewBidangStudiRepository(db)
	bidangStudiService := services.NewBidangStudiService(bidangStudiRepo)
	bidangStudiController := controllers.NewBidangStudiController(bidangStudiService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/bidang-studi")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create bidang studi
		protected.POST("/create-bidang-studi", bidangStudiController.Create)

		// Get all bidang studi
		protected.POST("/get-bidang-studi", bidangStudiController.GetAll)

		// Get bidang studi by ID
		protected.POST("/get-bidang-studi-by-id", bidangStudiController.GetByID)

		// Update bidang studi
		protected.POST("/update-bidang-studi", bidangStudiController.Update)

		// Delete bidang studi
		protected.POST("/delete-bidang-studi", bidangStudiController.Delete)
	}
}
