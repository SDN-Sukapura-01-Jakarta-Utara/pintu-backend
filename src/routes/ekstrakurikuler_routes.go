package routes

import (
	"pintu-backend/src/middleware"
	"pintu-backend/src/modules/controllers"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterEkstrakurikulerRoutes registers all ekstrakurikuler routes
func RegisterEkstrakurikulerRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repository, service, and controller
	ekstrakurikulerRepo := repositories.NewEkstrakurikulerRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	ekstrakurikulerService := services.NewEkstrakurikulerService(ekstrakurikulerRepo, kelasRepo)
	ekstrakurikulerController := controllers.NewEkstrakurikulerController(ekstrakurikulerService)

	// Initialize peserta didik ekstrakurikuler (registration)
	pesertaDidikEkskulRepo := repositories.NewPesertaDidikEkstrakurikulerRepository(db)
	pesertaDidikRombelRepo := repositories.NewPesertaDidikRombelRepository(db)
	pesertaDidikEkskulService := services.NewPesertaDidikEkstrakurikulerService(pesertaDidikEkskulRepo, pesertaDidikRombelRepo, ekstrakurikulerRepo)
	pesertaDidikEkskulController := controllers.NewPesertaDidikEkstrakurikulerController(pesertaDidikEkskulService)

	// Public routes (no authentication required)
	public := router.Group("/api/v1/public")
	{
		// Get total ekstrakurikuler with status "active"
		public.POST("/get-total-ekskul", ekstrakurikulerController.GetTotalEkskul)
	}

	// Protected routes (auth required)
	protected := router.Group("/api/v1/ekstrakurikuler")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create ekstrakurikuler
		protected.POST("/create-ekstrakurikuler", ekstrakurikulerController.Create)

		// Get all ekstrakurikuler
		protected.POST("/get-ekstrakurikuler", ekstrakurikulerController.GetAll)

		// Get ekstrakurikuler by ID
		protected.POST("/get-ekstrakurikuler-by-id", ekstrakurikulerController.GetByID)

		// Update ekstrakurikuler
		protected.POST("/update-ekstrakurikuler", ekstrakurikulerController.Update)

		// Delete ekstrakurikuler
		protected.POST("/delete-ekstrakurikuler", ekstrakurikulerController.Delete)

		// Register/update student ekstrakurikuler
		protected.POST("/register-ekskul-peserta-didik", pesertaDidikEkskulController.RegisterOrUpdateEkstrakurikuler)

		// Get student's ekstrakurikuler
		protected.POST("/get-ekskul-peserta-didik", pesertaDidikEkskulController.GetEkstrakurikulerByPesertaDidik)

		// Get all students' ekstrakurikuler by rombel
		protected.POST("/get-all-ekstrakurikuler-siswa", pesertaDidikEkskulController.GetAllEkstrakurikulerSiswa)

		// Bulk register/update ekstrakurikuler for multiple students
		protected.POST("/register-all-ekstrakurikuler-siswa", pesertaDidikEkskulController.RegisterAllEkstrakurikulerSiswa)

		// Get comprehensive statistics for monitoring
		protected.POST("/get-all-statistic-ekstrakurikuler-siswa", pesertaDidikEkskulController.GetStatistikEkstrakurikuler)

		// Get rekapitulasi data per ekstrakurikuler
		protected.POST("/rekapitulasi-data-per-ekskul", pesertaDidikEkskulController.GetRekapitulasiPerEkskul)

		// Get rekapitulasi data per rombel
		protected.POST("/rekapitulasi-data-per-rombel", pesertaDidikEkskulController.GetRekapitulasiPerRombel)
	}
}
