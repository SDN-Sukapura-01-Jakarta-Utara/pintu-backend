package dtos

// IDRequest represents a simple request with only ID
type IDRequest struct {
	ID uint `json:"id" binding:"required"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PesertaDidikSimple represents simplified peserta didik information
type PesertaDidikSimple struct {
	ID       uint   `json:"id"`
	NIS      string `json:"nis"`
	NISN     string `json:"nisn"`
	Nama     string `json:"nama"`
	Photo    string `json:"photo"`
	PhotoURL string `json:"photo_url"`
	NamaAyah string `json:"nama_ayah"`
	NamaIbu  string `json:"nama_ibu"`
	Status   string `json:"status"` // Status peserta didik (active, lulus, pindah, etc)
}
