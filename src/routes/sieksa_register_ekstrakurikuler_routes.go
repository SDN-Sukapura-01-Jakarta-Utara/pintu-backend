package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSieksaRegisterEkstrakurikulerRoutes registers SIEKSA register ekstrakurikuler routes
func RegisterSieksaRegisterEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories and services
	pesertaDidikEkskulRepo := repositories.NewPesertaDidikEkstrakurikulerRepository(db)
	pesertaDidikRombelRepo := repositories.NewPesertaDidikRombelRepository(db)
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	pesertaDidikEkskulService := services.NewPesertaDidikEkstrakurikulerService(pesertaDidikEkskulRepo, pesertaDidikRombelRepo, ekstrakurikulerRepo)
	pesertaDidikEkskulController := controllers.NewPesertaDidikEkstrakurikulerController(pesertaDidikEkskulService)

	// Protected routes for SIEKSA register ekstrakurikuler (auth required with SIEKSA token)
	registerGroup := router.Group("/api/sieksa/register-ekstrakurikuler")
	registerGroup.Use(middleware.AuthMiddleware())
	{
		// Register/update student ekstrakurikuler
		registerGroup.POST("/register-ekskul-peserta-didik", pesertaDidikEkskulController.RegisterOrUpdateEkstrakurikuler)

		// Get student's ekstrakurikuler
		registerGroup.POST("/get-ekskul-peserta-didik", pesertaDidikEkskulController.GetEkstrakurikulerByPesertaDidik)

		// Get all students' ekstrakurikuler by rombel
		registerGroup.POST("/get-all-ekstrakurikuler-siswa", pesertaDidikEkskulController.GetAllEkstrakurikulerSiswa)

		// Bulk register/update ekstrakurikuler for multiple students
		registerGroup.POST("/register-all-ekstrakurikuler-siswa", pesertaDidikEkskulController.RegisterAllEkstrakurikulerSiswa)
	}
}
