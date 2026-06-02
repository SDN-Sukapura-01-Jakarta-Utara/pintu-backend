package dtos

import "encoding/json"

// KelulusanCreateRequest represents the request for creating kelulusan data
type KelulusanCreateRequest struct {
	NomorPeserta string                 `json:"nomor_peserta" binding:"required"`
	NISN         string                 `json:"nisn" binding:"required"`
	Nama         string                 `json:"nama" binding:"required"`
	TanggalLahir string                 `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
	Nilai        map[string]interface{} `json:"nilai" binding:"required"`         // Dynamic: {"Matematika": 85, "Bahasa Indonesia": 90}
	Lulus        bool                   `json:"lulus" binding:"required"`
	MaxAttempts  int                    `json:"max_attempts"`                     // Optional: jumlah percobaan yang diperlukan (default: 0)
}

// KelulusanResponse represents the response for kelulusan data
type KelulusanResponse struct {
	ID            uint            `json:"id"`
	NomorPeserta  string          `json:"nomor_peserta"`
	NISN          string          `json:"nisn"`
	Nama          string          `json:"nama"`
	TanggalLahir  string          `json:"tanggal_lahir"`
	Nilai         json.RawMessage `json:"nilai"`
	RataRataNilai float64         `json:"rata_rata_nilai"` // Calculated average, 2 decimal places
	Lulus         bool            `json:"lulus"`
	SKL           string          `json:"skl,omitempty"`
	MaxAttempts   int             `json:"max_attempts"`
	AttemptCount  int             `json:"attempt_count"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CreatedByID   *uint           `json:"created_by_id,omitempty"`
	UpdatedByID   *uint           `json:"updated_by_id,omitempty"`
}

// KelulusanDownloadTemplateRequest represents the request for downloading template
type KelulusanDownloadTemplateRequest struct {
	MapelList []string `json:"mapel_list" binding:"required"` // List of mata pelajaran names
}

// ImportKelulusanResponse represents the response for import excel
type ImportKelulusanResponse struct {
	SuccessCount int                          `json:"success_count"`
	FailedCount  int                          `json:"failed_count"`
	Errors       []ImportKelulusanRowError    `json:"errors,omitempty"`
}

// ImportKelulusanRowError represents an error for a specific row during import
type ImportKelulusanRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// KelulusanGetAllRequest represents the request for getting all kelulusan with filters
type KelulusanGetAllRequest struct {
	Search struct {
		Nama         string `json:"nama"`
		NomorPeserta string `json:"nomor_peserta"`
		NISN         string `json:"nisn"`
		Lulus        *bool  `json:"lulus"` // nil = all, true = lulus, false = tidak lulus
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// KelulusanListWithPaginationResponse represents the response with pagination
type KelulusanListWithPaginationResponse struct {
	Data       []KelulusanResponse `json:"data"`
	Pagination PaginationInfo      `json:"pagination"`
}


// KelulusanUpdateRequest represents the request for updating kelulusan data
type KelulusanUpdateRequest struct {
	ID           uint                   `json:"id" binding:"required"`
	NomorPeserta string                 `json:"nomor_peserta" binding:"omitempty"`
	NISN         string                 `json:"nisn" binding:"omitempty"`
	Nama         string                 `json:"nama" binding:"omitempty"`
	TanggalLahir string                 `json:"tanggal_lahir" binding:"omitempty"` // Format: YYYY-MM-DD
	Nilai        map[string]interface{} `json:"nilai" binding:"omitempty"`
	Lulus        *bool                  `json:"lulus" binding:"omitempty"`
	MaxAttempts  *int                   `json:"max_attempts" binding:"omitempty"` // Optional: update max_attempts
	DeleteSKL    bool                   `json:"delete_skl" binding:"omitempty"`   // true = hapus file SKL
}

// CekNilaiKelulusanRequest represents the request for checking kelulusan by NISN and tanggal lahir
type CekNilaiKelulusanRequest struct {
	NISN         string `json:"nisn" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
}

// CekKelulusanRequest represents the request for checking full kelulusan data by NISN and tanggal lahir
type CekKelulusanRequest struct {
	NISN         string `json:"nisn" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
}

// DownloadLaporanNilaiKelulusanRequest represents the request for downloading laporan nilai kelulusan PDF
type DownloadLaporanNilaiKelulusanRequest struct {
	NISN         string `json:"nisn" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
}

// CekNilaiKelulusanResponse represents the response for public kelulusan check (without informasi_lulus)
type CekNilaiKelulusanResponse struct {
	ID            uint            `json:"id"`
	NomorPeserta  string          `json:"nomor_peserta"`
	NISN          string          `json:"nisn"`
	Nama          string          `json:"nama"`
	TanggalLahir  string          `json:"tanggal_lahir"`
	Nilai         json.RawMessage `json:"nilai"`
	RataRataNilai float64         `json:"rata_rata_nilai"`
	SKL           string          `json:"skl,omitempty"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}
