package dtos

import "time"

// StrukturOrganisasiCreateRequest represents the request payload for creating StrukturOrganisasi
type StrukturOrganisasiCreateRequest struct {
	PegawaiID         *uint  `json:"pegawai_id" binding:"omitempty"`
	NamaNonPegawai    string `json:"nama_non_pegawai" binding:"omitempty,max=255"`
	JabatanNonPegawai string `json:"jabatan_non_pegawai" binding:"omitempty,max=255"`
	Urutan            int    `json:"urutan" binding:"required"`
	Relasi            string `json:"relasi" binding:"required,min=1,max=100"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// StrukturOrganisasiUpdateRequest represents the request payload for updating StrukturOrganisasi
type StrukturOrganisasiUpdateRequest struct {
	ID                   uint    `json:"id" binding:"required"`
	PegawaiID            *uint   `json:"pegawai_id" binding:"omitempty"`
	PegawaiIDSet         bool    `json:"-"` // Internal flag to track if pegawai_id is explicitly set
	NamaNonPegawai       *string `json:"nama_non_pegawai" binding:"omitempty,max=255"`
	JabatanNonPegawai    *string `json:"jabatan_non_pegawai" binding:"omitempty,max=255"`
	Urutan               *int    `json:"urutan" binding:"omitempty"`
	Relasi               *string `json:"relasi" binding:"omitempty,min=1,max=100"`
	Status               *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// PegawaiSimpleResponse represents simple pegawai response
type PegawaiSimpleResponse struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Jabatan  string `json:"jabatan"`
	Status   string `json:"status"`
}

// PegawaiPublicDetailResponse represents detailed pegawai response for public API
type PegawaiPublicDetailResponse struct {
	NamaLengkap   string                        `json:"nama_lengkap"`
	NIP           string                        `json:"nip"`
	NKKI          string                        `json:"nkki"`
	Jabatan       string                        `json:"jabatan"`
	Kategori      string                        `json:"kategori"`
	BidangStudi   string                        `json:"bidang_studi,omitempty"`
	KelasMengajar []RombelWithKelasResponse     `json:"kelas_mengajar,omitempty"`
}

// RombelWithKelasResponse represents rombel with kelas name for public API
type RombelWithKelasResponse struct {
	ID        uint   `json:"id"`
	Rombel    string `json:"rombel"`
	NamaKelas string `json:"nama_kelas"`
	Status    string `json:"status"`
}

// GuruKelasGroupResponse represents grouped guru kelas by kelas
type GuruKelasGroupResponse struct {
	NamaKelas string                               `json:"nama_kelas"`
	Relasi    string                               `json:"relasi"`
	Guru      []PegawaiPublicDetailResponse        `json:"guru"`
}

// GuruMapelGroupResponse represents grouped guru mapel by bidang studi
type GuruMapelGroupResponse struct {
	BidangStudi string                               `json:"bidang_studi"`
	Relasi      string                               `json:"relasi"`
	Guru        []PegawaiPublicDetailResponse        `json:"guru"`
}

// JabatanGroupResponse represents grouped by jabatan for urutan 5+
type JabatanGroupResponse struct {
	Jabatan string                               `json:"jabatan"`
	Data    []StrukturOrganisasiPublicResponse  `json:"data"`
}

// StrukturOrganisasiResponse represents the response payload for StrukturOrganisasi
type StrukturOrganisasiResponse struct {
	ID                uint                     `json:"id"`
	PegawaiID         *uint                    `json:"pegawai_id"`
	Pegawai           *PegawaiSimpleResponse   `json:"pegawai,omitempty"`
	NamaNonPegawai    string                   `json:"nama_non_pegawai"`
	JabatanNonPegawai string                   `json:"jabatan_non_pegawai"`
	Urutan            int                      `json:"urutan"`
	Relasi            string                   `json:"relasi"`
	Status            string                   `json:"status"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	CreatedByID       *uint                    `json:"created_by_id"`
	UpdatedByID       *uint                    `json:"updated_by_id"`
}

// StrukturOrganisasiListResponse represents the response payload for listing StrukturOrganisasi
type StrukturOrganisasiListResponse struct {
	Data   []StrukturOrganisasiResponse `json:"data"`
	Limit  int                          `json:"limit"`
	Offset int                          `json:"offset"`
	Total  int64                        `json:"total"`
}

// StrukturOrganisasiGetAllRequest represents the request payload for getting all StrukturOrganisasi with filters
type StrukturOrganisasiGetAllRequest struct {
	Search struct {
		Nama    string `json:"nama" binding:"omitempty"`
		Urutan  int    `json:"urutan" binding:"omitempty"`
		Relasi  string `json:"relasi" binding:"omitempty"`
		Jabatan string `json:"jabatan" binding:"omitempty"`
		Status  string `json:"status" binding:"omitempty"`
	} `json:"search" binding:"omitempty"`
	Pagination struct {
		Limit int `json:"limit" binding:"omitempty"`
		Page  int `json:"page" binding:"omitempty"`
	} `json:"pagination" binding:"omitempty"`
}

// StrukturOrganisasiListWithPaginationResponse represents the response with pagination
type StrukturOrganisasiListWithPaginationResponse struct {
	Data       []StrukturOrganisasiResponse `json:"data"`
	Pagination PaginationInfo               `json:"pagination"`
}

// StrukturOrganisasiPublicResponse represents the public response payload for StrukturOrganisasi
type StrukturOrganisasiPublicResponse struct {
	Pegawai           *PegawaiPublicDetailResponse `json:"pegawai,omitempty"`
	NamaNonPegawai    string                       `json:"nama_non_pegawai,omitempty"`
	JabatanNonPegawai string                       `json:"jabatan_non_pegawai,omitempty"`
	Urutan            int                          `json:"urutan"`
	Relasi            string                       `json:"relasi"`
}

// StrukturOrganisasiGroupedResponse represents grouped response by urutan
type StrukturOrganisasiGroupedResponse struct {
	Urutan         int                                  `json:"urutan"`
	Data           []StrukturOrganisasiPublicResponse  `json:"data,omitempty"`
	GuruKelas      []GuruKelasGroupResponse             `json:"guru_kelas,omitempty"`
	GuruMapel      []GuruMapelGroupResponse             `json:"guru_mapel,omitempty"`
	ByJabatan      []JabatanGroupResponse               `json:"by_jabatan,omitempty"`
}
