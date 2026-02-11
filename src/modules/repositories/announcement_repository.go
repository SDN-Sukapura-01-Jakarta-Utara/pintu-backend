package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// AnnouncementRepository handles data operations for Announcement
type AnnouncementRepository interface {
	Create(data *models.Announcement) error
	GetByID(id uint) (*models.Announcement, error)
	GetAll(limit int, offset int) ([]models.Announcement, int64, error)
	Update(data *models.Announcement) error
	Delete(id uint) error
	DeleteByGambar(gambar string) error
}

type AnnouncementRepositoryImpl struct {
	db *gorm.DB
}

// NewAnnouncementRepository creates a new Announcement repository
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &AnnouncementRepositoryImpl{db: db}
}

// Create creates a new Announcement record
func (r *AnnouncementRepositoryImpl) Create(data *models.Announcement) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Announcement by ID
func (r *AnnouncementRepositoryImpl) GetByID(id uint) (*models.Announcement, error) {
	var data models.Announcement
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Announcement records with pagination
func (r *AnnouncementRepositoryImpl) GetAll(limit int, offset int) ([]models.Announcement, int64, error) {
	var data []models.Announcement
	var total int64

	// Get total count
	if err := r.db.Model(&models.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data, ordered by created_at descending
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates Announcement record
func (r *AnnouncementRepositoryImpl) Update(data *models.Announcement) error {
	return r.db.Save(data).Error
}

// Delete deletes Announcement record by ID
func (r *AnnouncementRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Announcement{}, id).Error
}

// DeleteByGambar deletes Announcement record by gambar key
func (r *AnnouncementRepositoryImpl) DeleteByGambar(gambar string) error {
	return r.db.Where("gambar = ?", gambar).Delete(&models.Announcement{}).Error
}
