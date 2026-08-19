package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

type AbsensiEkskulRepository struct {
	db *gorm.DB
}

func NewAbsensiEkskulRepository(db *gorm.DB) *AbsensiEkskulRepository {
	return &AbsensiEkskulRepository{db: db}
}

// CreateKegiatanEkskul creates a new ekstrakurikuler activity session
func (r *AbsensiEkskulRepository) CreateKegiatanEkskul(kegiatan *models.KegiatanEkskul) error {
	return r.db.Create(kegiatan).Error
}

// BulkCreateAbsensiEkskul creates multiple student attendance records
func (r *AbsensiEkskulRepository) BulkCreateAbsensiEkskul(absensi []models.AbsensiEkskul) error {
	if len(absensi) == 0 {
		return nil
	}
	return r.db.Create(&absensi).Error
}

// BulkCreateAbsensiPelatih creates multiple coach attendance records
func (r *AbsensiEkskulRepository) BulkCreateAbsensiPelatih(absensi []models.AbsensiPelatihEkskul) error {
	if len(absensi) == 0 {
		return nil
	}
	return r.db.Create(&absensi).Error
}

// CheckEkstrakurikulerExists checks if ekstrakurikuler exists
func (r *AbsensiEkskulRepository) CheckEkstrakurikulerExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Ekstrakurikuler{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count > 0, err
}

// CheckTahunPelajaranExists checks if tahun pelajaran exists
func (r *AbsensiEkskulRepository) CheckTahunPelajaranExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.TahunPelajaran{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count > 0, err
}

// CheckPesertaDidikRombelExists checks if peserta didik rombel exists
func (r *AbsensiEkskulRepository) CheckPesertaDidikRombelExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.PesertaDidikRombel{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count > 0, err
}

// CheckPelatihExists checks if pelatih exists
func (r *AbsensiEkskulRepository) CheckPelatihExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Pelatih{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count > 0, err
}

// CheckKegiatanExists checks if kegiatan already exists for same ekstrakurikuler, tahun pelajaran, and date
func (r *AbsensiEkskulRepository) CheckKegiatanExists(ekstrakurikulerID, tahunPelajaranID uint, tanggalKegiatan string) (bool, error) {
	var count int64
	err := r.db.Model(&models.KegiatanEkskul{}).
		Where("ekstrakurikuler_id = ? AND tahun_pelajaran_id = ? AND tanggal_kegiatan = ? AND deleted_at IS NULL", 
			ekstrakurikulerID, tahunPelajaranID, tanggalKegiatan).
		Count(&count).Error
	return count > 0, err
}

// GetEkstrakurikulerName gets ekstrakurikuler name by ID
func (r *AbsensiEkskulRepository) GetEkstrakurikulerName(id uint) (string, error) {
	var ekskul models.Ekstrakurikuler
	err := r.db.Select("name").Where("id = ? AND deleted_at IS NULL", id).First(&ekskul).Error
	if err != nil {
		return "", err
	}
	return ekskul.Name, nil
}

// GetTahunPelajaranName gets tahun pelajaran name by ID
func (r *AbsensiEkskulRepository) GetTahunPelajaranName(id uint) (string, error) {
	var tahunPelajaran models.TahunPelajaran
	err := r.db.Select("tahun_pelajaran").Where("id = ? AND deleted_at IS NULL", id).First(&tahunPelajaran).Error
	if err != nil {
		return "", err
	}
	return tahunPelajaran.TahunPelajaran, nil
}

