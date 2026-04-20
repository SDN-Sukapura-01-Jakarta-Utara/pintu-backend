package dtos

import (
	"time"
)

// FotoItemDTO represents a single photo in the foto array
type FotoItemDTO struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	Thumbnail string `json:"thumbnail"` // "active" or "inactive"
}

// AnggotaTimPrestasiDTO represents anggota tim prestasi details
type AnggotaTimPrestasiDTO struct {
	ID               uint                          `json:"id"`
	PrestasiID       uint                          `json:"prestasi_id"`
	PesertaDidikID   uint                          `json:"peserta_didik_id"`
	TahunPelajaranID uint                          `json:"tahun_pelajaran_id"`
	PesertaDidik     *PesertaDidikResponse         `json:"peserta_didik"`
	TahunPelajaran   *TahunPelajaranDetailResponse `json:"tahun_pelajaran"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
}

// EkstrakurikulerDetailDTO represents ekstrakurikuler details
type EkstrakurikulerDetailDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Kategori string `json:"kategori"`
	Status   string `json:"status"`
}

// PrestasiCreateRequest represents the request payload for creating Prestasi
type PrestasiCreateRequest struct {
	PesertaDidikID    *uint  `json:"peserta_didik_id" binding:"omitempty"`
	Jenis             string `json:"jenis" binding:"required"`
	NamaGrup          string `json:"nama_grup" binding:"omitempty"`
	NamaPrestasi      string `json:"nama_prestasi" binding:"required"`
	TingkatPrestasi   string `json:"tingkat_prestasi" binding:"omitempty"`
	Penyelenggara     string `json:"penyelenggara" binding:"omitempty"`
	TanggalLomba      string `json:"tanggal_lomba" binding:"required"` // Format: YYYY-MM-DD
	Juara             string `json:"juara" binding:"required"`
	Keterangan        string `json:"keterangan" binding:"omitempty"`
	EkstrakurikulerID *uint  `json:"ekstrakurikuler_id" binding:"omitempty"`
	TahunPelajaranID  uint   `json:"tahun_pelajaran_id" binding:"required"`
	// Anggota tim data
	AnggotaTim        []AnggotaTimCreateRequest `json:"anggota_tim" binding:"omitempty"`
}

// AnggotaTimCreateRequest represents anggota tim data for creating
type AnggotaTimCreateRequest struct {
	PesertaDidikID   uint `json:"peserta_didik_id" binding:"required"`
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
}

// PrestasiUpdateRequest represents the request payload for updating Prestasi
type PrestasiUpdateRequest struct {
	ID                uint                      `json:"id" binding:"required"`
	PesertaDidikID    *uint                     `json:"peserta_didik_id" binding:"omitempty"`
	Jenis             string                    `json:"jenis" binding:"omitempty"`
	NamaGrup          string                    `json:"nama_grup" binding:"omitempty"`
	NamaPrestasi      string                    `json:"nama_prestasi" binding:"omitempty"`
	TingkatPrestasi   string                    `json:"tingkat_prestasi" binding:"omitempty"`
	Penyelenggara     string                    `json:"penyelenggara" binding:"omitempty"`
	TanggalLomba      string                    `json:"tanggal_lomba" binding:"omitempty"` // Format: YYYY-MM-DD
	Juara             string                    `json:"juara" binding:"omitempty"`
	Keterangan        string                    `json:"keterangan" binding:"omitempty"`
	EkstrakurikulerID *uint                     `json:"ekstrakurikuler_id" binding:"omitempty"`
	TahunPelajaranID  uint                      `json:"tahun_pelajaran_id" binding:"omitempty"`
	FotoToDelete      []string                  `json:"foto_to_delete" binding:"omitempty"`
	// Anggota tim data
	AnggotaTim        []AnggotaTimUpdateRequest `json:"anggota_tim" binding:"omitempty"`
}

// AnggotaTimUpdateRequest represents anggota tim data for updating
type AnggotaTimUpdateRequest struct {
	ID               *uint  `json:"id" binding:"omitempty"` // If ID exists, update; if not, create new
	PesertaDidikID   uint   `json:"peserta_didik_id" binding:"required"`
	TahunPelajaranID uint   `json:"tahun_pelajaran_id" binding:"required"`
}

// PrestasiResponse represents the response payload for Prestasi
type PrestasiResponse struct {
	ID                 uint                        `json:"id"`
	PesertaDidikID     *uint                       `json:"peserta_didik_id"`
	PesertaDidik       *PesertaDidikResponse       `json:"peserta_didik"`
	Jenis              string                      `json:"jenis"`
	NamaGrup           string                      `json:"nama_grup"`
	NamaPrestasi       string                      `json:"nama_prestasi"`
	TingkatPrestasi    string                      `json:"tingkat_prestasi"`
	Penyelenggara      string                      `json:"penyelenggara"`
	TanggalLomba       time.Time                   `json:"tanggal_lomba"`
	Juara              string                      `json:"juara"`
	Keterangan         string                      `json:"keterangan"`
	Foto               []FotoItemDTO               `json:"foto"`
	EkstrakurikulerID  *uint                       `json:"ekstrakurikuler_id"`
	Ekstrakurikuler    *EkstrakurikulerDetailDTO   `json:"ekstrakurikuler"`
	TahunPelajaranID   uint                        `json:"tahun_pelajaran_id"`
	TahunPelajaran     *TahunPelajaranDetailResponse `json:"tahun_pelajaran"`
	AnggotaTimPrestasi []AnggotaTimPrestasiDTO     `json:"anggota_tim_prestasi"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
	CreatedByID        *uint                       `json:"created_by_id"`
	UpdatedByID        *uint                       `json:"updated_by_id"`
}

