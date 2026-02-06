package controllers

import (
	"net/http"

	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// RoleController handles HTTP requests for Role
type RoleController struct {
	service services.RoleService
}

// NewRoleController creates a new Role controller
func NewRoleController(service services.RoleService) *RoleController {
	return &RoleController{service: service}
}

// Create creates a new Role
func (c *RoleController) Create(ctx *gin.Context) {
	type CreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		System      string `json:"system"`
	}
	
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		System:      req.System,
	}

	if err := c.service.Create(role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": role})
}

// GetByID retrieves Role by ID
func (c *RoleController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all Roles
func (c *RoleController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Role
func (c *RoleController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID          uint   `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
		System      string `json:"system"`
	}
	
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Description != "" {
		data.Description = req.Description
	}
	if req.System != "" {
		data.System = req.System
	}

	if err := c.service.Update(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Role by ID
func (c *RoleController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
