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

// PertanyaanService handles business logic for Pertanyaan
type PertanyaanService interface {
	CreatePublic(files []*multipart.FileHeader, req *dtos.PertanyaanCreateRequest) (*dtos.PertanyaanResponse, error)
	TrackByIDTiket(idTiket string) (*dtos.PertanyaanTrackResponse, error)
	GetAllWithFilter(req *dtos.PertanyaanGetAllRequest) (*dtos.PertanyaanListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.PertanyaanResponse, error)
	SendReply(files []*multipart.FileHeader, req *dtos.PertanyaanSendReplyRequest, userID uint) (*dtos.PertanyaanResponse, error)
	ClosePertanyaan(id uint) (*dtos.PertanyaanResponse, error)
	DeletePertanyaan(id uint, userID uint) error
}

type PertanyaanServiceImpl struct {
	repository   repositories.PertanyaanRepository
	r2Storage    *utils.R2Storage
	emailService *utils.EmailService
}

// NewPertanyaanService creates a new Pertanyaan service
func NewPertanyaanService(repository repositories.PertanyaanRepository, r2Storage *utils.R2Storage) PertanyaanService {
	return &PertanyaanServiceImpl{
		repository:   repository,
		r2Storage:    r2Storage,
		emailService: utils.NewEmailService(),
	}
}

// generateTicketID generates unique ticket ID: PRT-YYYYMMDD-XXXX
func (s *PertanyaanServiceImpl) generateTicketID() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	
	// Generate 4 random alphanumeric characters
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, 4)
	for i := range random {
		random[i] = charset[rand.Intn(len(charset))]
	}
	
	return fmt.Sprintf("PRT-%s-%s", dateStr, string(random))
}