// PrestasiListResponse represents the response payload for listing Prestasi
type PrestasiListResponse struct {
	Data   []PrestasiResponse `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
	Total  int64              `json:"total"`
}

// PrestasiGetAllRequest represents the request payload for getting all prestasi with filters
type PrestasiGetAllRequest struct {
	Search struct {
		PesertaDidikID    *uint  `json:"peserta_didik_id"`
		NamaPesertaDidik  string `json:"nama_peserta_didik"`
		Jenis             string `json:"jenis"`
		NamaGrup          string `json:"nama_grup"`
		NamaPrestasi      string `json:"nama_prestasi"`
		TingkatPrestasi   string `json:"tingkat_prestasi"`
		Penyelenggara     string `json:"penyelenggara"`
		StartDate         string `json:"start_date"` // Format: YYYY-MM-DD
		EndDate           string `json:"end_date"`   // Format: YYYY-MM-DD
		Juara             string `json:"juara"`
		EkstrakurikulerID *uint  `json:"ekstrakurikuler_id"`
		TahunPelajaranID  *uint  `json:"tahun_pelajaran_id"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// PrestasiListWithPaginationResponse represents the response with pagination
type PrestasiListWithPaginationResponse struct {
	Data       []PrestasiResponse `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

// PrestasiPublicResponse represents the public response for prestasi
type PrestasiPublicResponse struct {
	ID               uint                        `json:"id"`
	Jenis            string                      `json:"jenis"`
	NamaPesertaDidik string                      `json:"nama_peserta_didik,omitempty"` // For individu
	NamaGrup         string                      `json:"nama_grup,omitempty"`          // For grup
	AnggotaTim       []AnggotaTimPublicDetail    `json:"anggota_tim,omitempty"`        // For grup
	NamaPrestasi     string                      `json:"nama_prestasi"`
	TingkatPrestasi  string                      `json:"tingkat_prestasi"`
	Juara            string                      `json:"juara"`
	TanggalLomba     time.Time                   `json:"tanggal_lomba"`
	FotoThumbnail    string                      `json:"foto_thumbnail"` // Active thumbnail only
}

// AnggotaTimPublicDetail represents anggota tim detail for public response
type AnggotaTimPublicDetail struct {
	Nama   string `json:"nama"`
	NIS    string `json:"nis"`
	Rombel string `json:"rombel"`
}

// PrestasiPublicListResponse represents the public list response
type PrestasiPublicListResponse struct {
	Data []PrestasiPublicResponse `json:"data"`
}
