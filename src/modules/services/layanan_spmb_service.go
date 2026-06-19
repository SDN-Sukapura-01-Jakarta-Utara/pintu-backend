package services

import (
	"fmt"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// LayananSPMBService handles business logic for Layanan SPMB
type LayananSPMBService interface {
	CreatePublic(req *dtos.LayananSPMBCreateRequest) (*dtos.LayananSPMBResponse, error)
	GetAllWithFilter(req *dtos.LayananSPMBGetAllRequest) (*dtos.LayananSPMBListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.LayananSPMBResponse, error)
	UpdateStatus(req *dtos.LayananSPMBUpdateStatusRequest) (*dtos.LayananSPMBResponse, error)
	DeleteLayananSPMB(id uint) error
	GetMonitoringPelayanan(req *dtos.MonitoringPelayananSPMBRequest) (*dtos.MonitoringPelayananSPMBResponse, error)
}

type LayananSPMBServiceImpl struct {
	repository repositories.LayananSPMBRepository
}

// NewLayananSPMBService creates a new Layanan SPMB service
func NewLayananSPMBService(repository repositories.LayananSPMBRepository) LayananSPMBService {
	return &LayananSPMBServiceImpl{
		repository: repository,
	}
}

// CreatePublic creates a new Layanan SPMB from public form
func (s *LayananSPMBServiceImpl) CreatePublic(req *dtos.LayananSPMBCreateRequest) (*dtos.LayananSPMBResponse, error) {
	// Create layanan_spmb record
	data := &models.LayananSPMB{
		NamaOrangTua:     req.NamaOrangTua,
		NomorTelepon:     req.NomorTelepon,
		Alamat:           req.Alamat,
		NamaLengkapMurid: req.NamaLengkapMurid,
		Keperluan:        req.Keperluan,
		Status:           "pending",
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// mapToResponse maps LayananSPMB model to response DTO
func (s *LayananSPMBServiceImpl) mapToResponse(data *models.LayananSPMB) *dtos.LayananSPMBResponse {
	return &dtos.LayananSPMBResponse{
		ID:               data.ID,
		NamaOrangTua:     data.NamaOrangTua,
		NomorTelepon:     data.NomorTelepon,
		Alamat:           data.Alamat,
		NamaLengkapMurid: data.NamaLengkapMurid,
		Keperluan:        data.Keperluan,
		TanggalLaporan:   data.TanggalLaporan.Format("2006-01-02 15:04:05"),
		Status:           data.Status,
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// GetAllWithFilter retrieves all Layanan SPMB with filters and pagination
func (s *LayananSPMBServiceImpl) GetAllWithFilter(req *dtos.LayananSPMBGetAllRequest) (*dtos.LayananSPMBListWithPaginationResponse, error) {
	// Set default pagination
	limit := 10
	page := 1
	if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
		limit = req.Pagination.Limit
	}
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	offset := (page - 1) * limit

	// Build filter params
	params := repositories.GetLayananSPMBParams{
		Filter: repositories.GetLayananSPMBFilter{
			StartDate:    req.Search.StartDate,
			EndDate:      req.Search.EndDate,
			NamaOrangTua: req.Search.NamaOrangTua,
			NamaMurid:    req.Search.NamaMurid,
			Status:       req.Search.Status,
		},
		Limit:  limit,
		Offset: offset,
	}

	// Get data from repository
	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	var responses []dtos.LayananSPMBResponse
	for _, item := range data {
		responses = append(responses, *s.mapToResponse(&item))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dtos.LayananSPMBListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationMeta{
			Limit:      limit,
			Offset:     offset,
			Page:       page,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetByID retrieves Layanan SPMB by ID
func (s *LayananSPMBServiceImpl) GetByID(id uint) (*dtos.LayananSPMBResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", id)
	}

	return s.mapToResponse(data), nil
}

// UpdateStatus updates status of Layanan SPMB
func (s *LayananSPMBServiceImpl) UpdateStatus(req *dtos.LayananSPMBUpdateStatusRequest) (*dtos.LayananSPMBResponse, error) {
	// Get existing data
	data, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", req.ID)
	}

	// Validate status value
	if req.Status != "pending" && req.Status != "selesai" {
		return nil, fmt.Errorf("status harus 'pending' atau 'selesai'")
	}

	// Update status
	data.Status = req.Status

	// Save to database
	if err := s.repository.Update(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// DeleteLayananSPMB soft deletes Layanan SPMB by ID
func (s *LayananSPMBServiceImpl) DeleteLayananSPMB(id uint) error {
	// Check if record exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("layanan SPMB dengan ID %d tidak ditemukan", id)
	}

	// Soft delete
	if err := s.repository.SoftDelete(id); err != nil {
		return err
	}

	return nil
}

// GetMonitoringPelayanan retrieves monitoring dashboard data
func (s *LayananSPMBServiceImpl) GetMonitoringPelayanan(req *dtos.MonitoringPelayananSPMBRequest) (*dtos.MonitoringPelayananSPMBResponse, error) {
	// Calculate date range based on custom range or view_type default
	params := repositories.MonitoringParams{
		ViewType: req.ViewType,
	}

	// Set default view_type if not provided
	if params.ViewType == "" {
		params.ViewType = "daily"
	}

	if req.StartDate != "" && req.EndDate != "" {
		// Custom date range provided
		params.StartDate = req.StartDate
		params.EndDate = req.EndDate
	} else {
		// Use default range based on view_type
		params.StartDate, params.EndDate = s.getDefaultDateRange(req.ViewType)
	}

	// Get monitoring data from repository
	data, err := s.repository.GetMonitoringData(params)
	if err != nil {
		return nil, err
	}

	// Calculate trend percentage (hari ini vs kemarin)
	var trendPercentage float64
	var trendDirection string
	
	if data.LayananKemarin > 0 {
		diff := float64(data.LayananHariIni - data.LayananKemarin)
		trendPercentage = (diff / float64(data.LayananKemarin)) * 100
	} else if data.LayananHariIni > 0 {
		trendPercentage = 100 // 100% increase if kemarin was 0
	}

	if trendPercentage > 0 {
		trendDirection = "up"
	} else if trendPercentage < 0 {
		trendDirection = "down"
	} else {
		trendDirection = "stable"
	}

	// Map StatusCounts to DTO
	var statusCounts []dtos.MonitoringStatusCount
	for _, sc := range data.StatusCounts {
		statusCounts = append(statusCounts, dtos.MonitoringStatusCount{
			Status: sc.Status,
			Count:  sc.Count,
		})
	}

	// Format Trend Data based on view_type
	trend := s.formatTrendData(data.TrendData, req.ViewType, params.StartDate, params.EndDate)

	// Map DetailLayanan to DTO
	var details []dtos.MonitoringDetailLayanan
	for _, d := range data.DetailLayanan {
		details = append(details, dtos.MonitoringDetailLayanan{
			ID:               d.ID,
			NamaOrangTua:     d.NamaOrangTua,
			NamaLengkapMurid: d.NamaLengkapMurid,
			Keperluan:        d.Keperluan,
			TanggalLaporan:   d.TanggalLaporan.Format("2006-01-02 15:04:05"),
			Status:           d.Status,
		})
	}

	// Build response
	response := &dtos.MonitoringPelayananSPMBResponse{
		Statistik: dtos.MonitoringStatistik{
			TotalLayanan:     data.TotalLayanan,
			LayananHariIni:   data.LayananHariIni,
			LayananKemarin:   data.LayananKemarin,
			TrendPercentage:  trendPercentage,
			TrendDirection:   trendDirection,
			LayananMingguIni: data.LayananMingguIni,
			LayananBulanIni:  data.LayananBulanIni,
			ByStatus:         statusCounts,
		},
		Trend:         trend,
		DetailLayanan: details,
	}

	return response, nil
}

// getDefaultDateRange returns default date range based on view_type
func (s *LayananSPMBServiceImpl) getDefaultDateRange(viewType string) (string, string) {
	now := time.Now()
	
	switch viewType {
	case "daily":
		// Default: Last 30 days
		start := now.AddDate(0, 0, -30)
		return start.Format("2006-01-02"), now.Format("2006-01-02")
		
	case "weekly":
		// Default: Last 12 weeks (84 days)
		start := now.AddDate(0, 0, -84)
		return start.Format("2006-01-02"), now.Format("2006-01-02")
		
	case "monthly":
		// Default: Last 12 months
		start := now.AddDate(-1, 0, 0)
		return start.Format("2006-01-02"), now.Format("2006-01-02")
		
	case "yearly":
		// Default: Last 5 years
		start := now.AddDate(-5, 0, 0)
		return start.Format("2006-01-02"), now.Format("2006-01-02")
		
	default:
		// Fallback: Last 30 days
		start := now.AddDate(0, 0, -30)
		return start.Format("2006-01-02"), now.Format("2006-01-02")
	}
}

// formatTrendData formats trend data based on view_type
func (s *LayananSPMBServiceImpl) formatTrendData(trendData []repositories.TrendPoint, viewType, startDate, endDate string) dtos.MonitoringTrend {
	trend := dtos.MonitoringTrend{
		ViewType: viewType,
		Data:     []dtos.MonitoringTrendData{},
	}

	// Set period description
	if startDate != "" && endDate != "" {
		trend.Period = fmt.Sprintf("%s to %s", startDate, endDate)
	} else {
		trend.Period = "Current Period"
	}

	// Create map for quick lookup - normalize date format
	dataMap := make(map[string]int64)
	for _, t := range trendData {
		// Try to parse and normalize the date
		normalizedDate := t.Date
		
		// Handle various date formats from database
		if parsed, err := time.Parse("2006-01-02", t.Date); err == nil {
			normalizedDate = parsed.Format("2006-01-02")
		} else if parsed, err := time.Parse("2006-01-02 15:04:05", t.Date); err == nil {
			normalizedDate = parsed.Format("2006-01-02")
		} else if parsed, err := time.Parse(time.RFC3339, t.Date); err == nil {
			normalizedDate = parsed.Format("2006-01-02")
		}
		
		dataMap[normalizedDate] = t.Count
	}

	// Generate complete date range based on view_type
	var allDates []time.Time
	now := time.Now()

	switch viewType {
	case "daily":
		// Generate all dates for last 30 days (or custom range)
		start := now.AddDate(0, 0, -30)
		if startDate != "" {
			if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
				start = parsed
			}
		}
		end := now
		if endDate != "" {
			if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
				end = parsed
			}
		}

		// Generate all dates from start to end
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			allDates = append(allDates, d)
		}

		// Format each date
		for _, date := range allDates {
			dateStr := date.Format("2006-01-02")
			count := dataMap[dateStr] // Will be 0 if not exists
			
			trend.Data = append(trend.Data, dtos.MonitoringTrendData{
				Label: date.Format("2 Jan"),
				Date:  dateStr,
				Count: count,
			})
		}

	case "weekly":
		// Generate all weeks for last 12 weeks (or custom range)
		start := now.AddDate(0, 0, -84) // 12 weeks ago
		if startDate != "" {
			if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
				start = parsed
			}
		}
		end := now
		if endDate != "" {
			if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
				end = parsed
			}
		}

		// Start from the beginning of week
		for start.Weekday() != time.Monday {
			start = start.AddDate(0, 0, -1)
		}

		// Generate all weeks
		for d := start; !d.After(end); d = d.AddDate(0, 0, 7) {
			weekStart := d
			weekEnd := d.AddDate(0, 0, 6)
			dateStr := weekStart.Format("2006-01-02")
			count := dataMap[dateStr]

			trend.Data = append(trend.Data, dtos.MonitoringTrendData{
				Label:     fmt.Sprintf("Week %s", weekStart.Format("2 Jan")),
				DateRange: fmt.Sprintf("%s - %s", weekStart.Format("2 Jan"), weekEnd.Format("2 Jan")),
				Date:      dateStr,
				Count:     count,
			})
		}

	case "monthly":
		// Generate all months for last 12 months (or custom range)
		start := time.Date(now.Year(), now.Month()-11, 1, 0, 0, 0, 0, now.Location())
		if startDate != "" {
			if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
				start = time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, parsed.Location())
			}
		}
		end := now
		if endDate != "" {
			if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
				end = parsed
			}
		}

		// Generate all months
		for d := start; d.Before(end) || d.Month() == end.Month(); d = d.AddDate(0, 1, 0) {
			monthStr := d.Format("2006-01")
			count := int64(0)
			
			// Check if we have data for this month - match by month prefix
			for key, val := range dataMap {
				if len(key) >= 7 && key[:7] == monthStr {
					count += val // Sum all data for this month
				}
			}

			trend.Data = append(trend.Data, dtos.MonitoringTrendData{
				Label: d.Format("Jan 2006"),
				Month: monthStr,
				Count: count,
			})

			// Break if we've reached the end month
			if d.Year() == end.Year() && d.Month() == end.Month() {
				break
			}
		}

	case "yearly":
		// Generate all years from 3 years ago to current
		startYear := now.Year() - 3
		if startDate != "" {
			if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
				startYear = parsed.Year()
			}
		}
		endYear := now.Year()
		if endDate != "" {
			if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
				endYear = parsed.Year()
			}
		}

		// Generate all years
		for year := startYear; year <= endYear; year++ {
			yearStr := fmt.Sprintf("%d", year)
			count := int64(0)
			
			// Check if we have data for this year - match by year prefix
			for key, val := range dataMap {
				if len(key) >= 4 && key[:4] == yearStr {
					count += val // Sum all data for this year
				}
			}

			trend.Data = append(trend.Data, dtos.MonitoringTrendData{
				Label: yearStr,
				Year:  yearStr,
				Count: count,
			})
		}

	default:
		// Fallback to daily for unknown view_type
		for _, t := range trendData {
			if parsedDate, err := parseDate(t.Date); err == nil {
				trend.Data = append(trend.Data, dtos.MonitoringTrendData{
					Label: parsedDate.Format("2 Jan"),
					Date:  parsedDate.Format("2006-01-02"),
					Count: t.Count,
				})
			}
		}
	}

	return trend
}

// parseDate helper to parse date string
func parseDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
