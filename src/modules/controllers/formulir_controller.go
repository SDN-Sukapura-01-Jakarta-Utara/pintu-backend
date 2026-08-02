package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// FormulirController handles HTTP requests for Formulir
type FormulirController struct {
	service services.FormulirService
}

// NewFormulirController creates a new Formulir controller
func NewFormulirController(service services.FormulirService) *FormulirController {
	return &FormulirController{service: service}
}

// Create creates a new Formulir with pertanyaan and optional dokumen uploads
// @Summary Create new Formulir
// @Description Create a new Formulir with questions/fields and optional document uploads
// @Tags formulir
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON string of formulir data"
// @Param dokumen_1 formData file false "Document for question with urutan 1"
// @Param dokumen_2 formData file false "Document for question with urutan 2"
// @Success 201 {object} gin.H{data=dtos.FormulirResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/create-formulir [post]
func (c *FormulirController) Create(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil { // 50 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.FormulirCreateRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Parse uploaded dokumen files
	// Files are expected with field name pattern: dokumen_{urutan}
	// Example: dokumen_1, dokumen_2, dokumen_3 (matching urutan in pertanyaan array)
	dokumenMap := make(map[int]*multipart.FileHeader)
	
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Check if field name starts with "dokumen_"
			if len(fieldName) > 8 && fieldName[:8] == "dokumen_" {
				// Extract urutan from field name
				urutanStr := fieldName[8:]
				urutan, err := strconv.Atoi(urutanStr)
				if err != nil {
					continue // Skip invalid field names
				}
				
				// Map urutan to file (take first file if multiple uploaded)
				if len(fileHeaders) > 0 {
					dokumenMap[urutan] = fileHeaders[0]
				}
			}
		}
	}

	// Call service
	data, err := c.service.Create(&req, dokumenMap, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}


