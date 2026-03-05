package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
)

// KepegawaianController handles HTTP requests for Kepegawaian
type KepegawaianController struct {
	service services.KepegawaianService
}

// NewKepegawaianController creates a new Kepegawaian controller
func NewKepegawaianController(service services.KepegawaianService) *KepegawaianController {
	return &KepegawaianController{service: service}
}

// Create creates a new Kepegawaian (JSON only, no file upload)
func (c *KepegawaianController) Create(ctx *gin.Context) {
	var req dtos.KepegawaianCreateRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service
	result, err := c.service.Create(&req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// GetByID retrieves a Kepegawaian by ID
func (c *KepegawaianController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "kepegawaian not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetByNIP retrieves a Kepegawaian by NIP
func (c *KepegawaianController) GetByNIP(ctx *gin.Context) {
	var req struct {
		NIP string `json:"nip" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GetByNIP(req.NIP)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "kepegawaian with NIP not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAll retrieves all Kepegawaian with pagination and filters
func (c *KepegawaianController) GetAll(ctx *gin.Context) {
	var req dtos.KepegawaianGetAllRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
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
	data, err := c.service.GetAllWithFilter(repositories.GetKepegawaianParams{
		Filter: repositories.GetKepegawaianFilter{
			Nama:     req.Search.Nama,
			Username: req.Search.Username,
			NIP:      req.Search.NIP,
			NKKI:     req.Search.NKKI,
			Kategori: req.Search.Kategori,
			Jabatan:  req.Search.Jabatan,
			RoleID:   req.Search.RoleID,
			Status:   req.Search.Status,
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

// Update updates a Kepegawaian
func (c *KepegawaianController) Update(ctx *gin.Context) {
	// Parse multipart form (max 100MB)
	if err := ctx.Request.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get optional fields
	nama := ctx.PostForm("nama")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	nip := ctx.PostForm("nip")
	nkki := ctx.PostForm("nkki")
	kategori := ctx.PostForm("kategori")
	jabatan := ctx.PostForm("jabatan")
	status := ctx.PostForm("status")

	// Get bidang_studi_id (optional)
	var bidangStudiID *uint
	if bidangStudiIDStr := ctx.PostForm("bidang_studi_id"); bidangStudiIDStr != "" {
		if bidangStudiIDUint, err := strconv.ParseUint(bidangStudiIDStr, 10, 32); err == nil {
			bidangID := uint(bidangStudiIDUint)
			bidangStudiID = &bidangID
		}
	}

	// Get rombel_guru_kelas_id (optional)
	var rombelGuruKelasID *uint
	if rombelIDStr := ctx.PostForm("rombel_guru_kelas_id"); rombelIDStr != "" {
		if rombelIDUint, err := strconv.ParseUint(rombelIDStr, 10, 32); err == nil {
			rombelID := uint(rombelIDUint)
			rombelGuruKelasID = &rombelID
		}
	}

	// Get rombel_bidang_studi (optional, as JSON string)
	var rombelBidangStudi []uint
	if rombelBidangStudiJSON := ctx.PostForm("rombel_bidang_studi"); rombelBidangStudiJSON != "" {
		var tempSlice []interface{}
		if err := json.Unmarshal([]byte(rombelBidangStudiJSON), &tempSlice); err == nil {
			for _, item := range tempSlice {
				if num, ok := item.(float64); ok {
					rombelBidangStudi = append(rombelBidangStudi, uint(num))
				}
			}
		}
	}

	// Get foto (optional)
	foto, _ := ctx.FormFile("foto")

	// Get document files
	docMap := make(map[string][]*multipart.FileHeader)
	docTypes := []string{"kk", "akta_lahir", "ktp", "ijazah_sd", "ijazah_smp", "ijazah_sma", 
		"ijazah_s1", "ijazah_s2", "ijazah_s3", "sertifikat_pendidik", "sertifikat_lainnya", 
		"sk", "dokumen_lainnya"}

	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		for _, docType := range docTypes {
			if uploadedFiles, exists := form.File[docType]; exists {
				docMap[docType] = uploadedFiles
			}
		}
	}

	// Get user ID from context
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Create request DTO
	req := &dtos.KepegawaianUpdateRequest{
		ID:                    uint(id),
		Nama:                  nama,
		Username:              username,
		Password:              password,
		NIP:                   nip,
		NKKI:                  nkki,
		Kategori:              kategori,
		Jabatan:               jabatan,
		BidangStudiID:         bidangStudiID,
		RombelGuruKelasID:     rombelGuruKelasID,
		RombelBidangStudi:     rombelBidangStudi,
		Status:                status,
	}

	// Call service
	result, err := c.service.Update(uint(id), foto, docMap, req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Delete deletes a Kepegawaian
func (c *KepegawaianController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "kepegawaian deleted successfully"})
}
