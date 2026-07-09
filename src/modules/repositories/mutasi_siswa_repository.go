package repositories

import (
	"fmt"
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// MutasiSiswaRepository handles data operations for Mutasi Siswa
type MutasiSiswaRepository interface {
	Create(data *models.MutasiSiswa) error
	GetLastRegistrationNumber(tahunPelajaranID, semester int) (string, error)
	GetAllWithFilter(params GetMutasiSiswaParams) ([]models.MutasiSiswa, int64, error)
	GetByID(id uint) (*models.MutasiSiswa, error)
	Update(data *models.MutasiSiswa) error
	Delete(id uint) error
}

// GetMutasiSiswaFilter represents filter parameters
type GetMutasiSiswaFilter struct {
	TahunPelajaranID *int
	Semester         *int
	StartDate        string
	EndDate          string
	NamaSiswa        string
	NISN             string
	TempatLahir      string
	JenisKelamin     string
	PindahanKelas    *int
}

// GetMutasiSiswaParams represents query parameters
type GetMutasiSiswaParams struct {
	Filter GetMutasiSiswaFilter
	Limit  int
	Offset int
}

type MutasiSiswaRepositoryImpl struct {
	db *gorm.DB
}

// NewMutasiSiswaRepository creates a new Mutasi Siswa repository
func NewMutasiSiswaRepository(db *gorm.DB) MutasiSiswaRepository {
	return &MutasiSiswaRepositoryImpl{db: db}
}

// Create creates a new Mutasi Siswa record
func (r *MutasiSiswaRepositoryImpl) Create(data *models.MutasiSiswa) error {
	return r.db.Create(data).Error
}

// GetLastRegistrationNumber retrieves the last registration number for given tahun pelajaran and semester
func (r *MutasiSiswaRepositoryImpl) GetLastRegistrationNumber(tahunPelajaranID, semester int) (string, error) {
	var lastRecord models.MutasiSiswa
	
	err := r.db.Where("tahun_pelajaran_id = ? AND semester = ?", tahunPelajaranID, semester).
		Order("nomor_pendaftaran DESC").
		First(&lastRecord).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No records found, return empty string
			return "", nil
		}
		return "", fmt.Errorf("failed to get last registration number: %w", err)
	}
	
	return lastRecord.NomorPendaftaran, nil
}

// GetAllWithFilter retrieves all Mutasi Siswa with filters, sorting, and pagination
func (r *MutasiSiswaRepositoryImpl) GetAllWithFilter(params GetMutasiSiswaParams) ([]models.MutasiSiswa, int64, error) {
	var data []models.MutasiSiswa
	var total int64

	query := r.db.Model(&models.MutasiSiswa{}).Preload("TahunPelajaran")

	// Apply filters
	if params.Filter.TahunPelajaranID != nil {
		query = query.Where("tahun_pelajaran_id = ?", *params.Filter.TahunPelajaranID)
	}
	if params.Filter.Semester != nil {
		query = query.Where("semester = ?", *params.Filter.Semester)
	}
	if params.Filter.StartDate != "" {
		query = query.Where("created_at >= ?", params.Filter.StartDate)
	}
	if params.Filter.EndDate != "" {
		query = query.Where("created_at <= ?", params.Filter.EndDate+" 23:59:59")
	}
	if params.Filter.NamaSiswa != "" {
		query = query.Where("nama_lengkap ILIKE ?", "%"+params.Filter.NamaSiswa+"%")
	}
	if params.Filter.NISN != "" {
		query = query.Where("nisn ILIKE ?", "%"+params.Filter.NISN+"%")
	}
	if params.Filter.TempatLahir != "" {
		query = query.Where("tempat_lahir ILIKE ?", "%"+params.Filter.TempatLahir+"%")
	}
	if params.Filter.JenisKelamin != "" {
		query = query.Where("jenis_kelamin = ?", params.Filter.JenisKelamin)
	}
	if params.Filter.PindahanKelas != nil {
		query = query.Where("pindahan_kelas = ?", *params.Filter.PindahanKelas)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting: Latest created first (DESC)
	query = query.Order("created_at DESC")

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit).Offset(params.Offset)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}


// GetByID retrieves Mutasi Siswa by ID
func (r *MutasiSiswaRepositoryImpl) GetByID(id uint) (*models.MutasiSiswa, error) {
	var data models.MutasiSiswa
	if err := r.db.Preload("TahunPelajaran").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}


// Update updates Mutasi Siswa record
func (r *MutasiSiswaRepositoryImpl) Update(data *models.MutasiSiswa) error {
	return r.db.Save(data).Error
}

// Delete deletes Mutasi Siswa record by ID
func (r *MutasiSiswaRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.MutasiSiswa{}, id).Error
}

// GetDB returns the database instance (helper for internal use)
func (r *MutasiSiswaRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}
