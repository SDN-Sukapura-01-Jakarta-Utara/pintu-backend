package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetPrestasiFilter represents filter parameters for GetAllWithFilter
type GetPrestasiFilter struct {
	PesertaDidikID    *uint
	NamaPesertaDidik  string
	Jenis             string
	NamaGrup          string
	NamaPrestasi      string
	TingkatPrestasi   string
	Penyelenggara     string
	StartDate         time.Time
	EndDate           time.Time
	Juara             string
	EkstrakurikulerID *uint
	TahunPelajaranID  *uint
	Status            string
}

// GetPrestasiParams represents parameters for GetAllWithFilter with filters
type GetPrestasiParams struct {
	Filter GetPrestasiFilter
	Limit  int
	Offset int
}

// PrestasiRepository handles data operations for Prestasi
type PrestasiRepository interface {
	Create(data *models.Prestasi) error
	GetByID(id uint) (*models.Prestasi, error)
	GetAll(limit int, offset int) ([]models.Prestasi, int64, error)
	GetAllWithFilter(params GetPrestasiParams) ([]models.Prestasi, int64, error)
	GetPublicLatest() ([]models.Prestasi, error)
	GetPublicList(sort string, offset int) ([]models.Prestasi, int64, error)
	GetPublicDetailByID(id uint) (*models.Prestasi, error)
	Update(data *models.Prestasi) error
	Delete(id uint) error
	// Anggota Tim methods
	CreateAnggotaTim(data *models.AnggotaTimPrestasi) error
	GetAnggotaTimByPrestasiID(prestasiID uint) ([]models.AnggotaTimPrestasi, error)
	UpdateAnggotaTim(data *models.AnggotaTimPrestasi) error
	DeleteAnggotaTim(id uint) error
	DeleteAnggotaTimByPrestasiID(prestasiID uint) error
}

type PrestasiRepositoryImpl struct {
	db *gorm.DB
}

// NewPrestasiRepository creates a new Prestasi repository
func NewPrestasiRepository(db *gorm.DB) PrestasiRepository {
	return &PrestasiRepositoryImpl{db: db}
}

