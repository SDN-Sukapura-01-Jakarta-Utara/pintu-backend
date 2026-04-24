package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterStrukturOrganisasiRoutes registers all struktur organisasi routes
func RegisterStrukturOrganisasiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories, service, and controller
	strukturOrganisasiRepo := repositories.NewStrukturOrganisasiRepository(db)
	kepegawaianRepo := repositories.NewKepegawaianRepository(db)
	strukturOrganisasiService := services.NewStrukturOrganisasiService(strukturOrganisasiRepo, kepegawaianRepo)
	strukturOrganisasiController := controllers.NewStrukturOrganisasiController(strukturOrganisasiService)

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/get-data-struktur-organisasi", strukturOrganisasiController.GetPublic)
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/struktur-organisasi")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create struktur organisasi
		protected.POST("/create-struktur-organisasi", strukturOrganisasiController.Create)

		// Get all struktur organisasi
		protected.POST("/get-struktur-organisasi", strukturOrganisasiController.GetAll)

		// Get struktur organisasi by ID
		protected.POST("/get-struktur-organisasi-by-id", strukturOrganisasiController.GetByID)

		// Update struktur organisasi
		protected.POST("/update-struktur-organisasi", strukturOrganisasiController.Update)

		// Delete struktur organisasi
		protected.POST("/delete-struktur-organisasi", strukturOrganisasiController.Delete)
	}
}
