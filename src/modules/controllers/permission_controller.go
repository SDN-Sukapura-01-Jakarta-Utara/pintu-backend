package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
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

	// Reload to get system data
	permissionData, _ := c.service.GetByID(permission.ID)

	// Map to response
	var systemResponse *dtos.SystemResponse
	if permissionData.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          permissionData.System.ID,
			Nama:        permissionData.System.Nama,
			Description: permissionData.System.Description,
			Status:      permissionData.System.Status,
			CreatedAt:   permissionData.System.CreatedAt,
			UpdatedAt:   permissionData.System.UpdatedAt,
			CreatedByID: permissionData.System.CreatedByID,
			UpdatedByID: permissionData.System.UpdatedByID,
		}
	}

	response := dtos.PermissionResponse{
		ID:          permissionData.ID,
		Name:        permissionData.Name,
		Description: permissionData.Description,
		GroupName:   permissionData.GroupName,
		SystemID:    permissionData.SystemID,
		System:      systemResponse,
		Status:      permissionData.Status,
		CreatedAt:   permissionData.CreatedAt,
		UpdatedAt:   permissionData.UpdatedAt,
		CreatedByID: permissionData.CreatedByID,
		UpdatedByID: permissionData.UpdatedByID,
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
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

	// Map to response
	var systemResponse *dtos.SystemResponse
	if data.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          data.System.ID,
			Nama:        data.System.Nama,
			Description: data.System.Description,
			Status:      data.System.Status,
			CreatedAt:   data.System.CreatedAt,
			UpdatedAt:   data.System.UpdatedAt,
			CreatedByID: data.System.CreatedByID,
			UpdatedByID: data.System.UpdatedByID,
		}
	}

	response := dtos.PermissionResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		GroupName:   data.GroupName,
		SystemID:    data.SystemID,
		System:      systemResponse,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// GetAll retrieves all Permissions with filters and pagination
func (c *PermissionController) GetAll(ctx *gin.Context) {
	var req dtos.PermissionGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Call service
	permissions, total, err := c.service.GetAllWithFilter(repositories.GetPermissionsParams{
		Filter: repositories.GetPermissionsFilter{
			Name:     req.Search.Name,
			GroupName: req.Search.GroupName,
			SystemID: req.Search.SystemID,
			Status:   req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to response
	var responseData []dtos.PermissionResponse
	for _, permission := range permissions {
		var systemResponse *dtos.SystemResponse
		if permission.System != nil {
			systemResponse = &dtos.SystemResponse{
				ID:          permission.System.ID,
				Nama:        permission.System.Nama,
				Description: permission.System.Description,
				Status:      permission.System.Status,
				CreatedAt:   permission.System.CreatedAt,
				UpdatedAt:   permission.System.UpdatedAt,
				CreatedByID: permission.System.CreatedByID,
				UpdatedByID: permission.System.UpdatedByID,
			}
		}

		responseData = append(responseData, dtos.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			GroupName:   permission.GroupName,
			SystemID:    permission.SystemID,
			System:      systemResponse,
			Status:      permission.Status,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
			CreatedByID: permission.CreatedByID,
			UpdatedByID: permission.UpdatedByID,
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	ctx.JSON(http.StatusOK, gin.H{
		"data": responseData,
		"pagination": gin.H{
			"limit":       limit,
			"offset":      offset,
			"page":        page,
			"total":       total,
			"total_pages": totalPages,
		},
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
		SystemID uint `json:"system_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetBySystem(req.SystemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to response
	var responseData []dtos.PermissionResponse
	for _, permission := range data {
		var systemResponse *dtos.SystemResponse
		if permission.System != nil {
			systemResponse = &dtos.SystemResponse{
				ID:          permission.System.ID,
				Nama:        permission.System.Nama,
				Description: permission.System.Description,
				Status:      permission.System.Status,
				CreatedAt:   permission.System.CreatedAt,
				UpdatedAt:   permission.System.UpdatedAt,
				CreatedByID: permission.System.CreatedByID,
				UpdatedByID: permission.System.UpdatedByID,
			}
		}

		responseData = append(responseData, dtos.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			GroupName:   permission.GroupName,
			SystemID:    permission.SystemID,
			System:      systemResponse,
			Status:      permission.Status,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
			CreatedByID: permission.CreatedByID,
			UpdatedByID: permission.UpdatedByID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
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

	// Reload to get system data
	permissionData, _ := c.service.GetByID(data.ID)

	// Map to response
	var systemResponse *dtos.SystemResponse
	if permissionData.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          permissionData.System.ID,
			Nama:        permissionData.System.Nama,
			Description: permissionData.System.Description,
			Status:      permissionData.System.Status,
			CreatedAt:   permissionData.System.CreatedAt,
			UpdatedAt:   permissionData.System.UpdatedAt,
			CreatedByID: permissionData.System.CreatedByID,
			UpdatedByID: permissionData.System.UpdatedByID,
		}
	}

	response := dtos.PermissionResponse{
		ID:          permissionData.ID,
		Name:        permissionData.Name,
		Description: permissionData.Description,
		GroupName:   permissionData.GroupName,
		SystemID:    permissionData.SystemID,
		System:      systemResponse,
		Status:      permissionData.Status,
		CreatedAt:   permissionData.CreatedAt,
		UpdatedAt:   permissionData.UpdatedAt,
		CreatedByID: permissionData.CreatedByID,
		UpdatedByID: permissionData.UpdatedByID,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
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

// Helper function to map Permission model to PermissionResponse DTO
func mapPermissionToResponse(permission *models.Permission) *dtos.PermissionResponse {
	if permission == nil {
		return nil
	}

	var systemResponse *dtos.SystemResponse
	if permission.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          permission.System.ID,
			Nama:        permission.System.Nama,
			Description: permission.System.Description,
			Status:      permission.System.Status,
			CreatedAt:   permission.System.CreatedAt,
			UpdatedAt:   permission.System.UpdatedAt,
			CreatedByID: permission.System.CreatedByID,
			UpdatedByID: permission.System.UpdatedByID,
		}
	}

	return &dtos.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		GroupName:   permission.GroupName,
		SystemID:    permission.SystemID,
		System:      systemResponse,
		Status:      permission.Status,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
		CreatedByID: permission.CreatedByID,
		UpdatedByID: permission.UpdatedByID,
	}
}
