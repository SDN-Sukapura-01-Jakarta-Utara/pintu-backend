package dtos

import "time"

// AbsensiEkskulStudentRequest represents individual student attendance
type AbsensiEkskulStudentRequest struct {
	PesertaDidikRombelID uint    `json:"peserta_didik_rombel_id" binding:"required"`
	Status               string  `json:"status" binding:"required,oneof=hadir sakit izin alpa alpha"`
	Keterangan           *string `json:"keterangan"`
}

// AbsensiPelatihRequest represents individual pelatih attendance
type AbsensiPelatihRequest struct {
	PelatihID uint `json:"pelatih_id" binding:"required"`
}

// AbsensiEkskulCreateRequest represents request to create ekstrakurikuler attendance
type AbsensiEkskulCreateRequest struct {
	// Kegiatan data
	EkstrakurikulerID uint      `json:"ekstrakurikuler_id" binding:"required"`
	TahunPelajaranID  uint      `json:"tahun_pelajaran_id" binding:"required"`
	TanggalKegiatan   time.Time `json:"tanggal_kegiatan" binding:"required"`
	WaktuMulai        *string   `json:"waktu_mulai"`
	WaktuSelesai      *string   `json:"waktu_selesai"`
	MateriKegiatan    string    `json:"materi_kegiatan" binding:"required"`
	
	// Bulk attendance data
	AbsensiSiswa   []AbsensiEkskulStudentRequest `json:"absensi_siswa" binding:"required,min=1,dive"`
	AbsensiPelatih []AbsensiPelatihRequest       `json:"absensi_pelatih" binding:"required,min=1,dive"`
}

// AbsensiEkskulResponse represents response after creating attendance
type AbsensiEkskulResponse struct {
	KegiatanEkskulID   uint   `json:"kegiatan_ekskul_id"`
	EkstrakurikulerID  uint   `json:"ekstrakurikuler_id"`
	TahunPelajaranID   uint   `json:"tahun_pelajaran_id"`
	TanggalKegiatan    string `json:"tanggal_kegiatan"`
	MateriKegiatan     string `json:"materi_kegiatan"`
	TotalSiswaHadir    int    `json:"total_siswa_hadir"`
	TotalSiswaSakit    int    `json:"total_siswa_sakit"`
	TotalSiswaIzin     int    `json:"total_siswa_izin"`
	TotalSiswaAlpha    int    `json:"total_siswa_alpha"`
	TotalPelatihHadir  int    `json:"total_pelatih_hadir"`
}

// AbsensiEkskulGetRequest represents request to get attendance by ekstrakurikuler and tahun pelajaran
type AbsensiEkskulGetRequest struct {
	EkstrakurikulerID uint   `json:"ekstrakurikuler_id" binding:"required"`
	TahunPelajaranID  uint   `json:"tahun_pelajaran_id" binding:"required"`
	Nama              string `json:"nama"` // Optional filter by student name
	RombelID          *uint  `json:"rombel_id"` // Optional filter by rombel
	Bulan             *int   `json:"bulan"` // Optional filter by month (1-12)
	Tahun             *int   `json:"tahun"` // Optional filter by year (e.g., 2026)
}

// AbsensiSiswaDetail represents individual student attendance detail
type AbsensiSiswaDetail struct {
	ID                   uint    `json:"id"`
	PesertaDidikRombelID uint    `json:"peserta_didik_rombel_id"`
	NamaSiswa            string  `json:"nama_siswa"`
	NIS                  string  `json:"nis"`
	NISN                 string  `json:"nisn"`
	NamaKelas            string  `json:"nama_kelas"`
	NamaRombel           string  `json:"nama_rombel"`
	Status               string  `json:"status"`
	Keterangan           *string `json:"keterangan"`
}

// AbsensiPelatihDetail represents individual coach attendance detail
type AbsensiPelatihDetail struct {
	ID           uint   `json:"id"`
	PelatihID    uint   `json:"pelatih_id"`
	NamaPelatih  string `json:"nama_pelatih"`
	Telepon      string `json:"telepon"`
}

