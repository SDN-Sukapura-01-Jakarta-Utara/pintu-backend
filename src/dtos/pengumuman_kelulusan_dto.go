package dtos

// PengumumanKelulusanConfigRequest represents the request for configuring pengumuman kelulusan
type PengumumanKelulusanConfigRequest struct {
	ID                         *uint  `json:"id"` // Optional: if provided, update; if not, create
	SambutanKelulusan          string `json:"sambutan_kelulusan" binding:"required"`
	TanggalPengumumanNilai     string `json:"tanggal_pengumuman_nilai" binding:"required"`     // Format: YYYY-MM-DD HH:MM:SS
	TanggalPengumumanKelulusan string `json:"tanggal_pengumuman_kelulusan" binding:"required"` // Format: YYYY-MM-DD HH:MM:SS
	NamaKepsek                 string `json:"nama_kepsek"`                                      // Optional
	DeleteFotoKepsek           bool   `json:"delete_foto_kepsek"`                              // true = hapus foto kepsek
	DeleteTtdKepsek            bool   `json:"delete_ttd_kepsek"`                               // true = hapus ttd kepsek
}

// PengumumanKelulusanResponse represents the response for pengumuman kelulusan
type PengumumanKelulusanResponse struct {
	ID                         uint   `json:"id"`
	SambutanKelulusan          string `json:"sambutan_kelulusan"`
	TanggalPengumumanNilai     string `json:"tanggal_pengumuman_nilai"`
	TanggalPengumumanKelulusan string `json:"tanggal_pengumuman_kelulusan"`
	FotoKepsek                 string `json:"foto_kepsek,omitempty"`
	TtdKepsek                  string `json:"ttd_kepsek,omitempty"`
	NamaKepsek                 string `json:"nama_kepsek,omitempty"`
	CreatedAt                  string `json:"created_at"`
	UpdatedAt                  string `json:"updated_at"`
	CreatedByID                *uint  `json:"created_by_id,omitempty"`
	UpdatedByID                *uint  `json:"updated_by_id,omitempty"`
}
