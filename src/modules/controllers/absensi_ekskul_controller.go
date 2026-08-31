package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"
	"pintu-backend/src/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AbsensiEkskulController struct {
	service *services.AbsensiEkskulService
}

func NewAbsensiEkskulController(service *services.AbsensiEkskulService) *AbsensiEkskulController {
	return &AbsensiEkskulController{service: service}
}

// Create handles creating ekstrakurikuler attendance
func (c *AbsensiEkskulController) Create(ctx *gin.Context) {
	var req dtos.AbsensiEkskulCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from auth context
	var createdByID *uint
	if userID, exists := ctx.Get("userID"); exists {
		if id, ok := userID.(uint); ok {
			createdByID = &id
		}
	}

	response, err := c.service.Create(&req, createdByID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Absensi ekstrakurikuler created successfully",
		"data":    response,
	})
}

// GetAbsensiSiswa handles getting attendance by ekstrakurikuler and tahun pelajaran
func (c *AbsensiEkskulController) GetAbsensiSiswa(ctx *gin.Context) {
	var req dtos.AbsensiEkskulGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetAbsensiSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Absensi ekstrakurikuler retrieved successfully",
		"data":    response,
	})
}

// GetAbsensiSiswaByID handles getting single absensi siswa by ID
func (c *AbsensiEkskulController) GetAbsensiSiswaByID(ctx *gin.Context) {
	var req dtos.AbsensiSiswaGetByIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetAbsensiSiswaByID(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Absensi siswa retrieved successfully",
		"data":    response,
	})
}

// UpdateAbsensiSiswa handles updating absensi siswa status and keterangan
func (c *AbsensiEkskulController) UpdateAbsensiSiswa(ctx *gin.Context) {
	var req dtos.AbsensiSiswaUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.UpdateAbsensiSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Absensi siswa updated successfully",
		"data":    response,
	})
}


// GetAbsensiPelatih handles getting pelatih attendance
func (c *AbsensiEkskulController) GetAbsensiPelatih(ctx *gin.Context) {
	var req dtos.AbsensiPelatihGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetAbsensiPelatih(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Absensi pelatih retrieved successfully",
		"data":    response,
	})
}


// GetKegiatanEkskul handles getting kegiatan list by ekstrakurikuler
func (c *AbsensiEkskulController) GetKegiatanEkskul(ctx *gin.Context) {
	var req dtos.KegiatanEkskulGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetKegiatanEkskul(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kegiatan ekstrakurikuler retrieved successfully",
		"data":    response,
	})
}

// UpdateAbsensiPelatih handles toggling pelatih attendance
func (c *AbsensiEkskulController) UpdateAbsensiPelatih(ctx *gin.Context) {
	var req dtos.AbsensiPelatihUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.UpdateAbsensiPelatih(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": response.Message,
		"data":    response,
	})
}


