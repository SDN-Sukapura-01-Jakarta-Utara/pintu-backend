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

// PrestasiController handles HTTP requests for Prestasi
type PrestasiController struct {
	service services.PrestasiService
}

// NewPrestasiController creates a new Prestasi controller
func NewPrestasiController(service services.PrestasiService) *PrestasiController {
	return &PrestasiController{service: service}
}

// Create creates a new Prestasi with foto uploads
// @Summary Create new Prestasi
// @Description Create a new Prestasi with foto upload to Cloudflare R2
// @Tags prestasi
// @Accept multipart/form-data
// @Produce json
// @Param peserta_didik_id formData uint false "Peserta Didik ID"
// @Param nama formData string true "Prestasi name"
// @Param jenis formData string true "Prestasi type"
// @Param nama_prestasi formData string true "Achievement name"
// @Param tingkat_prestasi formData string false "Achievement level"
// @Param penyelenggara formData string false "Organizer"
// @Param tanggal_lomba formData string true "Competition date (YYYY-MM-DD)"
// @Param juara formData string true "Rank/Position"
// @Param keterangan formData string false "Description"
// @Param ekstrakurikuler_id formData uint false "Ekstrakurikuler ID"
// @Param anggota_tim formData string false "Team members JSON array"
// @Param foto formData file false "Achievement photos - multiple files allowed (max 5MB each)"
// @Success 201 {object} gin.H{data=dtos.PrestasiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/prestasi/create-prestasi [post]
func (c *PrestasiController) Create(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get required fields
	jenis := ctx.PostForm("jenis")
	if jenis == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "jenis is required"})
		return
	}

	namaPrestasi := ctx.PostForm("nama_prestasi")
	if namaPrestasi == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nama_prestasi is required"})
		return
	}

	tanggalLomba := ctx.PostForm("tanggal_lomba")
	if tanggalLomba == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tanggal_lomba is required"})
		return
	}

	juara := ctx.PostForm("juara")
	if juara == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "juara is required"})
		return
	}

	tahunPelajaranID := ctx.PostForm("tahun_pelajaran_id")
	if tahunPelajaranID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tahun_pelajaran_id is required"})
		return
	}

	tahunPelajaranIDUint, err := strconv.ParseUint(tahunPelajaranID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tahun_pelajaran_id format"})
		return
	}

	// Get optional fields
	var pesertaDidikID *uint
	if pesertaDidikIDStr := ctx.PostForm("peserta_didik_id"); pesertaDidikIDStr != "" {
		if id, err := strconv.ParseUint(pesertaDidikIDStr, 10, 32); err == nil {
			pesertaDidikIDUint := uint(id)
			pesertaDidikID = &pesertaDidikIDUint
		}
	}

	namaGrup := ctx.PostForm("nama_grup")
	tingkatPrestasi := ctx.PostForm("tingkat_prestasi")
	penyelenggara := ctx.PostForm("penyelenggara")
	keterangan := ctx.PostForm("keterangan")
	status := ctx.PostForm("status")

	var ekstrakurikulerID *uint
	if ekstrakurikulerIDStr := ctx.PostForm("ekstrakurikuler_id"); ekstrakurikulerIDStr != "" {
		if id, err := strconv.ParseUint(ekstrakurikulerIDStr, 10, 32); err == nil {
			ekstrakurikulerIDUint := uint(id)
			ekstrakurikulerID = &ekstrakurikulerIDUint
		}
	}

	// Parse anggota tim JSON
	var anggotaTim []dtos.AnggotaTimCreateRequest
	if anggotaTimStr := ctx.PostForm("anggota_tim"); anggotaTimStr != "" {
		if err := json.Unmarshal([]byte(anggotaTimStr), &anggotaTim); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid anggota_tim JSON format"})
			return
		}
	}

	// Get foto (optional, multiple)
	foto := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["foto"]; exists {
			foto = uploadedFiles
		}
	}

	// Get foto thumbnail info (optional)
	var fotoThumbnails []string
	if fotoThumbnailStr := ctx.PostForm("foto_thumbnails"); fotoThumbnailStr != "" {
		if err := json.Unmarshal([]byte(fotoThumbnailStr), &fotoThumbnails); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid foto_thumbnails JSON format"})
			return
		}
	}

	// Create request DTO
	req := &dtos.PrestasiCreateRequest{
		PesertaDidikID:    pesertaDidikID,
		Jenis:             jenis,
		NamaGrup:          namaGrup,
		NamaPrestasi:      namaPrestasi,
		TingkatPrestasi:   tingkatPrestasi,
		Penyelenggara:     penyelenggara,
		TanggalLomba:      tanggalLomba,
		Juara:             juara,
		Keterangan:        keterangan,
		EkstrakurikulerID: ekstrakurikulerID,
		TahunPelajaranID:  uint(tahunPelajaranIDUint),
		Status:            status,
		AnggotaTim:        anggotaTim,
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Call service
	data, err := c.service.Create(foto, fotoThumbnails, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

// GetAll retrieves all prestasi with filters and pagination
// @Summary Get all Prestasi
// @Description Retrieve all Prestasi records with filters and pagination
// @Tags prestasi
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.PrestasiListWithPaginationResponse}
// @Failure 401 {object} gin.H{error=string}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/prestasi/get-prestasi [post]
func (c *PrestasiController) GetAll(ctx *gin.Context) {
	var req dtos.PrestasiGetAllRequest
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

	// Parse date filters
	var startDate, endDate time.Time
	if req.Search.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.StartDate); err == nil {
			startDate = parsed
		}
	}
	if req.Search.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.Search.EndDate); err == nil {
			// Set to end of day for inclusive range
			endDate = parsed.Add(time.Hour * 24).Add(-time.Nanosecond)
		}
	}

	// Call service with filters
	data, err := c.service.GetAllWithFilter(repositories.GetPrestasiParams{
		Filter: repositories.GetPrestasiFilter{
			PesertaDidikID:    req.Search.PesertaDidikID,
			NamaPesertaDidik:  req.Search.NamaPesertaDidik,
			Jenis:             req.Search.Jenis,
			NamaGrup:          req.Search.NamaGrup,
			NamaPrestasi:      req.Search.NamaPrestasi,
			TingkatPrestasi:   req.Search.TingkatPrestasi,
			Penyelenggara:     req.Search.Penyelenggara,
			StartDate:         startDate,
			EndDate:           endDate,
			Juara:             req.Search.Juara,
			EkstrakurikulerID: req.Search.EkstrakurikulerID,
			TahunPelajaranID:  req.Search.TahunPelajaranID,
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

// GetByID retrieves Prestasi by ID
// @Summary Get Prestasi by ID
// @Description Retrieve prestasi details by ID
// @Tags prestasi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.PrestasiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/prestasi/get-prestasi-by-id [post]
func (c *PrestasiController) GetByID(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := c.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Prestasi not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates Prestasi
// @Summary Update Prestasi
// @Description Update Prestasi details (all fields optional, including foto uploads)
// @Tags prestasi
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Prestasi ID"
// @Param peserta_didik_id formData uint false "Peserta Didik ID"
// @Param nama formData string false "Prestasi name"
// @Param jenis formData string false "Prestasi type"
// @Param nama_prestasi formData string false "Achievement name"
// @Param tingkat_prestasi formData string false "Achievement level"
// @Param penyelenggara formData string false "Organizer"
// @Param tanggal_lomba formData string false "Competition date (YYYY-MM-DD)"
// @Param juara formData string false "Rank/Position"
// @Param keterangan formData string false "Description"
// @Param ekstrakurikuler_id formData uint false "Ekstrakurikuler ID"
// @Param anggota_tim formData string false "Team members JSON array"
// @Param foto formData file false "Achievement photos - multiple files allowed (max 5MB each)"
// @Param foto_to_delete formData string false "JSON array of foto IDs to delete"
// @Success 200 {object} gin.H{data=dtos.PrestasiResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/prestasi/update-prestasi [post]
func (c *PrestasiController) Update(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID from form
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

	// Get optional fields
	var pesertaDidikID *uint
	if pesertaDidikIDStr := ctx.PostForm("peserta_didik_id"); pesertaDidikIDStr != "" {
		if parsedID, err := strconv.ParseUint(pesertaDidikIDStr, 10, 32); err == nil {
			pesertaDidikIDUint := uint(parsedID)
			pesertaDidikID = &pesertaDidikIDUint
		}
	}

	jenis := ctx.PostForm("jenis")
	namaGrup := ctx.PostForm("nama_grup")
	namaPrestasi := ctx.PostForm("nama_prestasi")
	tingkatPrestasi := ctx.PostForm("tingkat_prestasi")
	penyelenggara := ctx.PostForm("penyelenggara")
	tanggalLomba := ctx.PostForm("tanggal_lomba")
	juara := ctx.PostForm("juara")
	keterangan := ctx.PostForm("keterangan")
	status := ctx.PostForm("status")

	var tahunPelajaranID uint
	if tahunPelajaranIDStr := ctx.PostForm("tahun_pelajaran_id"); tahunPelajaranIDStr != "" {
		if parsedID, err := strconv.ParseUint(tahunPelajaranIDStr, 10, 32); err == nil {
			tahunPelajaranID = uint(parsedID)
		}
	}

	var ekstrakurikulerID *uint
	if ekstrakurikulerIDStr := ctx.PostForm("ekstrakurikuler_id"); ekstrakurikulerIDStr != "" {
		if parsedID, err := strconv.ParseUint(ekstrakurikulerIDStr, 10, 32); err == nil {
			ekstrakurikulerIDUint := uint(parsedID)
			ekstrakurikulerID = &ekstrakurikulerIDUint
		}
	}

	// Parse anggota tim JSON
	var anggotaTim []dtos.AnggotaTimUpdateRequest
	if anggotaTimStr := ctx.PostForm("anggota_tim"); anggotaTimStr != "" {
		if err := json.Unmarshal([]byte(anggotaTimStr), &anggotaTim); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid anggota_tim JSON format"})
			return
		}
	}

	// Get foto (optional, multiple)
	foto := []*multipart.FileHeader{}
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if uploadedFiles, exists := form.File["foto"]; exists {
			foto = uploadedFiles
		}
	}

	// Get foto thumbnail info (optional)
	var fotoThumbnails []string
	if fotoThumbnailStr := ctx.PostForm("foto_thumbnails"); fotoThumbnailStr != "" {
		if err := json.Unmarshal([]byte(fotoThumbnailStr), &fotoThumbnails); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid foto_thumbnails JSON format"})
			return
		}
	}

	// Parse foto_to_delete if provided
	var fotoToDelete []string
	if fotoDeleteStr := ctx.PostForm("foto_to_delete"); fotoDeleteStr != "" {
		_ = json.Unmarshal([]byte(fotoDeleteStr), &fotoToDelete)
	}

	// Create request DTO
	req := &dtos.PrestasiUpdateRequest{
		ID:                uint(id),
		PesertaDidikID:    pesertaDidikID,
		Jenis:             jenis,
		NamaGrup:          namaGrup,
		NamaPrestasi:      namaPrestasi,
		TingkatPrestasi:   tingkatPrestasi,
		Penyelenggara:     penyelenggara,
		TanggalLomba:      tanggalLomba,
		Juara:             juara,
		Keterangan:        keterangan,
		EkstrakurikulerID: ekstrakurikulerID,
		TahunPelajaranID:  tahunPelajaranID,
		Status:            status,
		FotoToDelete:      fotoToDelete,
		AnggotaTim:        anggotaTim,
	}

	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	data, err := c.service.Update(uint(id), foto, fotoThumbnails, req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete deletes Prestasi by ID
// @Summary Delete Prestasi
// @Description Delete Prestasi by ID (also deletes all foto from R2 and anggota tim)
// @Tags prestasi
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/prestasi/delete-prestasi [post]
func (c *PrestasiController) Delete(ctx *gin.Context) {
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
		"message": "Prestasi deleted successfully",
	})
}

// GetPublicLatest retrieves 10 latest prestasi for public display (no auth required)
// @Summary Get latest Prestasi for public
// @Description Retrieve 10 latest prestasi ordered by tanggal_lomba DESC (no authentication required)
// @Tags prestasi
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{data=dtos.PrestasiPublicListResponse}
// @Failure 500 {object} gin.H{error=string}
// @Router /api/v1/public/get-data-prestasi [post]
func (c *PrestasiController) GetPublicLatest(ctx *gin.Context) {
	data, err := c.service.GetPublicLatest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}