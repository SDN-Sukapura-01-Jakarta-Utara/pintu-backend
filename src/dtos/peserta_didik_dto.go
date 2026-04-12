package dtos

// PesertaDidikCreateRequest represents the request payload for creating PesertaDidik
type PesertaDidikCreateRequest struct {
	Nama               string `json:"nama" binding:"required"`
	NIS                string `json:"nis" binding:"required"`
	JenisKelamin       string `json:"jenis_kelamin" binding:"required"`
	NISN               string `json:"nisn" binding:"required"`
	TempatLahir        string `json:"tempat_lahir" binding:"omitempty"`
	TanggalLahir       string `json:"tanggal_lahir" binding:"omitempty"` // Format: YYYY-MM-DD
	NIK                string `json:"nik" binding:"omitempty"`
	Agama              string `json:"agama" binding:"omitempty"`
	Alamat             string `json:"alamat" binding:"omitempty"`
	RT                 string `json:"rt" binding:"omitempty"`
	RW                 string `json:"rw" binding:"omitempty"`
	Kelurahan          string `json:"kelurahan" binding:"omitempty"`
	Kecamatan          string `json:"kecamatan" binding:"omitempty"`
	KodePos            string `json:"kode_pos" binding:"omitempty"`
	NamaAyah           string `json:"nama_ayah" binding:"omitempty"`
	NamaIbu            string `json:"nama_ibu" binding:"omitempty"`
	RombelID           *uint  `json:"rombel_id" binding:"omitempty"`
	TahunPelajaranID   *uint  `json:"tahun_pelajaran_id" binding:"omitempty"`
	Status             string `json:"status" binding:"omitempty,oneof=active inactive"`
	Username           string `json:"username" binding:"omitempty"`
	Password           string `json:"password" binding:"omitempty,min=3"`
	RoleIDs            []uint `json:"role_ids" binding:"omitempty"`
}

// PesertaDidikUpdateRequest represents the request payload for updating PesertaDidik
type PesertaDidikUpdateRequest struct {
	ID                 uint   `json:"id" binding:"required"`
	Nama               string `json:"nama" binding:"omitempty"`
	NIS                string `json:"nis" binding:"omitempty"`
	JenisKelamin       string `json:"jenis_kelamin" binding:"omitempty"`
	NISN               string `json:"nisn" binding:"omitempty"`
	TempatLahir        string `json:"tempat_lahir" binding:"omitempty"`
	TanggalLahir       string `json:"tanggal_lahir" binding:"omitempty"` // Format: YYYY-MM-DD
	NIK                string `json:"nik" binding:"omitempty"`
	Agama              string `json:"agama" binding:"omitempty"`
	Alamat             string `json:"alamat" binding:"omitempty"`
	RT                 string `json:"rt" binding:"omitempty"`
	RW                 string `json:"rw" binding:"omitempty"`
	Kelurahan          string `json:"kelurahan" binding:"omitempty"`
	Kecamatan          string `json:"kecamatan" binding:"omitempty"`
	KodePos            string `json:"kode_pos" binding:"omitempty"`
	NamaAyah           string `json:"nama_ayah" binding:"omitempty"`
	NamaIbu            string `json:"nama_ibu" binding:"omitempty"`
	RombelID           *uint  `json:"rombel_id" binding:"omitempty"`
	TahunPelajaranID   *uint  `json:"tahun_pelajaran_id" binding:"omitempty"`
	Status             string `json:"status" binding:"omitempty,oneof=active inactive"`
	Username           string `json:"username" binding:"omitempty"`
	Password           string `json:"password" binding:"omitempty,min=3"`
	RoleIDs            []uint `json:"role_ids" binding:"omitempty"`
}

