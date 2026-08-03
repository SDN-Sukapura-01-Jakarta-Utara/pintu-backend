package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// PesertaDidikEkstrakurikulerRepository handles data operations
type PesertaDidikEkstrakurikulerRepository interface {
	Create(data *models.PesertaDidikEkstrakurikuler) error
	CreateBulk(data []models.PesertaDidikEkstrakurikuler) error
	GetByID(id uint) (*models.PesertaDidikEkstrakurikuler, error)
	GetByPesertaDidikRombelAndEkskul(pesertaDidikRombelID, ekstrakurikulerID uint) (*models.PesertaDidikEkstrakurikuler, error)
	GetAllByPesertaDidikRombel(pesertaDidikRombelID uint) ([]models.PesertaDidikEkstrakurikuler, error)
	GetAllByRombelAndTahunPelajaran(rombelID, tahunPelajaranID uint) ([]models.PesertaDidikEkstrakurikuler, error)
	GetAllByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidikEkstrakurikuler, error)
	GetStatistikByTahunPelajaran(tahunPelajaranID uint, rombelID *uint) ([]models.PesertaDidikEkstrakurikuler, error)
	GetRekapPerEkskul(tahunPelajaranID uint, nama, nis string, rombelID, ekstrakurikulerID uint, limit, offset int) ([]models.PesertaDidikEkstrakurikuler, int64, error)
	GetRekapPerRombel(rombelID, tahunPelajaranID uint, nama, nis string, limit, offset int) ([]models.PesertaDidikEkstrakurikuler, int64, error)
	Delete(id uint) error
	DeleteByPesertaDidikRombelAndEkskul(pesertaDidikRombelID, ekstrakurikulerID uint) error
}

type PesertaDidikEkstrakurikulerRepositoryImpl struct {
	db *gorm.DB
}

