package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// RombelRepository handles data operations for Rombel
type RombelRepository interface {
	Create(data *models.Rombel) error
	GetByID(id uint) (*models.Rombel, error)
	GetAll(limit int, offset int) ([]models.Rombel, int64, error)
	GetByName(name string) (*models.Rombel, error)
	GetByKelasID(kelasID uint) ([]models.Rombel, error)
	Update(data *models.Rombel) error
	Delete(id uint) error
}

type RombelRepositoryImpl struct {
	db *gorm.DB
}

// NewRombelRepository creates a new Rombel repository
func NewRombelRepository(db *gorm.DB) RombelRepository {
	return &RombelRepositoryImpl{db: db}
}

// Create creates a new Rombel record
func (r *RombelRepositoryImpl) Create(data *models.Rombel) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Rombel by ID
func (r *RombelRepositoryImpl) GetByID(id uint) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Rombel records with pagination
func (r *RombelRepositoryImpl) GetAll(limit int, offset int) ([]models.Rombel, int64, error) {
	var data []models.Rombel
	var total int64

	// Get total count
	if err := r.db.Model(&models.Rombel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves Rombel by name
func (r *RombelRepositoryImpl) GetByName(name string) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByKelasID retrieves all Rombel by Kelas ID
func (r *RombelRepositoryImpl) GetByKelasID(kelasID uint) ([]models.Rombel, error) {
	var data []models.Rombel
	if err := r.db.Where("kelas_id = ?", kelasID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Rombel record
func (r *RombelRepositoryImpl) Update(data *models.Rombel) error {
	return r.db.Save(data).Error
}

// Delete deletes Rombel record by ID
func (r *RombelRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Rombel{}, id).Error
}
