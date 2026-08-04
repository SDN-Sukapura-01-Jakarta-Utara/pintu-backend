package dtos

import "time"

// RegisterOrUpdateEkstrakurikulerRequest represents the request payload for registering/updating multiple ekstrakurikuler
type RegisterOrUpdateEkstrakurikulerRequest struct {
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id" binding:"required"`
	EkstrakurikulerIDs   []uint `json:"ekstrakurikuler_ids" binding:"required,min=1"`
}

// GetEkstrakurikulerPesertaDidikRequest represents the request payload for getting student's ekstrakurikuler
type GetEkstrakurikulerPesertaDidikRequest struct {
	PesertaDidikRombelID uint `json:"peserta_didik_rombel_id" binding:"required"`
}

// GetAllEkstrakurikulerSiswaRequest represents the request payload for getting all students' ekstrakurikuler by rombel
type GetAllEkstrakurikulerSiswaRequest struct {
	RombelID         uint `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
}

// SiswaEkstrakurikulerInput represents input for one student's ekstrakurikuler
type SiswaEkstrakurikulerInput struct {
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id" binding:"required"`
	EkstrakurikulerIDs   []uint `json:"ekstrakurikuler_ids" binding:"required"`
}

// RegisterAllEkstrakurikulerSiswaRequest represents the request payload for bulk register/update ekstrakurikuler
type RegisterAllEkstrakurikulerSiswaRequest struct {
	Siswa []SiswaEkstrakurikulerInput `json:"siswa" binding:"required,min=1"`
}

// PesertaDidikEkstrakurikulerResponse represents the response for a single registration
type PesertaDidikEkstrakurikulerResponse struct {
	ID                   uint                     `json:"id"`
	PesertaDidikRombelID uint                     `json:"peserta_didik_rombel_id"`
	EkstrakurikulerID    uint                     `json:"ekstrakurikuler_id"`
	Ekstrakurikuler      *EkstrakurikulerResponse `json:"ekstrakurikuler,omitempty"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
}

// RegisterOrUpdateEkstrakurikulerResponse represents the response for multiple registrations
type RegisterOrUpdateEkstrakurikulerResponse struct {
	Message       string                                `json:"message"`
	Registrations []PesertaDidikEkstrakurikulerResponse `json:"registrations"`
	Summary       struct {
		Added   int `json:"added"`
		Removed int `json:"removed"`
		Kept    int `json:"kept"`
	} `json:"summary"`
}

// GetEkstrakurikulerPesertaDidikResponse represents the response for getting student's ekstrakurikuler
type GetEkstrakurikulerPesertaDidikResponse struct {
	PesertaDidikRombelID uint                                  `json:"peserta_didik_rombel_id"`
	Ekstrakurikuler      []PesertaDidikEkstrakurikulerResponse `json:"ekstrakurikuler"`
	TotalEkskul          int                                   `json:"total_ekskul"`
}

// SiswaEkstrakurikuler represents student with their ekstrakurikuler
type SiswaEkstrakurikuler struct {
	PesertaDidikRombelID uint                                  `json:"peserta_didik_rombel_id"`
	PesertaDidikID       uint                                  `json:"peserta_didik_id"`
	NamaLengkap          string                                `json:"nama_lengkap"`
	NISN                 string                                `json:"nisn"`
	Ekstrakurikuler      []PesertaDidikEkstrakurikulerResponse `json:"ekstrakurikuler"`
	TotalEkskul          int                                   `json:"total_ekskul"`
}

// GetAllEkstrakurikulerSiswaResponse represents the response for getting all students' ekstrakurikuler
type GetAllEkstrakurikulerSiswaResponse struct {
	RombelID         uint                   `json:"rombel_id"`
	TahunPelajaranID uint                   `json:"tahun_pelajaran_id"`
	Siswa            []SiswaEkstrakurikuler `json:"siswa"`
	TotalSiswa       int                    `json:"total_siswa"`
}

// RegisterAllEkstrakurikulerSiswaResponse represents the response for bulk register/update
type RegisterAllEkstrakurikulerSiswaResponse struct {
	Message string `json:"message"`
	Summary struct {
		TotalSiswa       int `json:"total_siswa"`
		TotalAdded       int `json:"total_added"`
		TotalRemoved     int `json:"total_removed"`
		TotalKept        int `json:"total_kept"`
		SuccessCount     int `json:"success_count"`
		FailedCount      int `json:"failed_count"`
	} `json:"summary"`
	Details []struct {
		PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
		Status               string `json:"status"` // success or failed
		Added                int    `json:"added"`
		Removed              int    `json:"removed"`
		Kept                 int    `json:"kept"`
		Error                string `json:"error,omitempty"`
	} `json:"details"`
}


// GetStatistikEkstrakurikulerRequest represents the request for ekstrakurikuler statistics
type GetStatistikEkstrakurikulerRequest struct {
	TahunPelajaranID uint  `json:"tahun_pelajaran_id" binding:"required"`
	RombelID         *uint `json:"rombel_id"` // Optional, null = all rombel
}

// StatistikPerEkskul represents statistics per ekstrakurikuler
type StatistikPerEkskul struct {
	EkstrakurikulerID   uint   `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string `json:"nama_ekstrakurikuler"`
	Kategori            string `json:"kategori"`
	TotalSiswa          int    `json:"total_siswa"`
	Rombel              []struct {
		RombelID   uint   `json:"rombel_id"`
		NamaRombel string `json:"nama_rombel"`
		JumlahSiswa int   `json:"jumlah_siswa"`
	} `json:"rombel"`
}

