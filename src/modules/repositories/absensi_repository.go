package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
)

// AbsensiRepository handles data operations for Absensi
type AbsensiRepository interface {
	Create(data *models.RekapitulasiAbsensi) error
	GetByID(id uint) (*models.RekapitulasiAbsensi, error)
	GetPesertaDidikRombelID(pesertaDidikID, rombelID uint) (uint, error)
	GetPesertaDidikRombelByID(id uint) (*models.PesertaDidikRombel, error)
	CheckDuplicateGuruKelas(rombelID, tahunPelajaranID uint, semester int, tanggal time.Time) (*models.RekapitulasiAbsensi, error)
	CheckDuplicateGuruMapel(rombelID, tahunPelajaranID, bidangStudiID uint, semester, pertemuanKe int, bulan, tahun int) (*models.RekapitulasiAbsensi, error)
	GetByPesertaDidikTanggalMapel(pesertaDidikRombelID uint, tanggal time.Time, bidangStudiID *uint) (*models.RekapitulasiAbsensi, error)
	CheckPertemuanExists(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, semester int, bulan int, tahun int, pertemuanKe int) (*models.RekapitulasiAbsensi, error)
	GetPertemuanTanggal(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, bulan int, tahun int, pertemuanKe int) (*time.Time, error)
	BulkCreate(dataList []models.RekapitulasiAbsensi) error
	GetRekapAbsensi(tahunPelajaranID, rombelID uint, semester, bulan, tahun *int, tanggalMulai, tanggalSelesai *time.Time, bidangStudiID *uint) ([]models.RekapitulasiAbsensi, error)
	Update(data *models.RekapitulasiAbsensi) error
	GetDashboardSummary(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error)
	CountUniqueSiswa(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint) (int, error)
	GetPerbandinganRombel(tahunPelajaranID uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error)
	GetSiswaTerendah(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error)
	GetDashboardSiswa(pesertaDidikRombelID, tahunPelajaranID, rombelID uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error)
	GetAbsensiScanByDate(tanggal time.Time) ([]models.Absensi, error)
	GetAbsensiScanByMonth(bulan, tahun int) ([]models.Absensi, error)
	GetRekapByRombelTanggalBidangStudi(rombelID, tahunPelajaranID uint, tanggal time.Time, bidangStudiID *uint) ([]models.RekapitulasiAbsensi, error)
}

type AbsensiRepositoryImpl struct {
	db *gorm.DB
}

// NewAbsensiRepository creates a new Absensi repository
func NewAbsensiRepository(db *gorm.DB) AbsensiRepository {
	return &AbsensiRepositoryImpl{db: db}
}

// Create creates a new Absensi record
func (r *AbsensiRepositoryImpl) Create(data *models.RekapitulasiAbsensi) error {
	return r.db.Create(data).Error
}

// GetPesertaDidikRombelID gets peserta_didik_rombel_id by peserta_didik_id and rombel_id
func (r *AbsensiRepositoryImpl) GetPesertaDidikRombelID(pesertaDidikID, rombelID uint) (uint, error) {
	var pesertaDidikRombel models.PesertaDidikRombel
	err := r.db.Where("peserta_didik_id = ? AND rombel_id = ?", pesertaDidikID, rombelID).
		First(&pesertaDidikRombel).Error
	if err != nil {
		return 0, err
	}
	return pesertaDidikRombel.ID, nil
}

// GetPesertaDidikRombelByID gets peserta_didik_rombel by ID
func (r *AbsensiRepositoryImpl) GetPesertaDidikRombelByID(id uint) (*models.PesertaDidikRombel, error) {
	var pesertaDidikRombel models.PesertaDidikRombel
	err := r.db.First(&pesertaDidikRombel, id).Error
	if err != nil {
		return nil, err
	}
	return &pesertaDidikRombel, nil
}

// CheckDuplicateGuruKelas checks if absensi for guru kelas already exists on this date
func (r *AbsensiRepositoryImpl) CheckDuplicateGuruKelas(rombelID, tahunPelajaranID uint, semester int, tanggal time.Time) (*models.RekapitulasiAbsensi, error) {
	var data models.RekapitulasiAbsensi
	err := r.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND semester = ? AND tanggal = ? AND bidang_studi_id IS NULL", 
		rombelID, tahunPelajaranID, semester, tanggal).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// CheckDuplicateGuruMapel checks if absensi for guru mapel already exists with same pertemuan_ke in the same month