// CreatePublic creates a new Pertanyaan from public form
func (s *PertanyaanServiceImpl) CreatePublic(files []*multipart.FileHeader, req *dtos.PertanyaanCreateRequest) (*dtos.PertanyaanResponse, error) {
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

			// Upload file to R2 in layanan-umpan-balik/pertanyaan directory
			fileKey, err := s.r2Storage.UploadFile(file, "layanan-umpan-balik/pertanyaan")
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

	// Set prioritas from request, default to "Sedang" if empty
	prioritas := req.Prioritas
	if prioritas == "" {
		prioritas = "Sedang"
	}

	// Create pertanyaan record
	data := &models.Pertanyaan{
		IDTiket:        ticketID,
		Nama:           req.Nama,
		Email:          req.Email,
		Telepon:        req.Telepon,
		Kategori:       req.Kategori,
		Prioritas:      prioritas,
		Judul:          req.Judul,
		Deskripsi:      req.Deskripsi,
		FilePertanyaan: fileJSON,
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
func (s *PertanyaanServiceImpl) mapToResponse(data *models.Pertanyaan) *dtos.PertanyaanResponse {
	// Parse file_pertanyaan JSON
	var filePertanyaan []models.FileItem
	if len(data.FilePertanyaan) > 0 {
		_ = json.Unmarshal(data.FilePertanyaan, &filePertanyaan)
	}

	// Convert file URLs to full public URLs
	for i := range filePertanyaan {
		filePertanyaan[i].URL = s.r2Storage.GetPublicURL(filePertanyaan[i].URL)
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

	resp := &dtos.PertanyaanResponse{
		ID:               data.ID,
		IDTiket:          data.IDTiket,
		TanggalPengajuan: data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
		Nama:             data.Nama,
		Email:            data.Email,
		Telepon:          data.Telepon,
		Kategori:         data.Kategori,
		Prioritas:        data.Prioritas,
		Judul:            data.Judul,
		Deskripsi:        data.Deskripsi,
		FilePertanyaan:   filePertanyaan,
		JudulJawaban:     data.JudulJawaban,
		DeskripsiJawaban: data.DeskripsiJawaban,
		FileJawaban:      fileJawaban,
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

// TrackByIDTiket retrieves Pertanyaan tracking info by ID Tiket
func (s *PertanyaanServiceImpl) TrackByIDTiket(idTiket string) (*dtos.PertanyaanTrackResponse, error) {
	data, err := s.repository.GetByIDTiket(idTiket)
	if err != nil {
		return nil, fmt.Errorf("pertanyaan dengan ID Tiket %s tidak ditemukan", idTiket)
	}

	return &dtos.PertanyaanTrackResponse{
		IDTiket:          data.IDTiket,
		TanggalPengajuan: data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
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


// GetAllWithFilter retrieves all Pertanyaan with filters and pagination
func (s *PertanyaanServiceImpl) GetAllWithFilter(req *dtos.PertanyaanGetAllRequest) (*dtos.PertanyaanListWithPaginationResponse, error) {
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
	params := repositories.GetPertanyaanParams{
		Filter: repositories.GetPertanyaanFilter{
			IDTiket:   req.Search.IDTiket,
			StartDate: req.Search.StartDate,
			EndDate:   req.Search.EndDate,
			Nama:      req.Search.Nama,
			Email:     req.Search.Email,
			Kategori:  req.Search.Kategori,
			Prioritas: req.Search.Prioritas,
			Judul:     req.Search.Judul,
			Status:    req.Search.Status,
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
	var responses []dtos.PertanyaanResponse
	for _, item := range data {
		responses = append(responses, *s.mapToResponse(&item))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dtos.PertanyaanListWithPaginationResponse{
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


// GetByID retrieves Pertanyaan by ID
func (s *PertanyaanServiceImpl) GetByID(id uint) (*dtos.PertanyaanResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pertanyaan dengan ID %d tidak ditemukan", id)
	}

	return s.mapToResponse(data), nil
}


// SendReply sends email reply and updates pertanyaan record
func (s *PertanyaanServiceImpl) SendReply(files []*multipart.FileHeader, req *dtos.PertanyaanSendReplyRequest, userID uint) (*dtos.PertanyaanResponse, error) {
	// Get pertanyaan data
	data, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("pertanyaan tidak ditemukan")
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
			fileKey, err := s.r2Storage.UploadFile(file, "layanan-umpan-balik/pertanyaan/jawaban")
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
	// Parse file_pertanyaan for email
	var filePertanyaanItems []models.FileItem
	if len(data.FilePertanyaan) > 0 {
		_ = json.Unmarshal(data.FilePertanyaan, &filePertanyaanItems)
	}

	var filePertanyaanLinks []utils.FileLink
	for _, item := range filePertanyaanItems {
		filePertanyaanLinks = append(filePertanyaanLinks, utils.FileLink{
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

	emailData := utils.EmailData{
		IDTiket:             data.IDTiket,
		Nama:                data.Nama,
		Email:               data.Email,
		Telepon:             data.Telepon,
		TanggalPengajuan:    data.TanggalPengajuan.Format("2006-01-02 15:04:05"),
		Kategori:            data.Kategori,
		Prioritas:           data.Prioritas,
		JudulPertanyaan:     data.Judul,
		DeskripsiPertanyaan: data.Deskripsi,
		FilePertanyaan:      filePertanyaanLinks,
		JudulJawaban:        req.JudulJawaban,
		DeskripsiJawaban:    req.DeskripsiJawaban,
		FileJawaban:         fileJawabanLinks,
	}

	// Try to send email
	if err := s.emailService.SendPertanyaanReply(data.Email, emailData); err != nil {
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

// ClosePertanyaan closes pertanyaan and sets tanggal_selesai
func (s *PertanyaanServiceImpl) ClosePertanyaan(id uint) (*dtos.PertanyaanResponse, error) {
	// Get pertanyaan data
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pertanyaan tidak ditemukan")
	}

	// Update status and tanggal_selesai
	// Use Asia/Jakarta timezone (WIB - UTC+7)
	// Use FixedZone to ensure WIB timezone works even without timezone database
	wib := time.FixedZone("WIB", 7*60*60) // UTC+7
	now := time.Now().In(wib)
	data.Status = "closed"
	data.TanggalSelesai = &now

	if err := s.repository.Update(data); err != nil {
		return nil, fmt.Errorf("gagal menutup pertanyaan: %w", err)
	}

	return s.mapToResponse(data), nil
}
// DeletePertanyaan soft deletes pertanyaan by setting deleted_at and deleted_by_id
func (s *PertanyaanServiceImpl) DeletePertanyaan(id uint, userID uint) error {
	// Check if pertanyaan exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("pertanyaan tidak ditemukan")
	}

	// Soft delete with user tracking
	if err := s.repository.SoftDeleteWithUser(id, userID); err != nil {
		return fmt.Errorf("gagal menghapus pertanyaan: %w", err)
	}

	return nil
}