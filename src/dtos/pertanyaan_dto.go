package dtos

import "pintu-backend/src/modules/models"

// PertanyaanCreateRequest represents the request payload for creating Pertanyaan (public)
type PertanyaanCreateRequest struct {
	Nama      string `form:"nama" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
	Telepon   string `form:"telepon"`
	Kategori  string `form:"kategori" binding:"required"`
	Prioritas string `form:"prioritas"`
	Judul     string `form:"judul" binding:"required"`
	Deskripsi string `form:"deskripsi" binding:"required"`
}

// PertanyaanResponse represents the response payload for Pertanyaan
type PertanyaanResponse struct {
	ID               uint               `json:"id"`
	IDTiket          string             `json:"id_tiket"`
	TanggalPengajuan string             `json:"tanggal_pengajuan"`
	Nama             string             `json:"nama"`
	Email            string             `json:"email"`
	Telepon          string             `json:"telepon"`
	Kategori         string             `json:"kategori"`
	Prioritas        string             `json:"prioritas"`
	Judul            string             `json:"judul"`
	Deskripsi        string             `json:"deskripsi"`
	FilePertanyaan   []models.FileItem  `json:"file_pertanyaan"`
	JudulJawaban     *string            `json:"judul_jawaban"`
	DeskripsiJawaban *string            `json:"deskripsi_jawaban"`
	FileJawaban      []models.FileItem  `json:"file_jawaban"`
	TanggalProses    *string            `json:"tanggal_proses"`
	EmailTerkirim    bool               `json:"email_terkirim"`
	TanggalSelesai   *string            `json:"tanggal_selesai"`
	Status           string             `json:"status"`
	RepliedBy        *uint              `json:"replied_by"`
	CreatedAt        string             `json:"created_at"`
}

// PertanyaanTrackRequest represents the request payload for tracking Pertanyaan by ID Tiket
type PertanyaanTrackRequest struct {
	IDTiket string `json:"id_tiket" binding:"required"`
}

// PertanyaanTrackResponse represents the simplified response for tracking
type PertanyaanTrackResponse struct {
	IDTiket          string `json:"id_tiket"`
	TanggalPengajuan string `json:"tanggal_pengajuan"`
	Nama             string `json:"nama"`
	Email            string `json:"email"`
	Telepon          string `json:"telepon"`
	Kategori         string `json:"kategori"`
	Prioritas        string `json:"prioritas"`
	Judul            string `json:"judul"`
	Deskripsi        string `json:"deskripsi"`
	Status           string `json:"status"`
}
// PertanyaanGetAllRequest represents the request for getting all pertanyaan with filters
type PertanyaanGetAllRequest struct {
	Search struct {
		IDTiket   string `json:"id_tiket"`
		StartDate string `json:"start_date"` // YYYY-MM-DD
		EndDate   string `json:"end_date"`   // YYYY-MM-DD
		Nama      string `json:"nama"`
		Email     string `json:"email"`
		Kategori  string `json:"kategori"`
		Prioritas string `json:"prioritas"`
		Judul     string `json:"judul"`
		Status    string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// PertanyaanListWithPaginationResponse represents paginated list response
type PertanyaanListWithPaginationResponse struct {
	Data       []PertanyaanResponse `json:"data"`
	Pagination PaginationMeta       `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Page       int   `json:"page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PertanyaanSendReplyRequest represents the request for sending email reply
type PertanyaanSendReplyRequest struct {
	ID               uint   `json:"id" binding:"required"`
	JudulJawaban     string `json:"judul_jawaban" binding:"required"`
	DeskripsiJawaban string `json:"deskripsi_jawaban" binding:"required"`
}