// Update updates Formulir with pertanyaan and optional dokumen uploads
// @Summary Update Formulir
// @Description Update Formulir with questions/fields and optional document uploads
// @Tags formulir
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON string of formulir update data"
// @Param dokumen_1 formData file false "New document for question with urutan 1"
// @Param dokumen_2 formData file false "New document for question with urutan 2"
// @Success 200 {object} gin.H{data=dtos.FormulirResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/formulir/edit-formulir [post]
func (c *FormulirController) Update(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil { // 50 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gagal parse form data"})
		return
	}

	// Get JSON data from form field
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "field 'data' wajib diisi"})
		return
	}

	// Parse JSON data
	var req dtos.FormulirUpdateRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Parse uploaded dokumen files
	dokumenMap := make(map[int]*multipart.FileHeader)
	
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for fieldName, fileHeaders := range form.File {
			// Check if field name starts with "dokumen_"
			if len(fieldName) > 8 && fieldName[:8] == "dokumen_" {
				// Extract urutan from field name
				urutanStr := fieldName[8:]
				urutan, err := strconv.Atoi(urutanStr)
				if err != nil {
					continue // Skip invalid field names
				}
				
				// Map urutan to file (take first file if multiple uploaded)
				if len(fileHeaders) > 0 {
					dokumenMap[urutan] = fileHeaders[0]
				}
			}
		}
	}

	// Call service
	data, err := c.service.Update(req.ID, &req, dokumenMap, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all formulir with pagination and filters
// @Summary Get all formulir
// @Description Get all formulir with pagination and filters
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirGetAllRequest true "Filter parameters"
// @Success 200 {object} gin.H{data=dtos.FormulirListWithPaginationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-formulir [post]
func (c *FormulirController) GetAll(ctx *gin.Context) {
	var req dtos.FormulirGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Validate pagination
	if req.Pagination.Page < 1 {
		req.Pagination.Page = 1
	}
	if req.Pagination.Limit < 1 {
		req.Pagination.Limit = 10
	}
	if req.Pagination.Limit > 100 {
		req.Pagination.Limit = 100
	}

	// Calculate offset from page
	offset := (req.Pagination.Page - 1) * req.Pagination.Limit

	// Parse date filters
	var startDate, endDate time.Time
	var err error

	if req.Search.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.Search.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (use YYYY-MM-DD)"})
			return
		}
	}

	if req.Search.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.Search.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (use YYYY-MM-DD)"})
			return
		}
		// Set to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Handle access_type logic:
	// - null (not provided in JSON) → show all forms (both public and authenticated)
	// - "" (empty string) → show all forms (both public and authenticated)
	// - "public" or "authenticated" → filter by that specific type
	var accessTypeFilter *string
	if req.Search.AccessType != nil {
		if *req.Search.AccessType == "" {
			// Empty string = show all, set to nil
			accessTypeFilter = nil
		} else {
			// Use the provided value ("public" or "authenticated")
			accessTypeFilter = req.Search.AccessType
		}
	}
	// If nil (not provided), accessTypeFilter stays nil = show all

	// Build params
	params := repositories.GetFormulirParams{
		Filter: repositories.GetFormulirFilter{
			Judul:           req.Search.Judul,
			CreatedByID:     req.Search.CreatedByID,
			CreatedByRole:   req.Search.Role,
			StartDate:       startDate,
			EndDate:         endDate,
			AccessType:      accessTypeFilter,
			TargetUserTypes: req.Search.TargetUserTypes,
			RombelID:        req.Search.RombelID,
		},
		Limit:  req.Pagination.Limit,
		Offset: offset,
	}

	// Call service with userID for role-based filtering
	data, err := c.service.GetAllWithFilter(params, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetByID retrieves Formulir by ID
// @Summary Get Formulir by ID
// @Description Retrieve formulir details by ID with all pertanyaan
// @Tags formulir
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.FormulirResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-formulir-by-id [post]
func (c *FormulirController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Formulir not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetBySlug retrieves Formulir by slug
// @Summary Get Formulir by slug
// @Description Retrieve formulir details by slug with all pertanyaan
// @Tags formulir
// @Accept json
// @Produce json
// @Param body body dtos.FormulirGetBySlugRequest true "Request body with slug"
// @Success 200 {object} gin.H{data=dtos.FormulirResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-formulir-by-slug [post]
func (c *FormulirController) GetBySlug(ctx *gin.Context) {
	var req dtos.FormulirGetBySlugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	data, err := c.service.GetBySlug(req.Slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Formulir not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetBySlugPublic retrieves Formulir by slug (public, no auth)
// @Summary Get Formulir by slug (public)
// @Description Retrieve public formulir details by slug with all pertanyaan (no authentication required)
// @Tags formulir-public
// @Accept json
// @Produce json
// @Param body body dtos.FormulirGetBySlugRequest true "Request body with slug"
// @Success 200 {object} gin.H{data=dtos.FormulirResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/public/get-form-public [post]
func (c *FormulirController) GetBySlugPublic(ctx *gin.Context) {
	var req dtos.FormulirGetBySlugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	data, err := c.service.GetBySlug(req.Slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Formulir not found"})
		return
	}

	// Optional: Check if form is active (public users should only see active forms)
	if !data.IsActive {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Formulir not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes a formulir with cascade delete
// @Summary Delete Formulir
// @Description Delete formulir with all related data (responses, questions) and files
// @Tags formulir
// @Accept json
// @Produce json
// @Param body body dtos.FormulirDeleteRequest true "Request body with formulir_id and optional role"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/delete-formulir [post]
func (c *FormulirController) Delete(ctx *gin.Context) {
	// Get user ID from context (set by middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := c.service.Delete(req.FormulirID, req.Role, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Formulir berhasil dihapus"})
}

// GetFormulirByUser retrieves formulir filtered by user role and context
// @Summary Get formulir by user role
// @Description Get formulir filtered by user role (pendidik, tendik, murid) with pagination
// @Tags formulir
// @Accept json
// @Produce json
// @Param request body dtos.FormulirGetByUserRequest true "Filter parameters by user role"
// @Success 200 {object} gin.H{data=dtos.FormulirListWithPaginationResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/formulir/get-formulir-by-user [post]
func (c *FormulirController) GetFormulirByUser(ctx *gin.Context) {
	// Get user ID from context (set by middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak terautentikasi"})
		return
	}
	userIDUint := userID.(uint)

	var req dtos.FormulirGetByUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	data, err := c.service.GetFormulirByUser(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
