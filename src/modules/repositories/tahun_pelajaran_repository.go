package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetTahunPelajaranFilter represents filter parameters for GetAll
type GetTahunPelajaranFilter struct {
	TahunPelajaran string
	Status         string
}

// GetTahunPelajaranParams represents parameters for GetAll with filters
type GetTahunPelajaranParams struct {
	Filter GetTahunPelajaranFilter
	Limit  int
	Offset int
}

// TahunPelajaranRepository handles data operations for TahunPelajaran
type TahunPelajaranRepository interface {
	Create(data *models.TahunPelajaran) error
	GetByID(id uint) (*models.TahunPelajaran, error)
	GetAll(limit int, offset int) ([]models.TahunPelajaran, int64, error)
	GetAllWithFilter(params GetTahunPelajaranParams) ([]models.TahunPelajaran, int64, error)
	GetByTahunPelajaran(tahunPelajaran string) (*models.TahunPelajaran, error)
	Update(data *models.TahunPelajaran) error
	Delete(id uint) error
	UpdateAllStatusToInactive() error
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

// GetAll retrieves all TahunPelajaran records with pagination
func (r *TahunPelajaranRepositoryImpl) GetAll(limit int, offset int) ([]models.TahunPelajaran, int64, error) {
	var data []models.TahunPelajaran
	var total int64

	// Get total count
	if err := r.db.Model(&models.TahunPelajaran{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByTahunPelajaran retrieves TahunPelajaran by tahun_pelajaran
func (r *TahunPelajaranRepositoryImpl) GetByTahunPelajaran(tahunPelajaran string) (*models.TahunPelajaran, error) {
	var data models.TahunPelajaran
	if err := r.db.Where("tahun_pelajaran = ?", tahunPelajaran).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates TahunPelajaran record
func (r *TahunPelajaranRepositoryImpl) Update(data *models.TahunPelajaran) error {
	return r.db.Save(data).Error
}

// Delete deletes TahunPelajaran record by ID
func (r *TahunPelajaranRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.TahunPelajaran{}, id).Error
}

// GetAllWithFilter retrieves TahunPelajaran records with filters and pagination
func (r *TahunPelajaranRepositoryImpl) GetAllWithFilter(params GetTahunPelajaranParams) ([]models.TahunPelajaran, int64, error) {
	var data []models.TahunPelajaran
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.TahunPelajaran != "" {
		query = query.Where("tahun_pelajaran ILIKE ?", "%"+params.Filter.TahunPelajaran+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.TahunPelajaran{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// UpdateAllStatusToInactive updates all tahun_pelajaran status to inactive
func (r *TahunPelajaranRepositoryImpl) UpdateAllStatusToInactive() error {
	return r.db.Model(&models.TahunPelajaran{}).Where("status = ?", "active").Update("status", "inactive").Error
}
