package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// SieksaRoleController handles HTTP requests for Role in SIEKSA context
type SieksaRoleController struct {
	service services.RoleService
}

// NewSieksaRoleController creates a new SIEKSA Role controller
func NewSieksaRoleController(service services.RoleService) *SieksaRoleController {
	return &SieksaRoleController{service: service}
}

// GetAll retrieves all Roles with filters and pagination (filtered for SIEKSA - SystemID=2)
func (c *SieksaRoleController) GetAll(ctx *gin.Context) {
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

	// Force SystemID to 2 (SIEKSA)
	sieksaSystemID := uint(2)

	// Call service
	roles, total, err := c.service.GetAllWithFilter(repositories.GetRolesParams{
		Filter: repositories.GetRolesFilter{
			Name:     req.Search.Name,
			SystemID: sieksaSystemID, // Always filter by SIEKSA
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

// Helper function to map Role model with full Permission details to RoleResponse DTO
func (c *SieksaRoleController) mapRoleToResponseWithPermissionDetails(role *models.Role, permissions []models.Permission) *dtos.RoleResponse {
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
				Status:      perm.Status,
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
