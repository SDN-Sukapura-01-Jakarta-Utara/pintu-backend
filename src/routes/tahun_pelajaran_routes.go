package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterTahunPelajaranRoutes registers all tahun pelajaran routes
func RegisterTahunPelajaranRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	tahunPelajaranRepo := repositories.NewTahunPelajaranRepository(db)
	tahunPelajaranService := services.NewTahunPelajaranService(tahunPelajaranRepo)
	tahunPelajaranController := controllers.NewTahunPelajaranController(tahunPelajaranService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/tahun-pelajaran")
	protected.Use(middleware.AuthMiddleware())
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
