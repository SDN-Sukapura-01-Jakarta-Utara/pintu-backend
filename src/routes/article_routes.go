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

// RegisterArticleRoutes registers all article routes
func RegisterArticleRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	articleRepo := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepo, r2Storage)
	articleController := controllers.NewArticleController(articleService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/articles")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create article with gambar and files upload
		protected.POST("/create-article", articleController.Create)

		// Get all articles
		protected.POST("/get-articles", articleController.GetAll)

		// Get article by ID
		protected.POST("/get-article-by-id", articleController.GetByID)

		// Update article (handle update fields, add files, delete files)
		protected.POST("/update-article", articleController.Update)

		// Delete article
		protected.POST("/delete-article", articleController.Delete)
	}
}
