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
	lowerName := toLowerFirst(name)
	filepath := filepath.Join("src", "modules", "repositories", fmt.Sprintf("%s_repository.go", lowerName))

	content := fmt.Sprintf(`package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// %sRepository handles data operations for %s
type %sRepository interface {
	Create(data *models.%s) error
	GetByID(id uint) (*models.%s, error)
	GetAll() ([]models.%s, error)
	Update(data *models.%s) error
	Delete(id uint) error
}

type %sRepositoryImpl struct {
	db *gorm.DB
}

// New%sRepository creates a new %s repository
func New%sRepository(db *gorm.DB) %sRepository {
	return &%sRepositoryImpl{db: db}
}

// Create creates a new %s record
func (r *%sRepositoryImpl) Create(data *models.%s) error {
	return r.db.Create(data).Error
}

// GetByID retrieves %s by ID
func (r *%sRepositoryImpl) GetByID(id uint) (*models.%s, error) {
	var data models.%s
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all %s records
func (r *%sRepositoryImpl) GetAll() ([]models.%s, error) {
	var data []models.%s
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates %s record
func (r *%sRepositoryImpl) Update(data *models.%s) error {
	return r.db.Save(data).Error
}

// Delete deletes %s record by ID
func (r *%sRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.%s{}, id).Error
}
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name))

	return writeFile(filepath, content)
}

// createServiceFile generates service file
func createServiceFile(name string) error {
	lowerName := toLowerFirst(name)
	filepath := filepath.Join("src", "modules", "services", fmt.Sprintf("%s_service.go", lowerName))

	content := fmt.Sprintf(`package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// %sService handles business logic for %s
type %sService interface {
	Create(data *models.%s) error
	GetByID(id uint) (*models.%s, error)
	GetAll() ([]models.%s, error)
	Update(data *models.%s) error
	Delete(id uint) error
}

type %sServiceImpl struct {
	repository repositories.%sRepository
}

// New%sService creates a new %s service
func New%sService(repository repositories.%sRepository) %sService {
	return &%sServiceImpl{repository: repository}
}

// Create creates a new %s
func (s *%sServiceImpl) Create(data *models.%s) error {
	return s.repository.Create(data)
}

// GetByID retrieves %s by ID
func (s *%sServiceImpl) GetByID(id uint) (*models.%s, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all %s
func (s *%sServiceImpl) GetAll() ([]models.%s, error) {
	return s.repository.GetAll()
}

// Update updates %s
func (s *%sServiceImpl) Update(data *models.%s) error {
	return s.repository.Update(data)
}

// Delete deletes %s by ID
func (s *%sServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name))

	return writeFile(filepath, content)
}

// createControllerFile generates controller file
func createControllerFile(name string) error {
	lowerName := toLowerFirst(name)
	filepath := filepath.Join("src", "modules", "controllers", fmt.Sprintf("%s_controller.go", lowerName))

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
`, toUpperFirst(name), name, toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), toUpperFirst(name), name, name, name, lowerName, lowerName, toUpperFirst(name), toUpperFirst(name), name, name, lowerName, toUpperFirst(name), name, name, lowerName, toUpperFirst(name), name, name, lowerName, toUpperFirst(name), name, name, name, lowerName, toUpperFirst(name), toUpperFirst(name), name, name, lowerName, toUpperFirst(name), toUpperFirst(name), name, toUpperFirst(name))

	return writeFile(filepath, content)
}

// createDTOFile generates DTO file with Request and Response structs
func createDTOFile(name string) error {
	lowerName := toLowerFirst(name)
	filePath := filepath.Join("src", "dtos", fmt.Sprintf("%s_dto.go", lowerName))

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
