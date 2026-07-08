package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// KonfigurasiAbsensiRepository defines the interface for Konfigurasi Absensi repository
type KonfigurasiAbsensiRepository interface {
	GetByID(id uint) (*models.KonfigurasiAbsensi, error)
	Create(data *models.KonfigurasiAbsensi) error
	Update(data *models.KonfigurasiAbsensi) error
}

type KonfigurasiAbsensiRepositoryImpl struct {
	db *gorm.DB
}

// NewKonfigurasiAbsensiRepository creates a new Konfigurasi Absensi repository
func NewKonfigurasiAbsensiRepository(db *gorm.DB) KonfigurasiAbsensiRepository {
	return &KonfigurasiAbsensiRepositoryImpl{db: db}
}

// GetByID retrieves konfigurasi absensi by ID
func (r *KonfigurasiAbsensiRepositoryImpl) GetByID(id uint) (*models.KonfigurasiAbsensi, error) {
	var data models.KonfigurasiAbsensi
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Create creates a new konfigurasi absensi record
func (r *KonfigurasiAbsensiRepositoryImpl) Create(data *models.KonfigurasiAbsensi) error {
	return r.db.Create(data).Error
}

// Update updates an existing konfigurasi absensi record
func (r *KonfigurasiAbsensiRepositoryImpl) Update(data *models.KonfigurasiAbsensi) error {
	return r.db.Save(data).Error
}
