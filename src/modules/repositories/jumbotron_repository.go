package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// JumbotronRepository handles data operations for Jumbotron
type JumbotronRepository interface {
	Create(data *models.Jumbotron) error
	GetByID(id uint) (*models.Jumbotron, error)
	GetAll(limit int, offset int) ([]models.Jumbotron, int64, error)
	Update(data *models.Jumbotron) error
	Delete(id uint) error
	DeleteByFile(file string) error
}

type JumbotronRepositoryImpl struct {
	db *gorm.DB
}

// NewJumbotronRepository creates a new Jumbotron repository
func NewJumbotronRepository(db *gorm.DB) JumbotronRepository {
	return &JumbotronRepositoryImpl{db: db}
}

// Create creates a new Jumbotron record
func (r *JumbotronRepositoryImpl) Create(data *models.Jumbotron) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Jumbotron by ID
func (r *JumbotronRepositoryImpl) GetByID(id uint) (*models.Jumbotron, error) {
	var data models.Jumbotron
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Jumbotron records with pagination
func (r *JumbotronRepositoryImpl) GetAll(limit int, offset int) ([]models.Jumbotron, int64, error) {
	var data []models.Jumbotron
	var total int64

	// Get total count
	if err := r.db.Model(&models.Jumbotron{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates Jumbotron record
func (r *JumbotronRepositoryImpl) Update(data *models.Jumbotron) error {
	return r.db.Save(data).Error
}

// Delete deletes Jumbotron record by ID
func (r *JumbotronRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Jumbotron{}, id).Error
}

// DeleteByFile deletes Jumbotron record by file key
func (r *JumbotronRepositoryImpl) DeleteByFile(file string) error {
	return r.db.Where("file = ?", file).Delete(&models.Jumbotron{}).Error
}
