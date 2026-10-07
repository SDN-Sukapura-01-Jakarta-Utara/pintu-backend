package dtos

type AnakInklusiCreateRequest struct {
	PesertaDidikID         uint   `form:"peserta_didik_id" json:"peserta_didik_id" binding:"required"`
	JenisHambatan          string `form:"jenis_hambatan" json:"jenis_hambatan" binding:"required"`
	TanggalIdentifikasi    string `form:"tanggal_identifikasi" json:"tanggal_identifikasi"`
	TanggalDiagnosa        string `form:"tanggal_diagnosa" json:"tanggal_diagnosa"`
	TanggalKadaluarsaSurat string `form:"tanggal_kadaluarsa_surat" json:"tanggal_kadaluarsa_surat"`
	Status                 string `form:"status" json:"status"`
	Catatan                string `form:"catatan" json:"catatan"`
}

type AnakInklusiUpdateRequest struct {
	ID                     uint   `json:"id"`
	JenisHambatan          string `json:"jenis_hambatan"`
	TanggalIdentifikasi    string `json:"tanggal_identifikasi"`
	TanggalDiagnosa        string `json:"tanggal_diagnosa"`
	TanggalKadaluarsaSurat string `json:"tanggal_kadaluarsa_surat"`
	Status                 string `json:"status"`
	Catatan                string `json:"catatan"`
	DeleteFile             bool   `json:"delete_file"` // Set true to delete existing file
}

type AnakInklusiGetAllRequest struct {
	Search struct {
		NamaPesertaDidik string `json:"nama_peserta_didik"`
		NIS              string `json:"nis"`
		NISN             string `json:"nisn"`
		JenisHambatan    string `json:"jenis_hambatan"`
		StartDate        string `json:"start_date"`
		EndDate          string `json:"end_date"`
		Status           string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

type AnakInklusiResponse struct {
	ID                     uint                `json:"id"`
	PesertaDidikID         uint                `json:"peserta_didik_id"`
	JenisHambatan          string              `json:"jenis_hambatan"`
	TanggalIdentifikasi    *string             `json:"tanggal_identifikasi"`
	TanggalDiagnosa        *string             `json:"tanggal_diagnosa"`
	TanggalKadaluarsaSurat *string             `json:"tanggal_kadaluarsa_surat"`
	FileSuratDokter        string              `json:"file_surat_dokter"`
	FileSuratDokterURL     string              `json:"file_surat_dokter_url"`
	Status                 string              `json:"status"`
	Catatan                string              `json:"catatan"`
	PesertaDidik           *PesertaDidikSimple `json:"peserta_didik,omitempty"`
	CreatedAt              string              `json:"created_at"`
	UpdatedAt              string              `json:"updated_at"`
}

type AnakInklusiListWithPaginationResponse struct {
	Data       []AnakInklusiResponse `json:"data"`
	Pagination PaginationResponse    `json:"pagination"`
}
