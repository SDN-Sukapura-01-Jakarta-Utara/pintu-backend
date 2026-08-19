package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// PelatihController handles HTTP requests for Pelatih
type PelatihController struct {
	service services.PelatihService
}

// NewPelatihController creates a new controller
func NewPelatihController(service services.PelatihService) *PelatihController {
	return &PelatihController{service: service}
}

// Create creates a new Pelatih and assigns to ekstrakurikuler
// @Summary Create new Pelatih
// @Description Create a new pelatih/coach and assign to ekstrakurikuler
// @Tags pelatih
// @Accept json
// @Produce json
// @Param body body dtos.PelatihCreateRequest true "Request body"
// @Success 201 {object} gin.H{data=dtos.PelatihResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/sieksa/ekstrakurikuler/create-pelatih [post]
func (c *PelatihController) Create(ctx *gin.Context) {
	var req dtos.PelatihCreateRequest
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

// GetByID retrieves Pelatih by ID
// @Summary Get Pelatih by ID
// @Description Retrieve pelatih details by ID with ekstrakurikuler relations
// @Tags pelatih
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PelatihResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/sieksa/ekstrakurikuler/get-pelatih-by-id [post]
func (c *PelatihController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Pelatih
// @Summary Update Pelatih
// @Description Update pelatih details and ekstrakurikuler assignments with file uploads (form-data)
// @Tags pelatih
// @Accept multipart/form-data
// @Produce json
// @Param id formData int true "Pelatih ID"
// @Param nama formData string false "Nama"
// @Param username formData string false "Username"
// @Param password formData string false "Password"
// @Param telepon formData string false "Telepon"
// @Param alamat formData string false "Alamat"
// @Param keahlian formData string false "Keahlian"
// @Param status formData string false "Status"
// @Param foto_profil formData file false "Foto Profil (single file)"
// @Param sertifikat formData file false "Sertifikat (multiple files with name 'sertifikat')"
// @Param ekstrakurikuler_ids formData string false "Ekstrakurikuler IDs as JSON array, e.g., [1,2,3]"
// @Param sertifikat_to_delete formData string false "Sertifikat URLs to delete as JSON array, e.g., ['url1','url2']"
// @Success 200 {object} gin.H{data=dtos.PelatihResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/sieksa/ekstrakurikuler/update-pelatih [post]
func (c *PelatihController) Update(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID from form (required)
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	// Get optional fields from form
	nama := ctx.PostForm("nama")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	telepon := ctx.PostForm("telepon")
	alamat := ctx.PostForm("alamat")
	keahlian := ctx.PostForm("keahlian")
	status := ctx.PostForm("status")

	// Get foto_profil (optional, single file)
	fotoProfil, _ := ctx.FormFile("foto_profil")

	// Get sertifikat files (optional, multiple files)
	sertifikatFiles := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["sertifikat"]; exists {
			sertifikatFiles = uploadedFiles
		}
	}

	// Parse ekstrakurikuler_ids if provided (JSON array format)
	var ekstrakurikulerIDs []uint
	if ekskulIDsStr := ctx.PostForm("ekstrakurikuler_ids"); ekskulIDsStr != "" {
		if err := json.Unmarshal([]byte(ekskulIDsStr), &ekstrakurikulerIDs); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ekstrakurikuler_ids format, expected JSON array"})
			return
		}
	}

	// Parse sertifikat_to_delete if provided (JSON array of URLs to delete)
	var sertifikatToDelete []string
	if sertifikatDeleteStr := ctx.PostForm("sertifikat_to_delete"); sertifikatDeleteStr != "" {
		if err := json.Unmarshal([]byte(sertifikatDeleteStr), &sertifikatToDelete); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sertifikat_to_delete format, expected JSON array"})
			return
		}
	}

	// Create request DTO with pointers for optional fields
	req := dtos.PelatihUpdateRequest{
		ID:                  uint(id),
		EkstrakurikulerIDs:  ekstrakurikulerIDs,
		SertifikatToDelete:  sertifikatToDelete,
	}

	if nama != "" {
		req.Nama = &nama
	}
	if username != "" {
		req.Username = &username
	}
	if password != "" {
		req.Password = &password
	}
	if telepon != "" {
		req.Telepon = &telepon
	}
	if alamat != "" {
		req.Alamat = &alamat
	}
	if keahlian != "" {
		req.Keahlian = &keahlian
	}
	if status != "" {
		req.Status = &status
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service with files
	data, err := c.service.Update(&req, fotoProfil, sertifikatFiles, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Pelatih by ID
// @Summary Delete Pelatih
// @Description Delete pelatih by ID (soft delete)
// @Tags pelatih
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/sieksa/ekstrakurikuler/delete-pelatih [post]
func (c *PelatihController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pelatih deleted successfully",
	})
}

// GetAll retrieves all pelatih with filters and pagination
// @Summary Get all Pelatih
// @Description Retrieve all Pelatih records with filters and pagination
// @Tags pelatih
// @Accept json
// @Produce json
// @Param body body dtos.PelatihGetAllRequest true "Request body"
// @Success 200 {object} gin.H{data=[]dtos.PelatihResponse,pagination=dtos.PaginationInfo}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/sieksa/ekstrakurikuler/get-pelatih [post]
func (c *PelatihController) GetAll(ctx *gin.Context) {
	var req dtos.PelatihGetAllRequest
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
	data, err := c.service.GetAllWithFilter(repositories.GetPelatihParams{
		Filter: repositories.GetPelatihFilter{
			Nama:              req.Search.Nama,
			EkstrakurikulerID: req.Search.EkstrakurikulerID,
			Status:            req.Search.Status,
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
