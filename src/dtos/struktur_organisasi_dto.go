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
