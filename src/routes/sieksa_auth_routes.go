package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaAuthRoutes registers all SIEKSA authentication routes
func RegisterSieksaAuthRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	loginRepo := repositories.NewLoginRepository(db)
	loginService := services.NewLoginService(loginRepo)
	loginController := controllers.NewLoginController(loginService)

	// Public routes for SIEKSA (no auth required)
	public := router.Group("/api/sieksa/auth")
	{
		public.POST("/login", loginController.LoginSieksa)                   // SIEKSA login endpoint
		public.POST("/login/student", loginController.LoginStudentSieksa)    // SIEKSA student login endpoint
	}

	// Protected routes for SIEKSA (auth required)
	protected := router.Group("/api/sieksa/auth")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", loginController.LogoutSieksa)                   // SIEKSA logout endpoint
		protected.POST("/logout/student", loginController.LogoutStudentSieksa)    // SIEKSA student logout endpoint
	}
}
