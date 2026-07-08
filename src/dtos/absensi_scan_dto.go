package dtos

// AbsensiScanRequest represents the request for scanning attendance
type AbsensiScanRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}

// AbsensiScanResponse represents the response after scanning
type AbsensiScanResponse struct {
	Success       bool    `json:"success"`
	Message       string  `json:"message"`
	PesertaDidik  *PesertaDidikInfo `json:"peserta_didik,omitempty"`
	AbsensiInfo   *AbsensiInfo      `json:"absensi_info,omitempty"`
}

// PesertaDidikInfo contains basic student info
type PesertaDidikInfo struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	NISN     string `json:"nisn"`
}

// AbsensiInfo contains attendance information
type AbsensiInfo struct {
	Tanggal    string  `json:"tanggal"`
	JamDatang  *string `json:"jam_datang"`
	JamPulang  *string `json:"jam_pulang"`
	Status     string  `json:"status"` // "tepat_waktu", "terlambat"
	IsUpdate   bool    `json:"is_update"` // true jika update, false jika insert baru
}
