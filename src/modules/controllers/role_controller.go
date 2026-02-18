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
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
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

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		SystemID:    &req.SystemID,
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	if err := c.service.Create(role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to get system data
	roleData, _ := c.service.GetByID(role.ID)

	// Map to response
	var systemResponse *dtos.SystemResponse
	if roleData.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          roleData.System.ID,
			Nama:        roleData.System.Nama,
			Description: roleData.System.Description,
			Status:      roleData.System.Status,
			CreatedAt:   roleData.System.CreatedAt,
			UpdatedAt:   roleData.System.UpdatedAt,
			CreatedByID: roleData.System.CreatedByID,
			UpdatedByID: roleData.System.UpdatedByID,
		}
	}

	response := dtos.RoleResponse{
		ID:          roleData.ID,
		Name:        roleData.Name,
		Description: roleData.Description,
		SystemID:    roleData.SystemID,
		System:      systemResponse,
		Status:      roleData.Status,
		CreatedAt:   roleData.CreatedAt,
		UpdatedAt:   roleData.UpdatedAt,
		CreatedByID: roleData.CreatedByID,
		UpdatedByID: roleData.UpdatedByID,
	}

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

	response := dtos.RoleResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
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

		responseData = append(responseData, dtos.RoleResponse{
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

// Update updates Role
func (c *RoleController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID          uint   `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
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

	// Reload to get system data
	roleData, _ := c.service.GetByID(data.ID)

	// Map to response
	var systemResponse *dtos.SystemResponse
	if roleData.System != nil {
		systemResponse = &dtos.SystemResponse{
			ID:          roleData.System.ID,
			Nama:        roleData.System.Nama,
			Description: roleData.System.Description,
			Status:      roleData.System.Status,
			CreatedAt:   roleData.System.CreatedAt,
			UpdatedAt:   roleData.System.UpdatedAt,
			CreatedByID: roleData.System.CreatedByID,
			UpdatedByID: roleData.System.UpdatedByID,
		}
	}

	response := dtos.RoleResponse{
		ID:          roleData.ID,
		Name:        roleData.Name,
		Description: roleData.Description,
		SystemID:    roleData.SystemID,
		System:      systemResponse,
		Status:      roleData.Status,
		CreatedAt:   roleData.CreatedAt,
		UpdatedAt:   roleData.UpdatedAt,
		CreatedByID: roleData.CreatedByID,
		UpdatedByID: roleData.UpdatedByID,
	}

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
