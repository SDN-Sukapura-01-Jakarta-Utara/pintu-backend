package main

import (
	"log"
	"os"
	"time"

	"pintu-backend/src/middleware"
	"pintu-backend/src/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Force UTC timezone for entire application
	time.Local = time.UTC

	// Load .env
	godotenv.Load()

	// Database connection with timezone set to UTC
	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	// Create router
	router := gin.Default()

	// Define allowed origins
	allowedOrigins := []string{
		"http://localhost:3001",
		"http://localhost:3000",
		"http://localhost:4001",
		"https://sdnsukapura01.sch.id",
		"http://sdnsukapura01.sch.id",
		"https://www.sdnsukapura01.sch.id",
		"http://www.sdnsukapura01.sch.id",
	}

	// Handle OPTIONS globally BEFORE any middleware (for CORS preflight)
	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			origin := c.GetHeader("Origin")
			
			// Validate origin is in allowed list
			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isAllowed = true
					break
				}
			}
			
			if !isAllowed {
				c.AbortWithStatus(403)
				return
			}
			
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "43200")
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup CORS middleware - must be before other middlewares
	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Preflight cache for 12 hours
		AllowWildcard:    false,
		AllowAllOrigins:  false,
	}
	router.Use(cors.New(corsConfig))

	// Setup Prometheus middleware
	router.Use(middleware.PrometheusMiddleware())

	// Metrics endpoint (Prometheus)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PINTU Backend is running",
			"app":     "PINTU SDN Sukapura 01",
		})
	})

	// Register routes
	routes.RegisterAuthRoutes(router, db)
	routes.RegisterSieksaAuthRoutes(router, db) // SIEKSA authentication routes
	routes.RegisterSieksaKelasRoutes(router, db) // SIEKSA kelas routes
	routes.RegisterSieksaRoleRoutes(router, db) // SIEKSA role routes
	routes.RegisterSieksaEkstrakurikulerRoutes(router, db) // SIEKSA ekstrakurikuler CRUD routes
	routes.RegisterSieksaPelatihRoutes(router, db) // SIEKSA pelatih routes
	routes.RegisterSieksaRegisterEkstrakurikulerRoutes(router, db) // SIEKSA register ekstrakurikuler routes
	routes.RegisterSieksaRekapitulasiEkstrakurikulerRoutes(router, db) // SIEKSA rekapitulasi ekstrakurikuler routes
	routes.RegisterSieksaMonitoringEkstrakurikulerRoutes(router, db) // SIEKSA monitoring ekstrakurikuler routes
	routes.RegisterSieksaAbsensiEkstrakurikulerRoutes(router, db) // SIEKSA absensi ekstrakurikuler routes
	routes.RegisterSieksaKegiatanEkstrakurikulerRoutes(router, db) // SIEKSA kegiatan ekstrakurikuler routes
	routes.RegisterSieksaRombelRoutes(router, db) // SIEKSA rombel routes
	routes.RegisterSieksaTahunPelajaranRoutes(router, db) // SIEKSA tahun pelajaran routes
	routes.RegisterSieksaPesertaDidikRoutes(router, db) // SIEKSA peserta didik routes
	routes.RegisterSieksaPesertaDidikRombelRoutes(router, db) // SIEKSA peserta didik rombel routes
	routes.RegisterSystemRoutes(router, db)
	routes.RegisterPermissionRoutes(router, db)
	routes.RegisterRoleRoutes(router, db)
	routes.RegisterUserRoutes(router, db)
	routes.RegisterTahunPelajaranRoutes(router, db)
	routes.RegisterBidangStudiRoutes(router, db)
	routes.RegisterKelasRoutes(router, db)
	routes.RegisterRombelRoutes(router, db)
	routes.RegisterEkstrakurikulerRoutes(router, db)
	routes.RegisterJumbotronRoutes(router, db)
	routes.RegisterKutipanKepsekRoutes(router, db)
	routes.RegisterVisiMisiRoutes(router, db)
	routes.RegisterSaranaPrasaranaRoutes(router, db)
	routes.RegisterArticleRoutes(router, db)
	routes.RegisterAnnouncementRoutes(router, db)
	routes.RegisterActivityGalleryRoutes(router, db)
	routes.RegisterContactRoutes(router, db)
	routes.RegisterKepegawaianRoutes(router, db)
	routes.RegisterStrukturOrganisasiRoutes(router, db)
	routes.RegisterPesertaDidikRoutes(router, db)
	routes.RegisterPesertaDidikRombelRoutes(router, db)
	routes.RegisterPrestasiRoutes(router, db)
	routes.RegisterApplicationRoutes(router, db)
	routes.RegisterKritikSaranRoutes(router, db)
	routes.RegisterPertanyaanRoutes(router, db)
	routes.RegisterPengaduanRoutes(router, db)
	routes.RegisterAbsensiRoutes(router, db)
	routes.RegisterAbsensiScanRoutes(router, db)
	routes.RegisterKonfigurasiAbsensiRoutes(router, db)
	routes.RegisterKelulusanRoutes(router, db)
	routes.RegisterPengumumanKelulusanRoutes(router, db)
	routes.RegisterLayananSPMBRoutes(router, db)
	routes.RegisterMutasiSiswaRoutes(router, db)
	routes.RegisterFormulirRoutes(router, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}