// GetKegiatanByEkskulAndTahun gets all kegiatan with full relations by ekstrakurikuler and tahun pelajaran
func (r *AbsensiEkskulRepository) GetKegiatanByEkskulAndTahun(ekstrakurikulerID, tahunPelajaranID uint, nama string, rombelID *uint, bulan *int, tahun *int) ([]models.KegiatanEkskul, error) {
	var kegiatan []models.KegiatanEkskul
	
	// Build the base query for kegiatan
	query := r.db.
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AbsensiEkskul", func(db *gorm.DB) *gorm.DB {
			// Apply filters on absensi_ekskul if provided
			subQuery := db.Order("absensi_ekskul.id ASC")
			
			if nama != "" || rombelID != nil {
				// Join with peserta_didik_rombel and peserta_didik for filtering
				subQuery = subQuery.
					Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = absensi_ekskul.peserta_didik_rombel_id AND peserta_didik_rombel.deleted_at IS NULL").
					Joins("JOIN peserta_didik ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id AND peserta_didik.deleted_at IS NULL")
				
				if nama != "" {
					subQuery = subQuery.Where("peserta_didik.nama ILIKE ?", "%"+nama+"%")
				}
				if rombelID != nil {
					subQuery = subQuery.Where("peserta_didik_rombel.rombel_id = ?", *rombelID)
				}
			}
			
			return subQuery
		}).
		Preload("AbsensiEkskul.PesertaDidikRombel").
		Preload("AbsensiEkskul.PesertaDidikRombel.PesertaDidik").
		Preload("AbsensiEkskul.PesertaDidikRombel.Rombel").
		Preload("AbsensiEkskul.PesertaDidikRombel.Rombel.Kelas").
		Preload("AbsensiPelatih", func(db *gorm.DB) *gorm.DB {
			return db.Order("absensi_pelatih_ekskul.id ASC")
		}).
		Preload("AbsensiPelatih.Pelatih").
		Where("ekstrakurikuler_id = ? AND tahun_pelajaran_id = ?", ekstrakurikulerID, tahunPelajaranID)
	
	// Apply date filters if provided
	if bulan != nil && tahun != nil {
		// Filter by both month and year
		query = query.Where("EXTRACT(MONTH FROM tanggal_kegiatan) = ? AND EXTRACT(YEAR FROM tanggal_kegiatan) = ?", *bulan, *tahun)
	} else if bulan != nil {
		// Filter by month only
		query = query.Where("EXTRACT(MONTH FROM tanggal_kegiatan) = ?", *bulan)
	} else if tahun != nil {
		// Filter by year only
		query = query.Where("EXTRACT(YEAR FROM tanggal_kegiatan) = ?", *tahun)
	}
	
	query = query.Order("tanggal_kegiatan DESC")
	
	err := query.Find(&kegiatan).Error
	return kegiatan, err
}

// GetAbsensiSiswaByID gets single absensi siswa by ID with all relations
func (r *AbsensiEkskulRepository) GetAbsensiSiswaByID(id uint) (*models.AbsensiEkskul, error) {
	var absensi models.AbsensiEkskul
	err := r.db.
		Preload("KegiatanEkskul").
		Preload("KegiatanEkskul.Ekstrakurikuler").
		Preload("KegiatanEkskul.TahunPelajaran").
		Preload("PesertaDidikRombel").
		Preload("PesertaDidikRombel.PesertaDidik").
		Preload("PesertaDidikRombel.Rombel").
		Preload("PesertaDidikRombel.Rombel.Kelas").
		First(&absensi, id).Error
	if err != nil {
		return nil, err
	}
	return &absensi, nil
}

// UpdateAbsensiSiswa updates absensi siswa status and keterangan
func (r *AbsensiEkskulRepository) UpdateAbsensiSiswa(absensi *models.AbsensiEkskul) error {
	return r.db.Save(absensi).Error
}


// GetPelatihByID gets pelatih by ID
func (r *AbsensiEkskulRepository) GetPelatihByID(id uint) (*models.Pelatih, error) {
	var pelatih models.Pelatih
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pelatih).Error
	if err != nil {
		return nil, err
	}
	return &pelatih, nil
}

// GetKegiatanByPelatihAndTahun gets all kegiatan where pelatih is assigned
func (r *AbsensiEkskulRepository) GetKegiatanByPelatihAndTahun(pelatihID, tahunPelajaranID uint, ekstrakurikulerID *uint, bulan *int, tahun *int) ([]models.KegiatanEkskul, error) {
	var kegiatan []models.KegiatanEkskul
	
	// Build query - join with absensi_pelatih_ekskul OR get all kegiatan from pelatih's ekstrakurikuler
	// We need to get all kegiatan from ekstrakurikuler assigned to this pelatih
	query := r.db.
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AbsensiEkskul").
		Preload("AbsensiPelatih", "pelatih_id = ?", pelatihID).
		Joins("JOIN pelatih_ekstrakurikuler ON pelatih_ekstrakurikuler.ekstrakurikuler_id = kegiatan_ekskul.ekstrakurikuler_id AND pelatih_ekstrakurikuler.deleted_at IS NULL").
		Where("pelatih_ekstrakurikuler.pelatih_id = ? AND kegiatan_ekskul.tahun_pelajaran_id = ? AND kegiatan_ekskul.deleted_at IS NULL", pelatihID, tahunPelajaranID)
	
	// Apply ekstrakurikuler filter if provided
	if ekstrakurikulerID != nil {
		query = query.Where("kegiatan_ekskul.ekstrakurikuler_id = ?", *ekstrakurikulerID)
	}
	
	// Apply date filters if provided
	if bulan != nil && tahun != nil {
		query = query.Where("EXTRACT(MONTH FROM kegiatan_ekskul.tanggal_kegiatan) = ? AND EXTRACT(YEAR FROM kegiatan_ekskul.tanggal_kegiatan) = ?", *bulan, *tahun)
	} else if bulan != nil {
		query = query.Where("EXTRACT(MONTH FROM kegiatan_ekskul.tanggal_kegiatan) = ?", *bulan)
	} else if tahun != nil {
		query = query.Where("EXTRACT(YEAR FROM kegiatan_ekskul.tanggal_kegiatan) = ?", *tahun)
	}
	
	query = query.Order("kegiatan_ekskul.tanggal_kegiatan DESC")
	
	err := query.Find(&kegiatan).Error
	return kegiatan, err
}


