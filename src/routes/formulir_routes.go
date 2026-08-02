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

// RegisterFormulirRoutes registers all formulir routes
func RegisterFormulirRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repositories
	formulirRepo := repositories.NewFormulirRepository(db)
	responseRepo := repositories.NewFormulirResponseRepository(db)

	// Initialize services
	formulirService := services.NewFormulirService(formulirRepo, r2Storage)
	responseService := services.NewFormulirResponseService(formulirRepo, responseRepo, r2Storage)

	// Initialize controllers
	formulirController := controllers.NewFormulirController(formulirService)
	responseController := controllers.NewFormulirResponseController(responseService)

	// Protected routes (auth required)
	protected := router.Group("/api/v1/formulir")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create formulir with pertanyaan and dokumen uploads
		protected.POST("/create-formulir", formulirController.Create)

		// Update formulir with pertanyaan and dokumen uploads
		protected.POST("/edit-formulir", formulirController.Update)
		protected.POST("/update-formulir", formulirController.Update) // Alias

		// Get all formulir with pagination and filters
		protected.POST("/get-formulir", formulirController.GetAll)

		// Get formulir by ID
		protected.POST("/get-formulir-by-id", formulirController.GetByID)

		// Get formulir by slug
		protected.POST("/get-formulir-by-slug", formulirController.GetBySlug)

		// Get formulir by user role (pendidik, tendik, murid)
		protected.POST("/get-formulir-by-user", formulirController.GetFormulirByUser)

		// Delete formulir
		protected.POST("/delete-formulir", formulirController.Delete)

		// Submit form response (authenticated)
		protected.POST("/submit-response-authenticated", responseController.SubmitAuthenticated)

		// Edit form response
		protected.POST("/edit-response", responseController.EditResponse)

		// Get responses by slug (Google Forms style)
		protected.POST("/get-response-by-slug", responseController.GetResponsesBySlug)

		// Get statistics by slug
		protected.POST("/get-statistic-by-slug", responseController.GetStatisticsBySlug)

		// Get user's own response by slug
		protected.POST("/get-response-by-user", responseController.GetResponseByUser)

		// Delete response by ID
		protected.POST("/delete-response-by-id", responseController.DeleteResponse)

		// Reset all responses for a formulir
		protected.POST("/reset-response", responseController.ResetResponses)
	}

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		// Get form by slug (public access)
		public.POST("/get-form-public", formulirController.GetBySlugPublic)

		// Submit form response (public)
		public.POST("/submit-response-public", responseController.SubmitPublic)
	}
}
