package dtos

import (
	"time"
)

// KepegawaianCreateRequest represents the request payload for creating Kepegawaian
type KepegawaianCreateRequest struct {
	Nama              string `json:"nama" binding:"required"`
	Username          string `json:"username" binding:"required"`
	Password          string `json:"password" binding:"required,min=6"`
	NIP               string `json:"nip" binding:"omitempty"`
	NKKI              string `json:"nkki" binding:"omitempty"`
	Kategori          string `json:"kategori" binding:"omitempty"`
	Jabatan           string `json:"jabatan" binding:"omitempty"`
	BidangStudiID     *uint  `json:"bidang_studi_id" binding:"omitempty"`
	RombelGuruKelasID *uint  `json:"rombel_guru_kelas_id" binding:"omitempty"`
	RombelBidangStudi []uint `json:"rombel_bidang_studi" binding:"omitempty"`
	RoleIDs           []uint `json:"role_ids" binding:"omitempty"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// KepegawaianUpdateRequest represents the request payload for updating Kepegawaian
type KepegawaianUpdateRequest struct {
	ID                        uint     `json:"id" binding:"required"`
	Nama                      string   `json:"nama" binding:"omitempty"`
	Username                  string   `json:"username" binding:"omitempty"`
	Password                  string   `json:"password" binding:"omitempty,min=6"`
	NIP                       *string  `json:"nip" binding:"omitempty"`
	NKKI                      *string  `json:"nkki" binding:"omitempty"`
	Kategori                  string   `json:"kategori" binding:"omitempty"`
	Jabatan                   string   `json:"jabatan" binding:"omitempty"`
	BidangStudiID             *uint    `json:"bidang_studi_id" binding:"omitempty"`
	RombelGuruKelasID         *uint    `json:"rombel_guru_kelas_id" binding:"omitempty"`
	RombelBidangStudi         []uint   `json:"rombel_bidang_studi" binding:"omitempty"`
	RoleIDs                   []uint   `json:"role_ids" binding:"omitempty"`
	Status                    string   `json:"status" binding:"omitempty,oneof=active inactive"`
	FilesToDelete             []string `json:"files_to_delete" binding:"omitempty"` // e.g., ["kk", "ktp", "ijazah_s1"]
	SertifikatLainnyaToDelete []string `json:"sertifikat_lainnya_to_delete" binding:"omitempty"`
	DokumenLainnyaToDelete    []string `json:"dokumen_lainnya_to_delete" binding:"omitempty"`
}

// BidangStudiSimpleResponse represents simple bidang studi response
type BidangStudiSimpleResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RombelSimpleResponse represents simple rombel response
type RombelSimpleResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// KepegawaianResponse represents the response payload for Kepegawaian
type KepegawaianResponse struct {
	ID                  uint                       `json:"id"`
	Nama                string                     `json:"nama"`
	Username            string                     `json:"username"`
	NIP                 string                     `json:"nip"`
	NKKI                string                     `json:"nkki"`
	Foto                *string                    `json:"foto"`
	Kategori            string                     `json:"kategori"`
	Jabatan             string                     `json:"jabatan"`
	BidangStudiID       *uint                      `json:"bidang_studi_id"`
	BidangStudi         *BidangStudiSimpleResponse `json:"bidang_studi"`
	RombelGuruKelasID   *uint                      `json:"rombel_guru_kelas_id"`
	RombelGuruKelas     *RombelSimpleResponse      `json:"rombel_guru_kelas"`
	RombelBidangStudi   []uint                     `json:"rombel_bidang_studi"`
	KK                    *string        `json:"kk"`
	AktaLahir             *string        `json:"akta_lahir"`
	KTP                   *string        `json:"ktp"`
	IjazahSD              *string        `json:"ijazah_sd"`
	IjazahSMP             *string        `json:"ijazah_smp"`
	IjazahSMA             *string        `json:"ijazah_sma"`
	IjazahS1              *string        `json:"ijazah_s1"`
	IjazahS2              *string        `json:"ijazah_s2"`
	IjazahS3              *string        `json:"ijazah_s3"`
	SertifikatPendidik    *string        `json:"sertifikat_pendidik"`
	SertifikatLainnya     []string       `json:"sertifikat_lainnya"`
	SK                    *string        `json:"sk"`
	DokumenLainnya        []string       `json:"dokumen_lainnya"`
	Status                string         `json:"status"`
	Roles                 []RoleResponse `json:"roles"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedByID           *uint          `json:"created_by_id"`
	UpdatedByID           *uint          `json:"updated_by_id"`
}

// KepegawaianListResponse represents the response payload for listing Kepegawaian
type KepegawaianListResponse struct {
	Data   []KepegawaianResponse `json:"data"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Total  int64                 `json:"total"`
}

// KepegawaianGetAllRequest represents the request payload for getting all Kepegawaian with filters
type KepegawaianGetAllRequest struct {
	Search struct {
		Nama     string `json:"nama" binding:"omitempty"`
		Username string `json:"username" binding:"omitempty"`
		NIP      string `json:"nip" binding:"omitempty"`
		NKKI     string `json:"nkki" binding:"omitempty"`
		Kategori string `json:"kategori" binding:"omitempty"`
		Jabatan  string `json:"jabatan" binding:"omitempty"`
		RoleID   uint   `json:"role_id" binding:"omitempty"`
		Status   string `json:"status" binding:"omitempty"`
	} `json:"search" binding:"omitempty"`
	Pagination struct {
		Limit int `json:"limit" binding:"omitempty"`
		Page  int `json:"page" binding:"omitempty"`
	} `json:"pagination" binding:"omitempty"`
}

// KepegawaianListWithPaginationResponse represents the response with pagination
type KepegawaianListWithPaginationResponse struct {
	Data       []KepegawaianResponse `json:"data"`
	Pagination PaginationInfo        `json:"pagination"`
}

// TotalPendidikResponse represents the response for total pendidik count
type TotalPendidikResponse struct {
	TotalPendidik int64 `json:"total_pendidik"`
}

// TotalTendikResponse represents the response for total tenaga kependidikan count
type TotalTendikResponse struct {
	TotalTendik int64 `json:"total_tendik"`
}
