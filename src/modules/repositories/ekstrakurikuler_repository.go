package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// EkstrakurikulerRepository handles data operations for Ekstrakurikuler
type EkstrakurikulerRepository interface {
	Create(data *models.Ekstrakurikuler) error
	GetByID(id uint) (*models.Ekstrakurikuler, error)
	GetAll(limit int, offset int) ([]models.Ekstrakurikuler, int64, error)
	GetByName(name string) (*models.Ekstrakurikuler, error)
	GetByKategori(kategori string) ([]models.Ekstrakurikuler, error)
	Update(data *models.Ekstrakurikuler) error
	Delete(id uint) error
}

type EkstrakurikulerRepositoryImpl struct {
	db *gorm.DB
}

// NewEkstrakurikulerRepository creates a new Ekstrakurikuler repository
func NewEkstrakurikulerRepository(db *gorm.DB) EkstrakurikulerRepository {
	return &EkstrakurikulerRepositoryImpl{db: db}
}

// Create creates a new Ekstrakurikuler record
func (r *EkstrakurikulerRepositoryImpl) Create(data *models.Ekstrakurikuler) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Ekstrakurikuler by ID
func (r *EkstrakurikulerRepositoryImpl) GetByID(id uint) (*models.Ekstrakurikuler, error) {
	var data models.Ekstrakurikuler
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Ekstrakurikuler records with pagination
func (r *EkstrakurikulerRepositoryImpl) GetAll(limit int, offset int) ([]models.Ekstrakurikuler, int64, error) {
	var data []models.Ekstrakurikuler
	var total int64

	// Get total count
	if err := r.db.Model(&models.Ekstrakurikuler{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves Ekstrakurikuler by name
func (r *EkstrakurikulerRepositoryImpl) GetByName(name string) (*models.Ekstrakurikuler, error) {
	var data models.Ekstrakurikuler
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByKategori retrieves all Ekstrakurikuler by kategori
func (r *EkstrakurikulerRepositoryImpl) GetByKategori(kategori string) ([]models.Ekstrakurikuler, error) {
	var data []models.Ekstrakurikuler
	if err := r.db.Where("kategori = ?", kategori).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Ekstrakurikuler record
func (r *EkstrakurikulerRepositoryImpl) Update(data *models.Ekstrakurikuler) error {
	return r.db.Save(data).Error
}

// Delete deletes Ekstrakurikuler record by ID
func (r *EkstrakurikulerRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Ekstrakurikuler{}, id).Error
}
