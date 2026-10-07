package dtos

type PPIInklusiItem struct {
	AspekPembelajaran    string `json:"aspek_pembelajaran" binding:"required"`
	NamaProgram          string `json:"nama_program" binding:"required"`
	TujuanPembelajaran   string `json:"tujuan_pembelajaran" binding:"required"`
	StrategiPembelajaran string `json:"strategi_pembelajaran" binding:"required"`
	TargetWaktu          string `json:"target_waktu"`
}

type PPIInklusiCreateRequest struct {
	AnakInklusiRombelID uint             `json:"anak_inklusi_rombel_id" binding:"required"`
	TahunPelajaranID    uint             `json:"tahun_pelajaran_id" binding:"required"`
	GuruPembuatID       *uint            `json:"guru_pembuat_id"`
	ProgramList         []PPIInklusiItem `json:"program_list" binding:"required,min=1,dive"`
}

type PPIInklusiCreateResponse struct {
	TotalCreated int                  `json:"total_created"`
	Data         []PPIInklusiResponse `json:"data"`
}

type PPIInklusiResponse struct {
	ID                    uint                          `json:"id"`
	AnakInklusiRombelID   uint                          `json:"anak_inklusi_rombel_id"`
	TahunPelajaranID      uint                          `json:"tahun_pelajaran_id"`
	AspekPembelajaran     string                        `json:"aspek_pembelajaran"`
	NamaProgram           string                        `json:"nama_program"`
	TujuanPembelajaran    string                        `json:"tujuan_pembelajaran"`
	StrategiPembelajaran  string                        `json:"strategi_pembelajaran"`
	TargetWaktu           string                        `json:"target_waktu"`
	CapaianSaatIni        string                        `json:"capaian_saat_ini"`
	StatusCapaian         string                        `json:"status_capaian"`
	GuruPembuatID         *uint                         `json:"guru_pembuat_id"`
	AnakInklusiRombel     *AnakInklusiRombelResponse    `json:"anak_inklusi_rombel,omitempty"`
	TahunPelajaran        *TahunPelajaranSimple         `json:"tahun_pelajaran,omitempty"`
	GuruPembuat           *KepegawaianSimple            `json:"guru_pembuat,omitempty"`
	CreatedAt             string                        `json:"created_at"`
	UpdatedAt             string                        `json:"updated_at"`
}

type PPIInklusiGetAllRequest struct {
	Search struct {
		TahunPelajaranID    uint   `json:"tahun_pelajaran_id"`
		RombelID            uint   `json:"rombel_id"`
		GuruPembuatID       uint   `json:"guru_pembuat_id"`
		AnakInklusiRombelID uint   `json:"anak_inklusi_rombel_id"`
		AspekPembelajaran   string `json:"aspek_pembelajaran"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

type PPIInklusiListWithPaginationResponse struct {
	Data       []PPIInklusiResponse `json:"data"`
	Pagination PaginationResponse   `json:"pagination"`
}

type PPIInklusiUpdateRequest struct {
	ID                   uint   `json:"id" binding:"required"`
	AspekPembelajaran    string `json:"aspek_pembelajaran"`
	NamaProgram          string `json:"nama_program"`
	TujuanPembelajaran   string `json:"tujuan_pembelajaran"`
	StrategiPembelajaran string `json:"strategi_pembelajaran"`
	TargetWaktu          string `json:"target_waktu"`
	GuruPembuatID        *uint  `json:"guru_pembuat_id"`
	StatusCapaian        string `json:"status_capaian"`
}

type PPIInklusiSummaryRequest struct {
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
	Search struct {
		RombelID            uint   `json:"rombel_id"`
		GuruPembuatID       uint   `json:"guru_pembuat_id"`
		AnakInklusiRombelID uint   `json:"anak_inklusi_rombel_id"`
		AspekPembelajaran   string `json:"aspek_pembelajaran"`
	} `json:"search"`
}

type PPIInklusiSummaryResponse struct {
	TahunPelajaran          TahunPelajaranSimple           `json:"tahun_pelajaran"`
	TotalPPI                int                            `json:"total_ppi"`
	TotalAnakInklusi        int                            `json:"total_anak_inklusi"`
	StatusCapaianBreakdown  []StatusCapaianBreakdown       `json:"status_capaian_breakdown"`
	AspekPembelajaranBreakdown []AspekPembelajaranBreakdown `json:"aspek_pembelajaran_breakdown"`
}

type StatusCapaianBreakdown struct {
	StatusCapaian string `json:"status_capaian"`
	Total         int    `json:"total"`
	Percentage    float64 `json:"percentage"`
}

type AspekPembelajaranBreakdown struct {
	AspekPembelajaran string  `json:"aspek_pembelajaran"`
	Total             int     `json:"total"`
	Percentage        float64 `json:"percentage"`
}
