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

// RegisterKepegawaianRoutes registers all Kepegawaian routes
func RegisterKepegawaianRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	repository := repositories.NewKepegawaianRepository(db)
	service := services.NewKepegawaianService(repository, r2Storage)
	controller := controllers.NewKepegawaianController(service)

	// Public routes (no authentication required)
	public := router.Group("/api/v1/public")
	{
		// Get total pendidik with kategori "Pendidik" and status "active"
		public.POST("/get-total-pendidik", controller.GetTotalPendidik)
		
		// Get total tendik with kategori "Tenaga Kependidikan" and status "active"
		public.POST("/get-total-tendik", controller.GetTotalTendik)
		
		// Get public pendidik data (nama, nip, nkki, jabatan, foto) with kategori "Pendidik" and status "active"
		public.POST("/get-data-pendidik", controller.GetPublicPendidikData)
		
		// Get public tendik data (nama, nip, nkki, jabatan, foto) with kategori "Tenaga Kependidikan" and status "active"
		public.POST("/get-data-tendik", controller.GetPublicTendikData)
	}

	// Protected routes (require authentication)
	api := router.Group("/api/v1/kepegawaian")
	api.Use(middleware.AuthMiddleware())
	{
		// Create
		api.POST("/create-kepegawaian", controller.Create)

		// Read
		api.POST("/get-kepegawaian", controller.GetAll)
		api.POST("/get-kepegawaian-without-pagination", controller.GetAllWithoutPagination)
		api.POST("/get-kepegawaian-by-id", controller.GetByID)
		api.POST("/get-kepegawaian-by-nip", controller.GetByNIP)

		// Update
		api.POST("/update-kepegawaian", controller.Update)

		// Delete
		api.POST("/delete-kepegawaian", controller.Delete)
	}
}
