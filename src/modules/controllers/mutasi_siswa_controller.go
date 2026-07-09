package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// MutasiSiswaController handles HTTP requests for Mutasi Siswa
type MutasiSiswaController struct {
	service services.MutasiSiswaService
}

// NewMutasiSiswaController creates a new Mutasi Siswa controller
func NewMutasiSiswaController(service services.MutasiSiswaService) *MutasiSiswaController {
	return &MutasiSiswaController{service: service}
}

// CreatePublic creates a new Mutasi Siswa from public form (no auth required)
// @Summary Create new Mutasi Siswa (Public)
// @Description Create a new mutasi siswa dari form public untuk orang tua murid
// @Tags mutasi-siswa
// @Accept multipart/form-data
// @Produce json
// @Param tahun_pelajaran_id formData int true "Tahun Pelajaran ID"
// @Param semester formData int true "Semester"
// @Param nama_lengkap formData string true "Nama Lengkap"
// @Param nama_panggilan formData string false "Nama Panggilan"
// @Param nisn formData string false "NISN"
// @Param tempat_lahir formData string true "Tempat Lahir"
// @Param tanggal_lahir formData string true "Tanggal Lahir (YYYY-MM-DD)"
// @Param jenis_kelamin formData string true "Jenis Kelamin"
// @Param agama formData string true "Agama"
// @Param golongan_darah formData string false "Golongan Darah"
// @Param anak_ke formData int false "Anak Ke"
// @Param jumlah_saudara formData int false "Jumlah Saudara"
// @Param status_anak formData string false "Status Anak"
// @Param alamat formData string true "Alamat"
// @Param rt formData string false "RT"
// @Param rw formData string false "RW"
// @Param kelurahan formData string false "Kelurahan"
// @Param kecamatan formData string false "Kecamatan"
// @Param kota formData string false "Kota"
// @Param provinsi formData string false "Provinsi"
// @Param nama_ayah formData string false "Nama Ayah"
// @Param nama_ibu formData string false "Nama Ibu"
// @Param pendidikan_ayah formData string false "Pendidikan Ayah"
// @Param pendidikan_ibu formData string false "Pendidikan Ibu"
// @Param pekerjaan_ayah formData string false "Pekerjaan Ayah"
// @Param pekerjaan_ibu formData string false "Pekerjaan Ibu"
// @Param penghasilan_ayah formData number false "Penghasilan Ayah"
// @Param penghasilan_ibu formData number false "Penghasilan Ibu"
// @Param nomor_hp_ortu formData string false "Nomor HP Orang Tua"
// @Param nama_wali formData string false "Nama Wali"
// @Param pendidikan_wali formData string false "Pendidikan Wali"
// @Param hubungan_wali formData string false "Hubungan Wali"
// @Param pekerjaan_wali formData string false "Pekerjaan Wali"
// @Param nomor_hp_wali formData string false "Nomor HP Wali"
// @Param pindahan_kelas formData int false "Pindahan Kelas"
// @Param asal_sekolah formData string false "Asal Sekolah"
// @Param nama_asal_sekolah formData string false "Nama Asal Sekolah"
// @Param rapor formData file false "File Rapor"
// @Param akte_kelahiran formData file false "File Akte Kelahiran"
// @Param kartu_keluarga formData file false "File Kartu Keluarga"
// @Param sptjm formData file false "File SPTJM"
// @Success 201 {object} gin.H{data=dtos.MutasiSiswaResponse}
// @Failure 400 {object} gin.H{error=string}
// @Router /api/v1/public/create-mutasi-siswa [post]
func (c *MutasiSiswaController) CreatePublic(ctx *gin.Context) {
	var req dtos.MutasiSiswaCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract file uploads
	files := make(map[string]*multipart.FileHeader)
	if file, err := ctx.FormFile("rapor"); err == nil {
		files["rapor"] = file
	}
	if file, err := ctx.FormFile("akte_kelahiran"); err == nil {
		files["akte_kelahiran"] = file
	}
	if file, err := ctx.FormFile("kartu_keluarga"); err == nil {
		files["kartu_keluarga"] = file
	}
	if file, err := ctx.FormFile("sptjm"); err == nil {
		files["sptjm"] = file
	}

	data, err := c.service.CreatePublic(&req, files)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}


