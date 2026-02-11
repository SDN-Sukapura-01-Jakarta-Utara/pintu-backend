package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// KutipanKepsekRepository handles data operations for KutipanKepsek
type KutipanKepsekRepository interface {
	Create(data *models.KutipanKepsek) error
	GetByID(id uint) (*models.KutipanKepsek, error)
	GetAll(limit int, offset int) ([]models.KutipanKepsek, int64, error)
	Update(data *models.KutipanKepsek) error
	Delete(id uint) error
}

type KutipanKepsekRepositoryImpl struct {
	db *gorm.DB
}

// NewKutipanKepsekRepository creates a new KutipanKepsek repository
func NewKutipanKepsekRepository(db *gorm.DB) KutipanKepsekRepository {
	return &KutipanKepsekRepositoryImpl{db: db}
}

// Create creates a new KutipanKepsek record
func (r *KutipanKepsekRepositoryImpl) Create(data *models.KutipanKepsek) error {
	return r.db.Create(data).Error
}

// GetByID retrieves KutipanKepsek by ID
func (r *KutipanKepsekRepositoryImpl) GetByID(id uint) (*models.KutipanKepsek, error) {
	var data models.KutipanKepsek
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all KutipanKepsek records with pagination
func (r *KutipanKepsekRepositoryImpl) GetAll(limit int, offset int) ([]models.KutipanKepsek, int64, error) {
	var data []models.KutipanKepsek
	var total int64

	// Get total count
	if err := r.db.Model(&models.KutipanKepsek{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates KutipanKepsek record
func (r *KutipanKepsekRepositoryImpl) Update(data *models.KutipanKepsek) error {
	return r.db.Save(data).Error
}

// Delete deletes KutipanKepsek record by ID
func (r *KutipanKepsekRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.KutipanKepsek{}, id).Error
}
