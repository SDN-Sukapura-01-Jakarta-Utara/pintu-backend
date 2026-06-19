package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// LayananSPMBRepository handles data operations for Layanan SPMB
type LayananSPMBRepository interface {
	Create(data *models.LayananSPMB) error
	GetAllWithFilter(params GetLayananSPMBParams) ([]models.LayananSPMB, int64, error)
	GetByID(id uint) (*models.LayananSPMB, error)
	Update(data *models.LayananSPMB) error
	SoftDelete(id uint) error
	GetMonitoringData(params MonitoringParams) (*MonitoringData, error)
}

// GetLayananSPMBFilter represents filter parameters
type GetLayananSPMBFilter struct {
	StartDate    string
	EndDate      string
	NamaOrangTua string
	NamaMurid    string
	Status       string
}

// GetLayananSPMBParams represents query parameters
type GetLayananSPMBParams struct {
	Filter GetLayananSPMBFilter
	Limit  int
	Offset int
}

// MonitoringParams represents monitoring query parameters
type MonitoringParams struct {
	ViewType  string // "daily", "weekly", "monthly", "yearly"
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
}

// MonitoringData represents aggregated monitoring data from database
type MonitoringData struct {
	TotalLayanan     int64
	LayananHariIni   int64
	LayananKemarin   int64
	LayananMingguIni int64
	LayananBulanIni  int64
	StatusCounts     []StatusCount
	TrendData        []TrendPoint
	DetailLayanan    []models.LayananSPMB
}

// StatusCount represents count by status
type StatusCount struct {
	Status string
	Count  int64
}

// TrendPoint represents a data point in trend
type TrendPoint struct {
	Date  string
	Count int64
}

type LayananSPMBRepositoryImpl struct {
	db *gorm.DB
}

// NewLayananSPMBRepository creates a new Layanan SPMB repository
func NewLayananSPMBRepository(db *gorm.DB) LayananSPMBRepository {
	return &LayananSPMBRepositoryImpl{db: db}
}

// Create creates a new Layanan SPMB record
func (r *LayananSPMBRepositoryImpl) Create(data *models.LayananSPMB) error {
	return r.db.Create(data).Error
}

