package dtos

// KonfigurasiAbsensiRequest represents the request for setting konfigurasi absensi
type KonfigurasiAbsensiRequest struct {
	JamDatangMulai   string `json:"jam_datang_mulai" binding:"required"`   // Format: "HH:MM" atau "HH:MM:SS"
	JamMaxDatang     string `json:"jam_max_datang" binding:"required"`     // Format: "HH:MM" atau "HH:MM:SS"
	JamDatangSelesai string `json:"jam_datang_selesai" binding:"required"` // Format: "HH:MM" atau "HH:MM:SS"
	JamPulangMulai   string `json:"jam_pulang_mulai" binding:"required"`   // Format: "HH:MM" atau "HH:MM:SS"
	JamPulangSelesai string `json:"jam_pulang_selesai" binding:"required"` // Format: "HH:MM" atau "HH:MM:SS"
	NamaKepsek       string `json:"nama_kepsek"`
	NIPKepsek        string `json:"nip_kepsek"`
}

// KonfigurasiAbsensiResponse represents the response for konfigurasi absensi
type KonfigurasiAbsensiResponse struct {
	ID               uint    `json:"id"`
	JamDatangMulai   string  `json:"jam_datang_mulai"`
	JamMaxDatang     string  `json:"jam_max_datang"`
	JamDatangSelesai string  `json:"jam_datang_selesai"`
	JamPulangMulai   string  `json:"jam_pulang_mulai"`
	JamPulangSelesai string  `json:"jam_pulang_selesai"`
	NamaKepsek       *string `json:"nama_kepsek"`
	NIPKepsek        *string `json:"nip_kepsek"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
