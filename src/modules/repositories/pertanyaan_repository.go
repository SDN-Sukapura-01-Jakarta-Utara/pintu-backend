package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
)

// PertanyaanRepository handles data operations for Pertanyaan
type PertanyaanRepository interface {
	Create(data *models.Pertanyaan) error
	GetByID(id uint) (*models.Pertanyaan, error)
	GetByIDTiket(idTiket string) (*models.Pertanyaan, error)
	GetAllWithFilter(params GetPertanyaanParams) ([]models.Pertanyaan, int64, error)
	GetAll() ([]models.Pertanyaan, error)
	Update(data *models.Pertanyaan) error
	Delete(id uint) error
	SoftDeleteWithUser(id uint, userID uint) error
}

// GetPertanyaanFilter represents filter parameters
type GetPertanyaanFilter struct {
	IDTiket   string
	StartDate string
	EndDate   string
	Nama      string
	Email     string
	Kategori  string
	Prioritas string
	Judul     string
	Status    string
}

// GetPertanyaanParams represents query parameters
type GetPertanyaanParams struct {
	Filter GetPertanyaanFilter
	Limit  int
	Offset int
}

type PertanyaanRepositoryImpl struct {
	db *gorm.DB
}

// NewPertanyaanRepository creates a new Pertanyaan repository
func NewPertanyaanRepository(db *gorm.DB) PertanyaanRepository {
	return &PertanyaanRepositoryImpl{db: db}
}

// Create creates a new Pertanyaan record
func (r *PertanyaanRepositoryImpl) Create(data *models.Pertanyaan) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Pertanyaan by ID
func (r *PertanyaanRepositoryImpl) GetByID(id uint) (*models.Pertanyaan, error) {
	var data models.Pertanyaan
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByIDTiket retrieves Pertanyaan by ID Tiket
func (r *PertanyaanRepositoryImpl) GetByIDTiket(idTiket string) (*models.Pertanyaan, error) {
	var data models.Pertanyaan
	if err := r.db.Where("id_tiket = ?", idTiket).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Pertanyaan records
func (r *PertanyaanRepositoryImpl) GetAll() ([]models.Pertanyaan, error) {
	var data []models.Pertanyaan
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Pertanyaan record
func (r *PertanyaanRepositoryImpl) Update(data *models.Pertanyaan) error {
	return r.db.Save(data).Error
}

// Delete deletes Pertanyaan record by ID
func (r *PertanyaanRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Pertanyaan{}, id).Error
}

// SoftDeleteWithUser soft deletes Pertanyaan and sets deleted_by_id
func (r *PertanyaanRepositoryImpl) SoftDeleteWithUser(id uint, userID uint) error {
	return r.db.Model(&models.Pertanyaan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at":    time.Now(),
		"deleted_by_id": userID,
	}).Error
}


// GetAllWithFilter retrieves all Pertanyaan with filters, sorting, and pagination
func (r *PertanyaanRepositoryImpl) GetAllWithFilter(params GetPertanyaanParams) ([]models.Pertanyaan, int64, error) {
	var data []models.Pertanyaan
	var total int64

	query := r.db.Model(&models.Pertanyaan{})

	// Apply filters
	if params.Filter.IDTiket != "" {
		query = query.Where("id_tiket ILIKE ?", "%"+params.Filter.IDTiket+"%")
	}
	if params.Filter.StartDate != "" {
		query = query.Where("tanggal_pengajuan >= ?", params.Filter.StartDate)
	}
	if params.Filter.EndDate != "" {
		query = query.Where("tanggal_pengajuan <= ?", params.Filter.EndDate+" 23:59:59")
	}
	if params.Filter.Nama != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Filter.Nama+"%")
	}
	if params.Filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+params.Filter.Email+"%")
	}
	if params.Filter.Kategori != "" {
		query = query.Where("kategori ILIKE ?", "%"+params.Filter.Kategori+"%")
	}
	if params.Filter.Prioritas != "" {
		query = query.Where("prioritas = ?", params.Filter.Prioritas)
	}
	if params.Filter.Judul != "" {
		query = query.Where("judul ILIKE ?", "%"+params.Filter.Judul+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting: Status (pending > processed > closed) then Prioritas (Tinggi > Sedang > Rendah) then Tanggal Pengajuan (DESC)
	// For closed status, ignore prioritas and only sort by tanggal_pengajuan DESC
	query = query.Order(`
		CASE status
			WHEN 'pending' THEN 1
			WHEN 'processed' THEN 2
			WHEN 'closed' THEN 3
			ELSE 4
		END ASC,
		CASE 
			WHEN status = 'closed' THEN 0
			ELSE CASE prioritas
				WHEN 'Tinggi' THEN 1
				WHEN 'Sedang' THEN 2
				WHEN 'Rendah' THEN 3
				ELSE 4
			END
		END ASC,
		tanggal_pengajuan DESC
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
