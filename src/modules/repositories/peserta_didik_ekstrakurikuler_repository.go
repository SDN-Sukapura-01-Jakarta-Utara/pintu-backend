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