// KegiatanEkskulDetail represents activity session with attendance details
type KegiatanEkskulDetail struct {
	ID                 uint                   `json:"id"`
	TanggalKegiatan    string                 `json:"tanggal_kegiatan"`
	WaktuMulai         *string                `json:"waktu_mulai"`
	WaktuSelesai       *string                `json:"waktu_selesai"`
	MateriKegiatan     string                 `json:"materi_kegiatan"`
	FotoKegiatan       *string                `json:"foto_kegiatan"`
	AbsensiSiswa       []AbsensiSiswaDetail   `json:"absensi_siswa"`
	AbsensiPelatih     []AbsensiPelatihDetail `json:"absensi_pelatih"`
	TotalSiswaHadir    int                    `json:"total_siswa_hadir"`
	TotalSiswaSakit    int                    `json:"total_siswa_sakit"`
	TotalSiswaIzin     int                    `json:"total_siswa_izin"`
	TotalSiswaAlpha    int                    `json:"total_siswa_alpha"`
	TotalPelatihHadir  int                    `json:"total_pelatih_hadir"`
}

// AbsensiEkskulGetResponse represents response for get attendance
type AbsensiEkskulGetResponse struct {
	EkstrakurikulerID   uint                   `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string                 `json:"nama_ekstrakurikuler"`
	TahunPelajaranID    uint                   `json:"tahun_pelajaran_id"`
	TahunPelajaran      string                 `json:"tahun_pelajaran"`
	TotalKegiatan       int                    `json:"total_kegiatan"`
	Kegiatan            []KegiatanEkskulDetail `json:"kegiatan"`
}

// AbsensiSiswaGetByIDRequest represents request to get single absensi siswa by ID
type AbsensiSiswaGetByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

// AbsensiSiswaDetailResponse represents detailed response for single absensi siswa
type AbsensiSiswaDetailResponse struct {
	ID                   uint    `json:"id"`
	KegiatanEkskulID     uint    `json:"kegiatan_ekskul_id"`
	PesertaDidikRombelID uint    `json:"peserta_didik_rombel_id"`
	NamaSiswa            string  `json:"nama_siswa"`
	NISN                 string  `json:"nisn"`
	NamaKelas            string  `json:"nama_kelas"`
	NamaRombel           string  `json:"nama_rombel"`
	Status               string  `json:"status"`
	Keterangan           *string `json:"keterangan"`
	// Kegiatan info
	TanggalKegiatan      string  `json:"tanggal_kegiatan"`
	WaktuMulai           *string `json:"waktu_mulai"`
	WaktuSelesai         *string `json:"waktu_selesai"`
	MateriKegiatan       string  `json:"materi_kegiatan"`
	// Ekstrakurikuler info
	EkstrakurikulerID    uint   `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler  string `json:"nama_ekstrakurikuler"`
	TahunPelajaranID     uint   `json:"tahun_pelajaran_id"`
	TahunPelajaran       string `json:"tahun_pelajaran"`
}

// AbsensiSiswaUpdateRequest represents request to update absensi siswa
type AbsensiSiswaUpdateRequest struct {
	ID         uint    `json:"id" binding:"required"`
	Status     string  `json:"status" binding:"required,oneof=hadir sakit izin alpa alpha"`
	Keterangan *string `json:"keterangan"`
}

// AbsensiPelatihGetRequest represents request to get pelatih attendance
type AbsensiPelatihGetRequest struct {
	PelatihID         uint  `json:"pelatih_id" binding:"required"`
	TahunPelajaranID  uint  `json:"tahun_pelajaran_id" binding:"required"`
	EkstrakurikulerID *uint `json:"ekstrakurikuler_id"` // Optional filter by ekstrakurikuler
	Bulan             *int  `json:"bulan"` // Optional filter by month (1-12)
	Tahun             *int  `json:"tahun"` // Optional filter by year (e.g., 2026)
}

