package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles HTTP requests for User
type UserController struct {
	service services.UserService
}

// NewUserController creates a new User controller
func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// Create creates a new User
func (c *UserController) Create(ctx *gin.Context) {
	var req dtos.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	createdByID := userID.(uint)

	user := &models.User{
		Nama:        req.Nama,
		Username:    req.Username,
		Password:    string(hashedPassword),
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	if err := c.service.Create(user, req.RoleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload user to get roles
	user, _ = c.service.GetByID(user.ID)

	// Map to response
	response := mapUserToResponse(user)
	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetByID retrieves User by ID
func (c *UserController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := mapUserToResponse(data)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// GetAll retrieves all Users with filters and pagination
func (c *UserController) GetAll(ctx *gin.Context) {
	var req dtos.UserGetAllRequest
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
	users, total, err := c.service.GetAllWithFilter(repositories.GetUsersParams{
		Filter: repositories.GetUsersFilter{
			Nama:     req.Search.Nama,
			Username: req.Search.Username,
			RoleIDs:  req.Search.RoleIDs,
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
	var responseData []dtos.UserResponseDetail
	for _, user := range users {
		responseData = append(responseData, *mapUserToResponseDetail(&user))
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

// Update updates User
func (c *UserController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID       uint   `json:"id" binding:"required"`
		Nama     string `json:"nama"`
		Username string `json:"username"`
		Password string `json:"password"`
		RoleIDs  []uint `json:"role_ids"`
		Status   string `json:"status"`
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Nama != "" {
		data.Nama = req.Nama
	}
	if req.Username != "" {
		data.Username = req.Username
	}
	if req.Password != "" {
		// Hash password if provided
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		data.Password = string(hashedPassword)
	}
	if req.Status != "" {
		data.Status = req.Status
	}

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	updatedByID := userID.(uint)
	data.UpdatedByID = &updatedByID

	if err := c.service.Update(data, req.RoleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload user to get updated roles
	data, _ = c.service.GetByID(req.ID)
	response := mapUserToResponse(data)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdatePassword updates user password
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	type PasswordRequest struct {
		ID          uint   `json:"id" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	
	var req PasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	// Get user ID from JWT token
	userID, _ := ctx.Get("userID")
	updatedByID := userID.(uint)
	user.UpdatedByID = &updatedByID

	// Keep existing roles when updating password
	if err := c.service.Update(user, []uint{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Delete deletes User by ID
func (c *UserController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Helper function to map User model to UserResponse DTO
func mapUserToResponse(user *models.User) *dtos.UserResponse {
	if user == nil {
		return nil
	}

	roles := make([]dtos.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			SystemID:    role.SystemID,
			System:      system,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			CreatedByID: role.CreatedByID,
			UpdatedByID: role.UpdatedByID,
		}
	}

	return &dtos.UserResponse{
		ID:          user.ID,
		Nama:        user.Nama,
		Username:    user.Username,
		Roles:       roles,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		CreatedByID: user.CreatedByID,
		UpdatedByID: user.UpdatedByID,
	}
}

// Helper function to map User model to UserResponseDetail DTO
func mapUserToResponseDetail(user *models.User) *dtos.UserResponseDetail {
	if user == nil {
		return nil
	}

	roles := make([]dtos.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			SystemID:    role.SystemID,
			System:      system,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			CreatedByID: role.CreatedByID,
			UpdatedByID: role.UpdatedByID,
		}
	}

	return &dtos.UserResponseDetail{
		ID:          user.ID,
		Nama:        user.Nama,
		Username:    user.Username,
		Roles:       roles,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		CreatedByID: user.CreatedByID,
		UpdatedByID: user.UpdatedByID,
	}
}
