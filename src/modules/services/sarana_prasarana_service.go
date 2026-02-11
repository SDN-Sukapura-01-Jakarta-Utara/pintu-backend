package services

import (
	"errors"
	"mime/multipart"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

// SaranaPrasaranaService handles business logic for SaranaPrasarana
type SaranaPrasaranaService interface {
	Create(file *multipart.FileHeader, req *dtos.SaranaPrasaranaCreateRequest, userID uint) (*dtos.SaranaPrasaranaResponse, error)
	GetByID(id uint) (*dtos.SaranaPrasaranaResponse, error)
	GetAll(limit int, offset int) (*dtos.SaranaPrasaranaListResponse, error)
	UpdateWithFile(id uint, file *multipart.FileHeader, req *dtos.SaranaPrasaranaUpdateRequest, userID uint) (*dtos.SaranaPrasaranaResponse, error)
	Delete(id uint) error
}

type SaranaPrasaranaServiceImpl struct {
	repository repositories.SaranaPrasaranaRepository
	r2Storage  *utils.R2Storage
}

// NewSaranaPrasaranaService creates a new SaranaPrasarana service
func NewSaranaPrasaranaService(repository repositories.SaranaPrasaranaRepository, r2Storage *utils.R2Storage) SaranaPrasaranaService {
	return &SaranaPrasaranaServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new SaranaPrasarana with file upload
func (s *SaranaPrasaranaServiceImpl) Create(file *multipart.FileHeader, req *dtos.SaranaPrasaranaCreateRequest, userID uint) (*dtos.SaranaPrasaranaResponse, error) {
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

	// Upload file to R2 (sarpras directory)
	fileKey, err := s.r2Storage.UploadFile(file, "sarpras")
	if err != nil {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.SaranaPrasarana{
		Name:        req.Name,
		Foto:        fileKey,
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

// GetByID retrieves SaranaPrasarana by ID
func (s *SaranaPrasaranaServiceImpl) GetByID(id uint) (*dtos.SaranaPrasaranaResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all SaranaPrasarana
func (s *SaranaPrasaranaServiceImpl) GetAll(limit int, offset int) (*dtos.SaranaPrasaranaListResponse, error) {
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
	responses := make([]dtos.SaranaPrasaranaResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.SaranaPrasaranaListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// UpdateWithFile updates SaranaPrasarana with optional file upload
func (s *SaranaPrasaranaServiceImpl) UpdateWithFile(id uint, file *multipart.FileHeader, req *dtos.SaranaPrasaranaUpdateRequest, userID uint) (*dtos.SaranaPrasaranaResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldFile := existing.Foto

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
		newFileKey, err := s.r2Storage.UploadFile(file, "sarpras")
		if err != nil {
			return nil, err
		}

		// Delete old file from R2
		_ = s.r2Storage.DeleteFile(oldFile)

		// Update file in model
		existing.Foto = newFileKey
	}

	// Update name if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}

	// Update status if provided
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		// If DB save fails, delete the uploaded file
		if file != nil {
			_ = s.r2Storage.DeleteFile(existing.Foto)
		}
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes SaranaPrasarana by ID
func (s *SaranaPrasaranaServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	// Delete file from R2
	if err := s.r2Storage.DeleteFile(existing.Foto); err != nil {
		// Log error but continue with database deletion
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *SaranaPrasaranaServiceImpl) mapToResponse(data *models.SaranaPrasarana) *dtos.SaranaPrasaranaResponse {
	// Get public URL for the file
	publicURL := s.r2Storage.GetPublicURL(data.Foto)

	return &dtos.SaranaPrasaranaResponse{
		ID:          data.ID,
		Name:        data.Name,
		Foto:        publicURL,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}
