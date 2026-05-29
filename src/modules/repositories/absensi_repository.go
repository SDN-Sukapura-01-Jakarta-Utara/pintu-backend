package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
)

// AbsensiRepository handles data operations for Absensi
type AbsensiRepository interface {
	Create(data *models.Absensi) error
	GetByID(id uint) (*models.Absensi, error)
	GetByPesertaDidikTanggalMapel(pesertaDidikID uint, tanggal time.Time, bidangStudiID *uint) (*models.Absensi, error)
	CheckPertemuanExists(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, semester int, bulan int, tahun int, pertemuanKe int) (*models.Absensi, error)
	BulkCreate(dataList []models.Absensi) error
	GetRekapAbsensi(tahunPelajaranID, rombelID uint, semester, bulan, tahun *int, tanggalMulai, tanggalSelesai *time.Time, bidangStudiID *uint) ([]models.Absensi, error)
	Update(data *models.Absensi) error
}

type AbsensiRepositoryImpl struct {
	db *gorm.DB
}

// NewAbsensiRepository creates a new Absensi repository
func NewAbsensiRepository(db *gorm.DB) AbsensiRepository {
	return &AbsensiRepositoryImpl{db: db}
}

// Create creates a new Absensi record
func (r *AbsensiRepositoryImpl) Create(data *models.Absensi) error {
	return r.db.Create(data).Error
}

// GetByPesertaDidikTanggalMapel retrieves Absensi by peserta_didik_id, tanggal, and bidang_studi_id
func (r *AbsensiRepositoryImpl) GetByPesertaDidikTanggalMapel(pesertaDidikID uint, tanggal time.Time, bidangStudiID *uint) (*models.Absensi, error) {
	var data models.Absensi
	query := r.db.Where("peserta_didik_id = ? AND tanggal = ?", pesertaDidikID, tanggal)
	
	// Handle NULL vs NOT NULL for bidang_studi_id
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// BulkCreate creates multiple Absensi records in a single transaction
func (r *AbsensiRepositoryImpl) BulkCreate(dataList []models.Absensi) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, data := range dataList {
			if err := tx.Create(&data).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRekapAbsensi retrieves attendance records for recap with filters
func (r *AbsensiRepositoryImpl) GetRekapAbsensi(tahunPelajaranID, rombelID uint, semester, bulan, tahun *int, tanggalMulai, tanggalSelesai *time.Time, bidangStudiID *uint) ([]models.Absensi, error) {
	var data []models.Absensi
	
	query := r.db.Preload("PesertaDidik").Preload("Rombel").Preload("BidangStudi").Preload("DicatatOleh").
		Where("tahun_pelajaran_id = ? AND rombel_id = ?", tahunPelajaranID, rombelID)
	
	// Filter by semester if provided
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	// Filter by tanggal range if provided
	if tanggalMulai != nil && tanggalSelesai != nil {
		query = query.Where("tanggal BETWEEN ? AND ?", *tanggalMulai, *tanggalSelesai)
	} else if tanggalMulai != nil {
		query = query.Where("tanggal >= ?", *tanggalMulai)
	} else if tanggalSelesai != nil {
		query = query.Where("tanggal <= ?", *tanggalSelesai)
	}
	
	// Filter by bulan and tahun if provided (only if tanggal range not provided)
	if tanggalMulai == nil && tanggalSelesai == nil {
		if bulan != nil && tahun != nil {
			query = query.Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", *bulan, *tahun)
		} else if tahun != nil {
			query = query.Where("EXTRACT(YEAR FROM tanggal) = ?", *tahun)
		}
	}
	
	// Filter by bidang_studi_id
	// - If bidangStudiID is nil: Filter for NULL (guru kelas only)
	// - If bidangStudiID is not nil: Filter for specific bidang_studi_id (guru mapel)
	if bidangStudiID == nil {
		// Filter for NULL (guru kelas only)
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		// Filter for specific bidang_studi_id (guru mapel)
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	// Order by tanggal and peserta_didik_id
	query = query.Order("tanggal ASC, peserta_didik_id ASC")
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetByID retrieves Absensi by ID
func (r *AbsensiRepositoryImpl) GetByID(id uint) (*models.Absensi, error) {
	var data models.Absensi
	if err := r.db.Preload("PesertaDidik").Preload("Rombel").Preload("TahunPelajaran").Preload("BidangStudi").Preload("DicatatOleh").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates an Absensi record
func (r *AbsensiRepositoryImpl) Update(data *models.Absensi) error {
	return r.db.Save(data).Error
}

// CheckPertemuanExists checks if a pertemuan already exists in a specific month for guru mapel
// Returns the existing absensi record with preloaded Rombel and BidangStudi if found
func (r *AbsensiRepositoryImpl) CheckPertemuanExists(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, semester int, bulan int, tahun int, pertemuanKe int) (*models.Absensi, error) {
	var data models.Absensi
	
	err := r.db.Preload("Rombel").Preload("BidangStudi").
		Where("rombel_id = ?", rombelID).
		Where("bidang_studi_id = ?", bidangStudiID).
		Where("tahun_pelajaran_id = ?", tahunPelajaranID).
		Where("semester = ?", semester).
		Where("pertemuan_ke = ?", pertemuanKe).
		Where("EXTRACT(MONTH FROM tanggal) = ?", bulan).
		Where("EXTRACT(YEAR FROM tanggal) = ?", tahun).
		First(&data).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found, no error
		}
		return nil, err
	}
	
	return &data, nil
}