func (r *AbsensiRepositoryImpl) CheckDuplicateGuruMapel(rombelID, tahunPelajaranID, bidangStudiID uint, semester, pertemuanKe int, bulan, tahun int) (*models.RekapitulasiAbsensi, error) {
	var data models.RekapitulasiAbsensi
	err := r.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id = ? AND semester = ? AND pertemuan_ke = ? AND EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?",
		rombelID, tahunPelajaranID, bidangStudiID, semester, pertemuanKe, bulan, tahun).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByPesertaDidikTanggalMapel retrieves Absensi by peserta_didik_rombel_id, tanggal, and bidang_studi_id
func (r *AbsensiRepositoryImpl) GetByPesertaDidikTanggalMapel(pesertaDidikRombelID uint, tanggal time.Time, bidangStudiID *uint) (*models.RekapitulasiAbsensi, error) {
	var data models.RekapitulasiAbsensi
	query := r.db.Where("peserta_didik_rombel_id = ? AND tanggal = ?", pesertaDidikRombelID, tanggal)
	
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
func (r *AbsensiRepositoryImpl) BulkCreate(dataList []models.RekapitulasiAbsensi) error {
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
func (r *AbsensiRepositoryImpl) GetRekapAbsensi(tahunPelajaranID, rombelID uint, semester, bulan, tahun *int, tanggalMulai, tanggalSelesai *time.Time, bidangStudiID *uint) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Preload("PesertaDidikRombel.PesertaDidik").Preload("Rombel").Preload("BidangStudi").Preload("DicatatOleh").
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
	
	// Order by tanggal and peserta_didik_rombel_id
	query = query.Order("tanggal ASC, peserta_didik_rombel_id ASC")
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetByID retrieves Absensi by ID
func (r *AbsensiRepositoryImpl) GetByID(id uint) (*models.RekapitulasiAbsensi, error) {
	var data models.RekapitulasiAbsensi
	if err := r.db.Preload("PesertaDidikRombel.PesertaDidik").Preload("Rombel").Preload("TahunPelajaran").Preload("BidangStudi").Preload("DicatatOleh").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates an Absensi record
func (r *AbsensiRepositoryImpl) Update(data *models.RekapitulasiAbsensi) error {
	return r.db.Save(data).Error
}

// CheckPertemuanExists checks if a pertemuan already exists in a specific month for guru mapel
// Returns the existing absensi record with preloaded Rombel and BidangStudi if found
func (r *AbsensiRepositoryImpl) CheckPertemuanExists(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, semester int, bulan int, tahun int, pertemuanKe int) (*models.RekapitulasiAbsensi, error) {
	var data models.RekapitulasiAbsensi
	
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

// GetPertemuanTanggal gets the tanggal for a specific pertemuan_ke in a month for guru bidang studi
// Returns the date if exists, nil if not found
func (r *AbsensiRepositoryImpl) GetPertemuanTanggal(rombelID uint, bidangStudiID uint, tahunPelajaranID uint, bulan int, tahun int, pertemuanKe int) (*time.Time, error) {
	var data models.RekapitulasiAbsensi
	
	err := r.db.Select("tanggal").
		Where("rombel_id = ?", rombelID).
		Where("bidang_studi_id = ?", bidangStudiID).
		Where("tahun_pelajaran_id = ?", tahunPelajaranID).
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
	
	return &data.Tanggal, nil
}

// GetDashboardSummary retrieves attendance records for dashboard summary
func (r *AbsensiRepositoryImpl) GetDashboardSummary(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Where("tahun_pelajaran_id = ?", tahunPelajaranID)
	
	// Filter by rombel_id if provided
	if rombelID != nil {
		query = query.Where("rombel_id = ?", *rombelID)
	}
	
	// Filter by semester if provided
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	// Filter by bidang_studi_id
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	// Filter by tanggal range if provided
	if tanggalMulai != nil && tanggalSelesai != nil {
		query = query.Where("tanggal BETWEEN ? AND ?", *tanggalMulai, *tanggalSelesai)
	} else if tanggalMulai != nil {
		query = query.Where("tanggal >= ?", *tanggalMulai)
	} else if tanggalSelesai != nil {
		query = query.Where("tanggal <= ?", *tanggalSelesai)
	}
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// CountUniqueSiswa counts unique students in the filtered data
func (r *AbsensiRepositoryImpl) CountUniqueSiswa(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint) (int, error) {
	var count int64
	
	query := r.db.Model(&models.RekapitulasiAbsensi{}).
		Distinct("peserta_didik_rombel_id").
		Where("tahun_pelajaran_id = ?", tahunPelajaranID)
	
	if rombelID != nil {
		query = query.Where("rombel_id = ?", *rombelID)
	}
	
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	
	return int(count), nil
}

// GetPerbandinganRombel retrieves attendance records for rombel comparison (with Rombel preloaded)
func (r *AbsensiRepositoryImpl) GetPerbandinganRombel(tahunPelajaranID uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Preload("PesertaDidikRombel.PesertaDidik").Preload("Rombel").Where("tahun_pelajaran_id = ?", tahunPelajaranID)
	
	// Filter by semester if provided
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	// Filter by bidang_studi_id
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	// Filter by tanggal range if provided
	if tanggalMulai != nil && tanggalSelesai != nil {
		query = query.Where("tanggal BETWEEN ? AND ?", *tanggalMulai, *tanggalSelesai)
	} else if tanggalMulai != nil {
		query = query.Where("tanggal >= ?", *tanggalMulai)
	} else if tanggalSelesai != nil {
		query = query.Where("tanggal <= ?", *tanggalSelesai)
	}
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetSiswaTerendah retrieves attendance records for students with lowest attendance (with PesertaDidik preloaded)
func (r *AbsensiRepositoryImpl) GetSiswaTerendah(tahunPelajaranID uint, rombelID *uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Preload("PesertaDidikRombel.PesertaDidik").Where("tahun_pelajaran_id = ?", tahunPelajaranID)
	
	// Filter by rombel_id if provided
	if rombelID != nil {
		query = query.Where("rombel_id = ?", *rombelID)
	}
	
	// Filter by semester if provided
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	// Filter by bidang_studi_id
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	// Filter by tanggal range if provided
	if tanggalMulai != nil && tanggalSelesai != nil {
		query = query.Where("tanggal BETWEEN ? AND ?", *tanggalMulai, *tanggalSelesai)
	} else if tanggalMulai != nil {
		query = query.Where("tanggal >= ?", *tanggalMulai)
	} else if tanggalSelesai != nil {
		query = query.Where("tanggal <= ?", *tanggalSelesai)
	}
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetDashboardSiswa retrieves attendance records for a specific student (with PesertaDidik and Rombel preloaded)
func (r *AbsensiRepositoryImpl) GetDashboardSiswa(pesertaDidikRombelID, tahunPelajaranID, rombelID uint, semester *int, bidangStudiID *uint, tanggalMulai, tanggalSelesai *time.Time) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Preload("PesertaDidikRombel.PesertaDidik").Preload("Rombel").
		Where("peserta_didik_rombel_id = ?", pesertaDidikRombelID).
		Where("tahun_pelajaran_id = ?", tahunPelajaranID).
		Where("rombel_id = ?", rombelID)
	
	// Filter by semester if provided
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}
	
	// Filter by bidang_studi_id
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	// Filter by tanggal range if provided
	if tanggalMulai != nil && tanggalSelesai != nil {
		query = query.Where("tanggal BETWEEN ? AND ?", *tanggalMulai, *tanggalSelesai)
	} else if tanggalMulai != nil {
		query = query.Where("tanggal >= ?", *tanggalMulai)
	} else if tanggalSelesai != nil {
		query = query.Where("tanggal <= ?", *tanggalSelesai)
	}
	
	// Order by tanggal descending (newest first)
	query = query.Order("tanggal DESC")
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetAbsensiScanByDate retrieves absensi scan records by date
func (r *AbsensiRepositoryImpl) GetAbsensiScanByDate(tanggal time.Time) ([]models.Absensi, error) {
	var data []models.Absensi
	if err := r.db.Preload("PesertaDidik").Where("tanggal = ?", tanggal).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetAbsensiScanByMonth retrieves absensi scan records by month and year
func (r *AbsensiRepositoryImpl) GetAbsensiScanByMonth(bulan, tahun int) ([]models.Absensi, error) {
	var data []models.Absensi
	if err := r.db.Preload("PesertaDidik").
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", bulan, tahun).
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetRekapByRombelTanggalBidangStudi retrieves rekapitulasi by rombel, tanggal, and bidang_studi
func (r *AbsensiRepositoryImpl) GetRekapByRombelTanggalBidangStudi(rombelID, tahunPelajaranID uint, tanggal time.Time, bidangStudiID *uint) ([]models.RekapitulasiAbsensi, error) {
	var data []models.RekapitulasiAbsensi
	
	query := r.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND tanggal = ?", rombelID, tahunPelajaranID, tanggal)
	
	if bidangStudiID == nil {
		query = query.Where("bidang_studi_id IS NULL")
	} else {
		query = query.Where("bidang_studi_id = ?", *bidangStudiID)
	}
	
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	
	return data, nil
}
