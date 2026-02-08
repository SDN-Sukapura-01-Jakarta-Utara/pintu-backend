package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	loginRepo := repositories.NewLoginRepository(db)
	loginService := services.NewLoginService(loginRepo)
	loginController := controllers.NewLoginController(loginService)

	// Public routes (no auth required)
	public := router.Group("/api/v1/auth")
	{
		public.POST("/login", loginController.Login) // Login endpoint
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", loginController.Logout) // Logout endpoint
	}
}