// GetAll retrieves all mutasi siswa with filters and pagination (auth required)
// @Summary Get all Mutasi Siswa with filters
// @Description Retrieve all mutasi siswa dengan filters dan pagination
// @Tags mutasi-siswa
// @Accept json
// @Produce json
// @Param body body dtos.MutasiSiswaGetAllRequest true "Request body with filters and pagination"
// @Success 200 {object} dtos.MutasiSiswaListWithPaginationResponse
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/get-mutasi-siswa [post]
func (c *MutasiSiswaController) GetAll(ctx *gin.Context) {
	var req dtos.MutasiSiswaGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.GetAllWithFilter(&req)
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


// GetByID retrieves mutasi siswa by ID (auth required)
// @Summary Get Mutasi Siswa by ID
// @Description Retrieve mutasi siswa details by ID
// @Tags mutasi-siswa
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with ID"
// @Success 200 {object} gin.H{data=dtos.MutasiSiswaResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/get-mutasi-siswa-by-id [post]
func (c *MutasiSiswaController) GetByID(ctx *gin.Context) {
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


// Update updates mutasi siswa data (auth required)
// @Summary Update Mutasi Siswa
// @Description Update mutasi siswa details (all fields optional, including file uploads)
// @Tags mutasi-siswa
// @Accept multipart/form-data
// @Produce json
// @Param id formData uint true "Mutasi Siswa ID"
// @Param tahun_pelajaran_id formData int false "Tahun Pelajaran ID"
// @Param semester formData int false "Semester"
// @Param nama_lengkap formData string false "Nama Lengkap"
// @Param nama_panggilan formData string false "Nama Panggilan"
// @Param nisn formData string false "NISN"
// @Param tempat_lahir formData string false "Tempat Lahir"
// @Param tanggal_lahir formData string false "Tanggal Lahir (YYYY-MM-DD)"
// @Param jenis_kelamin formData string false "Jenis Kelamin"
// @Param agama formData string false "Agama"
// @Param golongan_darah formData string false "Golongan Darah"
// @Param anak_ke formData int false "Anak Ke"
// @Param jumlah_saudara formData int false "Jumlah Saudara"
// @Param status_anak formData string false "Status Anak"
// @Param alamat formData string false "Alamat"
// @Param rt formData string false "RT"
// @Param rw formData string false "RW"
// @Param kelurahan formData string false "Kelurahan"
// @Param kecamatan formData string false "Kecamatan"
// @Param kota formData string false "Kota"
// @Param provinsi formData string false "Provinsi"
// @Param nama_ayah formData string false "Nama Ayah"
// @Param nama_ibu formData string false "Nama Ibu"
// @Param pendidikan_ayah formData string false "Pendidikan Ayah"
// @Param pendidikan_ibu formData string false "Pendidikan Ibu"
// @Param pekerjaan_ayah formData string false "Pekerjaan Ayah"
// @Param pekerjaan_ibu formData string false "Pekerjaan Ibu"
// @Param penghasilan_ayah formData number false "Penghasilan Ayah"
// @Param penghasilan_ibu formData number false "Penghasilan Ibu"
// @Param nomor_hp_ortu formData string false "Nomor HP Orang Tua"
// @Param nama_wali formData string false "Nama Wali"
// @Param pendidikan_wali formData string false "Pendidikan Wali"
// @Param hubungan_wali formData string false "Hubungan Wali"
// @Param pekerjaan_wali formData string false "Pekerjaan Wali"
// @Param nomor_hp_wali formData string false "Nomor HP Wali"
// @Param pindahan_kelas formData int false "Pindahan Kelas"
// @Param asal_sekolah formData string false "Asal Sekolah"
// @Param nama_asal_sekolah formData string false "Nama Asal Sekolah"
// @Param rapor formData file false "File Rapor (replaces existing)"
// @Param akte_kelahiran formData file false "File Akte Kelahiran (replaces existing)"
// @Param kartu_keluarga formData file false "File Kartu Keluarga (replaces existing)"
// @Param sptjm formData file false "File SPTJM (replaces existing)"
// @Success 200 {object} gin.H{data=dtos.MutasiSiswaResponse}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/edit-mutasi-siswa [post]
func (c *MutasiSiswaController) Update(ctx *gin.Context) {
	// Parse multipart form (max 50MB)
	if err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var req dtos.MutasiSiswaUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract file uploads
	files := make(map[string]*multipart.FileHeader)
	if file, err := ctx.FormFile("rapor"); err == nil {
		files["rapor"] = file
	}
	if file, err := ctx.FormFile("akte_kelahiran"); err == nil {
		files["akte_kelahiran"] = file
	}
	if file, err := ctx.FormFile("kartu_keluarga"); err == nil {
		files["kartu_keluarga"] = file
	}
	if file, err := ctx.FormFile("sptjm"); err == nil {
		files["sptjm"] = file
	}

	data, err := c.service.Update(&req, files)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Mutasi siswa berhasil diupdate",
		"data":    data,
	})
}

