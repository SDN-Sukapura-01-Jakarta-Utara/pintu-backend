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

type ActivityGalleryService interface {
	Create(fotos []*multipart.FileHeader, req *dtos.ActivityGalleryCreateRequest, userID uint) (*dtos.ActivityGalleryResponse, error)
	GetByID(id uint) (*dtos.ActivityGalleryResponse, error)
	GetAll(limit int, offset int) (*dtos.ActivityGalleryListResponse, error)
	Update(id uint, fotos []*multipart.FileHeader, req *dtos.ActivityGalleryUpdateRequest, userID uint) (*dtos.ActivityGalleryResponse, error)
	Delete(id uint) error
}

type ActivityGalleryServiceImpl struct {
	repository repositories.ActivityGalleryRepository
	r2Storage  *utils.R2Storage
}

// NewActivityGalleryService creates a new ActivityGallery service
func NewActivityGalleryService(repository repositories.ActivityGalleryRepository, r2Storage *utils.R2Storage) ActivityGalleryService {
	return &ActivityGalleryServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new ActivityGallery with foto uploads to R2
func (s *ActivityGalleryServiceImpl) Create(fotos []*multipart.FileHeader, req *dtos.ActivityGalleryCreateRequest, userID uint) (*dtos.ActivityGalleryResponse, error) {
	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("invalid tanggal format (use YYYY-MM-DD)")
	}

	// Upload fotos if provided
	var fotoItems []models.FileItem
	if len(fotos) > 0 {
		for _, foto := range fotos {
			if foto == nil {
				continue
			}

			// Validate foto file size (max 10MB per file)
			if foto.Size > 10*1024*1024 { // 10MB
				return nil, errors.New("each foto must not exceed 10MB")
			}

			// Validate foto type - allow image files only
			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}
			contentType := foto.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("only image files are allowed for foto (jpeg, png, gif, webp)")
			}

			// Upload foto to R2 in galeri-kegiatan directory
			fileKey, err := s.r2Storage.UploadFile(foto, "galeri-kegiatan")
			if err != nil {
				return nil, err
			}

			// Generate unique file ID: file_timestamp_randomstring
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			fotoItems = append(fotoItems, models.FileItem{
				ID:       fileID,
				Filename: foto.Filename,
				URL:      fileKey,
				Size:     foto.Size,
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

	// Convert fotoItems to JSON
	fotosJSON, _ := json.Marshal(fotoItems)

	// Create activity gallery record
	data := &models.ActivityGallery{
		Judul:           req.Judul,
		Tanggal:         tanggal,
		Foto:            fotosJSON,
		StatusPublikasi: statusPublikasi,
		Status:          status,
		CreatedByID:     &userID,
	}

	if err := s.repository.Create(data); err != nil {
		// If database save fails, delete uploaded files
		for _, item := range fotoItems {
			_ = s.r2Storage.DeleteFile(item.URL)
		}
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves ActivityGallery by ID
func (s *ActivityGalleryServiceImpl) GetByID(id uint) (*dtos.ActivityGalleryResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all ActivityGalleries
func (s *ActivityGalleryServiceImpl) GetAll(limit int, offset int) (*dtos.ActivityGalleryListResponse, error) {
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
	responses := make([]dtos.ActivityGalleryResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.ActivityGalleryListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// Update updates ActivityGallery
func (s *ActivityGalleryServiceImpl) Update(id uint, fotos []*multipart.FileHeader, req *dtos.ActivityGalleryUpdateRequest, userID uint) (*dtos.ActivityGalleryResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("activity gallery not found")
	}

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
	if req.StatusPublikasi != "" {
		existing.StatusPublikasi = req.StatusPublikasi
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	// Delete fotos if specified
	if len(req.FotoToDelete) > 0 {
		var existingFotoItems []models.FileItem
		if err := json.Unmarshal(existing.Foto, &existingFotoItems); err == nil {
			// Build map of file IDs to delete
			deleteMap := make(map[string]bool)
			for _, fileID := range req.FotoToDelete {
				deleteMap[fileID] = true
			}

			// Filter out fotos to delete and delete from R2
			var remainingFotos []models.FileItem
			for _, foto := range existingFotoItems {
				if deleteMap[foto.ID] {
					// Delete from R2
					_ = s.r2Storage.DeleteFile(foto.URL)
				} else {
					remainingFotos = append(remainingFotos, foto)
				}
			}

			// Update fotos array
			fotosJSON, _ := json.Marshal(remainingFotos)
			existing.Foto = fotosJSON
		}
	}

	// Add new fotos if provided (fotos lama tetap)
	if len(fotos) > 0 {
		var existingFotoItems []models.FileItem
		// Get existing fotos
		if err := json.Unmarshal(existing.Foto, &existingFotoItems); err != nil && existing.Foto != nil && len(existing.Foto) > 0 {
			// Only log, don't fail
		}

		// Upload and add new fotos
		for _, foto := range fotos {
			if foto == nil {
				continue
			}

			// Validate foto file size (max 10MB per file)
			if foto.Size > 10*1024*1024 { // 10MB
				return nil, errors.New("each foto must not exceed 10MB")
			}

			// Validate foto type
			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}
			contentType := foto.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("only image files are allowed for foto (jpeg, png, gif, webp)")
			}

			// Upload foto to R2
			fileKey, err := s.r2Storage.UploadFile(foto, "galeri-kegiatan")
			if err != nil {
				return nil, err
			}

			// Generate unique file ID
			fileID := fmt.Sprintf("file_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			existingFotoItems = append(existingFotoItems, models.FileItem{
				ID:       fileID,
				Filename: foto.Filename,
				URL:      fileKey,
				Size:     foto.Size,
			})
		}

		// Convert updated fotos to JSON
		fotosJSON, _ := json.Marshal(existingFotoItems)
		existing.Foto = fotosJSON
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		// If DB save fails, delete the uploaded fotos
		if len(fotos) > 0 {
			for _, foto := range fotos {
				_ = s.r2Storage.DeleteFile(fmt.Sprintf("galeri-kegiatan/%s", foto.Filename))
			}
		}
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes ActivityGallery by ID
func (s *ActivityGalleryServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("activity gallery not found")
	}

	// Delete all fotos from R2
	var fotoItems []models.FileItem
	if err := json.Unmarshal(existing.Foto, &fotoItems); err == nil {
		for _, foto := range fotoItems {
			_ = s.r2Storage.DeleteFile(foto.URL)
		}
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *ActivityGalleryServiceImpl) mapToResponse(data *models.ActivityGallery) *dtos.ActivityGalleryResponse {
	// Map fotos from JSON
	var fotoItems []dtos.FileItemDTO
	var fotoModels []models.FileItem
	if err := json.Unmarshal(data.Foto, &fotoModels); err == nil {
		for _, foto := range fotoModels {
			fotoItems = append(fotoItems, dtos.FileItemDTO{
				ID:       foto.ID,
				Filename: foto.Filename,
				URL:      s.r2Storage.GetPublicURL(foto.URL),
				Size:     foto.Size,
			})
		}
	}

	return &dtos.ActivityGalleryResponse{
		ID:              data.ID,
		Judul:           data.Judul,
		Tanggal:         data.Tanggal,
		Foto:            fotoItems,
		StatusPublikasi: data.StatusPublikasi,
		Status:          data.Status,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
		CreatedByID:     data.CreatedByID,
		UpdatedByID:     data.UpdatedByID,
	}
}
