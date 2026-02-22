package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetBidangStudiFilter represents filter parameters for GetAll
type GetBidangStudiFilter struct {
	Name   string
	Status string
}

// GetBidangStudiParams represents parameters for GetAll with filters
type GetBidangStudiParams struct {
	Filter GetBidangStudiFilter
	Limit  int
	Offset int
}

// BidangStudiRepository handles data operations for BidangStudi
type BidangStudiRepository interface {
	Create(data *models.BidangStudi) error
	GetByID(id uint) (*models.BidangStudi, error)
	GetAll(limit int, offset int) ([]models.BidangStudi, int64, error)
	GetAllWithFilter(params GetBidangStudiParams) ([]models.BidangStudi, int64, error)
	GetByName(name string) (*models.BidangStudi, error)
	Update(data *models.BidangStudi) error
	Delete(id uint) error
}

type BidangStudiRepositoryImpl struct {
	db *gorm.DB
}

// NewBidangStudiRepository creates a new BidangStudi repository
func NewBidangStudiRepository(db *gorm.DB) BidangStudiRepository {
	return &BidangStudiRepositoryImpl{db: db}
}

// Create creates a new BidangStudi record
func (r *BidangStudiRepositoryImpl) Create(data *models.BidangStudi) error {
	return r.db.Create(data).Error
}

// GetByID retrieves BidangStudi by ID
func (r *BidangStudiRepositoryImpl) GetByID(id uint) (*models.BidangStudi, error) {
	var data models.BidangStudi
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all BidangStudi records with pagination
func (r *BidangStudiRepositoryImpl) GetAll(limit int, offset int) ([]models.BidangStudi, int64, error) {
	var data []models.BidangStudi
	var total int64

	// Get total count
	if err := r.db.Model(&models.BidangStudi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves BidangStudi by name
func (r *BidangStudiRepositoryImpl) GetByName(name string) (*models.BidangStudi, error) {
	var data models.BidangStudi
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates BidangStudi record
func (r *BidangStudiRepositoryImpl) Update(data *models.BidangStudi) error {
	return r.db.Save(data).Error
}

// Delete deletes BidangStudi record by ID
func (r *BidangStudiRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.BidangStudi{}, id).Error
}

// GetAllWithFilter retrieves BidangStudi records with filters and pagination
func (r *BidangStudiRepositoryImpl) GetAllWithFilter(params GetBidangStudiParams) ([]models.BidangStudi, int64, error) {
	var data []models.BidangStudi
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Filter.Name+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.BidangStudi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
