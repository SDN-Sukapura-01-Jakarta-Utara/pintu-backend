package services

import (
	"errors"
	"mime/multipart"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// JumbotronService handles business logic for Jumbotron
type JumbotronService interface {
	Create(file *multipart.FileHeader, req *dtos.JumbotronCreateRequest, userID uint) (*dtos.JumbotronResponse, error)
	GetByID(id uint) (*dtos.JumbotronResponse, error)
	GetAll(limit int, offset int) (*dtos.JumbotronListResponse, error)
	Update(id uint, req *dtos.JumbotronUpdateRequest, userID uint) (*dtos.JumbotronResponse, error)
	Delete(id uint) error
}

type JumbotronServiceImpl struct {
	repository repositories.JumbotronRepository
	r2Storage  *utils.R2Storage
}

// NewJumbotronService creates a new Jumbotron service
func NewJumbotronService(repository repositories.JumbotronRepository, r2Storage *utils.R2Storage) JumbotronService {
	return &JumbotronServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new Jumbotron with file upload to R2
func (s *JumbotronServiceImpl) Create(file *multipart.FileHeader, req *dtos.JumbotronCreateRequest, userID uint) (*dtos.JumbotronResponse, error) {
	// Validate file
	if file == nil {
		return nil, errors.New("file is required")
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
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

	// Upload file to R2
	fileKey, err := s.r2Storage.UploadFile(file, "jumbotron")
	if err != nil {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create jumbotron record
	data := &models.Jumbotron{
		File:        fileKey,
		Status:      status,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		// If database save fails, delete the uploaded file
		_ = s.r2Storage.DeleteFile(fileKey)
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Jumbotron by ID
func (s *JumbotronServiceImpl) GetByID(id uint) (*dtos.JumbotronResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Jumbotron
func (s *JumbotronServiceImpl) GetAll(limit int, offset int) (*dtos.JumbotronListResponse, error) {
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
	responses := make([]dtos.JumbotronResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.JumbotronListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// Update updates Jumbotron
func (s *JumbotronServiceImpl) Update(id uint, req *dtos.JumbotronUpdateRequest, userID uint) (*dtos.JumbotronResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Jumbotron by ID
func (s *JumbotronServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	// Delete file from R2
	if err := s.r2Storage.DeleteFile(existing.File); err != nil {
		// Log error but continue with database deletion
		// In production, you might want to handle this differently
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *JumbotronServiceImpl) mapToResponse(data *models.Jumbotron) *dtos.JumbotronResponse {
	// Get public URL for the file
	publicURL := s.r2Storage.GetPublicURL(data.File)

	return &dtos.JumbotronResponse{
		ID:          data.ID,
		File:        publicURL,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
