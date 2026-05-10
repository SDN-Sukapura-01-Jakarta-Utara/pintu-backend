package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// PengaduanService handles business logic for Pengaduan
type PengaduanService interface {
	CreatePublic(files []*multipart.FileHeader, req *dtos.PengaduanCreateRequest) (*dtos.PengaduanResponse, error)
	TrackByIDTiket(idTiket string) (*dtos.PengaduanTrackResponse, error)
	GetAllWithFilter(req *dtos.PengaduanGetAllRequest) (*dtos.PengaduanListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.PengaduanResponse, error)
	SendReply(files []*multipart.FileHeader, req *dtos.PengaduanSendReplyRequest, userID uint) (*dtos.PengaduanResponse, error)
	SaveTindakLanjut(files []*multipart.FileHeader, req *dtos.PengaduanSaveTindakLanjutRequest, userID uint) (*dtos.PengaduanResponse, error)
	ClosePengaduan(id uint) (*dtos.PengaduanResponse, error)
	DeletePengaduan(id uint, userID uint) error
}

type PengaduanServiceImpl struct {
	repository   repositories.PengaduanRepository
	r2Storage    *utils.R2Storage
	emailService *utils.EmailService
}

// NewPengaduanService creates a new Pengaduan service
func NewPengaduanService(repository repositories.PengaduanRepository, r2Storage *utils.R2Storage) PengaduanService {
	return &PengaduanServiceImpl{
		repository:   repository,
		r2Storage:    r2Storage,
		emailService: utils.NewEmailService(),
	}
}

// generateTicketID generates unique ticket ID: PGD-YYYYMMDD-XXXX
func (s *PengaduanServiceImpl) generateTicketID() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	
	// Generate 4 random alphanumeric characters
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, 4)
	for i := range random {
		random[i] = charset[rand.Intn(len(charset))]
	}
	
	return fmt.Sprintf("PGD-%s-%s", dateStr, string(random))
}

