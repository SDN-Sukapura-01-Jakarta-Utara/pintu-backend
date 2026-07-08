package repositories

import (
	"strings"
	
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetPesertaDidikRombelFilter represents filter parameters for GetAllWithFilter
type GetPesertaDidikRombelFilter struct {
	Nama             string
	RombelID         uint
	TahunPelajaranID uint
	Status           string
}

// GetPesertaDidikRombelParams represents parameters for GetAllWithFilter with filters
type GetPesertaDidikRombelParams struct {
	Filter GetPesertaDidikRombelFilter
	Limit  int
	Offset int
}

// PesertaDidikRombelRepository handles data operations for PesertaDidikRombel
type PesertaDidikRombelRepository interface {
	Create(data *models.PesertaDidikRombel) error
	GetByID(id uint) (*models.PesertaDidikRombel, error)
	GetByPesertaDidikAndTahunPelajaran(pesertaDidikID uint, tahunPelajaranID uint) (*models.PesertaDidikRombel, error)
	GetAllWithFilter(params GetPesertaDidikRombelParams) ([]models.PesertaDidikRombel, int64, error)
	Update(data *models.PesertaDidikRombel) error
	Delete(id uint) error
	DeleteByRombelID(rombelID uint) (int64, error)
	DeleteByTahunPelajaranID(tahunPelajaranID uint) (int64, error)
	DeleteByRombelAndTahunPelajaran(rombelID uint, tahunPelajaranID uint) (int64, error)
	CheckDuplicateMapping(pesertaDidikID uint, rombelID uint, tahunPelajaranID uint) (bool, error)
	CheckDuplicateMappingExcludingID(id uint, pesertaDidikID uint, rombelID uint, tahunPelajaranID uint) (bool, error)
	GetRombelByID(id uint) (*models.Rombel, error)
	GetRombelByName(name string) (*models.Rombel, error)
	GetTahunPelajaranByID(id uint) (*models.TahunPelajaran, error)
	GetTahunPelajaranByName(name string) (*models.TahunPelajaran, error)
	GetAllRombels() ([]models.Rombel, error)
	GetAllTahunPelajaran() ([]models.TahunPelajaran, error)
	CreateWithTransaction(fn func(tx interface{}) error) error
	CreateInTransaction(tx interface{}, data *models.PesertaDidikRombel) error
}

type PesertaDidikRombelRepositoryImpl struct {
	db *gorm.DB
}

// NewPesertaDidikRombelRepository creates a new PesertaDidikRombel repository
func NewPesertaDidikRombelRepository(db *gorm.DB) PesertaDidikRombelRepository {
	return &PesertaDidikRombelRepositoryImpl{db: db}
}

// Create creates a new PesertaDidikRombel record
func (r *PesertaDidikRombelRepositoryImpl) Create(data *models.PesertaDidikRombel) error {
	return r.db.Create(data).Error
}

// GetByID retrieves PesertaDidikRombel by ID with preloaded relations
func (r *PesertaDidikRombelRepositoryImpl) GetByID(id uint) (*models.PesertaDidikRombel, error) {
	var data models.PesertaDidikRombel
	if err := r.db.Preload("PesertaDidik").Preload("Rombel.Kelas").Preload("TahunPelajaran").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates PesertaDidikRombel record (only specific fields)
func (r *PesertaDidikRombelRepositoryImpl) Update(data *models.PesertaDidikRombel) error {
	// Use Updates with Select to update only specific fields and avoid association issues
	return r.db.Model(data).
		Select("RombelID", "TahunPelajaranID", "Status", "UpdatedByID", "UpdatedAt").
		Updates(map[string]interface{}{
			"rombel_id":         data.RombelID,
			"tahun_pelajaran_id": data.TahunPelajaranID,
			"status":            data.Status,
			"updated_by_id":     data.UpdatedByID,
		}).Error
}

// Delete deletes PesertaDidikRombel record by ID
func (r *PesertaDidikRombelRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.PesertaDidikRombel{}, id).Error
}

// DeleteByRombelID deletes all PesertaDidikRombel records by rombel_id
func (r *PesertaDidikRombelRepositoryImpl) DeleteByRombelID(rombelID uint) (int64, error) {
	result := r.db.Where("rombel_id = ?", rombelID).Delete(&models.PesertaDidikRombel{})
	return result.RowsAffected, result.Error
}

// DeleteByTahunPelajaranID deletes all PesertaDidikRombel records by tahun_pelajaran_id
func (r *PesertaDidikRombelRepositoryImpl) DeleteByTahunPelajaranID(tahunPelajaranID uint) (int64, error) {
	result := r.db.Where("tahun_pelajaran_id = ?", tahunPelajaranID).Delete(&models.PesertaDidikRombel{})
	return result.RowsAffected, result.Error
}

// DeleteByRombelAndTahunPelajaran deletes all PesertaDidikRombel records by rombel_id and tahun_pelajaran_id
func (r *PesertaDidikRombelRepositoryImpl) DeleteByRombelAndTahunPelajaran(rombelID uint, tahunPelajaranID uint) (int64, error) {
	result := r.db.Where("rombel_id = ? AND tahun_pelajaran_id = ?", rombelID, tahunPelajaranID).Delete(&models.PesertaDidikRombel{})
	return result.RowsAffected, result.Error
}

// GetByPesertaDidikAndTahunPelajaran retrieves PesertaDidikRombel by peserta_didik_id and tahun_pelajaran_id
func (r *PesertaDidikRombelRepositoryImpl) GetByPesertaDidikAndTahunPelajaran(pesertaDidikID uint, tahunPelajaranID uint) (*models.PesertaDidikRombel, error) {
	var data models.PesertaDidikRombel
	if err := r.db.Where("peserta_didik_id = ? AND tahun_pelajaran_id = ?", pesertaDidikID, tahunPelajaranID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// CheckDuplicateMapping checks if mapping already exists
func (r *PesertaDidikRombelRepositoryImpl) CheckDuplicateMapping(pesertaDidikID uint, rombelID uint, tahunPelajaranID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.PesertaDidikRombel{}).
		Where("peserta_didik_id = ? AND rombel_id = ? AND tahun_pelajaran_id = ?", pesertaDidikID, rombelID, tahunPelajaranID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CheckDuplicateMappingExcludingID checks if mapping already exists excluding specific ID (for update)
func (r *PesertaDidikRombelRepositoryImpl) CheckDuplicateMappingExcludingID(id uint, pesertaDidikID uint, rombelID uint, tahunPelajaranID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.PesertaDidikRombel{}).
		Where("id != ? AND peserta_didik_id = ? AND rombel_id = ? AND tahun_pelajaran_id = ?", id, pesertaDidikID, rombelID, tahunPelajaranID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// GetRombelByID retrieves Rombel by ID
func (r *PesertaDidikRombelRepositoryImpl) GetRombelByID(id uint) (*models.Rombel, error) {
	var rombel models.Rombel
	if err := r.db.First(&rombel, id).Error; err != nil {
		return nil, err
	}
	return &rombel, nil
}

// GetRombelByName retrieves Rombel by name (case-insensitive)
func (r *PesertaDidikRombelRepositoryImpl) GetRombelByName(name string) (*models.Rombel, error) {
	var rombel models.Rombel
	if err := r.db.Where("LOWER(name) = ?", strings.ToLower(name)).First(&rombel).Error; err != nil {
		return nil, err
	}
	return &rombel, nil
}

// GetTahunPelajaranByID retrieves TahunPelajaran by ID
func (r *PesertaDidikRombelRepositoryImpl) GetTahunPelajaranByID(id uint) (*models.TahunPelajaran, error) {
	var tahunPelajaran models.TahunPelajaran
	if err := r.db.First(&tahunPelajaran, id).Error; err != nil {
		return nil, err
	}
	return &tahunPelajaran, nil
}

// GetTahunPelajaranByName retrieves TahunPelajaran by tahun_pelajaran value
func (r *PesertaDidikRombelRepositoryImpl) GetTahunPelajaranByName(name string) (*models.TahunPelajaran, error) {
	var tahunPelajaran models.TahunPelajaran
	if err := r.db.Where("tahun_pelajaran = ?", name).First(&tahunPelajaran).Error; err != nil {
		return nil, err
	}
	return &tahunPelajaran, nil
}

// GetAllRombels retrieves all Rombel records (active and inactive) ordered by name
func (r *PesertaDidikRombelRepositoryImpl) GetAllRombels() ([]models.Rombel, error) {
	var rombels []models.Rombel
	if err := r.db.Order("name ASC").Find(&rombels).Error; err != nil {
		return nil, err
	}
	return rombels, nil
}

// GetAllTahunPelajaran retrieves all TahunPelajaran records (active and inactive) ordered by tahun_pelajaran
func (r *PesertaDidikRombelRepositoryImpl) GetAllTahunPelajaran() ([]models.TahunPelajaran, error) {
	var tahunPelajaranList []models.TahunPelajaran
	if err := r.db.Order("tahun_pelajaran DESC").Find(&tahunPelajaranList).Error; err != nil {
		return nil, err
	}
	return tahunPelajaranList, nil
}

// CreateWithTransaction executes a function within a database transaction
func (r *PesertaDidikRombelRepositoryImpl) CreateWithTransaction(fn func(tx interface{}) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// CreateInTransaction creates a record within a transaction
func (r *PesertaDidikRombelRepositoryImpl) CreateInTransaction(tx interface{}, data *models.PesertaDidikRombel) error {
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return gorm.ErrInvalidTransaction
	}
	return txDB.Create(data).Error
}

// GetAllWithFilter retrieves PesertaDidikRombel records with filters and pagination
func (r *PesertaDidikRombelRepositoryImpl) GetAllWithFilter(params GetPesertaDidikRombelParams) ([]models.PesertaDidikRombel, int64, error) {
	var data []models.PesertaDidikRombel
	var total int64

	query := r.db.Model(&models.PesertaDidikRombel{})

	// Join with peserta_didik for filtering by nama
	query = query.Joins("LEFT JOIN peserta_didik ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id")

	// Apply filters with table prefix to avoid ambiguity (all filters are optional)
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(peserta_didik.nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.RombelID != 0 {
		query = query.Where("peserta_didik_rombel.rombel_id = ?", params.Filter.RombelID)
	}
	if params.Filter.TahunPelajaranID != 0 {
		query = query.Where("peserta_didik_rombel.tahun_pelajaran_id = ?", params.Filter.TahunPelajaranID)
	}
	if params.Filter.Status != "" {
		query = query.Where("peserta_didik_rombel.status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with preloaded relations
	// Sorting hierarchy:
	// 1. Tahun pelajaran (created_at DESC - terbaru dulu)
	// 2. Rombel (nama ASC - alfabetis A-Z)
	// 3. Nama siswa (nama ASC - alfabetis A-Z)
	if err := query.
		Joins("LEFT JOIN tahun_pelajaran ON tahun_pelajaran.id = peserta_didik_rombel.tahun_pelajaran_id").
		Joins("LEFT JOIN rombel ON rombel.id = peserta_didik_rombel.rombel_id").
		Preload("PesertaDidik").
		Preload("Rombel.Kelas").
		Preload("TahunPelajaran").
		Order("tahun_pelajaran.created_at DESC").
		Order("rombel.name ASC").
		Order("peserta_didik.nama ASC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
