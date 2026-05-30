package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// KelulusanRepository handles data operations for Kelulusan
type KelulusanRepository interface {
	Create(data *models.Kelulusan) error
	GetByID(id uint) (*models.Kelulusan, error)
	GetByNomorPeserta(nomorPeserta string) (*models.Kelulusan, error)
	GetByNISNAndTanggalLahir(nisn string, tanggalLahir string) (*models.Kelulusan, error)
	GetAllWithFilter(params GetKelulusanParams) ([]models.Kelulusan, int64, error)
	Update(data *models.Kelulusan) error
	Delete(id uint) error
}

// GetKelulusanFilter represents filter parameters for GetAllWithFilter
type GetKelulusanFilter struct {
	Nama         string
	NomorPeserta string
	NISN         string
	Lulus        *bool // nil = all, true = lulus, false = tidak lulus
}

// GetKelulusanParams represents parameters for GetAllWithFilter with filters
type GetKelulusanParams struct {
	Filter GetKelulusanFilter
	Limit  int
	Offset int
}

type KelulusanRepositoryImpl struct {
	db *gorm.DB
}

// NewKelulusanRepository creates a new Kelulusan repository
func NewKelulusanRepository(db *gorm.DB) KelulusanRepository {
	return &KelulusanRepositoryImpl{db: db}
}

// Create creates a new Kelulusan record
func (r *KelulusanRepositoryImpl) Create(data *models.Kelulusan) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Kelulusan by ID
func (r *KelulusanRepositoryImpl) GetByID(id uint) (*models.Kelulusan, error) {
	var data models.Kelulusan
	if err := r.db.Preload("CreatedBy").Preload("UpdatedBy").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByNomorPeserta retrieves Kelulusan by nomor peserta
func (r *KelulusanRepositoryImpl) GetByNomorPeserta(nomorPeserta string) (*models.Kelulusan, error) {
	var data models.Kelulusan
	if err := r.db.Where("nomor_peserta = ?", nomorPeserta).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByNISNAndTanggalLahir retrieves Kelulusan by NISN and tanggal lahir
func (r *KelulusanRepositoryImpl) GetByNISNAndTanggalLahir(nisn string, tanggalLahir string) (*models.Kelulusan, error) {
	var data models.Kelulusan
	if err := r.db.Where("nisn = ? AND DATE(tanggal_lahir) = ?", nisn, tanggalLahir).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}


// GetAllWithFilter retrieves Kelulusan with filters and pagination
func (r *KelulusanRepositoryImpl) GetAllWithFilter(params GetKelulusanParams) ([]models.Kelulusan, int64, error) {
	var data []models.Kelulusan
	var total int64

	query := r.db.Model(&models.Kelulusan{})

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.NomorPeserta != "" {
		query = query.Where("LOWER(nomor_peserta) LIKE ?", "%"+strings.ToLower(params.Filter.NomorPeserta)+"%")
	}
	if params.Filter.NISN != "" {
		query = query.Where("LOWER(nisn) LIKE ?", "%"+strings.ToLower(params.Filter.NISN)+"%")
	}
	if params.Filter.Lulus != nil {
		query = query.Where("lulus = ?", *params.Filter.Lulus)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get data
	if err := query.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}


// Update updates a Kelulusan record
func (r *KelulusanRepositoryImpl) Update(data *models.Kelulusan) error {
	return r.db.Save(data).Error
}


// Delete soft deletes a Kelulusan record
func (r *KelulusanRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Kelulusan{}, id).Error
}