// ExportFormulirPDF exports formulir pendaftaran mutasi siswa to PDF (no auth required)
// @Summary Export Formulir Pendaftaran PDF
// @Description Export formulir pendaftaran mutasi siswa to PDF file
// @Tags mutasi-siswa
// @Accept json
// @Produce application/pdf
// @Param body body dtos.IDRequest true "Request body with Mutasi Siswa ID"
// @Success 200 {file} binary "PDF file"
// @Failure 400 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/public/export-pdf-formulir-mutasi-siswa [post]
func (c *MutasiSiswaController) ExportFormulirPDF(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	pdfBytes, err := c.service.ExportFormulirPDF(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=formulir_pendaftaran_%d.pdf", req.ID))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Delete deletes mutasi siswa and all associated files (auth required)
// @Summary Delete Mutasi Siswa
// @Description Delete mutasi siswa record and all associated files from storage
// @Tags mutasi-siswa
// @Accept json
// @Produce json
// @Param body body dtos.IDRequest true "Request body with Mutasi Siswa ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/delete-mutasi-siswa [post]
func (c *MutasiSiswaController) Delete(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Mutasi siswa berhasil dihapus"})
}

// ExportFormulirPDFAuth exports formulir pendaftaran mutasi siswa to PDF (auth required)
// @Summary Export Formulir Pendaftaran PDF (Admin)
// @Description Export formulir pendaftaran mutasi siswa to PDF file (requires authentication)
// @Tags mutasi-siswa
// @Accept json
// @Produce application/pdf
// @Param body body dtos.IDRequest true "Request body with Mutasi Siswa ID"
// @Success 200 {file} binary "PDF file"
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/export-pdf-formulir-mutasi-siswa [post]
func (c *MutasiSiswaController) ExportFormulirPDFAuth(ctx *gin.Context) {
	var req dtos.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	pdfBytes, err := c.service.ExportFormulirPDF(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=formulir_pendaftaran_%d.pdf", req.ID))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}


// ExportExcel exports mutasi siswa data to Excel (auth required)
// @Summary Export Data Mutasi Siswa to Excel
// @Description Export mutasi siswa data to Excel file by tahun pelajaran and semester
// @Tags mutasi-siswa
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param body body dtos.MutasiSiswaExportExcelRequest true "Request body with tahun_pelajaran_id and semester"
// @Success 200 {file} binary "Excel file"
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/export-excel-mutasi-siswa [post]
func (c *MutasiSiswaController) ExportExcel(ctx *gin.Context) {
	var req dtos.MutasiSiswaExportExcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	excelBytes, err := c.service.ExportExcel(&req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Set headers for Excel download
	filename := fmt.Sprintf("data_calon_murid_baru_semester_%d.xlsx", req.Semester)
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}


// ExportListPDF exports mutasi siswa list to PDF (auth required)
// @Summary Export Data Mutasi Siswa List to PDF
// @Description Export mutasi siswa list to PDF file by tahun pelajaran and semester
// @Tags mutasi-siswa
// @Accept json
// @Produce application/pdf
// @Param body body dtos.MutasiSiswaExportExcelRequest true "Request body with tahun_pelajaran_id and semester"
// @Success 200 {file} binary "PDF file"
// @Failure 400 {object} gin.H{error=string}
// @Failure 401 {object} gin.H{error=string}
// @Failure 404 {object} gin.H{error=string}
// @Router /api/v1/spmb-mutasi/export-pdf-mutasi-siswa [post]
func (c *MutasiSiswaController) ExportListPDF(ctx *gin.Context) {
	var req dtos.MutasiSiswaExportExcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := c.service.ExportListPDF(&req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("data_calon_murid_baru_semester_%d.pdf", req.Semester)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