// Create creates a new Prestasi record
func (r *PrestasiRepositoryImpl) Create(data *models.Prestasi) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Prestasi by ID with all relationships
func (r *PrestasiRepositoryImpl) GetByID(id uint) (*models.Prestasi, error) {
	var data models.Prestasi
	if err := r.db.Preload("PesertaDidik").
		Preload("PesertaDidik.Rombel").
		Preload("PesertaDidik.Rombel.Kelas").
		Preload("PesertaDidik.TahunPelajaran").
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel.Kelas").
		Preload("AnggotaTimPrestasi.PesertaDidik.TahunPelajaran").
		Preload("AnggotaTimPrestasi.TahunPelajaran").
		First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Prestasi records with pagination
func (r *PrestasiRepositoryImpl) GetAll(limit int, offset int) ([]models.Prestasi, int64, error) {
	var data []models.Prestasi
	var total int64

	// Get total count
	if err := r.db.Model(&models.Prestasi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with relationships
	if err := r.db.Preload("PesertaDidik").
		Preload("PesertaDidik.Rombel").
		Preload("PesertaDidik.Rombel.Kelas").
		Preload("PesertaDidik.TahunPelajaran").
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel.Kelas").
		Preload("AnggotaTimPrestasi.PesertaDidik.TahunPelajaran").
		Preload("AnggotaTimPrestasi.TahunPelajaran").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves Prestasi records with filters and pagination
func (r *PrestasiRepositoryImpl) GetAllWithFilter(params GetPrestasiParams) ([]models.Prestasi, int64, error) {
	var data []models.Prestasi
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.PesertaDidikID != nil {
		query = query.Where("(peserta_didik_id = ? OR id IN (SELECT prestasi_id FROM anggota_tim_prestasi WHERE peserta_didik_id = ? AND deleted_at IS NULL))", *params.Filter.PesertaDidikID, *params.Filter.PesertaDidikID)
	}
	if params.Filter.NamaPesertaDidik != "" {
		nameLike := "%" + strings.ToLower(params.Filter.NamaPesertaDidik) + "%"
		query = query.Where("(peserta_didik_id IN (SELECT id FROM peserta_didik WHERE LOWER(nama) LIKE ? AND deleted_at IS NULL) OR id IN (SELECT prestasi_id FROM anggota_tim_prestasi WHERE peserta_didik_id IN (SELECT id FROM peserta_didik WHERE LOWER(nama) LIKE ? AND deleted_at IS NULL) AND deleted_at IS NULL))", nameLike, nameLike)
	}
	if params.Filter.Jenis != "" {
		query = query.Where("LOWER(jenis) LIKE ?", "%"+strings.ToLower(params.Filter.Jenis)+"%")
	}
	if params.Filter.NamaGrup != "" {
		query = query.Where("LOWER(nama_grup) LIKE ?", "%"+strings.ToLower(params.Filter.NamaGrup)+"%")
	}
	if params.Filter.NamaPrestasi != "" {
		query = query.Where("LOWER(nama_prestasi) LIKE ?", "%"+strings.ToLower(params.Filter.NamaPrestasi)+"%")
	}
	if params.Filter.TingkatPrestasi != "" {
		query = query.Where("LOWER(tingkat_prestasi) LIKE ?", "%"+strings.ToLower(params.Filter.TingkatPrestasi)+"%")
	}
	if params.Filter.Penyelenggara != "" {
		query = query.Where("LOWER(penyelenggara) LIKE ?", "%"+strings.ToLower(params.Filter.Penyelenggara)+"%")
	}
	if !params.Filter.StartDate.IsZero() && !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal_lomba >= ? AND tanggal_lomba <= ?", params.Filter.StartDate, params.Filter.EndDate)
	} else if !params.Filter.StartDate.IsZero() {
		query = query.Where("tanggal_lomba >= ?", params.Filter.StartDate)
	} else if !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal_lomba <= ?", params.Filter.EndDate)
	}
	if params.Filter.Juara != "" {
		query = query.Where("LOWER(juara) LIKE ?", "%"+strings.ToLower(params.Filter.Juara)+"%")
	}
	if params.Filter.EkstrakurikulerID != nil {
		query = query.Where("ekstrakurikuler_id = ?", *params.Filter.EkstrakurikulerID)
	}
	if params.Filter.TahunPelajaranID != nil {
		query = query.Where("tahun_pelajaran_id = ?", *params.Filter.TahunPelajaranID)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.Prestasi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with relationships
	if err := query.Preload("PesertaDidik").
		Preload("PesertaDidik.Rombel").
		Preload("PesertaDidik.Rombel.Kelas").
		Preload("PesertaDidik.TahunPelajaran").
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel.Kelas").
		Preload("AnggotaTimPrestasi.PesertaDidik.TahunPelajaran").
		Preload("AnggotaTimPrestasi.TahunPelajaran").
		Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates Prestasi record
func (r *PrestasiRepositoryImpl) Update(data *models.Prestasi) error {
	return r.db.Save(data).Error
}

// Delete deletes Prestasi record by ID
func (r *PrestasiRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Prestasi{}, id).Error
}

// CreateAnggotaTim creates a new AnggotaTimPrestasi record
func (r *PrestasiRepositoryImpl) CreateAnggotaTim(data *models.AnggotaTimPrestasi) error {
	return r.db.Create(data).Error
}

// GetAnggotaTimByPrestasiID retrieves all anggota tim by prestasi ID
func (r *PrestasiRepositoryImpl) GetAnggotaTimByPrestasiID(prestasiID uint) ([]models.AnggotaTimPrestasi, error) {
	var data []models.AnggotaTimPrestasi
	if err := r.db.Preload("PesertaDidik").
		Preload("PesertaDidik.Rombel").
		Preload("PesertaDidik.Rombel.Kelas").
		Preload("PesertaDidik.TahunPelajaran").
		Preload("TahunPelajaran").
		Where("prestasi_id = ?", prestasiID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateAnggotaTim updates AnggotaTimPrestasi record
func (r *PrestasiRepositoryImpl) UpdateAnggotaTim(data *models.AnggotaTimPrestasi) error {
	return r.db.Save(data).Error
}

// DeleteAnggotaTim deletes AnggotaTimPrestasi record by ID
func (r *PrestasiRepositoryImpl) DeleteAnggotaTim(id uint) error {
	return r.db.Delete(&models.AnggotaTimPrestasi{}, id).Error
}

// DeleteAnggotaTimByPrestasiID deletes all AnggotaTimPrestasi records by prestasi ID
func (r *PrestasiRepositoryImpl) DeleteAnggotaTimByPrestasiID(prestasiID uint) error {
	return r.db.Where("prestasi_id = ?", prestasiID).Delete(&models.AnggotaTimPrestasi{}).Error
}

// GetPublicLatest retrieves 10 latest prestasi ordered by tanggal_lomba DESC
func (r *PrestasiRepositoryImpl) GetPublicLatest() ([]models.Prestasi, error) {
	var data []models.Prestasi
	if err := r.db.Where("status = ?", "active").
		Preload("PesertaDidik").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Order("tanggal_lomba DESC").
		Limit(10).
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}


// GetPublicList retrieves active prestasi with sorting and pagination (12 items per request)
func (r *PrestasiRepositoryImpl) GetPublicList(sort string, offset int) ([]models.Prestasi, int64, error) {
	var data []models.Prestasi
	var total int64

	query := r.db.Where("status = ?", "active")

	// Get total count
	if err := query.Model(&models.Prestasi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "tanggal_lomba DESC" // default: terbaru
	if sort == "terlama" {
		orderBy = "tanggal_lomba ASC"
	}

	// Get paginated data (12 items per request)
	if err := query.Preload("PesertaDidik").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Order(orderBy).
		Limit(12).
		Offset(offset).
		Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}


// GetPublicDetailByID retrieves prestasi detail by ID for public (only if active)
func (r *PrestasiRepositoryImpl) GetPublicDetailByID(id uint) (*models.Prestasi, error) {
	var data models.Prestasi
	if err := r.db.Where("id = ? AND status = ?", id, "active").
		Preload("PesertaDidik").
		Preload("PesertaDidik.Rombel").
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AnggotaTimPrestasi").
		Preload("AnggotaTimPrestasi.PesertaDidik").
		Preload("AnggotaTimPrestasi.PesertaDidik.Rombel").
		Preload("AnggotaTimPrestasi.TahunPelajaran").
		First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
