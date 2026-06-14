package dtos

// LayananSPMBCreateRequest represents the request payload for creating Layanan SPMB (public)
type LayananSPMBCreateRequest struct {
	NamaOrangTua      string `json:"nama_orang_tua" binding:"required"`
	NomorTelepon      string `json:"nomor_telepon" binding:"required"`
	Alamat            string `json:"alamat" binding:"required"`
	NamaLengkapMurid  string `json:"nama_lengkap_murid" binding:"required"`
	Keperluan         string `json:"keperluan" binding:"required"`
}

// LayananSPMBResponse represents the response payload for Layanan SPMB
type LayananSPMBResponse struct {
	ID               uint   `json:"id"`
	NamaOrangTua     string `json:"nama_orang_tua"`
	NomorTelepon     string `json:"nomor_telepon"`
	Alamat           string `json:"alamat"`
	NamaLengkapMurid string `json:"nama_lengkap_murid"`
	Keperluan        string `json:"keperluan"`
	TanggalLaporan   string `json:"tanggal_laporan"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// LayananSPMBGetAllRequest represents the request for getting all layanan SPMB with filters
type LayananSPMBGetAllRequest struct {
	Search struct {
		StartDate    string `json:"start_date"` // YYYY-MM-DD
		EndDate      string `json:"end_date"`   // YYYY-MM-DD
		NamaOrangTua string `json:"nama_orang_tua"`
		NamaMurid    string `json:"nama_murid"`
		Status       string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// LayananSPMBListWithPaginationResponse represents paginated list response
type LayananSPMBListWithPaginationResponse struct {
	Data       []LayananSPMBResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

// LayananSPMBUpdateStatusRequest represents the request for updating status
type LayananSPMBUpdateStatusRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// SettingLayananSPMBRequest represents the request for setting layanan SPMB
type SettingLayananSPMBRequest struct {
	NamaKepalaSekolah string `json:"nama_kepala_sekolah"`
	NIPKepalaSekolah  string `json:"nip_kepala_sekolah"`
	NamaKetuaPanitia  string `json:"nama_ketua_panitia"`
	NIPKetuaPanitia   string `json:"nip_ketua_panitia"`
	GrupWA            string `json:"grup_wa"`
}

// SettingLayananSPMBResponse represents the response for setting layanan SPMB
type SettingLayananSPMBResponse struct {
	ID                uint    `json:"id"`
	NamaKepalaSekolah *string `json:"nama_kepala_sekolah"`
	NIPKepalaSekolah  *string `json:"nip_kepala_sekolah"`
	NamaKetuaPanitia  *string `json:"nama_ketua_panitia"`
	NIPKetuaPanitia   *string `json:"nip_ketua_panitia"`
	GrupWA            *string `json:"grup_wa"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// GrupWASPMBResponse represents the response for public grup WA
type GrupWASPMBResponse struct {
	GrupWA *string `json:"grup_wa"`
}
