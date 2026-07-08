package dtos

// AbsensiSyncRequest represents request for synchronizing absensi from scanner to rekapitulasi
type AbsensiSyncRequest struct {
	TipeSync         string `json:"tipe_sync" binding:"required,oneof=tanggal bulan"`          // tanggal or bulan
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`                       // NULL = guru kelas, NOT NULL = guru mapel
	PertemuanKe      *int   `json:"pertemuan_ke" binding:"omitempty"`                          // For guru mapel only
	Tanggal          string `json:"tanggal" binding:"required_if=TipeSync tanggal,omitempty"` // YYYY-MM-DD (required if tipe_sync = tanggal)
	Bulan            int    `json:"bulan" binding:"required_if=TipeSync bulan,omitempty,min=1,max=12"` // 1-12 (required if tipe_sync = bulan, only for guru kelas)
	Tahun            int    `json:"tahun" binding:"required_if=TipeSync bulan,omitempty,min=2000"` // Year (required if tipe_sync = bulan, only for guru kelas)
}

// AbsensiSyncResponse represents response for synchronization
type AbsensiSyncResponse struct {
	TotalProcessed int                    `json:"total_processed"`
	TotalInserted  int                    `json:"total_inserted"`
	TotalUpdated   int                    `json:"total_updated"`
	TotalSkipped   int                    `json:"total_skipped"`
	Message        string                 `json:"message"`
	Details        []AbsensiSyncDetailItem `json:"details,omitempty"`
}

// AbsensiSyncDetailItem represents detail of each sync operation
type AbsensiSyncDetailItem struct {
	PesertaDidikID uint   `json:"peserta_didik_id"`
	NIS            string `json:"nis"`
	Nama           string `json:"nama"`
	Tanggal        string `json:"tanggal"`
	Action         string `json:"action"` // inserted, updated, skipped
	Reason         string `json:"reason,omitempty"`
}
