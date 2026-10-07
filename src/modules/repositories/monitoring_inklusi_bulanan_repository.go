package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

type MonitoringInklusiBulananRepository interface {
	Create(monitoring *models.MonitoringInklusiBulanan) error
	GetByID(id uint) (*models.MonitoringInklusiBulanan, error)
	GetPPIByID(id uint) (*models.PPIInklusi, error)
	UpdatePPIStatus(ppiID uint, status string) error
	GetByBulanTahunRombel(bulan int16, tahun int16, rombelID *uint, ppiInklusiID *uint) ([]models.MonitoringInklusiBulanan, error)
	GetRombelByID(id uint) (*models.Rombel, error)
	Update(monitoring *models.MonitoringInklusiBulanan) error
	GetSummaryByTahun(tahun int16, bulan *int16, rombelID *uint) ([]models.MonitoringInklusiBulanan, error)
	GetPPIsByTahun(tahun int16, rombelID *uint) ([]models.PPIInklusi, error)
	GetAnakInklusiByTahun(tahun int16, rombelID *uint) ([]models.AnakInklusi, error)
}

type monitoringInklusiBulananRepository struct {
	db *gorm.DB
}

func NewMonitoringInklusiBulananRepository(db *gorm.DB) MonitoringInklusiBulananRepository {
	return &monitoringInklusiBulananRepository{db: db}
}

func (r *monitoringInklusiBulananRepository) Create(monitoring *models.MonitoringInklusiBulanan) error {
	return r.db.Create(monitoring).Error
}

func (r *monitoringInklusiBulananRepository) GetByID(id uint) (*models.MonitoringInklusiBulanan, error) {
	var monitoring models.MonitoringInklusiBulanan
	if err := r.db.Preload("PPIInklusi").
		Preload("PPIInklusi.AnakInklusiRombel").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Preload("GuruPengisi").
		First(&monitoring, id).Error; err != nil {
		return nil, err
	}
	return &monitoring, nil
}

func (r *monitoringInklusiBulananRepository) GetPPIByID(id uint) (*models.PPIInklusi, error) {
	var ppi models.PPIInklusi
	if err := r.db.First(&ppi, id).Error; err != nil {
		return nil, err
	}
	return &ppi, nil
}

func (r *monitoringInklusiBulananRepository) UpdatePPIStatus(ppiID uint, status string) error {
	return r.db.Model(&models.PPIInklusi{}).
		Where("id = ?", ppiID).
		Update("status_capaian", status).Error
}

func (r *monitoringInklusiBulananRepository) GetRombelByID(id uint) (*models.Rombel, error) {
	var rombel models.Rombel
	if err := r.db.First(&rombel, id).Error; err != nil {
		return nil, err
	}
	return &rombel, nil
}

func (r *monitoringInklusiBulananRepository) GetByBulanTahunRombel(bulan int16, tahun int16, rombelID *uint, ppiInklusiID *uint) ([]models.MonitoringInklusiBulanan, error) {
	var monitoringList []models.MonitoringInklusiBulanan
	
	query := r.db.Preload("PPIInklusi").
		Preload("PPIInklusi.AnakInklusiRombel").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Preload("GuruPengisi").
		Joins("JOIN ppi_inklusi ppi ON ppi.id = monitoring_inklusi_bulanan.ppi_inklusi_id").
		Joins("JOIN anak_inklusi_rombel air ON air.id = ppi.anak_inklusi_rombel_id").
		Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
		Where("monitoring_inklusi_bulanan.bulan = ?", bulan).
		Where("monitoring_inklusi_bulanan.tahun = ?", tahun)
	
	// Optional filter by rombel
	if rombelID != nil && *rombelID > 0 {
		query = query.Where("pdr.rombel_id = ?", *rombelID)
	}
	
	// Optional filter by PPI
	if ppiInklusiID != nil && *ppiInklusiID > 0 {
		query = query.Where("monitoring_inklusi_bulanan.ppi_inklusi_id = ?", *ppiInklusiID)
	}
	
	err := query.Order("monitoring_inklusi_bulanan.created_at DESC").
		Find(&monitoringList).Error
	
	if err != nil {
		return nil, err
	}
	
	return monitoringList, nil
}

func (r *monitoringInklusiBulananRepository) Update(monitoring *models.MonitoringInklusiBulanan) error {
	return r.db.Model(monitoring).
		Select("deskripsi_perkembangan", "kendala_ditemui", "tindak_lanjut", "file_pendukung", "guru_pengisi_id", "updated_by_id", "updated_at").
		Updates(monitoring).Error
}

func (r *monitoringInklusiBulananRepository) GetSummaryByTahun(tahun int16, bulan *int16, rombelID *uint) ([]models.MonitoringInklusiBulanan, error) {
	var monitoringList []models.MonitoringInklusiBulanan
	
	query := r.db.Preload("PPIInklusi").
		Preload("PPIInklusi.AnakInklusiRombel").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi").
		Preload("PPIInklusi.AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel").
		Preload("PPIInklusi.AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Joins("JOIN ppi_inklusi ppi ON ppi.id = monitoring_inklusi_bulanan.ppi_inklusi_id").
		Joins("JOIN anak_inklusi_rombel air ON air.id = ppi.anak_inklusi_rombel_id").
		Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
		Where("monitoring_inklusi_bulanan.tahun = ?", tahun)
	
	// Optional filter by bulan
	if bulan != nil && *bulan >= 1 && *bulan <= 12 {
		query = query.Where("monitoring_inklusi_bulanan.bulan = ?", *bulan)
	}
	
	// Optional filter by rombel
	if rombelID != nil && *rombelID > 0 {
		query = query.Where("pdr.rombel_id = ?", *rombelID)
	}
	
	err := query.Order("monitoring_inklusi_bulanan.bulan ASC").
		Find(&monitoringList).Error
	
	return monitoringList, err
}

func (r *monitoringInklusiBulananRepository) GetPPIsByTahun(tahun int16, rombelID *uint) ([]models.PPIInklusi, error) {
	var ppiList []models.PPIInklusi
	
	query := r.db.Preload("AnakInklusiRombel").
		Preload("AnakInklusiRombel.AnakInklusi").
		Preload("AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("AnakInklusiRombel.PesertaDidikRombel").
		Preload("AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Preload("AnakInklusiRombel.PesertaDidikRombel.TahunPelajaran").
		Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
		Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
		Joins("JOIN tahun_pelajaran tp ON tp.id = pdr.tahun_pelajaran_id").
		Where("tp.tahun_mulai = ?", tahun)
	
	if rombelID != nil && *rombelID > 0 {
		query = query.Where("pdr.rombel_id = ?", *rombelID)
	}
	
	err := query.Find(&ppiList).Error
	return ppiList, err
}

func (r *monitoringInklusiBulananRepository) GetAnakInklusiByTahun(tahun int16, rombelID *uint) ([]models.AnakInklusi, error) {
	var anakInklusiList []models.AnakInklusi
	
	query := r.db.Preload("PesertaDidik").
		Joins("JOIN anak_inklusi_rombel air ON air.anak_inklusi_id = anak_inklusi.id").
		Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
		Joins("JOIN tahun_pelajaran tp ON tp.id = pdr.tahun_pelajaran_id").
		Where("tp.tahun_mulai = ?", tahun).
		Distinct()
	
	if rombelID != nil && *rombelID > 0 {
		query = query.Where("pdr.rombel_id = ?", *rombelID)
	}
	
	err := query.Find(&anakInklusiList).Error
	return anakInklusiList, err
}
