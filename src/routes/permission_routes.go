package routes

import (
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPermissionRoutes registers all permission routes
func RegisterPermissionRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	permissionRepo := repositories.NewPermissionRepository(db)
	permissionService := services.NewPermissionService(permissionRepo)
	permissionController := controllers.NewPermissionController(permissionService)

	// Group routes under /api/v1/permissions
	api := router.Group("/api/v1/permissions")
	{
		// CRUD operations - all POST with action-based routes
		api.POST("/create-permission", permissionController.Create)                    // Create permission
		api.POST("/get-permissions", permissionController.GetAll)                     // Get all permissions (with pagination)
		api.POST("/get-permission", permissionController.GetByID)                     // Get permission by ID
		api.POST("/update-permission", permissionController.Update)                   // Update permission
		api.POST("/delete-permission", permissionController.Delete)                   // Delete permission

		// Filter operations
		api.POST("/get-permissions-by-group", permissionController.GetByGroupName)    // Get by group name
		api.POST("/get-permissions-by-system", permissionController.GetBySystem)      // Get by system
	}
}
