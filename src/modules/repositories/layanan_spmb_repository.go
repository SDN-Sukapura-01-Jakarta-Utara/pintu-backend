package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// LayananSPMBRepository handles data operations for Layanan SPMB
type LayananSPMBRepository interface {
	Create(data *models.LayananSPMB) error
	GetAllWithFilter(params GetLayananSPMBParams) ([]models.LayananSPMB, int64, error)
	GetByID(id uint) (*models.LayananSPMB, error)
	Update(data *models.LayananSPMB) error
	SoftDelete(id uint) error
}

// GetLayananSPMBFilter represents filter parameters
type GetLayananSPMBFilter struct {
	StartDate    string
	EndDate      string
	NamaOrangTua string
	NamaMurid    string
	Status       string
}

// GetLayananSPMBParams represents query parameters
type GetLayananSPMBParams struct {
	Filter GetLayananSPMBFilter
	Limit  int
	Offset int
}

type LayananSPMBRepositoryImpl struct {
	db *gorm.DB
}

// NewLayananSPMBRepository creates a new Layanan SPMB repository
func NewLayananSPMBRepository(db *gorm.DB) LayananSPMBRepository {
	return &LayananSPMBRepositoryImpl{db: db}
}

// Create creates a new Layanan SPMB record
func (r *LayananSPMBRepositoryImpl) Create(data *models.LayananSPMB) error {
	return r.db.Create(data).Error
}

// GetAllWithFilter retrieves all Layanan SPMB with filters, sorting, and pagination
func (r *LayananSPMBRepositoryImpl) GetAllWithFilter(params GetLayananSPMBParams) ([]models.LayananSPMB, int64, error) {
	var data []models.LayananSPMB
	var total int64

	query := r.db.Model(&models.LayananSPMB{})

	// Apply filters
	if params.Filter.StartDate != "" {
		query = query.Where("tanggal_laporan >= ?", params.Filter.StartDate)
	}
	if params.Filter.EndDate != "" {
		query = query.Where("tanggal_laporan <= ?", params.Filter.EndDate+" 23:59:59")
	}
	if params.Filter.NamaOrangTua != "" {
		query = query.Where("nama_orang_tua ILIKE ?", "%"+params.Filter.NamaOrangTua+"%")
	}
	if params.Filter.NamaMurid != "" {
		query = query.Where("nama_lengkap_murid ILIKE ?", "%"+params.Filter.NamaMurid+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting: Status (pending first) then Tanggal Laporan (DESC)
	query = query.Order(`
		CASE status
			WHEN 'pending' THEN 1
			WHEN 'selesai' THEN 2
			ELSE 3
		END ASC,
		tanggal_laporan DESC
	`)

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit).Offset(params.Offset)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByID retrieves Layanan SPMB by ID
func (r *LayananSPMBRepositoryImpl) GetByID(id uint) (*models.LayananSPMB, error) {
	var data models.LayananSPMB
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Layanan SPMB record
func (r *LayananSPMBRepositoryImpl) Update(data *models.LayananSPMB) error {
	return r.db.Save(data).Error
}

// SoftDelete soft deletes Layanan SPMB by setting deleted_at
func (r *LayananSPMBRepositoryImpl) SoftDelete(id uint) error {
	return r.db.Delete(&models.LayananSPMB{}, id).Error
}
