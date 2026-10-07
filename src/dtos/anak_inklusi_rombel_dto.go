package dtos

type AnakInklusiRombelSyncRequest struct {
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
}

type AnakInklusiRombelGetAllRequest struct {
	TahunPelajaranID uint `json:"tahun_pelajaran_id" binding:"required"`
	Search struct {
		NamaPesertaDidik string `json:"nama_peserta_didik"`
		NIS              string `json:"nis"`
		NISN             string `json:"nisn"`
		JenisHambatan    string `json:"jenis_hambatan"`
		Status           string `json:"status"` // Status anak inklusi
		RombelID         uint   `json:"rombel_id"` // Filter by rombel ID
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

type AnakInklusiRombelResponse struct {
	ID                     uint                          `json:"id"`
	AnakInklusiID          uint                          `json:"anak_inklusi_id"`
	PesertaDidikRombelID   uint                          `json:"peserta_didik_rombel_id"`
	GuruPendampingKhususID *uint                         `json:"guru_pendamping_khusus_id"`
	GuruKelasID            *uint                         `json:"guru_kelas_id"`
	Catatan                string                        `json:"catatan"`
	AnakInklusi            *AnakInklusiResponse          `json:"anak_inklusi,omitempty"`
	PesertaDidikRombel     *PesertaDidikRombelSimple     `json:"peserta_didik_rombel,omitempty"`
	GuruPendampingKhusus   *KepegawaianSimple            `json:"guru_pendamping_khusus,omitempty"`
	GuruKelas              *KepegawaianSimple            `json:"guru_kelas,omitempty"`
	CreatedAt              string                        `json:"created_at"`
	UpdatedAt              string                        `json:"updated_at"`
}

type AnakInklusiRombelListWithPaginationResponse struct {
	Data       []AnakInklusiRombelResponse `json:"data"`
	Pagination PaginationResponse          `json:"pagination"`
}

type AnakInklusiRombelSyncResponse struct {
	TotalProcessed int                            `json:"total_processed"`
	TotalSynced    int                            `json:"total_synced"`
	TotalSkipped   int                            `json:"total_skipped"`
	Details        []AnakInklusiRombelSyncDetail  `json:"details"`
}

type AnakInklusiRombelSyncDetail struct {
	PesertaDidikID       uint   `json:"peserta_didik_id"`
	PesertaDidikNama     string `json:"peserta_didik_nama"`
	PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
	AnakInklusiID        uint   `json:"anak_inklusi_id"`
	Status               string `json:"status"` // "synced" or "skipped"
	Message              string `json:"message"`
}

type AnakInklusiRombelUpdateRequest struct {
	ID                     uint   `json:"id" binding:"required"`
	GuruPendampingKhususID *uint  `json:"guru_pendamping_khusus_id"`
	GuruKelasID            *uint  `json:"guru_kelas_id"`
	Catatan                string `json:"catatan"`
}

type KepegawaianSimple struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	NIP      string `json:"nip"`
	Foto     string `json:"foto"`
	FotoURL  string `json:"foto_url"`
	Kategori string `json:"kategori"`
	Jabatan  string `json:"jabatan"`
}

type RombelSimple struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type TahunPelajaranSimple struct {
	ID             uint   `json:"id"`
	TahunPelajaran string `json:"tahun_pelajaran"`
	Status         string `json:"status"`
}

type PesertaDidikRombelSimple struct {
	ID               uint                  `json:"id"`
	PesertaDidikID   uint                  `json:"peserta_didik_id"`
	RombelID         uint                  `json:"rombel_id"`
	TahunPelajaranID uint                  `json:"tahun_pelajaran_id"`
	Rombel           *RombelSimple         `json:"rombel,omitempty"`
	TahunPelajaran   *TahunPelajaranSimple `json:"tahun_pelajaran,omitempty"`
}
