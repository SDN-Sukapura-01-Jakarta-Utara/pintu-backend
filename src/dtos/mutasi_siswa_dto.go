package dtos

// MutasiSiswaCreateRequest represents the request payload for creating Mutasi Siswa (public)
type MutasiSiswaCreateRequest struct {
	TahunPelajaranID int     `form:"tahun_pelajaran_id" binding:"required"`
	Semester         int     `form:"semester" binding:"required"`
	NamaLengkap      string  `form:"nama_lengkap" binding:"required"`
	NamaPanggilan    *string `form:"nama_panggilan"`
	NISN             *string `form:"nisn"`
	TempatLahir      string  `form:"tempat_lahir" binding:"required"`
	TanggalLahir     string  `form:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
	JenisKelamin     string  `form:"jenis_kelamin" binding:"required"`
	Agama            string  `form:"agama" binding:"required"`
	GolonganDarah    *string `form:"golongan_darah"`
	AnakKe           *int    `form:"anak_ke"`
	JumlahSaudara    *int    `form:"jumlah_saudara"`
	StatusAnak       *string `form:"status_anak"`
	Alamat           string  `form:"alamat" binding:"required"`
	RT               *string `form:"rt"`
	RW               *string `form:"rw"`
	Kelurahan        *string `form:"kelurahan"`
	Kecamatan        *string `form:"kecamatan"`
	Kota             *string `form:"kota"`
	Provinsi         *string `form:"provinsi"`
	NamaAyah         *string `form:"nama_ayah"`
	NamaIbu          *string `form:"nama_ibu"`
	PendidikanAyah   *string `form:"pendidikan_ayah"`
	PendidikanIbu    *string `form:"pendidikan_ibu"`
	PekerjaanAyah    *string `form:"pekerjaan_ayah"`
	PekerjaanIbu     *string `form:"pekerjaan_ibu"`
	PenghasilanAyah  *float64 `form:"penghasilan_ayah"`
	PenghasilanIbu   *float64 `form:"penghasilan_ibu"`
	NomorHPOrtu      *string `form:"nomor_hp_ortu"`
	NamaWali         *string `form:"nama_wali"`
	PendidikanWali   *string `form:"pendidikan_wali"`
	HubunganWali     *string `form:"hubungan_wali"`
	PekerjaanWali    *string `form:"pekerjaan_wali"`
	NomorHPWali      *string `form:"nomor_hp_wali"`
	PindahanKelas    *int    `form:"pindahan_kelas"`
	AsalSekolah      *string `form:"asal_sekolah"`
	NamaAsalSekolah  *string `form:"nama_asal_sekolah"`
}

// MutasiSiswaGetAllRequest represents the request for getting all mutasi siswa with filters
type MutasiSiswaGetAllRequest struct {
	Search struct {
		TahunPelajaranID *int   `json:"tahun_pelajaran_id"`
		Semester         *int   `json:"semester"`
		StartDate        string `json:"start_date"` // YYYY-MM-DD
		EndDate          string `json:"end_date"`   // YYYY-MM-DD
		NamaSiswa        string `json:"nama_siswa"`
		NISN             string `json:"nisn"`
		TempatLahir      string `json:"tempat_lahir"`
		JenisKelamin     string `json:"jenis_kelamin"`
		PindahanKelas    *int   `json:"pindahan_kelas"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// MutasiSiswaResponse represents the response payload for Mutasi Siswa
type MutasiSiswaResponse struct {
	ID               uint   `json:"id"`
	NomorPendaftaran string `json:"nomor_pendaftaran"`
	TahunPelajaranID int    `json:"tahun_pelajaran_id"`
	TahunPelajaran   *struct {
		ID             uint   `json:"id"`
		TahunPelajaran string `json:"tahun_pelajaran"`
		Status         string `json:"status"`
	} `json:"tahun_pelajaran,omitempty"`
	Semester int `json:"semester"`
	NamaLengkap      string   `json:"nama_lengkap"`
	NamaPanggilan    *string  `json:"nama_panggilan"`
	NISN             *string  `json:"nisn"`
	TempatLahir      string   `json:"tempat_lahir"`
	TanggalLahir     string   `json:"tanggal_lahir"`
	JenisKelamin     string   `json:"jenis_kelamin"`
	Agama            string   `json:"agama"`
	GolonganDarah    *string  `json:"golongan_darah"`
	AnakKe           *int     `json:"anak_ke"`
	JumlahSaudara    *int     `json:"jumlah_saudara"`
	StatusAnak       *string  `json:"status_anak"`
	Alamat           string   `json:"alamat"`
	RT               *string  `json:"rt"`
	RW               *string  `json:"rw"`
	Kelurahan        *string  `json:"kelurahan"`
	Kecamatan        *string  `json:"kecamatan"`
	Kota             *string  `json:"kota"`
	Provinsi         *string  `json:"provinsi"`
	NamaAyah         *string  `json:"nama_ayah"`
	NamaIbu          *string  `json:"nama_ibu"`
	PendidikanAyah   *string  `json:"pendidikan_ayah"`
	PendidikanIbu    *string  `json:"pendidikan_ibu"`
	PekerjaanAyah    *string  `json:"pekerjaan_ayah"`
	PekerjaanIbu     *string  `json:"pekerjaan_ibu"`
	PenghasilanAyah  *float64 `json:"penghasilan_ayah"`
	PenghasilanIbu   *float64 `json:"penghasilan_ibu"`
	NomorHPOrtu      *string  `json:"nomor_hp_ortu"`
	NamaWali         *string  `json:"nama_wali"`
	PendidikanWali   *string  `json:"pendidikan_wali"`
	HubunganWali     *string  `json:"hubungan_wali"`
	PekerjaanWali    *string  `json:"pekerjaan_wali"`
	NomorHPWali      *string  `json:"nomor_hp_wali"`
	PindahanKelas    *int     `json:"pindahan_kelas"`
	AsalSekolah      *string  `json:"asal_sekolah"`
	NamaAsalSekolah  *string  `json:"nama_asal_sekolah"`
	Rapor            *string  `json:"rapor"`
	AkteKelahiran    *string  `json:"akte_kelahiran"`
	KartuKeluarga    *string  `json:"kartu_keluarga"`
	SPTJM            *string  `json:"sptjm"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

// MutasiSiswaListWithPaginationResponse represents paginated list response
type MutasiSiswaListWithPaginationResponse struct {
	Data       []MutasiSiswaResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}


// MutasiSiswaUpdateRequest represents the request payload for updating Mutasi Siswa
type MutasiSiswaUpdateRequest struct {
	ID               uint     `form:"id" binding:"required"`
	TahunPelajaranID *int     `form:"tahun_pelajaran_id"`
	Semester         *int     `form:"semester"`
	NamaLengkap      string   `form:"nama_lengkap"`
	NamaPanggilan    *string  `form:"nama_panggilan"`
	NISN             *string  `form:"nisn"`
	TempatLahir      string   `form:"tempat_lahir"`
	TanggalLahir     string   `form:"tanggal_lahir"` // Format: YYYY-MM-DD
	JenisKelamin     string   `form:"jenis_kelamin"`
	Agama            string   `form:"agama"`
	GolonganDarah    *string  `form:"golongan_darah"`
	AnakKe           *int     `form:"anak_ke"`
	JumlahSaudara    *int     `form:"jumlah_saudara"`
	StatusAnak       *string  `form:"status_anak"`
	Alamat           string   `form:"alamat"`
	RT               *string  `form:"rt"`
	RW               *string  `form:"rw"`
	Kelurahan        *string  `form:"kelurahan"`
	Kecamatan        *string  `form:"kecamatan"`
	Kota             *string  `form:"kota"`
	Provinsi         *string  `form:"provinsi"`
	NamaAyah         *string  `form:"nama_ayah"`
	NamaIbu          *string  `form:"nama_ibu"`
	PendidikanAyah   *string  `form:"pendidikan_ayah"`
	PendidikanIbu    *string  `form:"pendidikan_ibu"`
	PekerjaanAyah    *string  `form:"pekerjaan_ayah"`
	PekerjaanIbu     *string  `form:"pekerjaan_ibu"`
	PenghasilanAyah  *float64 `form:"penghasilan_ayah"`
	PenghasilanIbu   *float64 `form:"penghasilan_ibu"`
	NomorHPOrtu      *string  `form:"nomor_hp_ortu"`
	NamaWali         *string  `form:"nama_wali"`
	PendidikanWali   *string  `form:"pendidikan_wali"`
	HubunganWali     *string  `form:"hubungan_wali"`
	PekerjaanWali    *string  `form:"pekerjaan_wali"`
	NomorHPWali      *string  `form:"nomor_hp_wali"`
	PindahanKelas    *int     `form:"pindahan_kelas"`
	AsalSekolah      *string  `form:"asal_sekolah"`
	NamaAsalSekolah  *string  `form:"nama_asal_sekolah"`
}


// KonfigurasiMutasiSiswaRequest represents the request payload for setting konfigurasi mutasi siswa
type KonfigurasiMutasiSiswaRequest struct {
	TanggalBukaPendaftaran  string `form:"tanggal_buka_pendaftaran" binding:"required"`  // YYYY-MM-DD
	TanggalTutupPendaftaran string `form:"tanggal_tutup_pendaftaran" binding:"required"` // YYYY-MM-DD
	NamaKepalaSekolah       string `form:"nama_kepala_sekolah" binding:"required"`
	NIPKepalaSekolah        string `form:"nip_kepala_sekolah" binding:"required"`
	NamaKetuaPanitia        string `form:"nama_ketua_panitia" binding:"required"`
	NIPKetuaPanitia         string `form:"nip_ketua_panitia" binding:"required"`
	GrupWA                  string `form:"grup_wa"`
}

// KonfigurasiMutasiSiswaResponse represents the response payload for konfigurasi mutasi siswa
type KonfigurasiMutasiSiswaResponse struct {
	ID                      uint    `json:"id"`
	TanggalBukaPendaftaran  string  `json:"tanggal_buka_pendaftaran"`
	TanggalTutupPendaftaran string  `json:"tanggal_tutup_pendaftaran"`
	NamaKepalaSekolah       string  `json:"nama_kepala_sekolah"`
	NIPKepalaSekolah        string  `json:"nip_kepala_sekolah"`
	NamaKetuaPanitia        string  `json:"nama_ketua_panitia"`
	NIPKetuaPanitia         string  `json:"nip_ketua_panitia"`
	TemplateSPTJM           *string `json:"template_sptjm"`
	GrupWA                  *string `json:"grup_wa"`
	CreatedAt               string  `json:"created_at"`
	UpdatedAt               string  `json:"updated_at"`
}
