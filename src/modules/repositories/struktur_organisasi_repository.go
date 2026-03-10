package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// GetStrukturOrganisasiFilter represents filter parameters for GetAllWithFilter
type GetStrukturOrganisasiFilter struct {
	Nama           string
	Urutan         int
	Relasi         string
	Jabatan        string
	Status         string
}

// GetStrukturOrganisasiParams represents parameters for GetAllWithFilter with filters
type GetStrukturOrganisasiParams struct {
	Filter GetStrukturOrganisasiFilter
	Limit  int
	Offset int
}

// StrukturOrganisasiRepository handles data operations for StrukturOrganisasi
type StrukturOrganisasiRepository interface {
	Create(data *models.StrukturOrganisasi) error
	GetByID(id uint) (*models.StrukturOrganisasi, error)
	GetAll(limit int, offset int) ([]models.StrukturOrganisasi, int64, error)
	GetAllWithFilter(params GetStrukturOrganisasiParams) ([]models.StrukturOrganisasi, int64, error)
	Update(data *models.StrukturOrganisasi) error
	Delete(id uint) error
}

type StrukturOrganisasiRepositoryImpl struct {
	db *gorm.DB
}

// NewStrukturOrganisasiRepository creates a new StrukturOrganisasi repository
func NewStrukturOrganisasiRepository(db *gorm.DB) StrukturOrganisasiRepository {
	return &StrukturOrganisasiRepositoryImpl{db: db}
}

// Create creates a new StrukturOrganisasi record
func (r *StrukturOrganisasiRepositoryImpl) Create(data *models.StrukturOrganisasi) error {
	return r.db.Create(data).Error
}

// GetByID retrieves StrukturOrganisasi by ID
func (r *StrukturOrganisasiRepositoryImpl) GetByID(id uint) (*models.StrukturOrganisasi, error) {
	var data models.StrukturOrganisasi
	if err := r.db.Preload("Pegawai").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all StrukturOrganisasi records with pagination
func (r *StrukturOrganisasiRepositoryImpl) GetAll(limit int, offset int) ([]models.StrukturOrganisasi, int64, error) {
	var data []models.StrukturOrganisasi
	var total int64

	// Get total count
	if err := r.db.Model(&models.StrukturOrganisasi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC (newest first)
	if err := r.db.Preload("Pegawai").Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves StrukturOrganisasi records with filters and pagination
func (r *StrukturOrganisasiRepositoryImpl) GetAllWithFilter(params GetStrukturOrganisasiParams) ([]models.StrukturOrganisasi, int64, error) {
	var data []models.StrukturOrganisasi
	var total int64

	query := r.db

	// Check if we need to JOIN kepegawaian table (for nama or jabatan filters)
	needsJoin := params.Filter.Nama != "" || params.Filter.Jabatan != ""
	if needsJoin {
		query = query.Joins("LEFT JOIN kepegawaian ON struktur_organisasi.pegawai_id = kepegawaian.id")
	}

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(kepegawaian.nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.Urutan > 0 {
		query = query.Where("struktur_organisasi.urutan = ?", params.Filter.Urutan)
	}
	if params.Filter.Relasi != "" {
		query = query.Where("LOWER(struktur_organisasi.relasi) LIKE ?", "%"+strings.ToLower(params.Filter.Relasi)+"%")
	}
	if params.Filter.Jabatan != "" {
		query = query.Where("LOWER(kepegawaian.jabatan) LIKE ?", "%"+strings.ToLower(params.Filter.Jabatan)+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("struktur_organisasi.status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.StrukturOrganisasi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Preload relationships for data fetching
	query = query.Preload("Pegawai")

	// Get paginated data ordered by created_at DESC (newest first)
	if err := query.Order("struktur_organisasi.created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates StrukturOrganisasi record
func (r *StrukturOrganisasiRepositoryImpl) Update(data *models.StrukturOrganisasi) error {
	return r.db.Model(data).
		Select("pegawai_id", "nama_non_pegawai", "jabatan_non_pegawai", "urutan", "relasi", "status", "updated_by_id", "updated_at").
		Save(data).Error
}

// Delete deletes StrukturOrganisasi record by ID
func (r *StrukturOrganisasiRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.StrukturOrganisasi{}, id).Error
}
