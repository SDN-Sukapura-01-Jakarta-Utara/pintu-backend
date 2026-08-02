package controllers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// FormulirResponseController handles HTTP requests for form submissions
type FormulirResponseController struct {
	service services.FormulirResponseService
}

// NewFormulirResponseController creates a new controller
func NewFormulirResponseController(service services.FormulirResponseService) *FormulirResponseController {
	return &FormulirResponseController{service: service}
}

// SubmitAuthenticated handles authenticated form submission with file uploads
// @Summary Submit form response (authenticated)
// @Description Submit answers to an authenticated form with user ID tracking and file uploads
// @Tags formulir
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON string of FormulirSubmitRequest"
// @Param file_{question_id} formData file false "File upload for question (e.g., file_9, file_10)"
// @Success 201 {object} gin.H{data=dtos.FormulirSubmitResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/submit-response-authenticated [post]
func (c *FormulirResponseController) SubmitAuthenticated(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil { // 50MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	// Get JSON data from form field
	dataJSON := ctx.PostForm("data")
	if dataJSON == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "data field is required"})
		return
	}

	// Parse JSON request
	var req dtos.FormulirSubmitRequest
	if err := json.Unmarshal([]byte(dataJSON), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format in data field"})
		return
	}

	// Extract file uploads - files are named as file_{pertanyaan_id}
	files := make(map[uint][]*multipart.FileHeader)
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Parse field name like "file_9", "file_10"
			var pertanyaanID uint
			if _, err := fmt.Sscanf(fieldName, "file_%d", &pertanyaanID); err == nil {
				files[pertanyaanID] = fileHeaders
			}
		}
	}

	// Get IP address and user agent
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// Call service
	data, err := c.service.SubmitAuthenticated(&req, userIDUint, ipAddress, userAgent, files)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// SubmitPublic handles public form submission (no auth required) with file uploads
// @Summary Submit form response (public)
// @Description Submit answers to a public form without authentication and file uploads
// @Tags formulir-public
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON string of FormulirSubmitRequest"
// @Param file_{question_id} formData file false "File upload for question (e.g., file_9, file_10)"
// @Success 201 {object} gin.H{data=dtos.FormulirSubmitResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/submit-response-public [post]
func (c *FormulirResponseController) SubmitPublic(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil { // 50MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	// Get JSON data from form field
	dataJSON := ctx.PostForm("data")
	if dataJSON == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "data field is required"})
		return
	}

	// Parse JSON request
	var req dtos.FormulirSubmitRequest
	if err := json.Unmarshal([]byte(dataJSON), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format in data field"})
		return
	}

	// Extract file uploads - files are named as file_{pertanyaan_id}
	files := make(map[uint][]*multipart.FileHeader)
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Parse field name like "file_9", "file_10"
			var pertanyaanID uint
			if _, err := fmt.Sscanf(fieldName, "file_%d", &pertanyaanID); err == nil {
				files[pertanyaanID] = fileHeaders
			}
		}
	}

	// Get IP address and user agent
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// Call service
	data, err := c.service.SubmitPublic(&req, ipAddress, userAgent, files)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetResponsesBySlug gets all responses for a form (Google Forms style view)
// @Summary Get form responses by slug
// @Description Get all responses for a form with question details (like Google Forms response view)
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirResponseDetailRequest true "Slug request with optional role"
// @Success 200 {object} gin.H{data=dtos.FormulirResponsesDetailResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-response-by-slug [post]
func (c *FormulirResponseController) GetResponsesBySlug(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirResponseDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service with userID and role
	data, err := c.service.GetResponsesBySlug(req.Slug, req.Role, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// EditResponse handles editing existing form response with file uploads
// @Summary Edit form response
// @Description Edit an existing form response with updated answers and file uploads
// @Tags formulir
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON string of FormulirEditResponseRequest"
// @Param file_{question_id} formData file false "File upload for question (e.g., file_9, file_10)"
// @Success 200 {object} gin.H{data=dtos.FormulirSubmitResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/edit-response [post]
func (c *FormulirResponseController) EditResponse(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil { // 50MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	// Get JSON data from form field
	dataJSON := ctx.PostForm("data")
	if dataJSON == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "data field is required"})
		return
	}

	// Parse JSON request
	var req dtos.FormulirEditResponseRequest
	if err := json.Unmarshal([]byte(dataJSON), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format in data field"})
		return
	}

	// Extract file uploads - files are named as file_{pertanyaan_id}
	files := make(map[uint][]*multipart.FileHeader)
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Parse field name like "file_9", "file_10"
			var pertanyaanID uint
			if _, err := fmt.Sscanf(fieldName, "file_%d", &pertanyaanID); err == nil {
				files[pertanyaanID] = fileHeaders
			}
		}
	}

	// Get IP address and user agent
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// Call service
	data, err := c.service.EditResponse(&req, userIDUint, ipAddress, userAgent, files)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetResponseByUser gets a specific user's response for a form
// @Summary Get user's response by slug
// @Description Get a specific user's response for a form (for viewing own response)
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirResponseByUserRequest true "Request with slug and optional role"
// @Success 200 {object} gin.H{data=dtos.FormulirResponsesDetailResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-response-by-user [post]
func (c *FormulirResponseController) GetResponseByUser(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirResponseByUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service with userID from JWT token
	data, err := c.service.GetResponseByUser(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// DeleteResponse deletes a response and all its answers
// @Summary Delete response by ID
// @Description Delete a response and all its answers (admin or owner only)
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirDeleteResponseRequest true "Request with response_id"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/delete-response-by-id [post]
func (c *FormulirResponseController) DeleteResponse(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirDeleteResponseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := c.service.DeleteResponse(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Response berhasil dihapus"})
}

// ResetResponses deletes all responses for a formulir
// @Summary Reset all responses for a form
// @Description Delete all responses and answers for a formulir (form owner or admin only)
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirResetResponseRequest true "Request with formulir_id"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/reset-response [post]
func (c *FormulirResponseController) ResetResponses(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirResetResponseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := c.service.ResetResponses(req.FormulirID, req.Role, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Semua response berhasil dihapus"})
}

// GetStatisticsBySlug gets statistics for all questions in a form
// @Summary Get form statistics by slug
// @Description Get aggregated statistics for all questions (for charts and analytics)
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirStatisticRequest true "Request with slug and optional role"
// @Success 200 {object} gin.H{data=dtos.FormulirStatisticResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-statistic-by-slug [post]
func (c *FormulirResponseController) GetStatisticsBySlug(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirStatisticRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service with userID and role
	data, err := c.service.GetStatisticsBySlug(req.Slug, req.Role, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
