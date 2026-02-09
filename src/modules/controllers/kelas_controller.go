package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// KelasController handles HTTP requests for Kelas
type KelasController struct {
	service services.KelasService
}

// NewKelasController creates a new Kelas controller
func NewKelasController(service services.KelasService) *KelasController {
	return &KelasController{service: service}
}

// Create creates a new Kelas
// @Summary Create new Kelas
// @Description Create a new Kelas with name and status
// @Tags kelas
// @Accept json
// @Produce json
// @Param body body dtos.KelasCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.KelasResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/kelas/create-kelas [post]
func (c *KelasController) Create(ctx *gin.Context) {
	var req dtos.KelasCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Create(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all kelas
// @Summary Get all Kelas
// @Description Retrieve all Kelas records
// @Tags kelas
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.KelasResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/kelas/get-kelas [post]
func (c *KelasController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll(10, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data})
}

// GetByID retrieves Kelas by ID
// @Summary Get Kelas by ID
// @Description Retrieve kelas details by ID
// @Tags kelas
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.KelasResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/kelas/get-kelas-by-id [post]
func (c *KelasController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Kelas not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Kelas
// @Summary Update Kelas
// @Description Update Kelas details
// @Tags kelas
// @Accept json
// @Produce json
// @Param body body dtos.KelasUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.KelasResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/kelas/update-kelas [post]
func (c *KelasController) Update(ctx *gin.Context) {
	var req dtos.KelasUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Kelas by ID
// @Summary Delete Kelas
// @Description Delete Kelas by ID
// @Tags kelas
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/kelas/delete-kelas [post]
func (c *KelasController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Kelas not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kelas deleted successfully",
	})
}
