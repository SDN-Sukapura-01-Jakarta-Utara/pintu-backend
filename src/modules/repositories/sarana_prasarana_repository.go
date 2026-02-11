package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// SaranaPrasaranaRepository handles data operations for SaranaPrasarana
type SaranaPrasaranaRepository interface {
	Create(data *models.SaranaPrasarana) error
	GetByID(id uint) (*models.SaranaPrasarana, error)
	GetAll(limit int, offset int) ([]models.SaranaPrasarana, int64, error)
	Update(data *models.SaranaPrasarana) error
	Delete(id uint) error
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

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
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
