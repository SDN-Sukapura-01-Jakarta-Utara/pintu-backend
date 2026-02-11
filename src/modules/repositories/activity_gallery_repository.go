package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// ActivityGalleryRepository handles data operations for ActivityGallery
type ActivityGalleryRepository interface {
	Create(data *models.ActivityGallery) error
	GetByID(id uint) (*models.ActivityGallery, error)
	GetAll(limit int, offset int) ([]models.ActivityGallery, int64, error)
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

// Update updates ActivityGallery record
func (r *ActivityGalleryRepositoryImpl) Update(data *models.ActivityGallery) error {
	return r.db.Save(data).Error
}

// Delete deletes ActivityGallery record by ID
func (r *ActivityGalleryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.ActivityGallery{}, id).Error
}
