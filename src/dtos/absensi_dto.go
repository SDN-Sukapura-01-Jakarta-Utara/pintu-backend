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

// ============================================
// Dashboard Monitoring DTOs
// ============================================

// DashboardSummaryRequest represents request for dashboard summary
type DashboardSummaryRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         *uint  `json:"rombel_id" binding:"omitempty"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	TanggalMulai     string `json:"tanggal_mulai" binding:"omitempty"`
	TanggalSelesai   string `json:"tanggal_selesai" binding:"omitempty"`
}

// DashboardSummaryResponse represents response for dashboard summary
type DashboardSummaryResponse struct {
	TotalSiswa      int                    `json:"total_siswa"`
	TotalPertemuan  int                    `json:"total_pertemuan"`
	Summary         SummaryKehadiran       `json:"summary"`
	Trend           *TrendKehadiran        `json:"trend,omitempty"`
}

// SummaryKehadiran represents attendance summary
type SummaryKehadiran struct {
	TotalHadir          int     `json:"total_hadir"`
	TotalSakit          int     `json:"total_sakit"`
	TotalIzin           int     `json:"total_izin"`
	TotalAlpa           int     `json:"total_alpa"`
	PersentaseKehadiran float64 `json:"persentase_kehadiran"`
}

// TrendKehadiran represents attendance trend
type TrendKehadiran struct {
	HadirKemarin string `json:"hadir_kemarin"`
	HadirHariIni string `json:"hadir_hari_ini"`
	Perubahan    string `json:"perubahan"`
}

// GrafikKehadiranRequest represents request for attendance chart
type GrafikKehadiranRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         *uint  `json:"rombel_id" binding:"omitempty"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	Periode          string `json:"periode" binding:"required,oneof=harian mingguan bulanan"`
	TanggalMulai     string `json:"tanggal_mulai" binding:"required"`
	TanggalSelesai   string `json:"tanggal_selesai" binding:"required"`
}

// GrafikKehadiranResponse represents response for attendance chart
type GrafikKehadiranResponse struct {
	Labels   []string              `json:"labels"`
	Datasets []DatasetKehadiran    `json:"datasets"`
}

// DatasetKehadiran represents dataset for chart
type DatasetKehadiran struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

// SiswaTerendahRequest represents request for students with lowest attendance
type SiswaTerendahRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         *uint  `json:"rombel_id" binding:"omitempty"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	Limit            int    `json:"limit" binding:"omitempty,min=1,max=50"`
	TanggalMulai     string `json:"tanggal_mulai" binding:"omitempty"`
	TanggalSelesai   string `json:"tanggal_selesai" binding:"omitempty"`
}

// SiswaTerendahResponse represents response for students with lowest attendance
type SiswaTerendahResponse struct {
	Data []SiswaKehadiran `json:"data"`
}

// SiswaKehadiran represents student attendance summary
type SiswaKehadiran struct {
	PesertaDidikID  uint    `json:"peserta_didik_id"`
	NIS             string  `json:"nis"`
	Nama            string  `json:"nama"`
	TotalHadir      int     `json:"total_hadir"`
	TotalAbsen      int     `json:"total_absen"`
	TotalPertemuan  int     `json:"total_pertemuan"`
	PersentaseHadir float64 `json:"persentase_hadir"`
}

// PerbandinganRombelRequest represents request for class comparison
type PerbandinganRombelRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	TanggalMulai     string `json:"tanggal_mulai" binding:"omitempty"`
	TanggalSelesai   string `json:"tanggal_selesai" binding:"omitempty"`
}

// PerbandinganRombelResponse represents response for class comparison
type PerbandinganRombelResponse struct {
	Data []RombelKehadiran `json:"data"`
}

// RombelKehadiran represents class attendance summary
type RombelKehadiran struct {
	RombelID        uint    `json:"rombel_id"`
	RombelNama      string  `json:"rombel_nama"`
	TotalSiswa      int     `json:"total_siswa"`
	PersentaseHadir float64 `json:"persentase_hadir"`
	TotalHadir      int     `json:"total_hadir"`
	TotalSakit      int     `json:"total_sakit"`
	TotalIzin       int     `json:"total_izin"`
	TotalAlpa       int     `json:"total_alpa"`
}

// StatistikPerHariRequest represents request for daily statistics
type StatistikPerHariRequest struct {
	TahunPelajaranID uint  `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int  `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         *uint `json:"rombel_id" binding:"omitempty"`
	BidangStudiID    *uint `json:"bidang_studi_id" binding:"omitempty"`
	Bulan            int   `json:"bulan" binding:"required,min=1,max=12"`
	Tahun            int   `json:"tahun" binding:"required,min=2020"`
}

// StatistikPerHariResponse represents response for daily statistics
type StatistikPerHariResponse struct {
	Data []HariKehadiran `json:"data"`
}

// HariKehadiran represents attendance by day of week
type HariKehadiran struct {
	Hari            string  `json:"hari"`
	RataRataHadir   int     `json:"rata_rata_hadir"`
	RataRataAbsen   int     `json:"rata_rata_absen"`
	PersentaseHadir float64 `json:"persentase_hadir"`
}

// DashboardSiswaRequest represents request for student dashboard
type DashboardSiswaRequest struct {
	PesertaDidikID   uint   `json:"peserta_didik_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	Periode          string `json:"periode" binding:"required,oneof=harian mingguan bulanan"`
	TanggalMulai     string `json:"tanggal_mulai" binding:"omitempty"`
	TanggalSelesai   string `json:"tanggal_selesai" binding:"omitempty"`
}

