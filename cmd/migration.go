package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

// runMigrations executes all migration files in order
func runMigrations() error {
	// Load .env file
	godotenv.Load()

	// Get PostgreSQL credentials from environment
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	psqlPath := os.Getenv("PSQL_PATH")

	if host == "" || user == "" || dbName == "" || psqlPath == "" {
		return fmt.Errorf("missing database credentials or PSQL_PATH in environment variables")
	}

	if port == "" {
		port = "5432"
	}

	// Get all migration files
	migrationDir := filepath.Join("src", "database", "migrations")
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	if len(sqlFiles) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	sort.Strings(sqlFiles)

	// Run each migration
	for _, sqlFile := range sqlFiles {
		filePath := filepath.Join(migrationDir, sqlFile)
		fmt.Printf("Running migration: %s\n", sqlFile)

		// Build psql command
		cmd := exec.Command(
			psqlPath,
			"-h", host,
			"-p", port,
			"-U", user,
			"-d", dbName,
			"-f", filePath,
		)

		// Set password via environment variable
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

		// Execute command
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("migration %s failed: %w\nOutput: %s", sqlFile, err, string(output))
		}

		fmt.Printf("✓ %s completed\n", sqlFile)
	}

	return nil
}

// runMigrationFile executes a specific migration file
func runMigrationFile(filename string) error {
	// Load .env file
	godotenv.Load()

	// Get PostgreSQL credentials from environment
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	psqlPath := os.Getenv("PSQL_PATH")

	if host == "" || user == "" || dbName == "" || psqlPath == "" {
		return fmt.Errorf("missing database credentials or PSQL_PATH in environment variables")
	}

	if port == "" {
		port = "5432"
	}

	// Check if file exists
	filePath := filepath.Join("src", "database", "migrations", filename)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("migration file not found: %s", filename)
	}

	// Run the migration
	fmt.Printf("Running migration: %s\n", filename)
	cmd := exec.Command(
		psqlPath,
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbName,
		"-f", filePath,
	)

	// Set password via environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("migration failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ %s completed\n", filename)
	return nil
}
