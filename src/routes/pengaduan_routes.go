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

// RegisterPengaduanRoutes registers all pengaduan routes
func RegisterPengaduanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	pengaduanRepo := repositories.NewPengaduanRepository(db)
	pengaduanService := services.NewPengaduanService(pengaduanRepo, r2Storage)
	pengaduanController := controllers.NewPengaduanController(pengaduanService)

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/create-pengaduan", pengaduanController.CreatePublic)
		public.POST("/track-pengaduan", pengaduanController.TrackPengaduan)
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/pengaduan")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/get-pengaduan", pengaduanController.GetAll)
		protected.POST("/get-pengaduan-by-id", pengaduanController.GetByID)
		protected.POST("/send-reply", pengaduanController.SendReply)
		protected.POST("/save-tindak-lanjut", pengaduanController.SaveTindakLanjut)
		protected.POST("/close-pengaduan", pengaduanController.ClosePengaduan)
		protected.POST("/delete-pengaduan", pengaduanController.DeletePengaduan)
	}
}