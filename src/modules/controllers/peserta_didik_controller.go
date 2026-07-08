package controllers

import (
	"fmt"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
)

// PesertaDidikController handles HTTP requests for PesertaDidik
type PesertaDidikController struct {
	service services.PesertaDidikService
}

// NewPesertaDidikController creates a new PesertaDidik controller
func NewPesertaDidikController(service services.PesertaDidikService) *PesertaDidikController {
	return &PesertaDidikController{service: service}
}

// Create creates a new PesertaDidik
func (c *PesertaDidikController) Create(ctx *gin.Context) {
	var req dtos.PesertaDidikCreateRequest

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

// GetByID retrieves a PesertaDidik by ID
func (c *PesertaDidikController) GetByID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "peserta didik tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetByNIS retrieves a PesertaDidik by NIS
func (c *PesertaDidikController) GetByNIS(ctx *gin.Context) {
	var req struct {
		NIS string `json:"nis" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GetByNIS(req.NIS)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "peserta didik dengan NIS tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAll retrieves all PesertaDidik with pagination and filters
func (c *PesertaDidikController) GetAll(ctx *gin.Context) {
	var req dtos.PesertaDidikGetAllRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Build filter
	filter := repositories.GetPesertaDidikFilter{
		Nama:             req.Search.Nama,
		NIS:              req.Search.NIS,
		JenisKelamin:     req.Search.JenisKelamin,
		NISN:             req.Search.NISN,
		TempatLahir:      req.Search.TempatLahir,
		NIK:              req.Search.NIK,
		Agama:            req.Search.Agama,
		Status:           req.Search.Status,
	}

	// Set default pagination
	limit := req.Pagination.Limit
	page := req.Pagination.Page

	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// Call service with filter
	params := repositories.GetPesertaDidikParams{
		Filter: filter,
		Limit:  limit,
		Offset: offset,
	}

	result, err := c.service.GetAllWithFilter(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Update updates a PesertaDidik
func (c *PesertaDidikController) Update(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil { // 10MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get ID from form
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	// Get optional photo
	photo, _ := ctx.FormFile("photo")

	// Parse role_ids from form array - support multiple formats
	var roleIDs []uint
	hasRoleIDs := false
	if ctx.Request.MultipartForm != nil && ctx.Request.MultipartForm.Value != nil {
		// Debug: log all form values
		fmt.Printf("DEBUG - All form values: %+v\n", ctx.Request.MultipartForm.Value)
		
		// Try multiple possible formats
		var roleIDStrings []string
		
		// Format 1: role_ids[]
		if vals, ok := ctx.Request.MultipartForm.Value["role_ids[]"]; ok && len(vals) > 0 {
			roleIDStrings = vals
			fmt.Printf("DEBUG - Found role_ids[] format: %v\n", vals)
		}
		
		// Format 2: role_ids
		if len(roleIDStrings) == 0 {
			if vals, ok := ctx.Request.MultipartForm.Value["role_ids"]; ok && len(vals) > 0 {
				roleIDStrings = vals
				fmt.Printf("DEBUG - Found role_ids format: %v\n", vals)
			}
		}
		
		// Format 3: role_ids[0], role_ids[1], etc
		if len(roleIDStrings) == 0 {
			for key, vals := range ctx.Request.MultipartForm.Value {
				if key == "role_ids[0]" || key == "role_ids[1]" || key == "role_ids[2]" || key == "role_ids[3]" || key == "role_ids[4]" {
					roleIDStrings = append(roleIDStrings, vals...)
					fmt.Printf("DEBUG - Found %s format: %v\n", key, vals)
				}
			}
		}
		
		// Check if role_ids was sent
		if len(roleIDStrings) > 0 {
			hasRoleIDs = true
			// Parse the role IDs
			for _, roleIDStr := range roleIDStrings {
				var roleID uint
				if _, err := fmt.Sscanf(roleIDStr, "%d", &roleID); err == nil {
					roleIDs = append(roleIDs, roleID)
				}
			}
			fmt.Printf("DEBUG - Parsed roleIDs: %v, hasRoleIDs: %v\n", roleIDs, hasRoleIDs)
		}
	}

	// Build update request from form fields
	req := &dtos.PesertaDidikUpdateRequest{
		ID:           id,
		Nama:         ctx.PostForm("nama"),
		NIS:          ctx.PostForm("nis"),
		JenisKelamin: ctx.PostForm("jenis_kelamin"),
		NISN:         ctx.PostForm("nisn"),
		TempatLahir:  ctx.PostForm("tempat_lahir"),
		TanggalLahir: ctx.PostForm("tanggal_lahir"),
		NIK:          ctx.PostForm("nik"),
		Agama:        ctx.PostForm("agama"),
		Alamat:       ctx.PostForm("alamat"),
		RT:           ctx.PostForm("rt"),
		RW:           ctx.PostForm("rw"),
		Kelurahan:    ctx.PostForm("kelurahan"),
		Kecamatan:    ctx.PostForm("kecamatan"),
		KodePos:      ctx.PostForm("kode_pos"),
		NamaAyah:     ctx.PostForm("nama_ayah"),
		NamaIbu:      ctx.PostForm("nama_ibu"),
		Status:       ctx.PostForm("status"),
		Username:     ctx.PostForm("username"),
		Password:     ctx.PostForm("password"),
	}

	// Set RoleIDs only if provided
	if hasRoleIDs {
		req.RoleIDs = &roleIDs
	}

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	// Call service with photo
	result, err := c.service.Update(id, photo, req, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// Delete deletes a PesertaDidik
func (c *PesertaDidikController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "peserta didik berhasil dihapus"})
}

// ImportExcel imports peserta didik data from Excel file
func (c *PesertaDidikController) ImportExcel(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file excel wajib diunggah"})
		return
	}
	defer file.Close()

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	result, err := c.service.ImportExcel(file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// DownloadTemplate downloads the Excel template for peserta didik import
func (c *PesertaDidikController) DownloadTemplate(ctx *gin.Context) {
	f, err := c.service.DownloadTemplate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat template"})
		return
	}
	defer f.Close()

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_peserta_didik.xlsx")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ExportDataIndukSiswaExcel exports data induk siswa to Excel file
func (c *PesertaDidikController) ExportDataIndukSiswaExcel(ctx *gin.Context) {
	var req dtos.ExportDataIndukSiswaRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If error parsing, use empty status (get all)
		req.Status = ""
	}
	
	f, err := c.service.ExportDataIndukSiswaExcel(req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	filename := "data_induk_siswa.xlsx"
	if req.Status != "" {
		filename = fmt.Sprintf("data_induk_siswa_%s.xlsx", req.Status)
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ExportDataIndukSiswaPDF exports data induk siswa to PDF file
func (c *PesertaDidikController) ExportDataIndukSiswaPDF(ctx *gin.Context) {
	var req dtos.ExportDataIndukSiswaRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If error parsing, use empty status (get all)
		req.Status = ""
	}
	
	pdfBytes, err := c.service.ExportDataIndukSiswaPDF(req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := "data_induk_siswa.pdf"
	if req.Status != "" {
		filename = fmt.Sprintf("data_induk_siswa_%s.pdf", req.Status)
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportPemetaanRombelExcel exports pemetaan rombel to Excel file
func (c *PesertaDidikController) ExportPemetaanRombelExcel(ctx *gin.Context) {
	var req dtos.ExportPemetaanRombelRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If error parsing, use 0 (get all)
		req.RombelID = 0
		req.TahunPelajaranID = 0
	}
	
	f, err := c.service.ExportPemetaanRombelExcel(req.RombelID, req.TahunPelajaranID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	filename := "pemetaan_rombel.xlsx"
	if req.RombelID > 0 && req.TahunPelajaranID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_%d_%d.xlsx", req.RombelID, req.TahunPelajaranID)
	} else if req.RombelID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_%d.xlsx", req.RombelID)
	} else if req.TahunPelajaranID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_tp_%d.xlsx", req.TahunPelajaranID)
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ExportPemetaanRombelPDF exports pemetaan rombel to PDF file
func (c *PesertaDidikController) ExportPemetaanRombelPDF(ctx *gin.Context) {
	var req dtos.ExportPemetaanRombelRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If error parsing, use 0 (get all)
		req.RombelID = 0
		req.TahunPelajaranID = 0
	}
	
	pdfBytes, err := c.service.ExportPemetaanRombelPDF(req.RombelID, req.TahunPelajaranID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := "pemetaan_rombel.pdf"
	if req.RombelID > 0 && req.TahunPelajaranID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_%d_%d.pdf", req.RombelID, req.TahunPelajaranID)
	} else if req.RombelID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_%d.pdf", req.RombelID)
	} else if req.TahunPelajaranID > 0 {
		filename = fmt.Sprintf("pemetaan_rombel_tp_%d.pdf", req.TahunPelajaranID)
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetTotalSiswa retrieves total count of peserta didik with active tahun pelajaran (public endpoint)
func (c *PesertaDidikController) GetTotalSiswa(ctx *gin.Context) {
	result, err := c.service.GetTotalSiswa()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GenerateBarcodeAllPesertaDidik generates barcodes for all peserta didik
func (c *PesertaDidikController) GenerateBarcodeAllPesertaDidik(ctx *gin.Context) {
	result, err := c.service.GenerateBarcodeAllPesertaDidik()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GenerateBarcodePesertaDidikByID generates or regenerates barcode for a specific peserta didik by ID
func (c *PesertaDidikController) GenerateBarcodePesertaDidikByID(ctx *gin.Context) {
	var req dtos.GenerateBarcodeByIDRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := c.service.GenerateBarcodePesertaDidikByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// DownloadTemplateSiswaLulus downloads the Excel template for siswa lulus with nama and nis columns only
func (c *PesertaDidikController) DownloadTemplateSiswaLulus(ctx *gin.Context) {
	f, err := c.service.DownloadTemplateSiswaLulus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat template"})
		return
	}
	defer f.Close()

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_siswa_lulus.xlsx")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengirim file"})
		return
	}
}

// ImportSiswaLulus imports siswa lulus data from Excel file and updates status to "lulus"
func (c *PesertaDidikController) ImportSiswaLulus(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file excel wajib diunggah"})
		return
	}
	defer file.Close()

	// Get user ID from context (set by middleware)
	userID, _ := ctx.Get("userID")
	userIDUint := userID.(uint)

	result, err := c.service.ImportSiswaLulus(file, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "data": result})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// DownloadKartuPelajar downloads student cards as PDF
func (c *PesertaDidikController) DownloadKartuPelajar(ctx *gin.Context) {
	var req dtos.DownloadKartuPelajarRequest
	
	// Jika ada request body, bind JSON. Jika tidak ada atau kosong, download semua
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Jika error parsing, bisa jadi body kosong, download semua
		req.PesertaDidikIDs = []uint{}
	}
	
	pdfBytes, err := c.service.DownloadKartuPelajar(req.PesertaDidikIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=kartu_pelajar_SDN_Sukapura_01.pdf")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
