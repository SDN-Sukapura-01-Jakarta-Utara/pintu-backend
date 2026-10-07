package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

type PPIInklusiRepository interface {
	Create(ppiInklusi *models.PPIInklusi) error
	GetByID(id uint) (*models.PPIInklusi, error)
	GetAllWithFilter(params GetPPIInklusiParams) ([]models.PPIInklusi, int64, error)
	Update(ppiInklusi *models.PPIInklusi) error
	Delete(id uint) error
	GetSummaryByTahunPelajaran(filter GetPPIInklusiFilter) (*PPIInklusiSummaryData, error)
	GetTahunPelajaranByID(id uint) (*models.TahunPelajaran, error)
}

type PPIInklusiSummaryData struct {
	TotalPPI                       int
	TotalAnakInklusi               int
	StatusCapaianBreakdown         map[string]int
	AspekPembelajaranBreakdown     []AspekPembelajaranSummary
}

type AspekPembelajaranSummary struct {
	AspekPembelajaran string
	Total             int
}

type GetPPIInklusiFilter struct {
	TahunPelajaranID    string
	RombelID            string
	GuruPembuatID       string
	AnakInklusiRombelID string
	AspekPembelajaran   string
}

type GetPPIInklusiParams struct {
	Filter GetPPIInklusiFilter
	Limit  int
	Offset int
}

type ppiInklusiRepository struct {
	db *gorm.DB
}

func NewPPIInklusiRepository(db *gorm.DB) PPIInklusiRepository {
	return &ppiInklusiRepository{db: db}
}

func (r *ppiInklusiRepository) Create(ppiInklusi *models.PPIInklusi) error {
	return r.db.Create(ppiInklusi).Error
}

