package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoleRoutes registers all role routes
func RegisterRoleRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	roleRepo := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleController := controllers.NewRoleController(roleService)

	// Group routes under /api/v1/roles with auth middleware
	api := router.Group("/api/v1/roles")
	api.Use(middleware.AuthMiddleware()) // Require authentication
	{
		api.POST("/create-role", roleController.Create)      // Create role
		api.POST("/get-roles", roleController.GetAll)        // Get all roles
		api.POST("/get-role-by-id", roleController.GetByID)        // Get role by ID
		api.POST("/update-role", roleController.Update)      // Update role
		api.POST("/delete-role", roleController.Delete)      // Delete role
	}
}