// GetKegiatanListByEkskulAndTahun gets kegiatan list without detailed student attendance
func (r *AbsensiEkskulRepository) GetKegiatanListByEkskulAndTahun(ekstrakurikulerID, tahunPelajaranID uint, bulan *int, tahun *int) ([]models.KegiatanEkskul, error) {
	var kegiatan []models.KegiatanEkskul
	
	query := r.db.
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AbsensiEkskul").
		Preload("AbsensiPelatih.Pelatih"). // Nested preload to get pelatih details
		Where("ekstrakurikuler_id = ? AND tahun_pelajaran_id = ?", ekstrakurikulerID, tahunPelajaranID)
	
	// Apply date filters if provided
	if bulan != nil && tahun != nil {
		query = query.Where("EXTRACT(MONTH FROM tanggal_kegiatan) = ? AND EXTRACT(YEAR FROM tanggal_kegiatan) = ?", *bulan, *tahun)
	} else if bulan != nil {
		query = query.Where("EXTRACT(MONTH FROM tanggal_kegiatan) = ?", *bulan)
	} else if tahun != nil {
		query = query.Where("EXTRACT(YEAR FROM tanggal_kegiatan) = ?", *tahun)
	}
	
	query = query.Order("tanggal_kegiatan DESC")
	
	err := query.Find(&kegiatan).Error
	return kegiatan, err
}


// GetKegiatanEkskulByID gets kegiatan ekskul by ID
func (r *AbsensiEkskulRepository) GetKegiatanEkskulByID(id uint) (*models.KegiatanEkskul, error) {
	var kegiatan models.KegiatanEkskul
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&kegiatan).Error
	if err != nil {
		return nil, err
	}
	return &kegiatan, nil
}

// GetKegiatanEkskulByIDWithDetails gets kegiatan ekskul by ID with full relations
func (r *AbsensiEkskulRepository) GetKegiatanEkskulByIDWithDetails(id uint) (*models.KegiatanEkskul, error) {
	var kegiatan models.KegiatanEkskul
	err := r.db.
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AbsensiEkskul.PesertaDidikRombel.PesertaDidik").
		Preload("AbsensiEkskul.PesertaDidikRombel.Rombel.Kelas").
		Preload("AbsensiPelatih.Pelatih").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&kegiatan).Error
	if err != nil {
		return nil, err
	}
	return &kegiatan, nil
}

// CheckAbsensiPelatihExists checks if pelatih attendance record exists
func (r *AbsensiEkskulRepository) CheckAbsensiPelatihExists(kegiatanEkskulID, pelatihID uint) (*models.AbsensiPelatihEkskul, error) {
	var absensi models.AbsensiPelatihEkskul
	err := r.db.Where("kegiatan_ekskul_id = ? AND pelatih_id = ?", kegiatanEkskulID, pelatihID).First(&absensi).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Not found, return nil
	}
	if err != nil {
		return nil, err
	}
	return &absensi, nil
}

// CreateAbsensiPelatih creates pelatih attendance record
func (r *AbsensiEkskulRepository) CreateAbsensiPelatih(absensi *models.AbsensiPelatihEkskul) error {
	return r.db.Create(absensi).Error
}

// DeleteAbsensiPelatih deletes pelatih attendance record (soft delete)
func (r *AbsensiEkskulRepository) DeleteAbsensiPelatih(id uint) error {
	return r.db.Delete(&models.AbsensiPelatihEkskul{}, id).Error
}


// UpdateFotoKegiatan updates foto_kegiatan field in kegiatan_ekskul
func (r *AbsensiEkskulRepository) UpdateFotoKegiatan(kegiatanID uint, fotoKegiatan string) error {
	return r.db.Model(&models.KegiatanEkskul{}).
		Where("id = ? AND deleted_at IS NULL", kegiatanID).
		Update("foto_kegiatan", fotoKegiatan).Error
}

// UpdateKegiatanEkskul updates kegiatan ekstrakurikuler fields
func (r *AbsensiEkskulRepository) UpdateKegiatanEkskul(kegiatan *models.KegiatanEkskul) error {
	return r.db.Save(kegiatan).Error
}
