package dtos

// KegiatanEkskulDownloadWordRequest represents request to download documentation in Word format
type KegiatanEkskulDownloadWordRequest struct {
	EkstrakurikulerID uint `json:"ekstrakurikuler_id" binding:"required"`
	TahunPelajaranID  uint `json:"tahun_pelajaran_id" binding:"required"`
	Bulan             int  `json:"bulan" binding:"required,min=1,max=12"` // Month (1-12)
	Tahun             int  `json:"tahun" binding:"required,min=2020"`     // Year (e.g., 2026)
}
