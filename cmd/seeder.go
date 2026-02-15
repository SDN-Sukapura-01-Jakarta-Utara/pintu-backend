package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pintu-backend/src/database/seeders"
)

// runSeeders executes all seeders in order
func runSeeders() error {
	// Load .env file
	godotenv.Load()

	// Get database credentials
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || user == "" || dbName == "" {
		return fmt.Errorf("missing database credentials in environment variables")
	}

	if port == "" {
		port = "5432"
	}

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to database successfully!")

	// Run seeders in order
	seedersToRun := []struct {
		name   string
		seeder func(*gorm.DB) error
	}{
		{"System", runSystemSeeder},
		{"Permissions", runPermissionSeeder},
		{"Roles", runRoleSeeder},
		{"Role Permissions", runRolePermissionSeeder},
		{"Users", runUserSeeder},
	}

	for _, s := range seedersToRun {
		fmt.Printf("\n--- Seeding %s ---\n", s.name)
		if err := s.seeder(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", s.name, err)
		}
	}

	return nil
}

func runSystemSeeder(db *gorm.DB) error {
	return seeders.SystemSeeder(db)
}

func runPermissionSeeder(db *gorm.DB) error {
	seeder := seeders.NewPermissionSeeders(db)
	return seeder.Run()
}

func runRoleSeeder(db *gorm.DB) error {
	seeder := seeders.NewRoleSeeders(db)
	return seeder.Run()
}

func runRolePermissionSeeder(db *gorm.DB) error {
	seeder := seeders.NewRolePermissionSeeders(db)
	return seeder.Run()
}

func runUserSeeder(db *gorm.DB) error {
	seeder := seeders.NewUserSeeders(db)
	return seeder.Run()
}

// runSeedSpecific executes a specific seeder
func runSeedSpecific(seederName string) error {
	// Load .env file
	godotenv.Load()

	// Get database credentials
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || user == "" || dbName == "" {
		return fmt.Errorf("missing database credentials in environment variables")
	}

	if port == "" {
		port = "5432"
	}

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to database successfully!")

	// Run specific seeder
	switch seederName {
	case "permission":
		fmt.Printf("\n--- Seeding %s ---\n", seederName)
		return runPermissionSeeder(db)
	case "role":
		fmt.Printf("\n--- Seeding %s ---\n", seederName)
		return runRoleSeeder(db)
	case "role_permission":
		fmt.Printf("\n--- Seeding %s ---\n", seederName)
		return runRolePermissionSeeder(db)
	case "user":
		fmt.Printf("\n--- Seeding %s ---\n", seederName)
		return runUserSeeder(db)
	default:
		return fmt.Errorf("unknown seeder: %s. Available seeders: permission, role, role_permission, user", seederName)
	}
}
