package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// TahunPelajaranRepository handles data operations for TahunPelajaran
type TahunPelajaranRepository interface {
	Create(data *models.TahunPelajaran) error
	GetByID(id uint) (*models.TahunPelajaran, error)
	GetAll() ([]models.TahunPelajaran, error)
	Update(data *models.TahunPelajaran) error
	Delete(id uint) error
}

type TahunPelajaranRepositoryImpl struct {
	db *gorm.DB
}

// NewTahunPelajaranRepository creates a new TahunPelajaran repository
func NewTahunPelajaranRepository(db *gorm.DB) TahunPelajaranRepository {
	return &TahunPelajaranRepositoryImpl{db: db}
}

// Create creates a new TahunPelajaran record
func (r *TahunPelajaranRepositoryImpl) Create(data *models.TahunPelajaran) error {
	return r.db.Create(data).Error
}

// GetByID retrieves TahunPelajaran by ID
func (r *TahunPelajaranRepositoryImpl) GetByID(id uint) (*models.TahunPelajaran, error) {
	var data models.TahunPelajaran
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all TahunPelajaran records
func (r *TahunPelajaranRepositoryImpl) GetAll() ([]models.TahunPelajaran, error) {
	var data []models.TahunPelajaran
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates TahunPelajaran record
func (r *TahunPelajaranRepositoryImpl) Update(data *models.TahunPelajaran) error {
	return r.db.Save(data).Error
}

// Delete deletes TahunPelajaran record by ID
func (r *TahunPelajaranRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.TahunPelajaran{}, id).Error
}