func (r *ppiInklusiRepository) GetByID(id uint) (*models.PPIInklusi, error) {
	var ppiInklusi models.PPIInklusi
	if err := r.db.Preload("AnakInklusiRombel").
		Preload("AnakInklusiRombel.AnakInklusi").
		Preload("AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("AnakInklusiRombel.PesertaDidikRombel").
		Preload("AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Preload("AnakInklusiRombel.GuruPendampingKhusus").
		Preload("AnakInklusiRombel.GuruKelas").
		Preload("TahunPelajaran").
		Preload("GuruPembuat").
		First(&ppiInklusi, id).Error; err != nil {
		return nil, err
	}
	return &ppiInklusi, nil
}

func (r *ppiInklusiRepository) GetAllWithFilter(params GetPPIInklusiParams) ([]models.PPIInklusi, int64, error) {
	var ppiInklusiList []models.PPIInklusi
	var total int64

	query := r.db.Model(&models.PPIInklusi{}).
		Preload("AnakInklusiRombel").
		Preload("AnakInklusiRombel.AnakInklusi").
		Preload("AnakInklusiRombel.AnakInklusi.PesertaDidik").
		Preload("AnakInklusiRombel.PesertaDidikRombel").
		Preload("AnakInklusiRombel.PesertaDidikRombel.Rombel").
		Preload("AnakInklusiRombel.GuruPendampingKhusus").
		Preload("AnakInklusiRombel.GuruKelas").
		Preload("TahunPelajaran").
		Preload("GuruPembuat")

	// Filter by tahun_pelajaran_id
	if params.Filter.TahunPelajaranID != "" {
		query = query.Where("ppi_inklusi.tahun_pelajaran_id = ?", params.Filter.TahunPelajaranID)
	}

	// Filter by anak_inklusi_rombel_id
	if params.Filter.AnakInklusiRombelID != "" {
		query = query.Where("ppi_inklusi.anak_inklusi_rombel_id = ?", params.Filter.AnakInklusiRombelID)
	}

	// Filter by aspek_pembelajaran
	if params.Filter.AspekPembelajaran != "" {
		query = query.Where("ppi_inklusi.aspek_pembelajaran = ?", params.Filter.AspekPembelajaran)
	}

	// Filter by guru_pembuat_id
	if params.Filter.GuruPembuatID != "" {
		query = query.Where("ppi_inklusi.guru_pembuat_id = ?", params.Filter.GuruPembuatID)
	}

	// Filter by rombel_id (requires join with anak_inklusi_rombel and peserta_didik_rombel)
	if params.Filter.RombelID != "" {
		query = query.Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
			Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
			Where("pdr.rombel_id = ?", params.Filter.RombelID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Order("ppi_inklusi.created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&ppiInklusiList).Error; err != nil {
		return nil, 0, err
	}

	return ppiInklusiList, total, nil
}

func (r *ppiInklusiRepository) Update(ppiInklusi *models.PPIInklusi) error {
	return r.db.Model(ppiInklusi).
		Select("aspek_pembelajaran", "nama_program", "tujuan_pembelajaran", "strategi_pembelajaran", "target_waktu", "guru_pembuat_id", "status_capaian", "updated_by_id", "updated_at").
		Updates(ppiInklusi).Error
}

func (r *ppiInklusiRepository) Delete(id uint) error {
	return r.db.Delete(&models.PPIInklusi{}, id).Error
}

func (r *ppiInklusiRepository) GetSummaryByTahunPelajaran(filter GetPPIInklusiFilter) (*PPIInklusiSummaryData, error) {
	// Build base query
	query := r.db.Model(&models.PPIInklusi{})

	// Apply filters
	if filter.TahunPelajaranID != "" {
		query = query.Where("ppi_inklusi.tahun_pelajaran_id = ?", filter.TahunPelajaranID)
	}

	if filter.AnakInklusiRombelID != "" {
		query = query.Where("ppi_inklusi.anak_inklusi_rombel_id = ?", filter.AnakInklusiRombelID)
	}

	if filter.AspekPembelajaran != "" {
		query = query.Where("ppi_inklusi.aspek_pembelajaran = ?", filter.AspekPembelajaran)
	}

	if filter.GuruPembuatID != "" {
		query = query.Where("ppi_inklusi.guru_pembuat_id = ?", filter.GuruPembuatID)
	}

	if filter.RombelID != "" {
		query = query.Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
			Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
			Where("pdr.rombel_id = ?", filter.RombelID)
	}

	var totalPPI int64
	
	// Get total PPI with filters
	if err := query.Count(&totalPPI).Error; err != nil {
		return nil, err
	}

	// Get total unique anak_inklusi_rombel
	var totalAnakInklusi int64
	distinctQuery := r.db.Model(&models.PPIInklusi{}).Select("DISTINCT anak_inklusi_rombel_id")
	
	if filter.TahunPelajaranID != "" {
		distinctQuery = distinctQuery.Where("ppi_inklusi.tahun_pelajaran_id = ?", filter.TahunPelajaranID)
	}
	if filter.AnakInklusiRombelID != "" {
		distinctQuery = distinctQuery.Where("ppi_inklusi.anak_inklusi_rombel_id = ?", filter.AnakInklusiRombelID)
	}
	if filter.AspekPembelajaran != "" {
		distinctQuery = distinctQuery.Where("ppi_inklusi.aspek_pembelajaran = ?", filter.AspekPembelajaran)
	}
	if filter.GuruPembuatID != "" {
		distinctQuery = distinctQuery.Where("ppi_inklusi.guru_pembuat_id = ?", filter.GuruPembuatID)
	}
	if filter.RombelID != "" {
		distinctQuery = distinctQuery.Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
			Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
			Where("pdr.rombel_id = ?", filter.RombelID)
	}
	
	if err := distinctQuery.Count(&totalAnakInklusi).Error; err != nil {
		return nil, err
	}

	// Get breakdown by status_capaian with filters
	type StatusCapaianCount struct {
		StatusCapaian string
		Total         int64
	}
	var statusCapaianCounts []StatusCapaianCount
	
	statusQuery := r.db.Model(&models.PPIInklusi{}).
		Select("status_capaian, COUNT(*) as total").
		Group("status_capaian")
	
	if filter.TahunPelajaranID != "" {
		statusQuery = statusQuery.Where("ppi_inklusi.tahun_pelajaran_id = ?", filter.TahunPelajaranID)
	}
	if filter.AnakInklusiRombelID != "" {
		statusQuery = statusQuery.Where("ppi_inklusi.anak_inklusi_rombel_id = ?", filter.AnakInklusiRombelID)
	}
	if filter.AspekPembelajaran != "" {
		statusQuery = statusQuery.Where("ppi_inklusi.aspek_pembelajaran = ?", filter.AspekPembelajaran)
	}
	if filter.GuruPembuatID != "" {
		statusQuery = statusQuery.Where("ppi_inklusi.guru_pembuat_id = ?", filter.GuruPembuatID)
	}
	if filter.RombelID != "" {
		statusQuery = statusQuery.Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
			Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
			Where("pdr.rombel_id = ?", filter.RombelID)
	}
	
	if err := statusQuery.Scan(&statusCapaianCounts).Error; err != nil {
		return nil, err
	}

	statusBreakdown := make(map[string]int)
	for _, sc := range statusCapaianCounts {
		statusBreakdown[sc.StatusCapaian] = int(sc.Total)
	}

	// Get breakdown by aspek_pembelajaran with filters
	type AspekPembelajaranCount struct {
		AspekPembelajaran string
		Total             int64
	}
	var aspekPembelajaranCounts []AspekPembelajaranCount
	
	aspekQuery := r.db.Model(&models.PPIInklusi{}).
		Select("aspek_pembelajaran, COUNT(*) as total").
		Group("aspek_pembelajaran")
	
	if filter.TahunPelajaranID != "" {
		aspekQuery = aspekQuery.Where("ppi_inklusi.tahun_pelajaran_id = ?", filter.TahunPelajaranID)
	}
	if filter.AnakInklusiRombelID != "" {
		aspekQuery = aspekQuery.Where("ppi_inklusi.anak_inklusi_rombel_id = ?", filter.AnakInklusiRombelID)
	}
	if filter.AspekPembelajaran != "" {
		aspekQuery = aspekQuery.Where("ppi_inklusi.aspek_pembelajaran = ?", filter.AspekPembelajaran)
	}
	if filter.GuruPembuatID != "" {
		aspekQuery = aspekQuery.Where("ppi_inklusi.guru_pembuat_id = ?", filter.GuruPembuatID)
	}
	if filter.RombelID != "" {
		aspekQuery = aspekQuery.Joins("JOIN anak_inklusi_rombel air ON air.id = ppi_inklusi.anak_inklusi_rombel_id").
			Joins("JOIN peserta_didik_rombel pdr ON pdr.id = air.peserta_didik_rombel_id").
			Where("pdr.rombel_id = ?", filter.RombelID)
	}
	
	if err := aspekQuery.Scan(&aspekPembelajaranCounts).Error; err != nil {
		return nil, err
	}

	var aspekPembelajaranBreakdown []AspekPembelajaranSummary
	for _, asp := range aspekPembelajaranCounts {
		aspekPembelajaranBreakdown = append(aspekPembelajaranBreakdown, AspekPembelajaranSummary{
			AspekPembelajaran: asp.AspekPembelajaran,
			Total:             int(asp.Total),
		})
	}

	return &PPIInklusiSummaryData{
		TotalPPI:                   int(totalPPI),
		TotalAnakInklusi:           int(totalAnakInklusi),
		StatusCapaianBreakdown:     statusBreakdown,
		AspekPembelajaranBreakdown: aspekPembelajaranBreakdown,
	}, nil
}

func (r *ppiInklusiRepository) GetTahunPelajaranByID(id uint) (*models.TahunPelajaran, error) {
	var tahunPelajaran models.TahunPelajaran
	if err := r.db.First(&tahunPelajaran, id).Error; err != nil {
		return nil, err
	}
	return &tahunPelajaran, nil
}