// StatistikPerRombel represents statistics per rombel
type StatistikPerRombel struct {
	RombelID              uint    `json:"rombel_id"`
	NamaRombel            string  `json:"nama_rombel"`
	TotalSiswa            int     `json:"total_siswa"`
	SiswaIkutEkskul       int     `json:"siswa_ikut_ekskul"`
	SiswaTidakIkutEkskul  int     `json:"siswa_tidak_ikut_ekskul"`
	PersentaseIkutEkskul  float64 `json:"persentase_ikut_ekskul"`
	Ekstrakurikuler       []struct {
		EkstrakurikulerID   uint   `json:"ekstrakurikuler_id"`
		NamaEkstrakurikuler string `json:"nama_ekstrakurikuler"`
		JumlahSiswa         int    `json:"jumlah_siswa"`
	} `json:"ekstrakurikuler"`
}

// SiswaTidakIkutEkskul represents student who doesn't join any ekstrakurikuler
type SiswaTidakIkutEkskul struct {
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
	PesertaDidikID       uint   `json:"peserta_didik_id"`
	NamaLengkap          string `json:"nama_lengkap"`
	NISN                 string `json:"nisn"`
	RombelID             uint   `json:"rombel_id"`
	NamaRombel           string `json:"nama_rombel"`
}

// GetStatistikEkstrakurikulerResponse represents the complete statistics response
type GetStatistikEkstrakurikulerResponse struct {
	TahunPelajaranID uint   `json:"tahun_pelajaran_id"`
	RombelID         *uint  `json:"rombel_id"`
	NamaRombel       string `json:"nama_rombel,omitempty"`
	
	// Overall summary
	Summary struct {
		TotalSiswa               int     `json:"total_siswa"`
		TotalSiswaIkutEkskul     int     `json:"total_siswa_ikut_ekskul"`
		TotalSiswaTidakIkutEkskul int    `json:"total_siswa_tidak_ikut_ekskul"`
		PersentaseIkutEkskul     float64 `json:"persentase_ikut_ekskul"`
		TotalEkstrakurikuler     int     `json:"total_ekstrakurikuler"`
		TotalRombel              int     `json:"total_rombel"`
	} `json:"summary"`
	
	// Statistics per ekstrakurikuler
	StatistikPerEkskul []StatistikPerEkskul `json:"statistik_per_ekskul"`
	
	// Statistics per rombel
	StatistikPerRombel []StatistikPerRombel `json:"statistik_per_rombel"`
	
	// List of students not joining any ekstrakurikuler
	SiswaTidakIkutEkskul []SiswaTidakIkutEkskul `json:"siswa_tidak_ikut_ekskul"`
}


