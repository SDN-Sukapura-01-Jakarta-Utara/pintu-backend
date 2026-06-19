package dtos

// LayananSPMBCreateRequest represents the request payload for creating Layanan SPMB (public)
type LayananSPMBCreateRequest struct {
	NamaOrangTua      string `json:"nama_orang_tua" binding:"required"`
	NomorTelepon      string `json:"nomor_telepon" binding:"required"`
	Alamat            string `json:"alamat" binding:"required"`
	NamaLengkapMurid  string `json:"nama_lengkap_murid" binding:"required"`
	Keperluan         string `json:"keperluan" binding:"required"`
}

// LayananSPMBResponse represents the response payload for Layanan SPMB
type LayananSPMBResponse struct {
	ID               uint   `json:"id"`
	NamaOrangTua     string `json:"nama_orang_tua"`
	NomorTelepon     string `json:"nomor_telepon"`
	Alamat           string `json:"alamat"`
	NamaLengkapMurid string `json:"nama_lengkap_murid"`
	Keperluan        string `json:"keperluan"`
	TanggalLaporan   string `json:"tanggal_laporan"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// LayananSPMBGetAllRequest represents the request for getting all layanan SPMB with filters
type LayananSPMBGetAllRequest struct {
	Search struct {
		StartDate    string `json:"start_date"` // YYYY-MM-DD
		EndDate      string `json:"end_date"`   // YYYY-MM-DD
		NamaOrangTua string `json:"nama_orang_tua"`
		NamaMurid    string `json:"nama_murid"`
		Status       string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// LayananSPMBListWithPaginationResponse represents paginated list response
type LayananSPMBListWithPaginationResponse struct {
	Data       []LayananSPMBResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

// LayananSPMBUpdateStatusRequest represents the request for updating status
type LayananSPMBUpdateStatusRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// SettingLayananSPMBRequest represents the request for setting layanan SPMB
type SettingLayananSPMBRequest struct {
	NamaKepalaSekolah string `json:"nama_kepala_sekolah"`
	NIPKepalaSekolah  string `json:"nip_kepala_sekolah"`
	NamaKetuaPanitia  string `json:"nama_ketua_panitia"`
	NIPKetuaPanitia   string `json:"nip_ketua_panitia"`
	GrupWA            string `json:"grup_wa"`
}

// SettingLayananSPMBResponse represents the response for setting layanan SPMB
type SettingLayananSPMBResponse struct {
	ID                uint    `json:"id"`
	NamaKepalaSekolah *string `json:"nama_kepala_sekolah"`
	NIPKepalaSekolah  *string `json:"nip_kepala_sekolah"`
	NamaKetuaPanitia  *string `json:"nama_ketua_panitia"`
	NIPKetuaPanitia   *string `json:"nip_ketua_panitia"`
	GrupWA            *string `json:"grup_wa"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// GrupWASPMBResponse represents the response for public grup WA
type GrupWASPMBResponse struct {
	GrupWA *string `json:"grup_wa"`
}

// MonitoringPelayananSPMBRequest represents the request for monitoring dashboard
type MonitoringPelayananSPMBRequest struct {
	ViewType  string `json:"view_type"`  // "daily", "weekly", "monthly", "yearly" (default: "daily")
	StartDate string `json:"start_date"` // YYYY-MM-DD (optional, for custom range)
	EndDate   string `json:"end_date"`   // YYYY-MM-DD (optional, for custom range)
}

// MonitoringStatistik represents overview statistics
type MonitoringStatistik struct {
	TotalLayanan     int64                   `json:"total_layanan"`
	LayananHariIni   int64                   `json:"layanan_hari_ini"`
	LayananKemarin   int64                   `json:"layanan_kemarin"`
	TrendPercentage  float64                 `json:"trend_percentage"`  // Percentage change hari ini vs kemarin
	TrendDirection   string                  `json:"trend_direction"`   // "up", "down", "stable"
	LayananMingguIni int64                   `json:"layanan_minggu_ini"`
	LayananBulanIni  int64                   `json:"layanan_bulan_ini"`
	ByStatus         []MonitoringStatusCount `json:"by_status"`
}

// MonitoringStatusCount represents count by status
type MonitoringStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// MonitoringTrend represents trend data with view type context
type MonitoringTrend struct {
	ViewType string              `json:"view_type"` // "daily", "weekly", "monthly", "yearly"
	Period   string              `json:"period"`    // Description of period
	Data     []MonitoringTrendData `json:"data"`
}

// MonitoringTrendData represents individual trend data point
type MonitoringTrendData struct {
	Label     string `json:"label"`      // Display label (e.g., "1 Jun", "Week 1", "January")
	Date      string `json:"date"`       // YYYY-MM-DD (for daily)
	DateRange string `json:"date_range,omitempty"` // For weekly view (e.g., "1-7 Jun")
	Month     string `json:"month,omitempty"`      // YYYY-MM (for monthly)
	Year      string `json:"year,omitempty"`       // YYYY (for yearly)
	Count     int64  `json:"count"`
}

// MonitoringDetailLayanan represents simplified layanan info for monitoring
type MonitoringDetailLayanan struct {
	ID               uint   `json:"id"`
	NamaOrangTua     string `json:"nama_orang_tua"`
	NamaLengkapMurid string `json:"nama_lengkap_murid"`
	Keperluan        string `json:"keperluan"`
	TanggalLaporan   string `json:"tanggal_laporan"`
	Status           string `json:"status"`
}

// MonitoringPelayananSPMBResponse represents the complete monitoring response
type MonitoringPelayananSPMBResponse struct {
	Statistik     MonitoringStatistik       `json:"statistik"`
	Trend         MonitoringTrend           `json:"trend"`
	DetailLayanan []MonitoringDetailLayanan `json:"detail_layanan"`
}
