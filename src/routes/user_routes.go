package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUserRoutes registers all user routes
func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// Initialize service with role validation
	userService := services.NewUserServiceWithRole(userRepo, roleRepo)
	userController := controllers.NewUserController(userService)

	// Group routes under /api/v1/users with auth middleware
	api := router.Group("/api/v1/users")
	api.Use(middleware.AuthMiddleware()) // Require authentication
	{
		api.POST("/create-user", userController.Create)           // Create user
		api.POST("/get-users", userController.GetAll)             // Get all users
		api.POST("/get-user", userController.GetByID)             // Get user by ID
		api.POST("/update-user", userController.Update)           // Update user
		api.POST("/update-user-password", userController.UpdatePassword) // Update password
		api.POST("/delete-user", userController.Delete)           // Delete user
	}
}
