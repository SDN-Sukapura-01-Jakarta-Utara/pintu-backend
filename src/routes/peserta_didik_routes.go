package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPesertaDidikRoutes registers all PesertaDidik routes
func RegisterPesertaDidikRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	repository := repositories.NewPesertaDidikRepository(db)
	service := services.NewPesertaDidikService(repository)
	controller := controllers.NewPesertaDidikController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1/peserta-didik")
	api.Use(middleware.AuthMiddleware())
	{
		// Create
		api.POST("/create-peserta-didik", controller.Create)

		// Read
		api.POST("/get-peserta-didik", controller.GetAll)
		api.POST("/get-peserta-didik-by-id", controller.GetByID)
		api.POST("/get-peserta-didik-by-nis", controller.GetByNIS)

		// Update
		api.POST("/update-peserta-didik", controller.Update)

		// Delete
		api.POST("/delete-peserta-didik", controller.Delete)

		// Import Excel
		api.POST("/import-excel", controller.ImportExcel)

		// Download Template
		api.POST("/download-template", controller.DownloadTemplate)
	}
}