// DashboardSiswaResponse represents response for student dashboard
type DashboardSiswaResponse struct {
	Siswa           InfoSiswa              `json:"siswa"`
	Summary         SummarySiswa           `json:"summary"`
	Grafik          GrafikBulananSiswa     `json:"grafik"`
	RiwayatAbsensi  []RiwayatAbsensiSiswa  `json:"riwayat_absensi"`
}

// InfoSiswa represents student information
type InfoSiswa struct {
	PesertaDidikID uint   `json:"peserta_didik_id"`
	NIS            string `json:"nis"`
	Nama           string `json:"nama"`
	JenisKelamin   string `json:"jenis_kelamin"`
	RombelNama     string `json:"rombel_nama"`
	Foto           string `json:"foto,omitempty"`
}

// SummarySiswa represents student attendance summary
type SummarySiswa struct {
	TotalPertemuan  int     `json:"total_pertemuan"`
	TotalHadir      int     `json:"total_hadir"`
	TotalSakit      int     `json:"total_sakit"`
	TotalIzin       int     `json:"total_izin"`
	TotalAlpa       int     `json:"total_alpa"`
	PersentaseHadir float64 `json:"persentase_hadir"`
	StatusKehadiran string  `json:"status_kehadiran"`
}

// GrafikBulananSiswa represents monthly chart for student
type GrafikBulananSiswa struct {
	Labels []string `json:"labels"`
	Hadir  []int    `json:"hadir"`
	Sakit  []int    `json:"sakit"`
	Izin   []int    `json:"izin"`
	Alpa   []int    `json:"alpa"`
}

// RiwayatAbsensiSiswa represents student attendance history
type RiwayatAbsensiSiswa struct {
	Tanggal     string `json:"tanggal"`
	Hari        string `json:"hari"`
	Status      string `json:"status"`
	WaktuAbsen  string `json:"waktu_absen,omitempty"`
	MetodeInput string `json:"metode_input"`
	Keterangan  string `json:"keterangan"`
	FileSurat   string `json:"file_surat,omitempty"`
	PertemuanKe *int   `json:"pertemuan_ke,omitempty"`
}

