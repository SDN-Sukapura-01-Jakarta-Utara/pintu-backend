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

// RegisterPrestasiRoutes registers all Prestasi routes
func RegisterPrestasiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller
	repository := repositories.NewPrestasiRepository(db)
	service := services.NewPrestasiService(repository, r2Storage)
	controller := controllers.NewPrestasiController(service)

	// Protected routes (require authentication)
	api := router.Group("/api/v1")
	{
		PrestasiRoutes(api, controller)
	}

	// Public routes (no authentication required)
	publicAPI := router.Group("/api/v1/public")
	{
		PrestasiPublicRoutes(publicAPI, controller)
	}
}

// PrestasiRoutes sets up routes for prestasi endpoints
func PrestasiRoutes(router *gin.RouterGroup, controller *controllers.PrestasiController) {
	prestasiGroup := router.Group("/prestasi")
	prestasiGroup.Use(middleware.AuthMiddleware()) // Apply auth middleware to all prestasi routes
	{
		prestasiGroup.POST("/create-prestasi", controller.Create)
		prestasiGroup.POST("/get-prestasi", controller.GetAll)
		prestasiGroup.POST("/get-prestasi-by-id", controller.GetByID)
		prestasiGroup.POST("/update-prestasi", controller.Update)
		prestasiGroup.POST("/delete-prestasi", controller.Delete)
	}
}

// PrestasiPublicRoutes sets up public routes for prestasi endpoints (no auth required)
func PrestasiPublicRoutes(router *gin.RouterGroup, controller *controllers.PrestasiController) {
	router.POST("/get-data-prestasi", controller.GetPublicLatest)
}