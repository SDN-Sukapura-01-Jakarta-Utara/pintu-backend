package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaTahunPelajaranRoutes registers all SIEKSA tahun pelajaran routes
func RegisterSieksaTahunPelajaranRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller (reuse existing implementations)
	tahunPelajaranRepo := repositories.NewTahunPelajaranRepository(db)
	tahunPelajaranService := services.NewTahunPelajaranService(tahunPelajaranRepo)
	tahunPelajaranController := controllers.NewTahunPelajaranController(tahunPelajaranService)

	// Public routes for SIEKSA (no authentication required)
	publicAPI := router.Group("/api/sieksa/public")
	{
		// Get active tahun pelajaran
		publicAPI.POST("/get-tahun-pelajaran-aktif", tahunPelajaranController.GetActiveTahunPelajaran)
	}

	// Protected routes for SIEKSA (auth required with SIEKSA token)
	protected := router.Group("/api/sieksa/tahun-pelajaran")
	protected.Use(middleware.AuthMiddleware()) // Middleware will auto-detect SIEKSA from path
	{
		// Create tahun pelajaran
		protected.POST("/create-tahun-pelajaran", tahunPelajaranController.Create)

		// Get all tahun pelajaran
		protected.POST("/get-tahun-pelajaran", tahunPelajaranController.GetAll)

		// Get tahun pelajaran by ID
		protected.POST("/get-tahun-pelajaran-by-id", tahunPelajaranController.GetByID)

		// Update tahun pelajaran
		protected.POST("/update-tahun-pelajaran", tahunPelajaranController.Update)

		// Delete tahun pelajaran
		protected.POST("/delete-tahun-pelajaran", tahunPelajaranController.Delete)
	}
}
