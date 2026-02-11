package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type ArticleService interface {
	Create(gambar *multipart.FileHeader, files []*multipart.FileHeader, req *dtos.ArticleCreateRequest, userID uint) (*dtos.ArticleResponse, error)
	GetByID(id uint) (*dtos.ArticleResponse, error)
	GetAll(limit int, offset int) (*dtos.ArticleListResponse, error)
	Update(id uint, gambar *multipart.FileHeader, files []*multipart.FileHeader, req *dtos.ArticleUpdateRequest, userID uint) (*dtos.ArticleResponse, error)
	Delete(id uint) error
}

type ArticleServiceImpl struct {
	repository repositories.ArticleRepository
	r2Storage  *utils.R2Storage
}

// NewArticleService creates a new Article service
func NewArticleService(repository repositories.ArticleRepository, r2Storage *utils.R2Storage) ArticleService {
	return &ArticleServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new Article with file uploads to R2
func (s *ArticleServiceImpl) Create(gambar *multipart.FileHeader, files []*multipart.FileHeader, req *dtos.ArticleCreateRequest, userID uint) (*dtos.ArticleResponse, error) {
	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("invalid tanggal format (use YYYY-MM-DD)")
	}

	// Upload gambar if provided
	var gambarURL string
	if gambar != nil {
		// Validate gambar file
		if gambar.Size > 5*1024*1024 { // 5MB
			return nil, errors.New("gambar size must not exceed 5MB")
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
		contentType := gambar.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, errors.New("only image files are allowed for gambar (jpeg, png, gif, webp)")
		}

		// Upload gambar to R2
		fileKey, err := s.r2Storage.UploadFile(gambar, "artikel")
		if err != nil {
			return nil, err
		}
		gambarURL = fileKey
	}

	// Upload files if provided
	var fileItems []models.FileItem
	if len(files) > 0 {
		for _, file := range files {
			if file == nil {
				continue
			}

			// Validate file size (max 10MB per file)
			if file.Size > 10*1024*1024 { // 10MB
				return nil, errors.New("each file must not exceed 10MB")
			}

			// Validate file type - allow multiple file types
			allowedTypes := map[string]bool{
				"application/pdf":                          true,
				"application/msword":                       true,
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
				"application/vnd.ms-excel": true,
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
				"text/plain":               true,
				"image/jpeg":               true,
				"image/png":                true,
			}
			contentType := file.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("file type not allowed. Allowed: PDF, DOC, DOCX, XLS, XLSX, TXT, JPG, PNG")
			}

			// Upload file to R2
			fileKey, err := s.r2Storage.UploadFile(file, "artikel")
			if err != nil {
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

	// Set defaults
	statusPublikasi := req.StatusPublikasi
	if statusPublikasi == "" {
		statusPublikasi = "draft"
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	// Convert fileItems to JSON
	filesJSON, _ := json.Marshal(fileItems)

	// Create article record
	data := &models.Article{
		Judul:           req.Judul,
		Tanggal:         tanggal,
		Kategori:        req.Kategori,
		Deskripsi:       req.Deskripsi,
		Gambar:          gambarURL,
		Files:           filesJSON,
		Penulis:         req.Penulis,
		StatusPublikasi: statusPublikasi,
		Status:          status,
		CreatedByID:     &userID,
	}

	if err := s.repository.Create(data); err != nil {
		// If database save fails, delete uploaded files
		if gambarURL != "" {
			_ = s.r2Storage.DeleteFile(gambarURL)
		}
		for _, item := range fileItems {
			_ = s.r2Storage.DeleteFile(item.URL)
		}
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Article by ID
func (s *ArticleServiceImpl) GetByID(id uint) (*dtos.ArticleResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Articles
func (s *ArticleServiceImpl) GetAll(limit int, offset int) (*dtos.ArticleListResponse, error) {
	// Set default limit and offset
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	data, total, err := s.repository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.ArticleResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.ArticleListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// Update updates Article
func (s *ArticleServiceImpl) Update(id uint, gambar *multipart.FileHeader, files []*multipart.FileHeader, req *dtos.ArticleUpdateRequest, userID uint) (*dtos.ArticleResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	oldGambar := existing.Gambar

	// Update basic fields if provided
	if req.Judul != "" {
		existing.Judul = req.Judul
	}
	if req.Tanggal != "" {
		tanggal, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			return nil, errors.New("invalid tanggal format (use YYYY-MM-DD)")
		}
		existing.Tanggal = tanggal
	}
	if req.Kategori != "" {
		existing.Kategori = req.Kategori
	}
	if req.Deskripsi != "" {
		existing.Deskripsi = req.Deskripsi
	}
	if req.Penulis != "" {
		existing.Penulis = req.Penulis
	}
	if req.StatusPublikasi != "" {
		existing.StatusPublikasi = req.StatusPublikasi
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	// Update gambar if provided
	if gambar != nil {
		// Validate gambar file
		if gambar.Size > 5*1024*1024 { // 5MB
			return nil, errors.New("gambar size must not exceed 5MB")
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
		contentType := gambar.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, errors.New("only image files are allowed for gambar (jpeg, png, gif, webp)")
		}

		// Upload new gambar
		newFileKey, err := s.r2Storage.UploadFile(gambar, "artikel")
		if err != nil {
			return nil, err
		}

		// Delete old gambar if exists
		if oldGambar != "" {
			_ = s.r2Storage.DeleteFile(oldGambar)
		}

		existing.Gambar = newFileKey
	}

	// Delete files if specified
	if len(req.FilesToDelete) > 0 {
		var existingFileItems []models.FileItem
		if err := json.Unmarshal(existing.Files, &existingFileItems); err == nil {
			// Build map of file IDs to delete
			deleteMap := make(map[string]bool)
			for _, fileID := range req.FilesToDelete {
				deleteMap[fileID] = true
			}

			// Filter out files to delete and delete from R2
			var remainingFiles []models.FileItem
			for _, file := range existingFileItems {
				if deleteMap[file.ID] {
					// Delete from R2
					_ = s.r2Storage.DeleteFile(file.URL)
				} else {
					remainingFiles = append(remainingFiles, file)
				}
			}

			// Update files array
			filesJSON, _ := json.Marshal(remainingFiles)
			existing.Files = filesJSON
		}
	}

	// Add new files if provided (files lama tetap)
	if len(files) > 0 {
		var existingFileItems []models.FileItem
		// Get existing files
		if err := json.Unmarshal(existing.Files, &existingFileItems); err != nil && existing.Files != nil && len(existing.Files) > 0 {
			// Only log, don't fail
		}

		// Upload and add new files
		for _, file := range files {
			if file == nil {
				continue
			}

			// Validate file size (max 10MB per file)
			if file.Size > 10*1024*1024 { // 10MB
				return nil, errors.New("each file must not exceed 10MB")
			}

			// Validate file type
			allowedTypes := map[string]bool{
				"application/pdf":                          true,
				"application/msword":                       true,
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
				"application/vnd.ms-excel": true,
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
				"text/plain":               true,
				"image/jpeg":               true,
				"image/png":                true,
			}
			contentType := file.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("file type not allowed. Allowed: PDF, DOC, DOCX, XLS, XLSX, TXT, JPG, PNG")
			}

			// Upload file to R2
			fileKey, err := s.r2Storage.UploadFile(file, "artikel")
			if err != nil {
				return nil, err
			}

			// Generate unique file ID
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			existingFileItems = append(existingFileItems, models.FileItem{
				ID:       fileID,
				Filename: file.Filename,
				URL:      fileKey,
				Size:     file.Size,
			})
		}

		// Convert updated files to JSON
		filesJSON, _ := json.Marshal(existingFileItems)
		existing.Files = filesJSON
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		// If DB save fails, delete the uploaded files
		if gambar != nil {
			_ = s.r2Storage.DeleteFile(existing.Gambar)
		}
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Article by ID
func (s *ArticleServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	// Delete gambar from R2
	if existing.Gambar != "" {
		_ = s.r2Storage.DeleteFile(existing.Gambar)
	}

	// Delete all files from R2
	var fileItems []models.FileItem
	if err := json.Unmarshal(existing.Files, &fileItems); err == nil {
		for _, file := range fileItems {
			_ = s.r2Storage.DeleteFile(file.URL)
		}
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *ArticleServiceImpl) mapToResponse(data *models.Article) *dtos.ArticleResponse {
	// Map files from JSON
	var fileItems []dtos.FileItemDTO
	var fileModels []models.FileItem
	if err := json.Unmarshal(data.Files, &fileModels); err == nil {
		for _, file := range fileModels {
			fileItems = append(fileItems, dtos.FileItemDTO{
				ID:       file.ID,
				Filename: file.Filename,
				URL:      s.r2Storage.GetPublicURL(file.URL),
				Size:     file.Size,
			})
		}
	}

	return &dtos.ArticleResponse{
		ID:              data.ID,
		Judul:           data.Judul,
		Tanggal:         data.Tanggal,
		Kategori:        data.Kategori,
		Deskripsi:       data.Deskripsi,
		Gambar:          s.r2Storage.GetPublicURL(data.Gambar),
		Files:           fileItems,
		Penulis:         data.Penulis,
		StatusPublikasi: data.StatusPublikasi,
		Status:          data.Status,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
		CreatedByID:     data.CreatedByID,
		UpdatedByID:     data.UpdatedByID,
	}
}
