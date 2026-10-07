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

// RegisterMonitoringInklusiRoutes registers all monitoring inklusi routes
func RegisterMonitoringInklusiRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize R2 storage
	r2Storage := utils.NewR2Storage()

	// Initialize repository, service, and controller for anak inklusi
	anakInklusiRepo := repositories.NewAnakInklusiRepository(db)
	anakInklusiService := services.NewAnakInklusiService(anakInklusiRepo, r2Storage)
	anakInklusiController := controllers.NewAnakInklusiController(anakInklusiService)

	// Initialize repository, service, and controller for anak inklusi rombel
	anakInklusiRombelRepo := repositories.NewAnakInklusiRombelRepository(db)
	anakInklusiRombelService := services.NewAnakInklusiRombelService(anakInklusiRombelRepo, anakInklusiRepo, r2Storage)
	anakInklusiRombelController := controllers.NewAnakInklusiRombelController(anakInklusiRombelService)

	// Initialize repository, service, and controller for PPI inklusi
	ppiInklusiRepo := repositories.NewPPIInklusiRepository(db)
	ppiInklusiService := services.NewPPIInklusiService(ppiInklusiRepo, r2Storage)
	ppiInklusiController := controllers.NewPPIInklusiController(ppiInklusiService)

	// Initialize repository, service, and controller for monitoring inklusi bulanan
	monitoringInklusiBulananRepo := repositories.NewMonitoringInklusiBulananRepository(db)
	monitoringInklusiBulananService := services.NewMonitoringInklusiBulananService(monitoringInklusiBulananRepo, r2Storage)
	monitoringInklusiBulananController := controllers.NewMonitoringInklusiBulananController(monitoringInklusiBulananService)

	// Protected routes (auth required for PINTU)
	protected := router.Group("/api/v1/monitoring-inklusi")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create data induk inklusi
		protected.POST("/create-data-induk-inklusi", anakInklusiController.Create)

		// Get all data induk inklusi with filters and pagination
		protected.POST("/get-data-induk-inklusi", anakInklusiController.GetAll)

		// Get history data induk inklusi (non-active students)
		protected.POST("/get-history-data-induk-inklusi", anakInklusiController.GetHistory)

		// Get data induk inklusi by ID
		protected.POST("/get-data-induk-inklusi-by-id", anakInklusiController.GetByID)

		// Update data induk inklusi
		protected.POST("/update-data-induk-inklusi", anakInklusiController.Update)

		// Delete data induk inklusi
		protected.POST("/delete-data-induk-inklusi", anakInklusiController.Delete)

		// Sync data anak inklusi rombel by tahun pelajaran
		protected.POST("/sync-data-induk-inklusi-rombel", anakInklusiRombelController.SyncByTahunPelajaran)

		// Get all data anak inklusi rombel with filters and pagination
		protected.POST("/get-data-induk-inklusi-rombel", anakInklusiRombelController.GetAll)

		// Get data anak inklusi rombel by ID
		protected.POST("/get-data-induk-inklusi-rombel-by-id", anakInklusiRombelController.GetByID)

		// Update data anak inklusi rombel
		protected.POST("/update-data-induk-inklusi-rombel", anakInklusiRombelController.Update)

		// Delete data anak inklusi rombel
		protected.POST("/delete-data-induk-inklusi-rombel", anakInklusiRombelController.Delete)

		// Create PPI (Program Pembelajaran Individual)
		protected.POST("/create-ppi", ppiInklusiController.Create)

		// Get all PPI with filters and pagination
		protected.POST("/get-ppi", ppiInklusiController.GetAll)

		// Get PPI by ID
		protected.POST("/get-ppi-by-id", ppiInklusiController.GetByID)

		// Update PPI
		protected.POST("/update-ppi", ppiInklusiController.Update)

		// Delete PPI
		protected.POST("/delete-ppi", ppiInklusiController.Delete)

		// Get PPI summary
		protected.POST("/summary-ppi", ppiInklusiController.GetSummary)

		// Create monitoring inklusi bulanan
		protected.POST("/create-monitoring-bulanan", monitoringInklusiBulananController.Create)

		// Get monitoring inklusi bulanan
		protected.POST("/get-monitoring-bulanan", monitoringInklusiBulananController.GetByBulanTahunRombel)

		// Get monitoring inklusi bulanan by ID
		protected.POST("/get-monitoring-bulanan-by-id", monitoringInklusiBulananController.GetByID)

		// Update monitoring inklusi bulanan
		protected.POST("/update-monitoring-bulanan", monitoringInklusiBulananController.Update)

		// Get summary monitoring inklusi bulanan
		protected.POST("/get-summary-monitoring-bulanan", monitoringInklusiBulananController.GetSummary)
	}
}
