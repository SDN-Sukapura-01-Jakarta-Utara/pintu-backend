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
		RoleID:      &req.RoleID,
		Status:      req.Status,
		CreatedByID: &createdByID,
	}

	// Set accessible systems
	if len(req.AccessibleSystem) > 0 {
		user.SetAccessibleSystems(req.AccessibleSystem)
	}

	if err := c.service.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
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

	ctx.JSON(http.StatusOK, gin.H{"data": data})
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
			Nama:             req.Search.Nama,
			Username:         req.Search.Username,
			RoleID:           req.Search.RoleID,
			Status:           req.Search.Status,
			AccessibleSystem: req.Search.AccessibleSystem,
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
		systems, _ := user.AccessibleSystems()
		roleName := ""
		if user.Role != nil {
			roleName = user.Role.Name
		}

		responseData = append(responseData, dtos.UserResponseDetail{
			ID:               user.ID,
			Nama:             user.Nama,
			Username:         user.Username,
			RoleID:           *user.RoleID,
			RoleName:         roleName,
			AccessibleSystem: systems,
			Status:           user.Status,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			CreatedByID:      user.CreatedByID,
			UpdatedByID:      user.UpdatedByID,
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

// Update updates User
func (c *UserController) Update(ctx *gin.Context) {
	type UpdateRequest struct {
		ID               uint     `json:"id" binding:"required"`
		Nama             string   `json:"nama"`
		Username         string   `json:"username"`
		RoleID           *uint    `json:"role_id"`
		AccessibleSystem []string `json:"accessible_system"`
		Status           string   `json:"status"`
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
	if req.RoleID != nil && *req.RoleID > 0 {
		data.RoleID = req.RoleID
	}
	if len(req.AccessibleSystem) > 0 {
		data.SetAccessibleSystems(req.AccessibleSystem)
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

	if err := c.service.Update(user); err != nil {
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