// NewPesertaDidikEkstrakurikulerRepository creates a new repository
func NewPesertaDidikEkstrakurikulerRepository(db *gorm.DB) PesertaDidikEkstrakurikulerRepository {
	return &PesertaDidikEkstrakurikulerRepositoryImpl{db: db}
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) Create(data *models.PesertaDidikEkstrakurikuler) error {
	return r.db.Create(data).Error
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) CreateBulk(data []models.PesertaDidikEkstrakurikuler) error {
	return r.db.Create(&data).Error
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetByID(id uint) (*models.PesertaDidikEkstrakurikuler, error) {
	var data models.PesertaDidikEkstrakurikuler
	if err := r.db.Preload("Ekstrakurikuler").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetByPesertaDidikRombelAndEkskul(pesertaDidikRombelID, ekstrakurikulerID uint) (*models.PesertaDidikEkstrakurikuler, error) {
	var data models.PesertaDidikEkstrakurikuler
	if err := r.db.Where("peserta_didik_rombel_id = ? AND ekstrakurikuler_id = ?", pesertaDidikRombelID, ekstrakurikulerID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetAllByPesertaDidikRombel(pesertaDidikRombelID uint) ([]models.PesertaDidikEkstrakurikuler, error) {
	var data []models.PesertaDidikEkstrakurikuler
	if err := r.db.Preload("Ekstrakurikuler").Where("peserta_didik_rombel_id = ?", pesertaDidikRombelID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.PesertaDidikEkstrakurikuler{}, id).Error
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) DeleteByPesertaDidikRombelAndEkskul(pesertaDidikRombelID, ekstrakurikulerID uint) error {
	return r.db.Where("peserta_didik_rombel_id = ? AND ekstrakurikuler_id = ?", pesertaDidikRombelID, ekstrakurikulerID).Delete(&models.PesertaDidikEkstrakurikuler{}).Error
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetAllByRombelAndTahunPelajaran(rombelID, tahunPelajaranID uint) ([]models.PesertaDidikEkstrakurikuler, error) {
	var data []models.PesertaDidikEkstrakurikuler
	
	err := r.db.
		Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = peserta_didik_ekstrakurikuler.peserta_didik_rombel_id").
		Where("peserta_didik_rombel.rombel_id = ? AND peserta_didik_rombel.tahun_pelajaran_id = ?", rombelID, tahunPelajaranID).
		Preload("PesertaDidikRombel.PesertaDidik").
		Preload("Ekstrakurikuler").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	
	return data, nil
}


func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetAllByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidikEkstrakurikuler, error) {
	var data []models.PesertaDidikEkstrakurikuler
	
	err := r.db.
		Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = peserta_didik_ekstrakurikuler.peserta_didik_rombel_id").
		Where("peserta_didik_rombel.tahun_pelajaran_id = ?", tahunPelajaranID).
		Preload("PesertaDidikRombel.PesertaDidik").
		Preload("PesertaDidikRombel.Rombel").
		Preload("Ekstrakurikuler").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetStatistikByTahunPelajaran(tahunPelajaranID uint, rombelID *uint) ([]models.PesertaDidikEkstrakurikuler, error) {
	var data []models.PesertaDidikEkstrakurikuler
	
	query := r.db.
		Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = peserta_didik_ekstrakurikuler.peserta_didik_rombel_id").
		Where("peserta_didik_rombel.tahun_pelajaran_id = ?", tahunPelajaranID)
	
	// Add rombel filter if provided
	if rombelID != nil {
		query = query.Where("peserta_didik_rombel.rombel_id = ?", *rombelID)
	}
	
	err := query.
		Preload("PesertaDidikRombel.PesertaDidik").
		Preload("PesertaDidikRombel.Rombel").
		Preload("Ekstrakurikuler").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	
	return data, nil
}


func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetRekapPerEkskul(tahunPelajaranID uint, nama, nis string, rombelID, ekstrakurikulerID uint, limit, offset int) ([]models.PesertaDidikEkstrakurikuler, int64, error) {
	var data []models.PesertaDidikEkstrakurikuler
	var total int64
	
	query := r.db.
		Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = peserta_didik_ekstrakurikuler.peserta_didik_rombel_id").
		Joins("JOIN peserta_didik ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id").
		Joins("JOIN rombel ON rombel.id = peserta_didik_rombel.rombel_id").
		Where("peserta_didik_rombel.tahun_pelajaran_id = ?", tahunPelajaranID)
	
	// Apply filters
	if nama != "" {
		query = query.Where("peserta_didik.nama ILIKE ?", "%"+nama+"%")
	}
	if nis != "" {
		query = query.Where("peserta_didik.nis ILIKE ?", "%"+nis+"%")
	}
	if rombelID > 0 {
		query = query.Where("peserta_didik_rombel.rombel_id = ?", rombelID)
	}
	if ekstrakurikulerID > 0 {
		query = query.Where("peserta_didik_ekstrakurikuler.ekstrakurikuler_id = ?", ekstrakurikulerID)
	}
	
	// Count total
	if err := query.Model(&models.PesertaDidikEkstrakurikuler{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated data with sorting: ekstrakurikuler_id ASC, rombel alphabetically, then nama A-Z
	err := query.
		Preload("PesertaDidikRombel.PesertaDidik").
		Preload("PesertaDidikRombel.Rombel").
		Preload("Ekstrakurikuler").
		Order("ekstrakurikuler_id ASC, rombel.name ASC, peserta_didik.nama ASC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return data, total, nil
}

func (r *PesertaDidikEkstrakurikulerRepositoryImpl) GetRekapPerRombel(rombelID, tahunPelajaranID uint, nama, nis string, limit, offset int) ([]models.PesertaDidikEkstrakurikuler, int64, error) {
	var data []models.PesertaDidikEkstrakurikuler
	
	// Get all students in rombel first
	var allStudents []models.PesertaDidikRombel
	studentQuery := r.db.Model(&models.PesertaDidikRombel{}).
		Joins("JOIN peserta_didik ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id").
		Where("peserta_didik_rombel.rombel_id = ? AND peserta_didik_rombel.tahun_pelajaran_id = ?", rombelID, tahunPelajaranID).
		Where("peserta_didik_rombel.status = ?", "active")
	
	// Apply filters
	if nama != "" {
		studentQuery = studentQuery.Where("peserta_didik.nama ILIKE ?", "%"+nama+"%")
	}
	if nis != "" {
		studentQuery = studentQuery.Where("peserta_didik.nis ILIKE ?", "%"+nis+"%")
	}
	
	// Count total students
	var total int64
	if err := studentQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated students
	if err := studentQuery.
		Preload("PesertaDidik").
		Preload("Rombel").
		Order("peserta_didik.nama ASC").
		Limit(limit).
		Offset(offset).
		Find(&allStudents).Error; err != nil {
		return nil, 0, err
	}
	
	// Get ekstrakurikuler for these students
	var studentIDs []uint
	for _, student := range allStudents {
		studentIDs = append(studentIDs, student.ID)
	}
	
	if len(studentIDs) > 0 {
		if err := r.db.
			Where("peserta_didik_rombel_id IN ?", studentIDs).
			Preload("PesertaDidikRombel.PesertaDidik").
			Preload("PesertaDidikRombel.Rombel").
			Preload("Ekstrakurikuler").
			Order("peserta_didik_rombel_id ASC, ekstrakurikuler_id ASC").
			Find(&data).Error; err != nil {
			return nil, 0, err
		}
	}
	
	// Create dummy records for students without ekstrakurikuler
	ekskulMap := make(map[uint]bool)
	for _, reg := range data {
		ekskulMap[reg.PesertaDidikRombelID] = true
	}
	
	for _, student := range allStudents {
		if !ekskulMap[student.ID] {
			// Add dummy record so student appears in response
			data = append(data, models.PesertaDidikEkstrakurikuler{
				PesertaDidikRombelID: student.ID,
				PesertaDidikRombel:   &student,
			})
		}
	}
	
	return data, total, nil
}
