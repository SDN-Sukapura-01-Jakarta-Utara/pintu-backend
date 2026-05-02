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

// RegisterPertanyaanRoutes registers all pertanyaan routes
func RegisterPertanyaanRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	pertanyaanRepo := repositories.NewPertanyaanRepository(db)
	pertanyaanService := services.NewPertanyaanService(pertanyaanRepo, r2Storage)
	pertanyaanController := controllers.NewPertanyaanController(pertanyaanService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/pertanyaan")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/get-pertanyaan", pertanyaanController.GetAll)
		protected.POST("/get-pertanyaan-by-id", pertanyaanController.GetByID)
		protected.POST("/send-reply", pertanyaanController.SendReply)
		protected.POST("/close-pertanyaan", pertanyaanController.ClosePertanyaan)
		protected.POST("/delete-pertanyaan", pertanyaanController.DeletePertanyaan)
	}

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		public.POST("/create-pertanyaan", pertanyaanController.CreatePublic)
		public.POST("/track-pertanyaan", pertanyaanController.TrackPertanyaan)
	}
}
