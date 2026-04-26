package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetAnnouncementFilter represents filter parameters for GetAllWithFilter
type GetAnnouncementFilter struct {
	Judul            string
	StartDate        time.Time
	EndDate          time.Time
	Penulis          string
	StatusPublikasi  string
	Status           string
}

// GetAnnouncementParams represents parameters for GetAllWithFilter with filters
type GetAnnouncementParams struct {
	Filter GetAnnouncementFilter
	Limit  int
	Offset int
}

// AnnouncementRepository handles data operations for Announcement
type AnnouncementRepository interface {
	Create(data *models.Announcement) error
	GetByID(id uint) (*models.Announcement, error)
	GetAll(limit int, offset int) ([]models.Announcement, int64, error)
	GetAllWithFilter(params GetAnnouncementParams) ([]models.Announcement, int64, error)
	GetPublicLatest() (*models.Announcement, error)
	GetPublicNext3() ([]models.Announcement, error)
	GetPublicList(sort string, offset int) ([]models.Announcement, int64, error)
	GetPublicDetailByID(id uint) (*models.Announcement, error)
	GetPublicOtherAnnouncements(excludeID uint) ([]models.Announcement, error)
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

// GetAllWithFilter retrieves Announcement records with filters and pagination
func (r *AnnouncementRepositoryImpl) GetAllWithFilter(params GetAnnouncementParams) ([]models.Announcement, int64, error) {
	var data []models.Announcement
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
	if params.Filter.Penulis != "" {
		query = query.Where("LOWER(penulis) LIKE ?", "%"+strings.ToLower(params.Filter.Penulis)+"%")
	}
	if params.Filter.StatusPublikasi != "" {
		query = query.Where("status_publikasi = ?", params.Filter.StatusPublikasi)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
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

// GetPublicLatest retrieves the latest published and active announcement ordered by tanggal DESC
func (r *AnnouncementRepositoryImpl) GetPublicLatest() (*models.Announcement, error) {
	var data models.Announcement
	if err := r.db.Where("status = ? AND status_publikasi = ?", "active", "published").
		Order("tanggal DESC").
		First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetPublicNext3 retrieves 3 announcements (2nd to 4th latest) published and active ordered by tanggal DESC
func (r *AnnouncementRepositoryImpl) GetPublicNext3() ([]models.Announcement, error) {
	var data []models.Announcement
	if err := r.db.Where("status = ? AND status_publikasi = ?", "active", "published").
		Order("tanggal DESC").
		Offset(1). // Skip the first (latest) one
		Limit(3).  // Get next 3
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetPublicList retrieves published and active announcements with sorting and pagination (12 items per request)
func (r *AnnouncementRepositoryImpl) GetPublicList(sort string, offset int) ([]models.Announcement, int64, error) {
	var data []models.Announcement
	var total int64

	query := r.db.Where("status = ? AND status_publikasi = ?", "active", "published")

	// Get total count
	if err := query.Model(&models.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "tanggal DESC" // default: terbaru
	if sort == "terlama" {
		orderBy = "tanggal ASC"
	}

	// Get paginated data (12 items per request)
	if err := query.Order(orderBy).Limit(12).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetPublicDetailByID retrieves announcement detail by ID for public (only if active and published)
func (r *AnnouncementRepositoryImpl) GetPublicDetailByID(id uint) (*models.Announcement, error) {
	var data models.Announcement
	if err := r.db.Where("id = ? AND status = ? AND status_publikasi = ?", id, "active", "published").
		First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetPublicOtherAnnouncements retrieves 5 latest published and active announcements excluding the specified ID
func (r *AnnouncementRepositoryImpl) GetPublicOtherAnnouncements(excludeID uint) ([]models.Announcement, error) {
	var data []models.Announcement
	if err := r.db.Where("status = ? AND status_publikasi = ? AND id != ?", "active", "published", excludeID).
		Order("tanggal DESC").
		Limit(5).
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
