package dtos

// AbsensiSiswaItem represents a single student attendance item in bulk input
type AbsensiSiswaItem struct {
	PesertaDidikID uint   `json:"peserta_didik_id" binding:"required"`
	Status         string `json:"status" binding:"required,oneof=hadir sakit izin alpa"`
	Keterangan     string `json:"keterangan" binding:"omitempty"`
}

// AbsensiManualCreateRequest represents the request for bulk manual attendance input
type AbsensiManualCreateRequest struct {
	RombelID         uint               `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint               `json:"tahun_pelajaran_id" binding:"required"`
	Semester         int                `json:"semester" binding:"required,oneof=1 2"`
	Tanggal          string             `json:"tanggal" binding:"required"` // Format: YYYY-MM-DD
	BidangStudiID    *uint              `json:"bidang_studi_id"` // NULL = guru kelas, NOT NULL = guru mapel
	PertemuanKe      *int               `json:"pertemuan_ke"`    // NULL = guru kelas, NOT NULL = guru mapel
	WaktuAbsen       string             `json:"waktu_absen" binding:"omitempty"` // Format: YYYY-MM-DD HH:MM:SS
	AbsensiList      []AbsensiSiswaItem `json:"absensi_list" binding:"required,min=1"`
}

// AbsensiManualCreateResponse represents the response for bulk manual attendance input
type AbsensiManualCreateResponse struct {
	TotalSuccess int                       `json:"total_success"`
	TotalFailed  int                       `json:"total_failed"`
	Message      string                    `json:"message"`
	Errors       []AbsensiCreateErrorItem  `json:"errors,omitempty"`
}

// AbsensiCreateErrorItem represents an error for a specific student during bulk input
type AbsensiCreateErrorItem struct {
	PesertaDidikID uint   `json:"peserta_didik_id"`
	Message        string `json:"message"`
}

// AbsensiResponse represents the response for a single absensi record
type AbsensiResponse struct {
	ID               uint   `json:"id"`
	PesertaDidikID   uint   `json:"peserta_didik_id"`
	PesertaDidikNama string `json:"peserta_didik_nama,omitempty"`
	RombelID         *uint  `json:"rombel_id"`
	RombelNama       string `json:"rombel_nama,omitempty"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id"`
	Semester         int    `json:"semester"`
	Tanggal          string `json:"tanggal"`
	BidangStudiID    *uint  `json:"bidang_studi_id"`
	BidangStudiNama  string `json:"bidang_studi_nama,omitempty"`
	PertemuanKe      *int   `json:"pertemuan_ke"`
	Status           string `json:"status"`
	WaktuAbsen       string `json:"waktu_absen,omitempty"`
	MetodeInput      string `json:"metode_input"`
	Keterangan       string `json:"keterangan,omitempty"`
	FileSurat        string `json:"file_surat,omitempty"`
	DicatatOlehID    *uint  `json:"dicatat_oleh_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// AbsensiRekapRequest represents the request for attendance recap
type AbsensiRekapRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`   // Optional: 1 atau 2
	Bulan            *int   `json:"bulan" binding:"omitempty,min=1,max=12"`   // Optional: 1-12
	Tahun            *int   `json:"tahun" binding:"omitempty,min=2020"`       // Optional: tahun kalender
	TanggalMulai     string `json:"tanggal_mulai" binding:"omitempty"`        // Optional: YYYY-MM-DD
	TanggalSelesai   string `json:"tanggal_selesai" binding:"omitempty"`      // Optional: YYYY-MM-DD
	BidangStudiID    *uint  `json:"bidang_studi_id"`                          // Optional: NULL = guru kelas, NOT NULL = guru mapel
}

// AbsensiRekapSiswa represents attendance recap for a single student
type AbsensiRekapSiswa struct {
	PesertaDidikID   uint                    `json:"peserta_didik_id"`
	NIS              string                  `json:"nis"`
	Nama             string                  `json:"nama"`
	JenisKelamin     string                  `json:"jenis_kelamin"`
	TotalHadir       int                     `json:"total_hadir"`
	TotalSakit       int                     `json:"total_sakit"`
	TotalIzin        int                     `json:"total_izin"`
	TotalAlpa        int                     `json:"total_alpa"`
	TotalAbsen       int                     `json:"total_absen"` // sakit + izin + alpa
	TotalPertemuan   int                     `json:"total_pertemuan"`
	PersentaseHadir  float64                 `json:"persentase_hadir"`
	DetailPerTanggal []AbsensiDetailTanggal  `json:"detail_per_tanggal"`
}

// AbsensiDetailTanggal represents attendance detail for a specific date
type AbsensiDetailTanggal struct {
	ID            uint   `json:"id"`
	Tanggal       string `json:"tanggal"`
	PertemuanKe   *int   `json:"pertemuan_ke,omitempty"` // Only for guru mapel (bidang_studi_id not null)
	Status        string `json:"status"`
	WaktuAbsen    string `json:"waktu_absen,omitempty"`
	MetodeInput   string `json:"metode_input"`
	Keterangan    string `json:"keterangan"`
	FileSurat     string `json:"file_surat,omitempty"`
	DicatatOleh   string `json:"dicatat_oleh"`
	DicatatOlehID *uint  `json:"dicatat_oleh_id,omitempty"`
}

// AbsensiRekapResponse represents the response for attendance recap
type AbsensiRekapResponse struct {
	TahunPelajaranID uint                `json:"tahun_pelajaran_id"`
	RombelID         uint                `json:"rombel_id"`
	RombelNama       string              `json:"rombel_nama"`
	Semester         *int                `json:"semester,omitempty"`
	Bulan            *int                `json:"bulan,omitempty"`
	Tahun            *int                `json:"tahun,omitempty"`
	BidangStudiID    *uint               `json:"bidang_studi_id,omitempty"`
	BidangStudiNama  string              `json:"bidang_studi_nama,omitempty"`
	TotalSiswa       int                 `json:"total_siswa"`
	DataSiswa        []AbsensiRekapSiswa `json:"data_siswa"`
}

// AbsensiUpdateRequest represents the request for updating a single absensi record
type AbsensiUpdateRequest struct {
	ID              uint   `json:"id" binding:"required"`
	Status          string `json:"status" binding:"required,oneof=hadir sakit izin alpa"`
	Keterangan      string `json:"keterangan" binding:"omitempty"`
	DeleteFileSurat bool   `json:"delete_file_surat" binding:"omitempty"` // true = hapus file, false/tidak ada = tidak hapus
}

// AbsensiUpdateResponse represents the response for updating absensi
type AbsensiUpdateResponse struct {
	Message string           `json:"message"`
	Data    *AbsensiResponse `json:"data,omitempty"`
}
