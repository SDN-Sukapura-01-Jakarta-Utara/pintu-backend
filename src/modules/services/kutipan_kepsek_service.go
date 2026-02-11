package services

import (
	"errors"
	"mime/multipart"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// KutipanKepsekService handles business logic for KutipanKepsek
type KutipanKepsekService interface {
	Create(file *multipart.FileHeader, req *dtos.KutipanKepsekCreateRequest, userID uint) (*dtos.KutipanKepsekResponse, error)
	GetByID(id uint) (*dtos.KutipanKepsekResponse, error)
	GetAll(limit int, offset int) (*dtos.KutipanKepsekListResponse, error)
	UpdateWithFile(id uint, file *multipart.FileHeader, req *dtos.KutipanKepsekUpdateRequest, userID uint) (*dtos.KutipanKepsekResponse, error)
	Delete(id uint) error
}

type KutipanKepsekServiceImpl struct {
	repository repositories.KutipanKepsekRepository
	r2Storage  *utils.R2Storage
}

// NewKutipanKepsekService creates a new KutipanKepsek service
func NewKutipanKepsekService(repository repositories.KutipanKepsekRepository, r2Storage *utils.R2Storage) KutipanKepsekService {
	return &KutipanKepsekServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new KutipanKepsek with file upload
func (s *KutipanKepsekServiceImpl) Create(file *multipart.FileHeader, req *dtos.KutipanKepsekCreateRequest, userID uint) (*dtos.KutipanKepsekResponse, error) {
	// Validate file
	if file == nil {
		return nil, errors.New("file is required")
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return nil, errors.New("file size must not exceed 5MB")
	}

	// Validate file type (only images)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, errors.New("only image files are allowed (jpeg, png, gif, webp)")
	}

	// Upload file to R2 (kepsek directory)
	fileKey, err := s.r2Storage.UploadFile(file, "kepsek")
	if err != nil {
		return nil, err
	}

	data := &models.KutipanKepsek{
		NamaKepsek:    req.NamaKepsek,
		FotoKepsek:    fileKey,
		KutipanKepsek: req.KutipanKepsek,
		CreatedByID:   &userID,
	}

	if err := s.repository.Create(data); err != nil {
		// If database save fails, delete the uploaded file
		_ = s.r2Storage.DeleteFile(fileKey)
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves KutipanKepsek by ID
func (s *KutipanKepsekServiceImpl) GetByID(id uint) (*dtos.KutipanKepsekResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all KutipanKepsek
func (s *KutipanKepsekServiceImpl) GetAll(limit int, offset int) (*dtos.KutipanKepsekListResponse, error) {
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
	responses := make([]dtos.KutipanKepsekResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.KutipanKepsekListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// UpdateWithFile updates KutipanKepsek with optional file upload
func (s *KutipanKepsekServiceImpl) UpdateWithFile(id uint, file *multipart.FileHeader, req *dtos.KutipanKepsekUpdateRequest, userID uint) (*dtos.KutipanKepsekResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldFile := existing.FotoKepsek

	// If file provided, validate and upload
	if file != nil {
		// Validate file
		if file.Size > 5*1024*1024 {
			return nil, errors.New("file size must not exceed 5MB")
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
		contentType := file.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, errors.New("only image files are allowed (jpeg, png, gif, webp)")
		}

		// Upload new file
		newFileKey, err := s.r2Storage.UploadFile(file, "kepsek")
		if err != nil {
			return nil, err
		}

		// Delete old file from R2
		_ = s.r2Storage.DeleteFile(oldFile)

		// Update file in model
		existing.FotoKepsek = newFileKey
	}

	// Update nama if provided
	if req.NamaKepsek != nil {
		existing.NamaKepsek = *req.NamaKepsek
	}

	// Update kutipan if provided
	if req.KutipanKepsek != nil {
		existing.KutipanKepsek = *req.KutipanKepsek
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		// If DB save fails, delete the uploaded file
		if file != nil {
			_ = s.r2Storage.DeleteFile(existing.FotoKepsek)
		}
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes KutipanKepsek by ID
func (s *KutipanKepsekServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	// Delete file from R2
	if err := s.r2Storage.DeleteFile(existing.FotoKepsek); err != nil {
		// Log error but continue with database deletion
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *KutipanKepsekServiceImpl) mapToResponse(data *models.KutipanKepsek) *dtos.KutipanKepsekResponse {
	// Get public URL for the file
	publicURL := s.r2Storage.GetPublicURL(data.FotoKepsek)

	return &dtos.KutipanKepsekResponse{
		ID:            data.ID,
		NamaKepsek:    data.NamaKepsek,
		FotoKepsek:    publicURL,
		KutipanKepsek: data.KutipanKepsek,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
		CreatedByID:   data.CreatedByID,
		UpdatedByID:   data.UpdatedByID,
	}
}
