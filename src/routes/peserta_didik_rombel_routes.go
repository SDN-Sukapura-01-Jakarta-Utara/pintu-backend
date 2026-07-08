package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPesertaDidikRombelRoutes registers all PesertaDidikRombel routes
func RegisterPesertaDidikRombelRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	pesertaDidikRombelRepo := repositories.NewPesertaDidikRombelRepository(db)
	pesertaDidikRepo := repositories.NewPesertaDidikRepository(db)

	// Initialize service
	service := services.NewPesertaDidikRombelService(pesertaDidikRombelRepo, pesertaDidikRepo)

	// Initialize controller
	controller := controllers.NewPesertaDidikRombelController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1/peserta-didik")
	api.Use(middleware.AuthMiddleware())
	{
		// Bulk create pemetaan rombel
		api.POST("/create-pemetaan-rombel", controller.BulkCreate)
		
		// Get pemetaan rombel with filters
		api.POST("/get-pemetaan-rombel", controller.GetAll)
		
		// Get pemetaan rombel by ID
		api.POST("/get-pemetaan-rombel-by-id", controller.GetByID)
		
		// Update pemetaan rombel
		api.POST("/edit-pemetaan-rombel-by-id", controller.Update)
		
		// Delete pemetaan rombel
		api.POST("/delete-pemetaan-rombel-by-id", controller.Delete)
		
		// Download Template
		api.POST("/download-template-pemetaan-rombel", controller.DownloadTemplate)
		
		// Import Excel
		api.POST("/import-excel-pemetaan-rombel", controller.ImportExcel)
		
		// Reset pemetaan rombel
		api.POST("/reset-pemetaan-rombel", controller.Reset)
	}
}
