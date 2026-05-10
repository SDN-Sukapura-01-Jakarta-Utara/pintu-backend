package dtos

import "pintu-backend/src/modules/models"

// PengaduanCreateRequest represents the request payload for creating Pengaduan (public)
type PengaduanCreateRequest struct {
	TipePelapor string `form:"tipe_pelapor"`
	Nama        string `form:"nama"`
	Email       string `form:"email"`
	Telepon     string `form:"telepon"`
	Kategori    string `form:"kategori" binding:"required"`
	Prioritas   string `form:"prioritas"`
	Judul       string `form:"judul" binding:"required"`
	Deskripsi   string `form:"deskripsi" binding:"required"`
}

// PengaduanResponse represents the response payload for Pengaduan
type PengaduanResponse struct {
	ID               uint               `json:"id"`
	IDTiket          string             `json:"id_tiket"`
	TanggalPengajuan string             `json:"tanggal_pengajuan"`
	TipePelapor      string             `json:"tipe_pelapor"`
	Nama             *string            `json:"nama"`
	Email            *string            `json:"email"`
	Telepon          *string            `json:"telepon"`
	Kategori         string             `json:"kategori"`
	Prioritas        string             `json:"prioritas"`
	Judul            string             `json:"judul"`
	Deskripsi        string             `json:"deskripsi"`
	FilePengaduan    []models.FileItem  `json:"file_pengaduan"`
	JudulJawaban     *string            `json:"judul_jawaban"`
	DeskripsiJawaban *string            `json:"deskripsi_jawaban"`
	FileJawaban      []models.FileItem  `json:"file_jawaban"`
	TindakLanjut     *string            `json:"tindak_lanjut"`
	FileTindakLanjut []models.FileItem  `json:"file_tindak_lanjut"`
	TanggalProses    *string            `json:"tanggal_proses"`
	EmailTerkirim    bool               `json:"email_terkirim"`
	TanggalSelesai   *string            `json:"tanggal_selesai"`
	Status           string             `json:"status"`
	RepliedBy        *uint              `json:"replied_by"`
	CreatedAt        string             `json:"created_at"`
}

// PengaduanTrackRequest represents the request payload for tracking Pengaduan by ID Tiket
type PengaduanTrackRequest struct {
	IDTiket string `json:"id_tiket" binding:"required"`
}

// PengaduanTrackResponse represents the simplified response for tracking
type PengaduanTrackResponse struct {
	IDTiket          string  `json:"id_tiket"`
	TanggalPengajuan string  `json:"tanggal_pengajuan"`
	TipePelapor      string  `json:"tipe_pelapor"`
	Nama             *string `json:"nama"`
	Email            *string `json:"email"`
	Telepon          *string `json:"telepon"`
	Kategori         string  `json:"kategori"`
	Prioritas        string  `json:"prioritas"`
	Judul            string  `json:"judul"`
	Deskripsi        string  `json:"deskripsi"`
	Status           string  `json:"status"`
}

// PengaduanGetAllRequest represents the request for getting all pengaduan with filters
type PengaduanGetAllRequest struct {
	Search struct {
		IDTiket     string `json:"id_tiket"`
		StartDate   string `json:"start_date"` // YYYY-MM-DD
		EndDate     string `json:"end_date"`   // YYYY-MM-DD
		TipePelapor string `json:"tipe_pelapor"`
		Nama        string `json:"nama"`
		Email       string `json:"email"`
		Kategori    string `json:"kategori"`
		Prioritas   string `json:"prioritas"`
		Judul       string `json:"judul"`
		Status      string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// PengaduanListWithPaginationResponse represents paginated list response
type PengaduanListWithPaginationResponse struct {
	Data       []PengaduanResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

// PengaduanSendReplyRequest represents the request for sending email reply
type PengaduanSendReplyRequest struct {
	ID               uint   `json:"id" binding:"required"`
	JudulJawaban     string `json:"judul_jawaban" binding:"required"`
	DeskripsiJawaban string `json:"deskripsi_jawaban" binding:"required"`
}

// PengaduanSaveTindakLanjutRequest represents the request for saving tindak lanjut
type PengaduanSaveTindakLanjutRequest struct {
	ID            uint     `form:"id" binding:"required"`
	TindakLanjut  string   `form:"tindak_lanjut" binding:"required"`
	FilesToDelete []string `form:"files_to_delete"` // Array of file IDs to delete
}