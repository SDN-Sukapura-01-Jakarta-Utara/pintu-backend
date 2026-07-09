package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// KonfigurasiMutasiSiswaRepository handles data operations for Konfigurasi Mutasi Siswa
type KonfigurasiMutasiSiswaRepository interface {
	GetByID(id uint) (*models.KonfigurasiMutasiSiswa, error)
	Create(data *models.KonfigurasiMutasiSiswa) error
	Update(data *models.KonfigurasiMutasiSiswa) error
}

type KonfigurasiMutasiSiswaRepositoryImpl struct {
	db *gorm.DB
}

// NewKonfigurasiMutasiSiswaRepository creates a new Konfigurasi Mutasi Siswa repository
func NewKonfigurasiMutasiSiswaRepository(db *gorm.DB) KonfigurasiMutasiSiswaRepository {
	return &KonfigurasiMutasiSiswaRepositoryImpl{db: db}
}

// GetByID retrieves Konfigurasi Mutasi Siswa by ID
func (r *KonfigurasiMutasiSiswaRepositoryImpl) GetByID(id uint) (*models.KonfigurasiMutasiSiswa, error) {
	var data models.KonfigurasiMutasiSiswa
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Create creates a new Konfigurasi Mutasi Siswa record
func (r *KonfigurasiMutasiSiswaRepositoryImpl) Create(data *models.KonfigurasiMutasiSiswa) error {
	return r.db.Create(data).Error
}

// Update updates Konfigurasi Mutasi Siswa record
func (r *KonfigurasiMutasiSiswaRepositoryImpl) Update(data *models.KonfigurasiMutasiSiswa) error {
	return r.db.Save(data).Error
}