// KegiatanPelatihDetail represents activity session from pelatih perspective
type KegiatanPelatihDetail struct {
	ID                  uint    `json:"id"`
	EkstrakurikulerID   uint    `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string  `json:"nama_ekstrakurikuler"`
	TanggalKegiatan     string  `json:"tanggal_kegiatan"`
	WaktuMulai          *string `json:"waktu_mulai"`
	WaktuSelesai        *string `json:"waktu_selesai"`
	MateriKegiatan      string  `json:"materi_kegiatan"`
	FotoKegiatan        *string `json:"foto_kegiatan"`
	TotalSiswaHadir     int     `json:"total_siswa_hadir"`
	TotalSiswaSakit     int     `json:"total_siswa_sakit"`
	TotalSiswaIzin      int     `json:"total_siswa_izin"`
	TotalSiswaAlpa      int     `json:"total_siswa_alpa"`
	IsHadir             bool    `json:"is_hadir"` // Pelatih hadir atau tidak
}

// AbsensiPelatihGetResponse represents response for get pelatih attendance
type AbsensiPelatihGetResponse struct {
	PelatihID        uint                    `json:"pelatih_id"`
	NamaPelatih      string                  `json:"nama_pelatih"`
	Telepon          string                  `json:"telepon"`
	TahunPelajaranID uint                    `json:"tahun_pelajaran_id"`
	TahunPelajaran   string                  `json:"tahun_pelajaran"`
	TotalKegiatan    int                     `json:"total_kegiatan"`
	TotalHadir       int                     `json:"total_hadir"`
	TotalTidakHadir  int                     `json:"total_tidak_hadir"`
	Kegiatan         []KegiatanPelatihDetail `json:"kegiatan"`
}

// KegiatanEkskulGetRequest represents request to get kegiatan by ekstrakurikuler
type KegiatanEkskulGetRequest struct {
	EkstrakurikulerID uint `json:"ekstrakurikuler_id" binding:"required"`
	TahunPelajaranID  uint `json:"tahun_pelajaran_id" binding:"required"`
	Bulan             *int `json:"bulan"` // Optional filter by month (1-12)
	Tahun             *int `json:"tahun"` // Optional filter by year (e.g., 2026)
}

// KegiatanEkskulItem represents single kegiatan item in list
type KegiatanEkskulItem struct {
	ID                uint     `json:"id"`
	TanggalKegiatan   string   `json:"tanggal_kegiatan"`
	WaktuMulai        *string  `json:"waktu_mulai"`
	WaktuSelesai      *string  `json:"waktu_selesai"`
	MateriKegiatan    string   `json:"materi_kegiatan"`
	FotoKegiatan      *string  `json:"foto_kegiatan"`
	TotalSiswa        int      `json:"total_siswa"`
	TotalSiswaHadir   int      `json:"total_siswa_hadir"`
	TotalSiswaSakit   int      `json:"total_siswa_sakit"`
	TotalSiswaIzin    int      `json:"total_siswa_izin"`
	TotalSiswaAlpa    int      `json:"total_siswa_alpa"`
	TotalPelatihHadir int      `json:"total_pelatih_hadir"`
	PelatihHadir      []string `json:"pelatih_hadir"` // Array of pelatih names
}

// KegiatanEkskulGetResponse represents response for get kegiatan
type KegiatanEkskulGetResponse struct {
	EkstrakurikulerID   uint                 `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string               `json:"nama_ekstrakurikuler"`
	TahunPelajaranID    uint                 `json:"tahun_pelajaran_id"`
	TahunPelajaran      string               `json:"tahun_pelajaran"`
	TotalKegiatan       int                  `json:"total_kegiatan"`
	Kegiatan            []KegiatanEkskulItem `json:"kegiatan"`
}

// AbsensiPelatihUpdateRequest represents request to update/toggle pelatih attendance
type AbsensiPelatihUpdateRequest struct {
	KegiatanEkskulID uint `json:"kegiatan_ekskul_id" binding:"required"`
	PelatihID        uint `json:"pelatih_id" binding:"required"`
	Status           bool `json:"status"` // true = hadir (create), false = tidak hadir (delete)
}

// AbsensiPelatihUpdateResponse represents response after updating pelatih attendance
type AbsensiPelatihUpdateResponse struct {
	KegiatanEkskulID uint   `json:"kegiatan_ekskul_id"`
	PelatihID        uint   `json:"pelatih_id"`
	NamaPelatih      string `json:"nama_pelatih"`
	IsHadir          bool   `json:"is_hadir"`
	Message          string `json:"message"`
}