// UpdateKegiatan handles updating kegiatan ekstrakurikuler including photos and other fields
func (c *AbsensiEkskulController) UpdateKegiatan(ctx *gin.Context) {
	// Get kegiatan_ekskul_id from form data
	kegiatanIDStr := ctx.PostForm("id")
	if kegiatanIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Convert to uint
	var kegiatanID uint
	if _, err := fmt.Sscanf(kegiatanIDStr, "%d", &kegiatanID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	// Get optional text fields from form data
	tanggalKegiatan := ctx.PostForm("tanggal_kegiatan")
	waktuMulai := ctx.PostForm("waktu_mulai")
	waktuSelesai := ctx.PostForm("waktu_selesai")
	materiKegiatan := ctx.PostForm("materi_kegiatan")

	// Convert to pointers
	var tanggalPtr, waktuMulaiPtr, waktuSelesaiPtr, materiPtr *string
	if tanggalKegiatan != "" {
		tanggalPtr = &tanggalKegiatan
	}
	if waktuMulai != "" {
		waktuMulaiPtr = &waktuMulai
	}
	if waktuSelesai != "" {
		waktuSelesaiPtr = &waktuSelesai
	}
	if materiKegiatan != "" {
		materiPtr = &materiKegiatan
	}

	// Get multiple files from form data (optional for upload)
	form, err := ctx.MultipartForm()
	if err != nil {
		// If no multipart form, that's okay - might be just updating text fields
		form = nil
	}

	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["foto"]
	}

	// Get foto URLs to delete from form data (optional)
	fotoToDeleteStr := ctx.PostFormArray("foto_to_delete")

	// Initialize SIEKSA R2 storage
	storage := utils.NewSieksaR2Storage()

	// Upload new files and collect URLs
	var fotoUrls []string
	for _, file := range files {
		// Validate file type (only images)
		contentType := file.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png" && contentType != "image/webp" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid file type for %s, only JPEG, PNG, and WebP are allowed", file.Filename)})
			return
		}

		// Upload to R2
		fileKey, err := storage.UploadFile(file, "kegiatan-ekskul")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload %s: %v", file.Filename, err)})
			return
		}

		// Get public URL
		publicURL := storage.GetPublicURL(fileKey)
		fotoUrls = append(fotoUrls, publicURL)
	}

	// Delete files from R2 storage
	for _, fotoURL := range fotoToDeleteStr {
		if fotoURL == "" {
			continue
		}

		// Extract file key from URL
		// URL format: https://sieksa-storage.sdnsukapura01.sch.id/kegiatan-ekskul/123456-filename.jpg
		// We need: kegiatan-ekskul/123456-filename.jpg
		
		// Split by domain to get the path
		parts := strings.Split(fotoURL, "/")
		if len(parts) >= 2 {
			// Get last 2 parts: folder/filename
			fileKey := strings.Join(parts[len(parts)-2:], "/")
			
			// Delete from R2
			if err := storage.DeleteFile(fileKey); err != nil {
				fmt.Printf("Warning: failed to delete file %s from R2: %v\n", fileKey, err)
			} else {
				fmt.Printf("Successfully deleted file %s from R2\n", fileKey)
			}
		}
	}

	// Update kegiatan in database
	response, err := c.service.UpdateKegiatanEkskul(kegiatanID, fotoUrls, fotoToDeleteStr, tanggalPtr, waktuMulaiPtr, waktuSelesaiPtr, materiPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": response.Message,
		"data":    response,
	})
}



// GetKegiatanEkskulByID handles getting kegiatan ekstrakurikuler detail by ID
func (c *AbsensiEkskulController) GetKegiatanEkskulByID(ctx *gin.Context) {
	var req dtos.KegiatanEkskulGetByIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetKegiatanEkskulByID(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kegiatan ekstrakurikuler retrieved successfully",
		"data":    response,
	})
}


// DownloadExcelAbsensiSiswa handles downloading Excel file for student attendance
func (c *AbsensiEkskulController) DownloadExcelAbsensiSiswa(ctx *gin.Context) {
	var req dtos.AbsensiEkskulGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate Excel file
	file, filename, err := c.service.DownloadExcelAbsensiSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write file to response
	if err := file.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}
}

// DownloadPDFAbsensiSiswa handles downloading PDF file for student attendance
func (c *AbsensiEkskulController) DownloadPDFAbsensiSiswa(ctx *gin.Context) {
	var req dtos.AbsensiEkskulGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF file
	pdfBytes, filename, err := c.service.DownloadPDFAbsensiSiswa(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write PDF to response
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// DownloadExcelAbsensiPelatih handles downloading Excel file for pelatih attendance
func (c *AbsensiEkskulController) DownloadExcelAbsensiPelatih(ctx *gin.Context) {
	var req dtos.AbsensiPelatihExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate Excel file
	file, filename, err := c.service.DownloadExcelAbsensiPelatih(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write file to response
	if err := file.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}
}

// DownloadPDFAbsensiPelatih handles downloading PDF file for pelatih attendance
func (c *AbsensiEkskulController) DownloadPDFAbsensiPelatih(ctx *gin.Context) {
	var req dtos.AbsensiPelatihExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF file
	pdfBytes, filename, err := c.service.DownloadPDFAbsensiPelatih(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write PDF to response
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// DownloadWordDokumentasiEkskul handles downloading Word file for ekstrakurikuler documentation
func (c *AbsensiEkskulController) DownloadWordDokumentasiEkskul(ctx *gin.Context) {
	var req dtos.KegiatanEkskulDownloadWordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate Word file
	wordBytes, filename, err := c.service.DownloadWordDokumentasiEkskul(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write Word to response
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", wordBytes)
}

// DownloadPDFDokumentasiEkskul handles downloading PDF file for ekstrakurikuler documentation
func (c *AbsensiEkskulController) DownloadPDFDokumentasiEkskul(ctx *gin.Context) {
	var req dtos.KegiatanEkskulDownloadWordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF file
	pdfBytes, filename, err := c.service.DownloadPDFDokumentasiEkskul(&req)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error generating PDF: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write PDF to response
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