// KelasDetailResponse represents kelas details in response
type KelasDetailResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RombelDetailResponse represents rombel details with kelas in response
type RombelDetailResponse struct {
	ID        uint                  `json:"id"`
	Name      string                `json:"name"`
	Status    string                `json:"status"`
	KelasID   uint                  `json:"kelas_id"`
	Kelas     *KelasDetailResponse  `json:"kelas"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

// TahunPelajaranDetailResponse represents tahun pelajaran details in response
type TahunPelajaranDetailResponse struct {
	ID               uint   `json:"id"`
	TahunPelajaran   string `json:"tahun_pelajaran"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	CreatedByID      *uint  `json:"created_by_id"`
	UpdatedByID      *uint  `json:"updated_by_id"`
}

// PesertaDidikResponse represents the response payload for PesertaDidik
type PesertaDidikResponse struct {
	ID               uint                          `json:"id"`
	Nama             string                        `json:"nama"`
	NIS              string                        `json:"nis"`
	JenisKelamin     string                        `json:"jenis_kelamin"`
	NISN             string                        `json:"nisn"`
	TempatLahir      string                        `json:"tempat_lahir"`
	TanggalLahir     string                        `json:"tanggal_lahir"`
	NIK              string                        `json:"nik"`
	Agama            string                        `json:"agama"`
	Alamat           string                        `json:"alamat"`
	RT               string                        `json:"rt"`
	RW               string                        `json:"rw"`
	Kelurahan        string                        `json:"kelurahan"`
	Kecamatan        string                        `json:"kecamatan"`
	KodePos          string                        `json:"kode_pos"`
	NamaAyah         string                        `json:"nama_ayah"`
	NamaIbu          string                        `json:"nama_ibu"`
	RombelID         *uint                         `json:"rombel_id"`
	Rombel           *RombelDetailResponse         `json:"rombel"`
	TahunPelajaranID *uint                         `json:"tahun_pelajaran_id"`
	TahunPelajaran   *TahunPelajaranDetailResponse `json:"tahun_pelajaran"`
	Status           string                        `json:"status"`
	Username         string                        `json:"username"`
	Roles            []RoleResponse                `json:"roles"`
	CreatedAt        string                        `json:"created_at"`
	UpdatedAt        string                        `json:"updated_at"`
	CreatedByID      *uint                         `json:"created_by_id"`
	UpdatedByID      *uint                         `json:"updated_by_id"`
}

// PesertaDidikListResponse represents list response
type PesertaDidikListResponse struct {
	Data   []PesertaDidikResponse `json:"data"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
	Total  int64                  `json:"total"`
}

// PesertaDidikListWithPaginationResponse represents list response with pagination info
type PesertaDidikListWithPaginationResponse struct {
	Data       []PesertaDidikResponse `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
}

// PesertaDidikGetAllRequest represents the request payload for getting all PesertaDidik with filters
type PesertaDidikGetAllRequest struct {
	Search struct {
		TahunPelajaranID *uint  `json:"tahun_pelajaran_id" binding:"omitempty"`
		RombelID         *uint  `json:"rombel_id" binding:"omitempty"`
		Nama             string `json:"nama" binding:"omitempty"`
		NIS              string `json:"nis" binding:"omitempty"`
		JenisKelamin     string `json:"jenis_kelamin" binding:"omitempty"`
		NISN             string `json:"nisn" binding:"omitempty"`
		TempatLahir      string `json:"tempat_lahir" binding:"omitempty"`
		NIK              string `json:"nik" binding:"omitempty"`
		Agama            string `json:"agama" binding:"omitempty"`
		Status           string `json:"status" binding:"omitempty"`
	} `json:"search" binding:"omitempty"`
	Pagination struct {
		Limit int `json:"limit" binding:"omitempty"`
		Page  int `json:"page" binding:"omitempty"`
	} `json:"pagination" binding:"omitempty"`
}

// ImportExcelResponse represents the response for import Excel operation
type ImportExcelResponse struct {
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
	Errors       []ImportExcelRowError `json:"errors"`
}

// ImportExcelRowError represents an error for a specific row during import
type ImportExcelRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}
