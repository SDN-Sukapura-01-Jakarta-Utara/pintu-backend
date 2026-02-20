package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
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
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		SystemID      uint   `json:"system_id" binding:"required"`
		Status        string `json:"status"`
		PermissionIDs []uint `json:"permission_ids"`
	}
	
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	createdByID := userID.(uint)

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		SystemID:    &req.SystemID,
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	// Create role with permissions
	if err := c.service.CreateWithPermissions(role, req.PermissionIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to get system data and permission details
	roleData, _ := c.service.GetByID(role.ID)
	_, permissionDetails, _ := c.service.GetRoleWithPermissionDetails(role.ID)

	response := c.mapRoleToResponseWithPermissionDetails(roleData, permissionDetails)

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
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

	// Get role with permission details
	_, permissionDetails, _ := c.service.GetRoleWithPermissionDetails(req.ID)

	response := c.mapRoleToResponseWithPermissionDetails(data, permissionDetails)

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// GetAll retrieves all Roles with filters and pagination
func (c *RoleController) GetAll(ctx *gin.Context) {
	var req dtos.RoleGetAllRequest
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
	roles, total, err := c.service.GetAllWithFilter(repositories.GetRolesParams{
		Filter: repositories.GetRolesFilter{
			Name:     req.Search.Name,
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
	var responseData []dtos.RoleResponse
	for _, role := range roles {
		// Get role with permission details
		_, permissionDetails, _ := c.service.GetRoleWithPermissionDetails(role.ID)
		
		response := c.mapRoleToResponseWithPermissionDetails(&role, permissionDetails)
		responseData = append(responseData, *response)
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

// Update updates Role
func (c *RoleController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID            uint   `json:"id" binding:"required"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		SystemID      *uint  `json:"system_id"`
		Status        string `json:"status"`
		PermissionIDs []uint `json:"permission_ids"`
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

	// Update permissions if permission_ids is provided
	// Always update if the field is sent (even if empty array means clear all permissions)
	if err := c.service.AssignPermissions(req.ID, req.PermissionIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to get system data and permission details
	roleData, _ := c.service.GetByID(data.ID)
	_, permissionDetails, _ := c.service.GetRoleWithPermissionDetails(data.ID)

	response := c.mapRoleToResponseWithPermissionDetails(roleData, permissionDetails)

	ctx.JSON(http.StatusOK, gin.H{"data": response})
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

// Helper function to map Role model to RoleResponse DTO
func mapRoleToResponse(role *models.Role) *dtos.RoleResponse {
	if role == nil {
		return nil
	}

	var systemResponse *dtos.SystemResponse
	if role.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          role.System.ID,
			Nama:        role.System.Nama,
			Description: role.System.Description,
			Status:      role.System.Status,
			CreatedAt:   role.System.CreatedAt,
			UpdatedAt:   role.System.UpdatedAt,
			CreatedByID: role.System.CreatedByID,
			UpdatedByID: role.System.UpdatedByID,
		}
	}

	return &dtos.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		SystemID:    role.SystemID,
		System:      systemResponse,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		CreatedByID: role.CreatedByID,
		UpdatedByID: role.UpdatedByID,
	}
}

// Helper function to map Role model with full Permission details to RoleResponse DTO
func (c *RoleController) mapRoleToResponseWithPermissionDetails(role *models.Role, permissions []models.Permission) *dtos.RoleResponse {
	if role == nil {
		return nil
	}

	var systemResponse *dtos.SystemResponse
	if role.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          role.System.ID,
			Nama:        role.System.Nama,
			Description: role.System.Description,
			Status:      role.System.Status,
			CreatedAt:   role.System.CreatedAt,
			UpdatedAt:   role.System.UpdatedAt,
			CreatedByID: role.System.CreatedByID,
			UpdatedByID: role.System.UpdatedByID,
		}
	}

	// Map permissions to permission data
	var permissionData []dtos.PermissionData
	if len(permissions) > 0 {
		permissionData = make([]dtos.PermissionData, len(permissions))
		for i, perm := range permissions {
			permissionData[i] = dtos.PermissionData{
				ID:          perm.ID,
				Name:        perm.Name,
				Description: perm.Description,
				GroupName:   perm.GroupName,
				System:      "", // Will be populated if System data exists
			}
			if perm.System != nil {
				permissionData[i].System = perm.System.Nama
			}
		}
	}

	return &dtos.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		SystemID:    role.SystemID,
		System:      systemResponse,
		Status:      role.Status,
		Permissions: permissionData,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		CreatedByID: role.CreatedByID,
		UpdatedByID: role.UpdatedByID,
	}
}