// FotoKegiatanUploadResponse represents response after uploading foto kegiatan
type FotoKegiatanUploadResponse struct {
	KegiatanEkskulID uint     `json:"kegiatan_ekskul_id"`
	FotoUrls         []string `json:"foto_urls"`
	TotalFoto        int      `json:"total_foto"`
	UploadedFoto     int      `json:"uploaded_foto"`
	DeletedFoto      int      `json:"deleted_foto"`
	Message          string   `json:"message"`
}

// KegiatanEkskulUpdateRequest represents request to update kegiatan ekstrakurikuler
type KegiatanEkskulUpdateRequest struct {
	ID              uint     `json:"id"`
	TanggalKegiatan *string  `json:"tanggal_kegiatan"` // Format: YYYY-MM-DD
	WaktuMulai      *string  `json:"waktu_mulai"`      // Format: HH:MM:SS
	WaktuSelesai    *string  `json:"waktu_selesai"`    // Format: HH:MM:SS
	MateriKegiatan  *string  `json:"materi_kegiatan"`
	FotoToDelete    []string `json:"foto_to_delete"`   // URLs to delete
}

// KegiatanEkskulUpdateResponse represents response after updating kegiatan
type KegiatanEkskulUpdateResponse struct {
	ID              uint     `json:"id"`
	TanggalKegiatan string   `json:"tanggal_kegiatan"`
	WaktuMulai      *string  `json:"waktu_mulai"`
	WaktuSelesai    *string  `json:"waktu_selesai"`
	MateriKegiatan  string   `json:"materi_kegiatan"`
	FotoUrls        []string `json:"foto_urls"`
	TotalFoto       int      `json:"total_foto"`
	UploadedFoto    int      `json:"uploaded_foto"`
	DeletedFoto     int      `json:"deleted_foto"`
	Message         string   `json:"message"`
}

// KegiatanEkskulGetByIDRequest represents request to get kegiatan by ID
type KegiatanEkskulGetByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

// KegiatanEkskulDetailResponse represents detailed response for single kegiatan
type KegiatanEkskulDetailResponse struct {
	ID                  uint                   `json:"id"`
	EkstrakurikulerID   uint                   `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string                 `json:"nama_ekstrakurikuler"`
	TahunPelajaranID    uint                   `json:"tahun_pelajaran_id"`
	TahunPelajaran      string                 `json:"tahun_pelajaran"`
	TanggalKegiatan     string                 `json:"tanggal_kegiatan"`
	WaktuMulai          *string                `json:"waktu_mulai"`
	WaktuSelesai        *string                `json:"waktu_selesai"`
	MateriKegiatan      string                 `json:"materi_kegiatan"`
	FotoKegiatan        []string               `json:"foto_kegiatan"`
	AbsensiSiswa        []AbsensiSiswaDetail   `json:"absensi_siswa"`
	AbsensiPelatih      []AbsensiPelatihDetail `json:"absensi_pelatih"`
	TotalSiswa          int                    `json:"total_siswa"`
	TotalSiswaHadir     int                    `json:"total_siswa_hadir"`
	TotalSiswaSakit     int                    `json:"total_siswa_sakit"`
	TotalSiswaIzin      int                    `json:"total_siswa_izin"`
	TotalSiswaAlpa      int                    `json:"total_siswa_alpa"`
	TotalPelatihHadir   int                    `json:"total_pelatih_hadir"`
}

// AbsensiPelatihExportRequest represents request for exporting pelatih attendance (all pelatih)
type AbsensiPelatihExportRequest struct {
	TahunPelajaranID  uint  `json:"tahun_pelajaran_id" binding:"required"`
	PelatihID         *uint `json:"pelatih_id"`         // Optional filter by specific pelatih
	EkstrakurikulerID *uint `json:"ekstrakurikuler_id"` // Optional filter by ekstrakurikuler
	Bulan             *int  `json:"bulan"`              // Optional filter by month (1-12)
	Tahun             *int  `json:"tahun"`              // Optional filter by year (e.g., 2026)
}
