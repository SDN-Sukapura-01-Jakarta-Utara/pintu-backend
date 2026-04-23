package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// GetSaranaPrasaranaFilter represents filter parameters for GetAllWithFilter
type GetSaranaPrasaranaFilter struct {
	Name   string
	Status string
}

// GetSaranaPrasaranaParams represents parameters for GetAllWithFilter with filters
type GetSaranaPrasaranaParams struct {
	Filter GetSaranaPrasaranaFilter
	Limit  int
	Offset int
}

// SaranaPrasaranaRepository handles data operations for SaranaPrasarana
type SaranaPrasaranaRepository interface {
	Create(data *models.SaranaPrasarana) error
	GetByID(id uint) (*models.SaranaPrasarana, error)
	GetAll(limit int, offset int) ([]models.SaranaPrasarana, int64, error)
	GetAllWithFilter(params GetSaranaPrasaranaParams) ([]models.SaranaPrasarana, int64, error)
	Update(data *models.SaranaPrasarana) error
	Delete(id uint) error
	GetAllPublic() ([]models.SaranaPrasarana, error)
}

type SaranaPrasaranaRepositoryImpl struct {
	db *gorm.DB
}

// NewSaranaPrasaranaRepository creates a new SaranaPrasarana repository
func NewSaranaPrasaranaRepository(db *gorm.DB) SaranaPrasaranaRepository {
	return &SaranaPrasaranaRepositoryImpl{db: db}
}

// Create creates a new SaranaPrasarana record
func (r *SaranaPrasaranaRepositoryImpl) Create(data *models.SaranaPrasarana) error {
	return r.db.Create(data).Error
}

// GetByID retrieves SaranaPrasarana by ID
func (r *SaranaPrasaranaRepositoryImpl) GetByID(id uint) (*models.SaranaPrasarana, error) {
	var data models.SaranaPrasarana
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all SaranaPrasarana records with pagination
func (r *SaranaPrasaranaRepositoryImpl) GetAll(limit int, offset int) ([]models.SaranaPrasarana, int64, error) {
	var data []models.SaranaPrasarana
	var total int64

	// Get total count
	if err := r.db.Model(&models.SaranaPrasarana{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves SaranaPrasarana records with filters and pagination
func (r *SaranaPrasaranaRepositoryImpl) GetAllWithFilter(params GetSaranaPrasaranaParams) ([]models.SaranaPrasarana, int64, error) {
	var data []models.SaranaPrasarana
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Filter.Name)+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.SaranaPrasarana{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates SaranaPrasarana record
func (r *SaranaPrasaranaRepositoryImpl) Update(data *models.SaranaPrasarana) error {
	return r.db.Save(data).Error
}

// Delete deletes SaranaPrasarana record by ID
func (r *SaranaPrasaranaRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.SaranaPrasarana{}, id).Error
}

// GetAllPublic retrieves all active SaranaPrasarana records for public display
func (r *SaranaPrasaranaRepositoryImpl) GetAllPublic() ([]models.SaranaPrasarana, error) {
	var data []models.SaranaPrasarana
	if err := r.db.Where("status = ?", "active").Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