// GetAllWithFilter retrieves all Layanan SPMB with filters, sorting, and pagination
func (r *LayananSPMBRepositoryImpl) GetAllWithFilter(params GetLayananSPMBParams) ([]models.LayananSPMB, int64, error) {
	var data []models.LayananSPMB
	var total int64

	query := r.db.Model(&models.LayananSPMB{})

	// Apply filters
	if params.Filter.StartDate != "" {
		query = query.Where("tanggal_laporan >= ?", params.Filter.StartDate)
	}
	if params.Filter.EndDate != "" {
		query = query.Where("tanggal_laporan <= ?", params.Filter.EndDate+" 23:59:59")
	}
	if params.Filter.NamaOrangTua != "" {
		query = query.Where("nama_orang_tua ILIKE ?", "%"+params.Filter.NamaOrangTua+"%")
	}
	if params.Filter.NamaMurid != "" {
		query = query.Where("nama_lengkap_murid ILIKE ?", "%"+params.Filter.NamaMurid+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting: Status (pending first) then Tanggal Laporan (DESC)
	query = query.Order(`
		CASE status
			WHEN 'pending' THEN 1
			WHEN 'selesai' THEN 2
			ELSE 3
		END ASC,
		tanggal_laporan DESC
	`)

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit).Offset(params.Offset)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByID retrieves Layanan SPMB by ID
func (r *LayananSPMBRepositoryImpl) GetByID(id uint) (*models.LayananSPMB, error) {
	var data models.LayananSPMB
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Layanan SPMB record
func (r *LayananSPMBRepositoryImpl) Update(data *models.LayananSPMB) error {
	return r.db.Save(data).Error
}

// SoftDelete soft deletes Layanan SPMB by setting deleted_at
func (r *LayananSPMBRepositoryImpl) SoftDelete(id uint) error {
	return r.db.Delete(&models.LayananSPMB{}, id).Error
}

// GetMonitoringData retrieves all monitoring statistics and data
func (r *LayananSPMBRepositoryImpl) GetMonitoringData(params MonitoringParams) (*MonitoringData, error) {
	result := &MonitoringData{}

	// Build base query with date range
	baseQuery := func() *gorm.DB {
		q := r.db.Model(&models.LayananSPMB{})
		if params.StartDate != "" && params.EndDate != "" {
			q = q.Where("tanggal_laporan >= ?", params.StartDate).
				Where("tanggal_laporan <= ?", params.EndDate+" 23:59:59")
		}
		return q
	}

	// 1. Get Total Layanan (in date range)
	if err := baseQuery().Count(&result.TotalLayanan).Error; err != nil {
		return nil, err
	}

	// 2. Get Layanan Hari Ini
	if err := r.db.Model(&models.LayananSPMB{}).
		Where("DATE(tanggal_laporan) = CURRENT_DATE").
		Count(&result.LayananHariIni).Error; err != nil {
		return nil, err
	}

	// 3. Get Layanan Kemarin
	if err := r.db.Model(&models.LayananSPMB{}).
		Where("DATE(tanggal_laporan) = CURRENT_DATE - INTERVAL '1 day'").
		Count(&result.LayananKemarin).Error; err != nil {
		return nil, err
	}

	// 4. Get Layanan Minggu Ini (Monday to Sunday)
	if err := r.db.Model(&models.LayananSPMB{}).
		Where("tanggal_laporan >= DATE_TRUNC('week', CURRENT_DATE)").
		Where("tanggal_laporan < DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '1 week'").
		Count(&result.LayananMingguIni).Error; err != nil {
		return nil, err
	}

	// 5. Get Layanan Bulan Ini
	if err := r.db.Model(&models.LayananSPMB{}).
		Where("DATE_TRUNC('month', tanggal_laporan) = DATE_TRUNC('month', CURRENT_DATE)").
		Count(&result.LayananBulanIni).Error; err != nil {
		return nil, err
	}

	// 6. Get Count By Status (in date range)
	statusQuery := baseQuery().
		Select("status, COUNT(*) as count").
		Group("status").
		Order("status ASC")
	
	var statusCounts []StatusCount
	if err := statusQuery.Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	result.StatusCounts = statusCounts

	// 7. Get Trend Data based on view_type
	var trendPoints []TrendPoint
	
	switch params.ViewType {
	case "daily":
		// Daily trend: group by date
		trendQuery := baseQuery().
			Select("DATE(tanggal_laporan) as date, COUNT(*) as count").
			Group("DATE(tanggal_laporan)").
			Order("date ASC")
		if err := trendQuery.Scan(&trendPoints).Error; err != nil {
			return nil, err
		}

	case "weekly":
		// Weekly trend: group by week
		trendQuery := baseQuery().
			Select("DATE_TRUNC('week', tanggal_laporan) as date, COUNT(*) as count").
			Group("DATE_TRUNC('week', tanggal_laporan)").
			Order("date ASC")
		if err := trendQuery.Scan(&trendPoints).Error; err != nil {
			return nil, err
		}

	case "monthly":
		// Monthly trend: group by month
		trendQuery := baseQuery().
			Select("DATE_TRUNC('month', tanggal_laporan) as date, COUNT(*) as count").
			Group("DATE_TRUNC('month', tanggal_laporan)").
			Order("date ASC")
		if err := trendQuery.Scan(&trendPoints).Error; err != nil {
			return nil, err
		}

	case "yearly":
		// Yearly trend: group by year
		trendQuery := baseQuery().
			Select("DATE_TRUNC('year', tanggal_laporan) as date, COUNT(*) as count").
			Group("DATE_TRUNC('year', tanggal_laporan)").
			Order("date ASC")
		if err := trendQuery.Scan(&trendPoints).Error; err != nil {
			return nil, err
		}

	default:
		// Default to daily
		trendQuery := baseQuery().
			Select("DATE(tanggal_laporan) as date, COUNT(*) as count").
			Group("DATE(tanggal_laporan)").
			Order("date ASC")
		if err := trendQuery.Scan(&trendPoints).Error; err != nil {
			return nil, err
		}
	}

	result.TrendData = trendPoints

	// 8. Get Detail Layanan (limited to last 50 in date range, sorted by status and date)
	detailQuery := baseQuery().
		Order(`CASE status WHEN 'pending' THEN 1 WHEN 'selesai' THEN 2 ELSE 3 END ASC, tanggal_laporan DESC`).
		Limit(50)
	
	var details []models.LayananSPMB
	if err := detailQuery.Find(&details).Error; err != nil {
		return nil, err
	}
	result.DetailLayanan = details

	return result, nil
}
