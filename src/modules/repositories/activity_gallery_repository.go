package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetActivityGalleryFilter represents filter parameters for GetAllWithFilter
type GetActivityGalleryFilter struct {
	Judul            string
	StartDate        time.Time
	EndDate          time.Time
	StatusPublikasi  string
	Status           string
}

// GetActivityGalleryParams represents parameters for GetAllWithFilter with filters
type GetActivityGalleryParams struct {
	Filter GetActivityGalleryFilter
	Limit  int
	Offset int
}

// ActivityGalleryRepository handles data operations for ActivityGallery
type ActivityGalleryRepository interface {
	Create(data *models.ActivityGallery) error
	GetByID(id uint) (*models.ActivityGallery, error)
	GetAll(limit int, offset int) ([]models.ActivityGallery, int64, error)
	GetAllWithFilter(params GetActivityGalleryParams) ([]models.ActivityGallery, int64, error)
	Update(data *models.ActivityGallery) error
	Delete(id uint) error
}

type ActivityGalleryRepositoryImpl struct {
	db *gorm.DB
}

// NewActivityGalleryRepository creates a new ActivityGallery repository
func NewActivityGalleryRepository(db *gorm.DB) ActivityGalleryRepository {
	return &ActivityGalleryRepositoryImpl{db: db}
}

// Create creates a new ActivityGallery record
func (r *ActivityGalleryRepositoryImpl) Create(data *models.ActivityGallery) error {
	return r.db.Create(data).Error
}

// GetByID retrieves ActivityGallery by ID
func (r *ActivityGalleryRepositoryImpl) GetByID(id uint) (*models.ActivityGallery, error) {
	var data models.ActivityGallery
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all ActivityGallery records with pagination
func (r *ActivityGalleryRepositoryImpl) GetAll(limit int, offset int) ([]models.ActivityGallery, int64, error) {
	var data []models.ActivityGallery
	var total int64

	// Get total count
	if err := r.db.Model(&models.ActivityGallery{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data, ordered by created_at descending
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves ActivityGallery records with filters and pagination
func (r *ActivityGalleryRepositoryImpl) GetAllWithFilter(params GetActivityGalleryParams) ([]models.ActivityGallery, int64, error) {
	var data []models.ActivityGallery
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Judul != "" {
		query = query.Where("LOWER(judul) LIKE ?", "%"+strings.ToLower(params.Filter.Judul)+"%")
	}
	if !params.Filter.StartDate.IsZero() && !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal >= ? AND tanggal <= ?", params.Filter.StartDate, params.Filter.EndDate)
	} else if !params.Filter.StartDate.IsZero() {
		query = query.Where("tanggal >= ?", params.Filter.StartDate)
	} else if !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal <= ?", params.Filter.EndDate)
	}
	if params.Filter.StatusPublikasi != "" {
		query = query.Where("status_publikasi = ?", params.Filter.StatusPublikasi)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.ActivityGallery{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates ActivityGallery record
func (r *ActivityGalleryRepositoryImpl) Update(data *models.ActivityGallery) error {
	return r.db.Save(data).Error
}

// Delete deletes ActivityGallery record by ID
func (r *ActivityGalleryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.ActivityGallery{}, id).Error
}
