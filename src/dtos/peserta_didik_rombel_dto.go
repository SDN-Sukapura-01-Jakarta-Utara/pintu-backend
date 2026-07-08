package dtos

// PesertaDidikRombelCreateRequest represents the request payload for bulk creating PesertaDidikRombel
type PesertaDidikRombelCreateRequest struct {
	PesertaDidikIDs  []uint `json:"peserta_didik_ids" binding:"required,min=1"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Status           string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// PesertaDidikRombelUpdateRequest represents the request payload for updating PesertaDidikRombel
type PesertaDidikRombelUpdateRequest struct {
	ID               uint   `json:"id" binding:"required"`
	RombelID         uint   `json:"rombel_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
	Status           string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// PesertaDidikRombelGetAllRequest represents the request payload for getting all PesertaDidikRombel with filters
type PesertaDidikRombelGetAllRequest struct {
	Search struct {
		Nama             string `json:"nama" binding:"omitempty"`
		RombelID         uint   `json:"rombel_id" binding:"omitempty"`
		TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"omitempty"`
		Status           string `json:"status" binding:"omitempty"`
	} `json:"search" binding:"omitempty"`
	Pagination struct {
		Limit int `json:"limit" binding:"omitempty"`
		Page  int `json:"page" binding:"omitempty"`
	} `json:"pagination" binding:"omitempty"`
}

// PesertaDidikRombelResponse represents the response payload for PesertaDidikRombel
type PesertaDidikRombelResponse struct {
	ID               uint                          `json:"id"`
	PesertaDidikID   uint                          `json:"peserta_didik_id"`
	PesertaDidik     *PesertaDidikResponse         `json:"peserta_didik,omitempty"`
	RombelID         uint                          `json:"rombel_id"`
	Rombel           *RombelDetailResponse         `json:"rombel,omitempty"`
	TahunPelajaranID uint                          `json:"tahun_pelajaran_id"`
	TahunPelajaran   *TahunPelajaranDetailResponse `json:"tahun_pelajaran,omitempty"`
	Status           string                        `json:"status"`
	CreatedAt        string                        `json:"created_at"`
	UpdatedAt        string                        `json:"updated_at"`
	CreatedByID      *uint                         `json:"created_by_id"`
	UpdatedByID      *uint                         `json:"updated_by_id"`
}

// PesertaDidikRombelBulkCreateResponse represents the response for bulk create
type PesertaDidikRombelBulkCreateResponse struct {
	SuccessCount int                          `json:"success_count"`
	FailedCount  int                          `json:"failed_count"`
	Data         []PesertaDidikRombelResponse `json:"data"`
	Errors       []BulkCreateError            `json:"errors,omitempty"`
}

// BulkCreateError represents an error for a specific peserta didik during bulk create
type BulkCreateError struct {
	PesertaDidikID uint   `json:"peserta_didik_id"`
	Message        string `json:"message"`
}

// PesertaDidikRombelResetRequest represents the request payload for reset pemetaan rombel
type PesertaDidikRombelResetRequest struct {
	RombelID         uint `json:"rombel_id" binding:"omitempty"`
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"omitempty"`
}

// PesertaDidikRombelResetResponse represents the response for reset operation
type PesertaDidikRombelResetResponse struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// PesertaDidikRombelListWithPaginationResponse represents list response with pagination info
type PesertaDidikRombelListWithPaginationResponse struct {
	Data       []PesertaDidikRombelResponse `json:"data"`
	Pagination PaginationInfo               `json:"pagination"`
}
