package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

type AnakInklusiRombelRepository interface {
	Create(anakInklusiRombel *models.AnakInklusiRombel) error
	CheckExists(anakInklusiID uint, pesertaDidikRombelID uint) (bool, error)
	GetPesertaDidikRombelByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidikRombel, error)
	GetAllWithFilter(params GetAnakInklusiRombelParams) ([]models.AnakInklusiRombel, int64, error)
	GetByID(id uint) (*models.AnakInklusiRombel, error)
	Update(anakInklusiRombel *models.AnakInklusiRombel) error
	Delete(id uint) error
}

type GetAnakInklusiRombelFilter struct {
	TahunPelajaranID string
	NamaPesertaDidik string
	NIS              string
	NISN             string
	JenisHambatan    string
	Status           string // Status anak inklusi
	RombelID         string // Filter by rombel ID
}

type GetAnakInklusiRombelParams struct {
	Filter GetAnakInklusiRombelFilter
	Limit  int
	Offset int
}

type anakInklusiRombelRepository struct {
	db *gorm.DB
}

func NewAnakInklusiRombelRepository(db *gorm.DB) AnakInklusiRombelRepository {
	return &anakInklusiRombelRepository{db: db}
}

func (r *anakInklusiRombelRepository) Create(anakInklusiRombel *models.AnakInklusiRombel) error {
	return r.db.Create(anakInklusiRombel).Error
}

func (r *anakInklusiRombelRepository) CheckExists(anakInklusiID uint, pesertaDidikRombelID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AnakInklusiRombel{}).
		Where("anak_inklusi_id = ? AND peserta_didik_rombel_id = ?", anakInklusiID, pesertaDidikRombelID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *anakInklusiRombelRepository) GetPesertaDidikRombelByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidikRombel, error) {
	var pesertaDidikRombelList []models.PesertaDidikRombel
	err := r.db.Preload("PesertaDidik").
		Where("tahun_pelajaran_id = ?", tahunPelajaranID).
		Find(&pesertaDidikRombelList).Error
	
	if err != nil {
		return nil, err
	}
	
	return pesertaDidikRombelList, nil
}

