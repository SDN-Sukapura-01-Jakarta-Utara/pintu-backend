package dtos

type MonitoringInklusiBulananCreateRequest struct {
	PPIInklusiID          uint     `json:"ppi_inklusi_id" binding:"required"`
	Bulan                 int16    `json:"bulan" binding:"required,min=1,max=12"`
	Tahun                 int16    `json:"tahun" binding:"required"`
	DeskripsiPerkembangan string   `json:"deskripsi_perkembangan" binding:"required"`
	KendalaDitemui        string   `json:"kendala_ditemui"`
	TindakLanjut          string   `json:"tindak_lanjut"`
	FilePendukung         []string `json:"file_pendukung"` // Array of base64 encoded files
	GuruPengisiID         *uint    `json:"guru_pengisi_id"`
}

type MonitoringInklusiBulananResponse struct {
	ID                     uint               `json:"id"`
	PPIInklusiID           uint               `json:"ppi_inklusi_id"`
	Bulan                  int16              `json:"bulan"`
	Tahun                  int16              `json:"tahun"`
	DeskripsiPerkembangan  string             `json:"deskripsi_perkembangan"`
	KendalaDitemui         string             `json:"kendala_ditemui"`
	TindakLanjut           string             `json:"tindak_lanjut"`
	FilePendukung          []string           `json:"file_pendukung"`
	FilePendukungURL       []string           `json:"file_pendukung_url"`
	GuruPengisiID          *uint              `json:"guru_pengisi_id"`
	PPIInklusi             *PPIInklusiDetail  `json:"ppi_inklusi,omitempty"`
	GuruPengisi            *KepegawaianSimple `json:"guru_pengisi,omitempty"`
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
}

type PPIInklusiDetail struct {
	ID                uint   `json:"id"`
	AspekPembelajaran string `json:"aspek_pembelajaran"`
	NamaProgram       string `json:"nama_program"`
}

type MonitoringInklusiBulananGetRequest struct {
	Bulan        int16  `json:"bulan" binding:"required,min=1,max=12"`
	Tahun        int16  `json:"tahun" binding:"required"`
	RombelID     *uint  `json:"rombel_id"` // Optional: filter by rombel
	PPIInklusiID *uint  `json:"ppi_inklusi_id"` // Optional: filter by specific PPI
}

type MonitoringInklusiBulananGetAllResponse struct {
	Bulan           int16                              `json:"bulan"`
	Tahun           int16                              `json:"tahun"`
	RombelID        *uint                              `json:"rombel_id,omitempty"`
	RombelName      string                             `json:"rombel_name,omitempty"`
	PPIInklusiID    *uint                              `json:"ppi_inklusi_id,omitempty"`
	TotalMonitoring int                                `json:"total_monitoring"`
	Data            []MonitoringInklusiBulananResponse `json:"data"`
}

type MonitoringInklusiBulananUpdateRequest struct {
	ID                    uint     `json:"id" binding:"required"`
	DeskripsiPerkembangan string   `json:"deskripsi_perkembangan" binding:"required"`
	KendalaDitemui        string   `json:"kendala_ditemui"`
	TindakLanjut          string   `json:"tindak_lanjut"`
	FilesToDelete         []string `json:"files_to_delete"` // Array of file paths to delete
	GuruPengisiID         *uint    `json:"guru_pengisi_id"`
}

type MonitoringInklusiBulananSummaryRequest struct {
	Tahun    int16  `json:"tahun" binding:"required"`
	Bulan    *int16 `json:"bulan"` // Optional: filter by specific month (1-12)
	RombelID *uint  `json:"rombel_id"` // Optional: filter by rombel
}

type MonitoringInklusiBulananSummaryResponse struct {
	Tahun                    int16                                  `json:"tahun"`
	Bulan                    *int16                                 `json:"bulan,omitempty"`
	RombelID                 *uint                                  `json:"rombel_id,omitempty"`
	RombelName               string                                 `json:"rombel_name,omitempty"`
	TotalAnakInklusi         int                                    `json:"total_anak_inklusi"`
	TotalPPI                 int                                    `json:"total_ppi"`
	TotalMonitoring          int                                    `json:"total_monitoring"`
	RataRataMonitoringPerPPI float64                                `json:"rata_rata_monitoring_per_ppi"`
	MonitoringPerBulan       []MonitoringBulananPerBulan            `json:"monitoring_per_bulan"`
	MonitoringPerAspek       []MonitoringBulananPerAspek            `json:"monitoring_per_aspek"`
	MonitoringPerSiswa       []MonitoringBulananPerSiswa            `json:"monitoring_per_siswa"`
	PPIAktifPerStatus        PPIStatusBreakdown                     `json:"ppi_aktif_per_status"`
}

type MonitoringBulananPerBulan struct {
	Bulan           int16  `json:"bulan"`
	NamaBulan       string `json:"nama_bulan"`
	TotalMonitoring int    `json:"total_monitoring"`
	TotalPPIAktif   int    `json:"total_ppi_aktif"`
}

type MonitoringBulananPerAspek struct {
	AspekPembelajaran string `json:"aspek_pembelajaran"`
	TotalPPI          int    `json:"total_ppi"`
	TotalMonitoring   int    `json:"total_monitoring"`
	RataRataPerPPI    float64 `json:"rata_rata_per_ppi"`
}

type MonitoringBulananPerSiswa struct {
	AnakInklusiID     uint                 `json:"anak_inklusi_id"`
	PesertaDidik      *PesertaDidikSimple  `json:"peserta_didik"`
	JenisHambatan     string               `json:"jenis_hambatan"`
	TotalPPI          int                  `json:"total_ppi"`
	TotalMonitoring   int                  `json:"total_monitoring"`
	MonitoringDetails []MonitoringDetail   `json:"monitoring_details"`
}

type MonitoringDetail struct {
	PPIInklusiID      uint   `json:"ppi_inklusi_id"`
	AspekPembelajaran string `json:"aspek_pembelajaran"`
	NamaProgram       string `json:"nama_program"`
	TotalMonitoring   int    `json:"total_monitoring"`
	BulanTerakhir     int16  `json:"bulan_terakhir"`
}

type PPIStatusBreakdown struct {
	BelumDimulai      int `json:"belum_dimulai"`
	SedangBerlangsung int `json:"sedang_berlangsung"`
	Selesai           int `json:"selesai"`
	Dihentikan        int `json:"dihentikan"`
}
