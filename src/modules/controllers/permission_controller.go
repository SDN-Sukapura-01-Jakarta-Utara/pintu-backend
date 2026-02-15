package controllers

import (
	"net/http"
	"strconv"

	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// PermissionController handles HTTP requests for Permission
type PermissionController struct {
	service services.PermissionService
}

// NewPermissionController creates a new Permission controller
func NewPermissionController(service services.PermissionService) *PermissionController {
	return &PermissionController{service: service}
}

// Create creates a new Permission
func (c *PermissionController) Create(ctx *gin.Context) {
	type CreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		GroupName   string `json:"group_name"`
		SystemID    uint   `json:"system_id" binding:"required"`
		Status      string `json:"status"`
	}
	
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	createdByID := userID.(uint)

	permission := &models.Permission{
		Name:        req.Name,
		Description: req.Description,
		GroupName:   req.GroupName,
		SystemID:    &req.SystemID,
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	if err := c.service.Create(permission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": permission})
}

// GetByID retrieves Permission by ID
func (c *PermissionController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all Permissions
func (c *PermissionController) GetAll(ctx *gin.Context) {
	limit := 10
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	data, total, err := c.service.GetAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetByGroupName retrieves permissions by group name
func (c *PermissionController) GetByGroupName(ctx *gin.Context) {
	var req struct {
		GroupName string `json:"group_name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByGroupName(req.GroupName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetBySystem retrieves permissions by system
func (c *PermissionController) GetBySystem(ctx *gin.Context) {
	var req struct {
		System string `json:"system" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetBySystem(req.System)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Permission
func (c *PermissionController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID          uint   `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
		GroupName   string `json:"group_name"`
		SystemID    *uint  `json:"system_id"`
		Status      string `json:"status"`
	}
	
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Description != "" {
		data.Description = req.Description
	}
	if req.GroupName != "" {
		data.GroupName = req.GroupName
	}
	if req.SystemID != nil && *req.SystemID > 0 {
		data.SystemID = req.SystemID
	}
	if req.Status != "" {
		data.Status = req.Status
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	updatedByID := userID.(uint)
	data.UpdatedByID = &updatedByID

	if err := c.service.Update(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Permission by ID
func (c *PermissionController) Delete(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Permission deleted successfully"})
}
