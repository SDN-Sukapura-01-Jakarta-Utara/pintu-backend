package main

import (
	"log"
	"os"

	"pintu-backend/src/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	godotenv.Load()

	// Database connection
	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

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

	// Setup CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3001", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

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

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}
