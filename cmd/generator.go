package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// Helper function to convert string to lowercase first letter
func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// Helper function to convert string to uppercase first letter
func toUpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Helper function to convert PascalCase to snake_case
// e.g., TahunPelajaran -> tahun_pelajaran
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// createMigrationFile generates migration file
func createMigrationFile(name string) error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := filepath.Join("src", "database", "migrations", filename)

	content := fmt.Sprintf(`-- Migration: %s
-- Created: %s

BEGIN;

-- Add your migration SQL here

COMMIT;
`, name, time.Now().Format("2006-01-02 15:04:05"))

	return writeFile(filepath, content)
}

// createSeederFile generates seeder file
func createSeederFile(name string) error {
	lowerName := toLowerFirst(name)
	filepath := filepath.Join("src", "database", "seeders", fmt.Sprintf("%s_seeder.go", lowerName))

	content := fmt.Sprintf(`package seeders

import (
	"fmt"
	"gorm.io/gorm"
)

// %sSeeder handles seeding data for %s
type %sSeeder struct {
	db *gorm.DB
}

// New%sSeeders creates a new %s seeder
func New%sSeeders(db *gorm.DB) *%sSeeder {
	return &%sSeeder{db: db}
}

// Run executes the seeder
func (s *%sSeeder) Run() error {
	fmt.Println("Seeding %s...")

	// Add your seeding logic here

	fmt.Println("%s seeding completed")
	return nil
}
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, name)

	return writeFile(filepath, content)
}

// createModelFile generates model file
func createModelFile(name string) error {
	lowerName := toLowerFirst(name)
	filepath := filepath.Join("src", "modules", "models", fmt.Sprintf("%s.go", lowerName))

	content := fmt.Sprintf(`package models

import (
	"time"

	"gorm.io/gorm"
)

// %s represents the %s model
type %s struct {
	ID        uint            ` + "`" + `gorm:"primaryKey"` + "`" + `
	// Add fields here
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index"` + "`" + `
}

// TableName specifies the table name for %s
func (m *%s) TableName() string {
	return "%s"
}
`, toUpperFirst(name), name, toUpperFirst(name), name, toUpperFirst(name), strings.ToLower(name)+"s")

	return writeFile(filepath, content)
}

// createRepositoryFile generates repository file
func createRepositoryFile(name string) error {
	snakeName := toSnakeCase(name)
	modelName := toUpperFirst(name)
	filepath := filepath.Join("src", "modules", "repositories", fmt.Sprintf("%s_repository.go", snakeName))

	template := `package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// ` + modelName + `Repository handles data operations for ` + modelName + `
type ` + modelName + `Repository interface {
	Create(data *models.` + modelName + `) error
	GetByID(id uint) (*models.` + modelName + `, error)
	GetAll() ([]models.` + modelName + `, error)
	Update(data *models.` + modelName + `) error
	Delete(id uint) error
}

type ` + modelName + `RepositoryImpl struct {
	db *gorm.DB
}

// New` + modelName + `Repository creates a new ` + modelName + ` repository
func New` + modelName + `Repository(db *gorm.DB) ` + modelName + `Repository {
	return &` + modelName + `RepositoryImpl{db: db}
}

// Create creates a new ` + modelName + ` record
func (r *` + modelName + `RepositoryImpl) Create(data *models.` + modelName + `) error {
	return r.db.Create(data).Error
}

// GetByID retrieves ` + modelName + ` by ID
func (r *` + modelName + `RepositoryImpl) GetByID(id uint) (*models.` + modelName + `, error) {
	var data models.` + modelName + `
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all ` + modelName + ` records
func (r *` + modelName + `RepositoryImpl) GetAll() ([]models.` + modelName + `, error) {
	var data []models.` + modelName + `
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates ` + modelName + ` record
func (r *` + modelName + `RepositoryImpl) Update(data *models.` + modelName + `) error {
	return r.db.Save(data).Error
}

// Delete deletes ` + modelName + ` record by ID
func (r *` + modelName + `RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.` + modelName + `{}, id).Error
}
`

	return writeFile(filepath, template)
}

// createServiceFile generates service file
func createServiceFile(name string) error {
	snakeName := toSnakeCase(name)
	modelName := toUpperFirst(name)
	filepath := filepath.Join("src", "modules", "services", fmt.Sprintf("%s_service.go", snakeName))

	template := `package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// ` + modelName + `Service handles business logic for ` + modelName + `
type ` + modelName + `Service interface {
	Create(data *models.` + modelName + `) error
	GetByID(id uint) (*models.` + modelName + `, error)
	GetAll() ([]models.` + modelName + `, error)
	Update(data *models.` + modelName + `) error
	Delete(id uint) error
}

type ` + modelName + `ServiceImpl struct {
	repository repositories.` + modelName + `Repository
}

// New` + modelName + `Service creates a new ` + modelName + ` service
func New` + modelName + `Service(repository repositories.` + modelName + `Repository) ` + modelName + `Service {
	return &` + modelName + `ServiceImpl{repository: repository}
}

// Create creates a new ` + modelName + `
func (s *` + modelName + `ServiceImpl) Create(data *models.` + modelName + `) error {
	return s.repository.Create(data)
}

// GetByID retrieves ` + modelName + ` by ID
func (s *` + modelName + `ServiceImpl) GetByID(id uint) (*models.` + modelName + `, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all ` + modelName + `
func (s *` + modelName + `ServiceImpl) GetAll() ([]models.` + modelName + `, error) {
	return s.repository.GetAll()
}

// Update updates ` + modelName + `
func (s *` + modelName + `ServiceImpl) Update(data *models.` + modelName + `) error {
	return s.repository.Update(data)
}

// Delete deletes ` + modelName + ` by ID
func (s *` + modelName + `ServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
`

	return writeFile(filepath, template)
}

// createControllerFile generates controller file
func createControllerFile(name string) error {
	snakeName := toSnakeCase(name)
	filepath := filepath.Join("src", "modules", "controllers", fmt.Sprintf("%s_controller.go", snakeName))

	content := fmt.Sprintf(`package controllers

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// %sController handles HTTP requests for %s
type %sController struct {
	service services.%sService
}

// New%sController creates a new %s controller
func New%sController(service services.%sService) *%sController {
	return &%sController{service: service}
}

// Create creates a new %s
// @Summary Create new %s
// @Description Create a new %s
// @Tags %s
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Router /%s [post]
func (c *%sController) Create(ctx *gin.Context) {
	var req models.%s
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": req})
}

// GetByID retrieves %s by ID
// @Summary Get %s by ID
// @Description Retrieve %s details by ID
// @Tags %s
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /%s/{id} [get]
func (c *%sController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all %s
// @Summary Get all %s
// @Description Retrieve all %s records
// @Tags %s
// @Produce json
// @Success 200
// @Failure 500
// @Router /%s [get]
func (c *%sController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates %s
// @Summary Update %s
// @Description Update %s details
// @Tags %s
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /%s/{id} [put]
func (c *%sController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.%s
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = uint(id)
	if err := c.service.Update(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": req})
}

// Delete deletes %s by ID
// @Summary Delete %s
// @Description Delete %s by ID
// @Tags %s
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /%s/{id} [delete]
func (c *%sController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, name, name, snakeName, snakeName, toUpperFirst(name), toUpperFirst(name), name, name, snakeName, toUpperFirst(name), name, name, snakeName, toUpperFirst(name), name, name, snakeName, toUpperFirst(name), name, name, name, snakeName, toUpperFirst(name), toUpperFirst(name), name, name, snakeName, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name))

	return writeFile(filepath, content)
}

// createDTOFile generates DTO file with Request and Response structs
func createDTOFile(name string) error {
	snakeName := toSnakeCase(name)
	filePath := filepath.Join("src", "dtos", fmt.Sprintf("%s_dto.go", snakeName))

	content := fmt.Sprintf(`package dtos

// %sCreateRequest represents the request payload for creating %s
type %sCreateRequest struct {
	// Add fields here
}

// %sUpdateRequest represents the request payload for updating %s
type %sUpdateRequest struct {
	// Add fields here
}

// %sResponse represents the response payload for %s
type %sResponse struct {
	ID uint ` + "`" + `json:"id"` + "`" + `
	// Add fields here
}

// %sListResponse represents the response payload for listing %s
type %sListResponse struct {
	Data []%sResponse ` + "`" + `json:"data"` + "`" + `
	Total int64 ` + "`" + `json:"total"` + "`" + `
}
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name))

	return writeFile(filePath, content)
}

// writeFile writes content to a file
func writeFile(filePath string, content string) error {
	// Create directory if not exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
