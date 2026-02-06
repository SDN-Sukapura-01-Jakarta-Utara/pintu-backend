package controllers

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Summary Create new User
// @Description Create a new User
// @Tags user
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Router /user [post]
func (c *UserController) Create(ctx *gin.Context) {
	var req models.User
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

// GetByID retrieves User by ID
// @Summary Get User by ID
// @Description Retrieve user details by ID
// @Tags User
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /User/{id} [get]
func (c *UserController) GetByID(ctx *gin.Context) {
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

// GetAll retrieves all user
// @Summary Get all User
// @Description Retrieve all User records
// @Tags User
// @Produce json
// @Success 200
// @Failure 500
// @Router /user [get]
func (c *UserController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates User
// @Summary Update User
// @Description Update User details
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /User/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.User
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

// Delete deletes User by ID
// @Summary Delete user
// @Description Delete User by ID
// @Tags User
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /User/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
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
