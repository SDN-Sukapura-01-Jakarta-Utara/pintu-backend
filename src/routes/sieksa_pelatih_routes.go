package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaPelatihRoutes registers SIEKSA pelatih routes
func RegisterSieksaPelatihRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services
	pelatihRepo := repositories.NewPelatihRepository(db)
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	pelatihService := services.NewPelatihService(pelatihRepo, ekstrakurikulerRepo)
	pelatihController := controllers.NewPelatihController(pelatihService)

	// Protected routes for SIEKSA pelatih (auth required with SIEKSA token)
	pelatihGroup := router.Group("/api/sieksa/pelatih")
	pelatihGroup.Use(middleware.AuthMiddleware())
	{
		// Create pelatih
		pelatihGroup.POST("/create-pelatih", pelatihController.Create)

		// Get all pelatih
		pelatihGroup.POST("/get-pelatih", pelatihController.GetAll)

		// Get pelatih by ID
		pelatihGroup.POST("/get-pelatih-by-id", pelatihController.GetByID)

		// Update pelatih
		pelatihGroup.POST("/update-pelatih", pelatihController.Update)

		// Delete pelatih
		pelatihGroup.POST("/delete-pelatih", pelatihController.Delete)
	}
}
