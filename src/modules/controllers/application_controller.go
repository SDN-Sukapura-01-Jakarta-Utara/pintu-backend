package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// ApplicationController handles HTTP requests for Application
type ApplicationController struct {
	service services.ApplicationService
}

// NewApplicationController creates a new Application controller
func NewApplicationController(service services.ApplicationService) *ApplicationController {
	return &ApplicationController{service: service}
}

// Create creates a new Application
// @Summary Create new Application
// @Description Create a new Application with nama, link, show_in_jumbotron, and status
// @Tags application
// @Accept json
// @Produce json
// @Param body body dtos.ApplicationCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.ApplicationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/application/create-application [post]
func (c *ApplicationController) Create(ctx *gin.Context) {
	var req dtos.ApplicationCreateRequest
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

// GetAll retrieves all applications
// @Summary Get all Applications
// @Description Retrieve all Application records
// @Tags application
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=[]dtos.ApplicationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/application/get-application [post]
func (c *ApplicationController) GetAll(ctx *gin.Context) {
	var req dtos.ApplicationGetAllRequest
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

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetApplicationParams{
		Filter: repositories.GetApplicationFilter{
			Nama:            req.Search.Nama,
			Status:          req.Search.Status,
			ShowInJumbotron: req.Search.ShowInJumbotron,
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data.Data,
		"pagination": gin.H{
			"limit":       data.Pagination.Limit,
			"offset":      data.Pagination.Offset,
			"page":        data.Pagination.Page,
			"total":       data.Pagination.Total,
			"total_pages": data.Pagination.TotalPages,
		},
	})
}

// GetByID retrieves Application by ID
// @Summary Get Application by ID
// @Description Retrieve application details by ID
// @Tags application
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.ApplicationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/application/get-application-by-id [post]
func (c *ApplicationController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Application
// @Summary Update Application
// @Description Update Application details
// @Tags application
// @Accept json
// @Produce json
// @Param body body dtos.ApplicationUpdateRequest true "Request body"
// @Success 200 {object} gin.H{data=dtos.ApplicationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/application/update-application [post]
func (c *ApplicationController) Update(ctx *gin.Context) {
	var req dtos.ApplicationUpdateRequest
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

// Delete deletes Application by ID
// @Summary Delete Application
// @Description Delete Application by ID
// @Tags application
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/application/delete-application [post]
func (c *ApplicationController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Application deleted successfully",
	})
}

// GetPublicList retrieves all active applications for public display (no auth required)
// @Summary Get all active applications for public
// @Description Retrieve all active applications (no authentication required, sorted from oldest to newest)
// @Tags application
// @Accept json
// @Produce json
// @Param body body dtos.ApplicationPublicRequest true "Request body with filter"
// @Success 200 {object} dtos.ApplicationPublicListResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/public/get-data-aplikasi-sekolah [post]
func (c *ApplicationController) GetPublicList(ctx *gin.Context) {
	var req dtos.ApplicationPublicRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetPublicList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