// CreatePublic creates a new Pengaduan from public form
func (s *PengaduanServiceImpl) CreatePublic(files []*multipart.FileHeader, req *dtos.PengaduanCreateRequest) (*dtos.PengaduanResponse, error) {
	// Upload files if provided
	var fileItems []models.FileItem
	if len(files) > 0 {
		for _, file := range files {
			if file == nil {
				continue
			}

			// Validate file size (max 10MB per file)
			if file.Size > 10*1024*1024 {
				return nil, fmt.Errorf("each file must not exceed 10MB")
			}

			// Upload file to R2 in layanan-umpan-balik/pengaduan directory
			fileKey, err := s.r2Storage.UploadFile(file, "layanan-umpan-balik/pengaduan")
			if err != nil {
				// Cleanup already uploaded files on error
				for _, item := range fileItems {
					_ = s.r2Storage.DeleteFile(item.URL)
				}
				return nil, err
			}

			// Generate unique file ID: file_timestamp_randomstring
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			fileItems = append(fileItems, models.FileItem{
				ID:       fileID,
				Filename: file.Filename,
				URL:      fileKey,
				Size:     file.Size,
			})
		}
	}

	// Generate unique ticket ID
	ticketID := s.generateTicketID()

	// Convert fileItems to JSON
	fileJSON, _ := json.Marshal(fileItems)

	// Set default values
	tipePelapor := req.TipePelapor
	if tipePelapor == "" {
		tipePelapor = "anonim"
	}

	prioritas := req.Prioritas
	if prioritas == "" {
		prioritas = "Sedang"
	}

	// Handle nullable fields
	var nama, email, telepon *string
	if req.Nama != "" {
		nama = &req.Nama
	}
	if req.Email != "" {
		email = &req.Email
	}
	if req.Telepon != "" {
		telepon = &req.Telepon
	}

	// Create pengaduan record
	data := &models.Pengaduan{
		IDTiket:        ticketID,
		TipePelapor:    tipePelapor,
		Nama:           nama,
		Email:          email,
		Telepon:        telepon,
		Kategori:       req.Kategori,
		Prioritas:      prioritas,
		Judul:          req.Judul,
		Deskripsi:      req.Deskripsi,
		FilePengaduan:  fileJSON,
		Status:         "pending",
		EmailTerkirim:  false,
	}

	if err := s.repository.Create(data); err != nil {
		// Cleanup uploaded files on database error
		for _, item := range fileItems {
			_ = s.r2Storage.DeleteFile(item.URL)
		}
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// mapToResponse converts model to response DTO
func (s *PengaduanServiceImpl) mapToResponse(data *models.Pengaduan) *dtos.PengaduanResponse {
	// Parse file_pengaduan JSON
	var filePengaduan []models.FileItem
	if len(data.FilePengaduan) > 0 {
		_ = json.Unmarshal(data.FilePengaduan, &filePengaduan)
	}

	// Convert file URLs to full public URLs
	for i := range filePengaduan {
		filePengaduan[i].URL = s.r2Storage.GetPublicURL(filePengaduan[i].URL)
	}

	// Parse file_jawaban JSON
	var fileJawaban []models.FileItem
	if len(data.FileJawaban) > 0 {
		_ = json.Unmarshal(data.FileJawaban, &fileJawaban)
	}

	// Convert file URLs to full public URLs
	for i := range fileJawaban {
		fileJawaban[i].URL = s.r2Storage.GetPublicURL(fileJawaban[i].URL)
	}

	// Parse file_tindak_lanjut JSON
	var fileTindakLanjut []models.FileItem
	if len(data.FileTindakLanjut) > 0 {
		_ = json.Unmarshal(data.FileTindakLanjut, &fileTindakLanjut)
	}

	// Convert file URLs to full public URLs
	for i := range fileTindakLanjut {
		fileTindakLanjut[i].URL = s.r2Storage.GetPublicURL(fileTindakLanjut[i].URL)
	}

	resp := &dtos.PengaduanResponse{
		ID:               data.ID,
		IDTiket:          data.IDTiket,
		TanggalPengajuan: data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
		TipePelapor:      data.TipePelapor,
		Nama:             data.Nama,
		Email:            data.Email,
		Telepon:          data.Telepon,
		Kategori:         data.Kategori,
		Prioritas:        data.Prioritas,
		Judul:            data.Judul,
		Deskripsi:        data.Deskripsi,
		FilePengaduan:    filePengaduan,
		JudulJawaban:     data.JudulJawaban,
		DeskripsiJawaban: data.DeskripsiJawaban,
		FileJawaban:      fileJawaban,
		TindakLanjut:     data.TindakLanjut,
		FileTindakLanjut: fileTindakLanjut,
		EmailTerkirim:    data.EmailTerkirim,
		Status:           data.Status,
		RepliedBy:        data.RepliedBy,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if data.TanggalProses != nil {
		tanggalProses := data.TanggalProses.Format("2006-01-02 15:04:05")
		resp.TanggalProses = &tanggalProses
	}

	if data.TanggalSelesai != nil {
		tanggalSelesai := data.TanggalSelesai.Format("2006-01-02 15:04:05")
		resp.TanggalSelesai = &tanggalSelesai
	}

	return resp
}

// TrackByIDTiket retrieves Pengaduan tracking info by ID Tiket
func (s *PengaduanServiceImpl) TrackByIDTiket(idTiket string) (*dtos.PengaduanTrackResponse, error) {
	data, err := s.repository.GetByIDTiket(idTiket)
	if err != nil {
		return nil, fmt.Errorf("pengaduan dengan ID Tiket %s tidak ditemukan", idTiket)
	}

	return &dtos.PengaduanTrackResponse{
		IDTiket:          data.IDTiket,
		TanggalPengajuan: data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
		TipePelapor:      data.TipePelapor,
		Nama:             data.Nama,
		Email:            data.Email,
		Telepon:          data.Telepon,
		Kategori:         data.Kategori,
		Prioritas:        data.Prioritas,
		Judul:            data.Judul,
		Deskripsi:        data.Deskripsi,
		Status:           data.Status,
	}, nil
}

// GetAllWithFilter retrieves all Pengaduan with filters and pagination
func (s *PengaduanServiceImpl) GetAllWithFilter(req *dtos.PengaduanGetAllRequest) (*dtos.PengaduanListWithPaginationResponse, error) {
	// Set default pagination
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Build filter params
	params := repositories.GetPengaduanParams{
		Filter: repositories.GetPengaduanFilter{
			IDTiket:     req.Search.IDTiket,
			StartDate:   req.Search.StartDate,
			EndDate:     req.Search.EndDate,
			TipePelapor: req.Search.TipePelapor,
			Nama:        req.Search.Nama,
			Email:       req.Search.Email,
			Kategori:    req.Search.Kategori,
			Prioritas:   req.Search.Prioritas,
			Judul:       req.Search.Judul,
			Status:      req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	}

	// Get data from repository
	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	var responses []dtos.PengaduanResponse
	for _, item := range data {
		responses = append(responses, *s.mapToResponse(&item))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dtos.PengaduanListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationMeta{
			Limit:      limit,
			Offset:     offset,
			Page:       page,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetByID retrieves Pengaduan by ID
func (s *PengaduanServiceImpl) GetByID(id uint) (*dtos.PengaduanResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pengaduan dengan ID %d tidak ditemukan", id)
	}

	return s.mapToResponse(data), nil
}

// SendReply sends email reply and updates pengaduan record
func (s *PengaduanServiceImpl) SendReply(files []*multipart.FileHeader, req *dtos.PengaduanSendReplyRequest, userID uint) (*dtos.PengaduanResponse, error) {
	// Get pengaduan data
	data, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("pengaduan tidak ditemukan")
	}

	// Check if email is available for sending reply
	if data.Email == nil || *data.Email == "" {
		return nil, fmt.Errorf("tidak dapat mengirim email karena pengaduan ini tidak memiliki alamat email")
	}

	// Upload file jawaban if provided
	var fileItems []models.FileItem
	if len(files) > 0 {
		for _, file := range files {
			if file == nil {
				continue
			}

			// Validate file size (max 10MB per file)
			if file.Size > 10*1024*1024 {
				return nil, fmt.Errorf("each file must not exceed 10MB")
			}

			// Upload file to R2
			fileKey, err := s.r2Storage.UploadFile(file, "layanan-umpan-balik/pengaduan/jawaban")
			if err != nil {
				// Cleanup already uploaded files on error
				for _, item := range fileItems {
					_ = s.r2Storage.DeleteFile(item.URL)
				}
				return nil, err
			}

			// Generate unique file ID
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			fileItems = append(fileItems, models.FileItem{
				ID:       fileID,
				Filename: file.Filename,
				URL:      fileKey,
				Size:     file.Size,
			})
		}
	}

	// Convert fileItems to JSON
	fileJSON, _ := json.Marshal(fileItems)

	// Update database FIRST before sending email
	// Use Asia/Jakarta timezone (WIB - UTC+7)
	// Use FixedZone to ensure WIB timezone works even without timezone database
	wib := time.FixedZone("WIB", 7*60*60) // UTC+7
	now := time.Now().In(wib)
	judulJawaban := req.JudulJawaban
	deskripsiJawaban := req.DeskripsiJawaban

	data.JudulJawaban = &judulJawaban
	data.DeskripsiJawaban = &deskripsiJawaban
	data.FileJawaban = fileJSON
	data.TanggalProses = &now
	data.EmailTerkirim = false // Set to false first, will update after email sent
	data.Status = "processed"
	data.RepliedBy = &userID

	// Save to database first
	if err := s.repository.Update(data); err != nil {
		// Cleanup uploaded files if database update fails
		for _, item := range fileItems {
			_ = s.r2Storage.DeleteFile(item.URL)
		}
		return nil, fmt.Errorf("gagal menyimpan data ke database: %w", err)
	}

	// Database saved successfully, now prepare and send email
	// Parse file_pengaduan for email
	var filePengaduanItems []models.FileItem
	if len(data.FilePengaduan) > 0 {
		_ = json.Unmarshal(data.FilePengaduan, &filePengaduanItems)
	}

	var filePengaduanLinks []utils.FileLink
	for _, item := range filePengaduanItems {
		filePengaduanLinks = append(filePengaduanLinks, utils.FileLink{
			Name: item.Filename,
			URL:  s.r2Storage.GetPublicURL(item.URL),
		})
	}

	var fileJawabanLinks []utils.FileLink
	for _, item := range fileItems {
		fileJawabanLinks = append(fileJawabanLinks, utils.FileLink{
			Name: item.Filename,
			URL:  s.r2Storage.GetPublicURL(item.URL),
		})
	}

	// Handle nullable fields for email
	nama := "Anonim"
	if data.Nama != nil {
		nama = *data.Nama
	}

	telepon := ""
	if data.Telepon != nil {
		telepon = *data.Telepon
	}

	emailData := utils.EmailData{
		IDTiket:             data.IDTiket,
		Nama:                nama,
		Email:               *data.Email,
		Telepon:             telepon,
		TanggalPengajuan:    data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
		Kategori:            data.Kategori,
		Prioritas:           data.Prioritas,
		JudulPertanyaan:     data.Judul,
		DeskripsiPertanyaan: data.Deskripsi,
		FilePertanyaan:      filePengaduanLinks,
		JudulJawaban:        req.JudulJawaban,
		DeskripsiJawaban:    req.DeskripsiJawaban,
		FileJawaban:         fileJawabanLinks,
	}

	// Try to send email
	if err := s.emailService.SendPengaduanReply(*data.Email, emailData); err != nil {
		// Email failed but data already saved - return success with warning
		// Keep email_terkirim as false to indicate email wasn't sent
		return s.mapToResponse(data), nil
	}

	// Email sent successfully, update email_terkirim flag
	data.EmailTerkirim = true
	if err := s.repository.Update(data); err != nil {
		// Failed to update email flag but email was sent - not critical
		// Return success anyway since main data is saved
	}

	return s.mapToResponse(data), nil
}

// SaveTindakLanjut saves tindak lanjut for pengaduan
func (s *PengaduanServiceImpl) SaveTindakLanjut(files []*multipart.FileHeader, req *dtos.PengaduanSaveTindakLanjutRequest, userID uint) (*dtos.PengaduanResponse, error) {
	// Get pengaduan data
	data, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("pengaduan tidak ditemukan")
	}

	// Parse existing file_tindak_lanjut
	var existingFiles []models.FileItem
	if len(data.FileTindakLanjut) > 0 {
		_ = json.Unmarshal(data.FileTindakLanjut, &existingFiles)
	}

	// Debug log
	fmt.Printf("DEBUG - Files to delete: %v\n", req.FilesToDelete)
	fmt.Printf("DEBUG - Existing files before deletion: %+v\n", existingFiles)

	// Handle file deletion if files_to_delete is provided
	if len(req.FilesToDelete) > 0 {
		// Build map of file IDs to delete
		deleteMap := make(map[string]bool)
		for _, fileID := range req.FilesToDelete {
			deleteMap[fileID] = true
		}

		// Filter out files to delete and delete from R2
		var remainingFiles []models.FileItem
		for _, file := range existingFiles {
			if deleteMap[file.ID] {
				// Delete from R2
				fmt.Printf("DEBUG - Deleting file from R2: ID='%s', URL='%s'\n", file.ID, file.URL)
				if err := s.r2Storage.DeleteFile(file.URL); err != nil {
					fmt.Printf("ERROR - Failed to delete file from R2: %v\n", err)
				} else {
					fmt.Printf("SUCCESS - File deleted from R2: %s\n", file.URL)
				}
			} else {
				remainingFiles = append(remainingFiles, file)
			}
		}
		existingFiles = remainingFiles
		fmt.Printf("DEBUG - Files after deletion: %+v\n", existingFiles)
	}

	// Upload new files if provided
	if len(files) > 0 {
		for _, file := range files {
			if file == nil {
				continue
			}

			// Validate file size (max 10MB per file)
			if file.Size > 10*1024*1024 {
				return nil, fmt.Errorf("each file must not exceed 10MB")
			}

			// Upload file to R2
			fileKey, err := s.r2Storage.UploadFile(file, "layanan-umpan-balik/pengaduan/tindak-lanjut")
			if err != nil {
				// Cleanup already uploaded files on error
				return nil, err
			}

			// Generate unique file ID
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			existingFiles = append(existingFiles, models.FileItem{
				ID:       fileID,
				Filename: file.Filename,
				URL:      fileKey,
				Size:     file.Size,
			})
		}
	}

	// Convert fileItems to JSON
	fileJSON, _ := json.Marshal(existingFiles)

	// Update tindak lanjut and file tindak lanjut
	tindakLanjut := req.TindakLanjut
	data.TindakLanjut = &tindakLanjut
	data.FileTindakLanjut = fileJSON

	// Update tanggal_proses only if tipe_pelapor is "anonim" and tanggal_proses is still nil
	if data.TipePelapor == "anonim" && data.TanggalProses == nil {
		now := time.Now()
		data.TanggalProses = &now
	}

	// Update status to "processed"
	data.Status = "processed"

	// Save to database
	if err := s.repository.Update(data); err != nil {
		return nil, fmt.Errorf("gagal menyimpan tindak lanjut: %w", err)
	}

	return s.mapToResponse(data), nil
}

// ClosePengaduan closes pengaduan and sets tanggal_selesai
func (s *PengaduanServiceImpl) ClosePengaduan(id uint) (*dtos.PengaduanResponse, error) {
	// Get pengaduan data
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pengaduan tidak ditemukan")
	}

	// Update status and tanggal_selesai
	// Use Asia/Jakarta timezone (WIB - UTC+7)
	wib := time.FixedZone("WIB", 7*60*60) // UTC+7
	now := time.Now().In(wib)
	data.Status = "closed"
	data.TanggalSelesai = &now

	if err := s.repository.Update(data); err != nil {
		return nil, fmt.Errorf("gagal menutup pengaduan: %w", err)
	}

	return s.mapToResponse(data), nil
}

// DeletePengaduan soft deletes pengaduan by setting deleted_at and deleted_by_id
func (s *PengaduanServiceImpl) DeletePengaduan(id uint, userID uint) error {
	// Check if pengaduan exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("pengaduan tidak ditemukan")
	}

	// Soft delete with user tracking
	if err := s.repository.SoftDeleteWithUser(id, userID); err != nil {
		return fmt.Errorf("gagal menghapus pengaduan: %w", err)
	}

	return nil
}
