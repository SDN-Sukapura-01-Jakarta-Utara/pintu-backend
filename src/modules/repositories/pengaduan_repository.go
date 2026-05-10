package repositories

import (
	"time"

	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// PengaduanRepository handles data operations for Pengaduan
type PengaduanRepository interface {
	Create(data *models.Pengaduan) error
	GetByID(id uint) (*models.Pengaduan, error)
	GetByIDTiket(idTiket string) (*models.Pengaduan, error)
	GetAllWithFilter(params GetPengaduanParams) ([]models.Pengaduan, int64, error)
	Update(data *models.Pengaduan) error
	SoftDeleteWithUser(id uint, userID uint) error
}

// GetPengaduanFilter represents filter parameters
type GetPengaduanFilter struct {
	IDTiket     string
	StartDate   string
	EndDate     string
	TipePelapor string
	Nama        string
	Email       string
	Kategori    string
	Prioritas   string
	Judul       string
	Status      string
}

// GetPengaduanParams represents query parameters
type GetPengaduanParams struct {
	Filter GetPengaduanFilter
	Limit  int
	Offset int
}

type PengaduanRepositoryImpl struct {
	db *gorm.DB
}

// NewPengaduanRepository creates a new Pengaduan repository
func NewPengaduanRepository(db *gorm.DB) PengaduanRepository {
	return &PengaduanRepositoryImpl{db: db}
}

// Create creates a new Pengaduan record
func (r *PengaduanRepositoryImpl) Create(data *models.Pengaduan) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Pengaduan by ID
func (r *PengaduanRepositoryImpl) GetByID(id uint) (*models.Pengaduan, error) {
	var data models.Pengaduan
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByIDTiket retrieves Pengaduan by ID Tiket
func (r *PengaduanRepositoryImpl) GetByIDTiket(idTiket string) (*models.Pengaduan, error) {
	var data models.Pengaduan
	if err := r.db.Where("id_tiket = ?", idTiket).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Pengaduan record
func (r *PengaduanRepositoryImpl) Update(data *models.Pengaduan) error {
	return r.db.Save(data).Error
}

// SoftDeleteWithUser soft deletes Pengaduan and sets deleted_by_id
func (r *PengaduanRepositoryImpl) SoftDeleteWithUser(id uint, userID uint) error {
	return r.db.Model(&models.Pengaduan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at":    time.Now(),
		"deleted_by_id": userID,
	}).Error
}

// GetAllWithFilter retrieves all Pengaduan with filters, sorting, and pagination
func (r *PengaduanRepositoryImpl) GetAllWithFilter(params GetPengaduanParams) ([]models.Pengaduan, int64, error) {
	var data []models.Pengaduan
	var total int64

	query := r.db.Model(&models.Pengaduan{})

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
	if params.Filter.TipePelapor != "" {
		query = query.Where("tipe_pelapor = ?", params.Filter.TipePelapor)
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