// RekapitulasiPerEkskulRequest represents the request for rekap data per ekstrakurikuler
type RekapitulasiPerEkskulRequest struct {
	EkstrakurikulerID uint `json:"ekstrakurikuler_id"`
	Search            struct {
		Nama             string `json:"nama"`
		NIS              string `json:"nis"`
		RombelID         uint   `json:"rombel_id"`
		TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// SiswaPerEkskul represents student data in ekstrakurikuler rekap
type SiswaPerEkskul struct {
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
	PesertaDidikID       uint   `json:"peserta_didik_id"`
	Nama                 string `json:"nama"`
	NIS                  string `json:"nis"`
	NISN                 string `json:"nisn"`
	RombelID             uint   `json:"rombel_id"`
	NamaRombel           string `json:"nama_rombel"`
	TanggalDaftar        string `json:"tanggal_daftar"`
}

// EkskulWithSiswa represents ekstrakurikuler with list of students
type EkskulWithSiswa struct {
	EkstrakurikulerID   uint             `json:"ekstrakurikuler_id"`
	NamaEkstrakurikuler string           `json:"nama_ekstrakurikuler"`
	Kategori            string           `json:"kategori"`
	TotalSiswa          int              `json:"total_siswa"`
	Siswa               []SiswaPerEkskul `json:"siswa"`
}

// RekapitulasiPerEkskulResponse represents the response for rekap per ekstrakurikuler
type RekapitulasiPerEkskulResponse struct {
	TahunPelajaranID    uint                `json:"tahun_pelajaran_id"`
	Ekstrakurikuler     []EkskulWithSiswa   `json:"ekstrakurikuler"`
	TotalEkstrakurikuler int                `json:"total_ekstrakurikuler"`
	Pagination          PaginationInfo      `json:"pagination"`
}

// RekapitulasiPerRombelRequest represents the request for rekap data per rombel
type RekapitulasiPerRombelRequest struct {
	RombelID         uint   `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Search           struct {
		Nama string `json:"nama"`
		NIS  string `json:"nis"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// SiswaWithEkskul represents student with their ekstrakurikuler list
type SiswaWithEkskul struct {
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
	PesertaDidikID       uint   `json:"peserta_didik_id"`
	Nama                 string `json:"nama"`
	NIS                  string `json:"nis"`
	NISN                 string `json:"nisn"`
	Ekstrakurikuler      []struct {
		EkstrakurikulerID   uint   `json:"ekstrakurikuler_id"`
		NamaEkstrakurikuler string `json:"nama_ekstrakurikuler"`
		Kategori            string `json:"kategori"`
		TanggalDaftar       string `json:"tanggal_daftar"`
	} `json:"ekstrakurikuler"`
	TotalEkskul int `json:"total_ekskul"`
}

// RekapitulasiPerRombelResponse represents the response for rekap per rombel
type RekapitulasiPerRombelResponse struct {
	RombelID         uint              `json:"rombel_id"`
	NamaRombel       string            `json:"nama_rombel"`
	TahunPelajaranID uint              `json:"tahun_pelajaran_id"`
	Siswa            []SiswaWithEkskul `json:"siswa"`
	TotalSiswa       int               `json:"total_siswa"`
	Pagination       PaginationInfo    `json:"pagination"`
}

// ExportExcelPerEkskulRequest represents the request for export excel per ekstrakurikuler
type ExportExcelPerEkskulRequest struct {
	TahunPelajaranID  uint `json:"tahun_pelajaran_id" binding:"required"`
	EkstrakurikulerID uint `json:"ekstrakurikuler_id"`
}

// ExportExcelPerRombelRequest represents the request for export excel per rombel
type ExportExcelPerRombelRequest struct {
	RombelID         uint `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
}
