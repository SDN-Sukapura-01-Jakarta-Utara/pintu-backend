package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
)

type AnakInklusiRepository interface {
	Create(anakInklusi *models.AnakInklusi) error
	GetByID(id uint) (*models.AnakInklusi, error)
	GetByPesertaDidikID(pesertaDidikID uint) (*models.AnakInklusi, error)
	GetAllWithFilter(params GetAnakInklusiParams) ([]models.AnakInklusi, int64, error)
	Update(anakInklusi *models.AnakInklusi) error
	Delete(id uint) error
}

type anakInklusiRepository struct {
	db *gorm.DB
}

type GetAnakInklusiFilter struct {
	NamaPesertaDidik       string
	NIS                    string
	NISN                   string
	JenisHambatan          string
	StartDate              time.Time
	EndDate                time.Time
	Status                 string
	PesertaDidikStatusOnly string // "active" or "non-active" or empty for all
}

type GetAnakInklusiParams struct {
	Filter GetAnakInklusiFilter
	Limit  int
	Offset int
}

func NewAnakInklusiRepository(db *gorm.DB) AnakInklusiRepository {
	return &anakInklusiRepository{db: db}
}

func (r *anakInklusiRepository) Create(anakInklusi *models.AnakInklusi) error {
	return r.db.Create(anakInklusi).Error
}

func (r *anakInklusiRepository) GetByID(id uint) (*models.AnakInklusi, error) {
	var anakInklusi models.AnakInklusi
	if err := r.db.Preload("PesertaDidik").First(&anakInklusi, id).Error; err != nil {
		return nil, err
	}
	return &anakInklusi, nil
}

func (r *anakInklusiRepository) GetByPesertaDidikID(pesertaDidikID uint) (*models.AnakInklusi, error) {
	var anakInklusi models.AnakInklusi
	if err := r.db.Preload("PesertaDidik").Where("peserta_didik_id = ?", pesertaDidikID).First(&anakInklusi).Error; err != nil {
		return nil, err
	}
	return &anakInklusi, nil
}

func (r *anakInklusiRepository) GetAllWithFilter(params GetAnakInklusiParams) ([]models.AnakInklusi, int64, error) {
	var anakInklusiList []models.AnakInklusi
	var total int64

	query := r.db.Model(&models.AnakInklusi{}).Preload("PesertaDidik")

	// Apply filters
	if params.Filter.NamaPesertaDidik != "" {
		query = query.Joins("JOIN peserta_didik ON peserta_didik.id = anak_inklusi.peserta_didik_id").
			Where("peserta_didik.nama ILIKE ?", "%"+params.Filter.NamaPesertaDidik+"%")
	}

	if params.Filter.NIS != "" {
		if params.Filter.NamaPesertaDidik == "" {
			query = query.Joins("JOIN peserta_didik ON peserta_didik.id = anak_inklusi.peserta_didik_id")
		}
		query = query.Where("peserta_didik.nis ILIKE ?", "%"+params.Filter.NIS+"%")
	}

	if params.Filter.NISN != "" {
		if params.Filter.NamaPesertaDidik == "" && params.Filter.NIS == "" {
			query = query.Joins("JOIN peserta_didik ON peserta_didik.id = anak_inklusi.peserta_didik_id")
		}
		query = query.Where("peserta_didik.nisn ILIKE ?", "%"+params.Filter.NISN+"%")
	}

	if params.Filter.JenisHambatan != "" {
		query = query.Where("anak_inklusi.jenis_hambatan ILIKE ?", "%"+params.Filter.JenisHambatan+"%")
	}

	if !params.Filter.StartDate.IsZero() {
		query = query.Where("anak_inklusi.tanggal_diagnosa >= ?", params.Filter.StartDate)
	}

	if !params.Filter.EndDate.IsZero() {
		query = query.Where("anak_inklusi.tanggal_diagnosa <= ?", params.Filter.EndDate)
	}

	if params.Filter.Status != "" {
		query = query.Where("anak_inklusi.status = ?", params.Filter.Status)
	}

	// Filter by peserta_didik status
	if params.Filter.PesertaDidikStatusOnly == "active" {
		// Ensure join exists
		if params.Filter.NamaPesertaDidik == "" && params.Filter.NIS == "" && params.Filter.NISN == "" {
			query = query.Joins("JOIN peserta_didik ON peserta_didik.id = anak_inklusi.peserta_didik_id")
		}
		query = query.Where("peserta_didik.status = ?", "active")
	} else if params.Filter.PesertaDidikStatusOnly == "non-active" {
		// Ensure join exists
		if params.Filter.NamaPesertaDidik == "" && params.Filter.NIS == "" && params.Filter.NISN == "" {
			query = query.Joins("JOIN peserta_didik ON peserta_didik.id = anak_inklusi.peserta_didik_id")
		}
		query = query.Where("peserta_didik.status != ?", "active")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Order("anak_inklusi.created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&anakInklusiList).Error; err != nil {
		return nil, 0, err
	}

	return anakInklusiList, total, nil
}

func (r *anakInklusiRepository) Update(anakInklusi *models.AnakInklusi) error {
	return r.db.Save(anakInklusi).Error
}

func (r *anakInklusiRepository) Delete(id uint) error {
	return r.db.Delete(&models.AnakInklusi{}, id).Error
}
