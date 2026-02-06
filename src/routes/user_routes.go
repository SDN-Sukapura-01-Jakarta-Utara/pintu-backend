package routes

import (
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUserRoutes registers all user routes
func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Group routes under /api/v1/users - all POST with action-based routes
	api := router.Group("/api/v1/users")
	{
		api.POST("/create-user", userController.Create)           // Create user
		api.POST("/get-users", userController.GetAll)             // Get all users
		api.POST("/get-user", userController.GetByID)             // Get user by ID
		api.POST("/update-user", userController.Update)           // Update user
		api.POST("/update-user-password", userController.UpdatePassword) // Update password
		api.POST("/delete-user", userController.Delete)           // Delete user
	}
}
