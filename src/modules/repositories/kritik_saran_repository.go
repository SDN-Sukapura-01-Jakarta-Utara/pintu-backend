package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
)

// GetKritikSaranFilter represents filter parameters for GetAllWithFilter
type GetKritikSaranFilter struct {
	StartDate time.Time
	EndDate   time.Time
}

// GetKritikSaranParams represents parameters for GetAllWithFilter with filters
type GetKritikSaranParams struct {
	Filter GetKritikSaranFilter
	Limit  int
	Offset int
}

// KritikSaranRepository handles data operations for KritikSaran
type KritikSaranRepository interface {
	Create(data *models.KritikSaran) error
	GetByID(id uint) (*models.KritikSaran, error)
	GetAll(limit int, offset int) ([]models.KritikSaran, int64, error)
	GetAllWithFilter(params GetKritikSaranParams) ([]models.KritikSaran, int64, error)
	Update(data *models.KritikSaran) error
	Delete(id uint) error
}

type KritikSaranRepositoryImpl struct {
	db *gorm.DB
}

// NewKritikSaranRepository creates a new KritikSaran repository
func NewKritikSaranRepository(db *gorm.DB) KritikSaranRepository {
	return &KritikSaranRepositoryImpl{db: db}
}

// Create creates a new KritikSaran record
func (r *KritikSaranRepositoryImpl) Create(data *models.KritikSaran) error {
	return r.db.Create(data).Error
}

// GetByID retrieves KritikSaran by ID
func (r *KritikSaranRepositoryImpl) GetByID(id uint) (*models.KritikSaran, error) {
	var data models.KritikSaran
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all KritikSaran records with pagination
func (r *KritikSaranRepositoryImpl) GetAll(limit int, offset int) ([]models.KritikSaran, int64, error) {
	var data []models.KritikSaran
	var total int64

	// Get total count
	if err := r.db.Model(&models.KritikSaran{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data, ordered by created_at descending
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves KritikSaran records with filters and pagination
func (r *KritikSaranRepositoryImpl) GetAllWithFilter(params GetKritikSaranParams) ([]models.KritikSaran, int64, error) {
	var data []models.KritikSaran
	var total int64

	query := r.db

	// Apply date filters on created_at
	if !params.Filter.StartDate.IsZero() && !params.Filter.EndDate.IsZero() {
		query = query.Where("created_at >= ? AND created_at <= ?", params.Filter.StartDate, params.Filter.EndDate)
	} else if !params.Filter.StartDate.IsZero() {
		query = query.Where("created_at >= ?", params.Filter.StartDate)
	} else if !params.Filter.EndDate.IsZero() {
		query = query.Where("created_at <= ?", params.Filter.EndDate)
	}

	// Get total count
	if err := query.Model(&models.KritikSaran{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates KritikSaran record
func (r *KritikSaranRepositoryImpl) Update(data *models.KritikSaran) error {
	return r.db.Save(data).Error
}

// Delete deletes KritikSaran record by ID
func (r *KritikSaranRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.KritikSaran{}, id).Error
}
