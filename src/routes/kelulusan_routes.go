package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterKelulusanRoutes registers all Kelulusan routes
func RegisterKelulusanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewKelulusanRepository(db)
	service := services.NewKelulusanService(repository)
	controller := controllers.NewKelulusanController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1/kelulusan")
	api.Use(middleware.AuthMiddleware())
	{
		// Create data kelulusan
		api.POST("/create-data-kelulusan", controller.CreateKelulusan)
		
		// Get all data kelulusan with filters
		api.POST("/get-data-kelulusan", controller.GetAll)
		
		// Get by ID
		api.POST("/get-data-kelulusan-by-id", controller.GetByID)
		
		// Update data kelulusan
		api.POST("/update-data-kelulusan", controller.Update)
		
		// Delete data kelulusan
		api.POST("/delete-data-kelulusan", controller.Delete)
		
		// Download template
		api.POST("/download-template", controller.DownloadTemplate)
		
		// Import Excel
		api.POST("/import-excel", controller.ImportExcel)
	}
}