func (r *anakInklusiRombelRepository) GetAllWithFilter(params GetAnakInklusiRombelParams) ([]models.AnakInklusiRombel, int64, error) {
	var anakInklusiRombelList []models.AnakInklusiRombel
	var total int64

	query := r.db.Model(&models.AnakInklusiRombel{}).
		Preload("AnakInklusi").
		Preload("AnakInklusi.PesertaDidik").
		Preload("PesertaDidikRombel").
		Preload("PesertaDidikRombel.Rombel").
		Preload("PesertaDidikRombel.TahunPelajaran").
		Preload("GuruPendampingKhusus").
		Preload("GuruKelas")

	// Track which joins have been added
	hasPesertaDidikRombelJoin := false
	hasPesertaDidikJoin := false
	hasAnakInklusiJoin := false
	hasRombelJoin := false

	// Required filter: tahun_pelajaran_id
	if params.Filter.TahunPelajaranID != "" {
		query = query.Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = anak_inklusi_rombel.peserta_didik_rombel_id").
			Where("peserta_didik_rombel.tahun_pelajaran_id = ?", params.Filter.TahunPelajaranID)
		hasPesertaDidikRombelJoin = true
	}

	// Filter by nama peserta didik
	if params.Filter.NamaPesertaDidik != "" {
		if !hasPesertaDidikRombelJoin {
			query = query.Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = anak_inklusi_rombel.peserta_didik_rombel_id")
			hasPesertaDidikRombelJoin = true
		}
		query = query.Joins("JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id").
			Joins("JOIN peserta_didik pd2 ON pd2.id = ai2.peserta_didik_id").
			Where("pd2.nama ILIKE ?", "%"+params.Filter.NamaPesertaDidik+"%")
		hasPesertaDidikJoin = true
		hasAnakInklusiJoin = true
	}

	// Filter by NIS
	if params.Filter.NIS != "" {
		if !hasPesertaDidikRombelJoin {
			query = query.Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = anak_inklusi_rombel.peserta_didik_rombel_id")
			hasPesertaDidikRombelJoin = true
		}
		if !hasPesertaDidikJoin {
			query = query.Joins("JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id").
				Joins("JOIN peserta_didik pd2 ON pd2.id = ai2.peserta_didik_id")
			hasPesertaDidikJoin = true
			hasAnakInklusiJoin = true
		}
		query = query.Where("pd2.nis ILIKE ?", "%"+params.Filter.NIS+"%")
	}

	// Filter by NISN
	if params.Filter.NISN != "" {
		if !hasPesertaDidikRombelJoin {
			query = query.Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = anak_inklusi_rombel.peserta_didik_rombel_id")
			hasPesertaDidikRombelJoin = true
		}
		if !hasPesertaDidikJoin {
			query = query.Joins("JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id").
				Joins("JOIN peserta_didik pd2 ON pd2.id = ai2.peserta_didik_id")
			hasPesertaDidikJoin = true
			hasAnakInklusiJoin = true
		}
		query = query.Where("pd2.nisn ILIKE ?", "%"+params.Filter.NISN+"%")
	}

	// Filter by jenis_hambatan
	if params.Filter.JenisHambatan != "" {
		if !hasAnakInklusiJoin {
			query = query.Joins("JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id")
			hasAnakInklusiJoin = true
		}
		query = query.Where("ai2.jenis_hambatan ILIKE ?", "%"+params.Filter.JenisHambatan+"%")
	}

	// Filter by status (anak_inklusi status)
	if params.Filter.Status != "" {
		if !hasAnakInklusiJoin {
			query = query.Joins("JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id")
			hasAnakInklusiJoin = true
		}
		query = query.Where("ai2.status = ?", params.Filter.Status)
	}

	// Filter by rombel_id
	if params.Filter.RombelID != "" {
		if !hasPesertaDidikRombelJoin {
			query = query.Joins("JOIN peserta_didik_rombel ON peserta_didik_rombel.id = anak_inklusi_rombel.peserta_didik_rombel_id")
			hasPesertaDidikRombelJoin = true
		}
		query = query.Where("peserta_didik_rombel.rombel_id = ?", params.Filter.RombelID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// For ordering, ensure we have the necessary joins
	if !hasRombelJoin {
		query = query.Joins("LEFT JOIN rombel ON rombel.id = peserta_didik_rombel.rombel_id")
		hasRombelJoin = true
	}
	if !hasAnakInklusiJoin {
		query = query.Joins("LEFT JOIN anak_inklusi ai2 ON ai2.id = anak_inklusi_rombel.anak_inklusi_id")
		hasAnakInklusiJoin = true
	}
	if !hasPesertaDidikJoin {
		query = query.Joins("LEFT JOIN peserta_didik pd2 ON pd2.id = ai2.peserta_didik_id")
		hasPesertaDidikJoin = true
	}

	// Get data with pagination and sorting
	// Sort by rombel name first, then by peserta didik name alphabetically
	if err := query.
		Order("rombel.name ASC").
		Order("pd2.nama ASC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&anakInklusiRombelList).Error; err != nil {
		return nil, 0, err
	}

	return anakInklusiRombelList, total, nil
}

func (r *anakInklusiRombelRepository) GetByID(id uint) (*models.AnakInklusiRombel, error) {
	var anakInklusiRombel models.AnakInklusiRombel
	if err := r.db.Preload("AnakInklusi").
		Preload("AnakInklusi.PesertaDidik").
		Preload("PesertaDidikRombel").
		Preload("PesertaDidikRombel.Rombel").
		Preload("PesertaDidikRombel.TahunPelajaran").
		Preload("GuruPendampingKhusus").
		Preload("GuruKelas").
		First(&anakInklusiRombel, id).Error; err != nil {
		return nil, err
	}
	return &anakInklusiRombel, nil
}

func (r *anakInklusiRombelRepository) Update(anakInklusiRombel *models.AnakInklusiRombel) error {
	// Use Updates with Select to update all fields including nil pointers
	return r.db.Model(anakInklusiRombel).
		Select("GuruPendampingKhususID", "GuruKelasID", "Catatan", "UpdatedByID", "UpdatedAt").
		Updates(anakInklusiRombel).Error
}

func (r *anakInklusiRombelRepository) Delete(id uint) error {
	return r.db.Delete(&models.AnakInklusiRombel{}, id).Error
}