// PerbandinganSiswaRequest represents request for student comparison
type PerbandinganSiswaRequest struct {
	PesertaDidikID   uint  `json:"peserta_didik_id" binding:"required"`
	TahunPelajaranID uint  `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int  `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         uint  `json:"rombel_id" binding:"required"`
	BidangStudiID    *uint `json:"bidang_studi_id" binding:"omitempty"`
}

// PerbandinganSiswaResponse represents response for student comparison
type PerbandinganSiswaResponse struct {
	Siswa        InfoPerbandinganSiswa `json:"siswa"`
	Kelas        InfoKelas             `json:"kelas"`
	Perbandingan Perbandingan          `json:"perbandingan"`
}

// InfoPerbandinganSiswa represents student comparison info
type InfoPerbandinganSiswa struct {
	Nama            string  `json:"nama"`
	PersentaseHadir float64 `json:"persentase_hadir"`
	Ranking         int     `json:"ranking"`
}

// InfoKelas represents class info
type InfoKelas struct {
	TotalSiswa          int     `json:"total_siswa"`
	RataRataKehadiran   float64 `json:"rata_rata_kehadiran"`
	Tertinggi           float64 `json:"tertinggi"`
	Terendah            float64 `json:"terendah"`
}

// Perbandingan represents comparison result
type Perbandingan struct {
	SelisihDenganRataRata float64 `json:"selisih_dengan_rata_rata"`
	Status                string  `json:"status"`
}

// TrendSiswaRequest represents request for student trend
type TrendSiswaRequest struct {
	PesertaDidikID   uint   `json:"peserta_didik_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	Periode          string `json:"periode" binding:"required,oneof=harian mingguan bulanan"`
}

// TrendSiswaResponse represents response for student trend
type TrendSiswaResponse struct {
	Labels     []string           `json:"labels"`
	Data       []TrendDataSiswa   `json:"data"`
	Trend      string             `json:"trend"`
	Perubahan  float64            `json:"perubahan"`
}

// TrendDataSiswa represents trend data point
type TrendDataSiswa struct {
	Periode         string  `json:"periode"`
	PersentaseHadir float64 `json:"persentase_hadir"`
	TotalHadir      int     `json:"total_hadir"`
	TotalPertemuan  int     `json:"total_pertemuan"`
}

// DaftarSiswaRequest represents request for student list
type DaftarSiswaRequest struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Semester         *int   `json:"semester" binding:"omitempty,oneof=1 2"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	BidangStudiID    *uint  `json:"bidang_studi_id" binding:"omitempty"`
	Search           string `json:"search" binding:"omitempty"`
	FilterStatus     string `json:"filter_status" binding:"omitempty,oneof=tinggi sedang rendah"`
	SortBy           string `json:"sort_by" binding:"omitempty,oneof=nama nis persentase_hadir"`
	SortOrder        string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
	Limit            int    `json:"limit" binding:"omitempty,min=1,max=100"`
	Offset           int    `json:"offset" binding:"omitempty,min=0"`
}

// DaftarSiswaResponse represents response for student list
type DaftarSiswaResponse struct {
	Data   []DaftarSiswaItem `json:"data"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

// DaftarSiswaItem represents student list item
type DaftarSiswaItem struct {
	PesertaDidikID  uint             `json:"peserta_didik_id"`
	NIS             string           `json:"nis"`
	Nama            string           `json:"nama"`
	JenisKelamin    string           `json:"jenis_kelamin"`
	Foto            string           `json:"foto,omitempty"`
	TotalHadir      int              `json:"total_hadir"`
	TotalAbsen      int              `json:"total_absen"`
	TotalPertemuan  int              `json:"total_pertemuan"`
	PersentaseHadir float64          `json:"persentase_hadir"`
	StatusKehadiran string           `json:"status_kehadiran"`
	AbsenTerakhir   *AbsenTerakhir   `json:"absen_terakhir,omitempty"`
}

// AbsenTerakhir represents last absence info
type AbsenTerakhir struct {
	Tanggal    string `json:"tanggal"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
}
