package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaRoleRoutes registers SIEKSA role routes
func RegisterSieksaRoleRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and SIEKSA-specific controller
	roleRepo := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	sieksaRoleController := controllers.NewSieksaRoleController(roleService)

	// Group routes under /api/sieksa/roles with auth middleware
	api := router.Group("/api/sieksa/roles")
	api.Use(middleware.AuthMiddleware()) // Require authentication
	{
		api.POST("/get-roles", sieksaRoleController.GetAll) // Get all roles (filtered for SIEKSA)
	}
}
