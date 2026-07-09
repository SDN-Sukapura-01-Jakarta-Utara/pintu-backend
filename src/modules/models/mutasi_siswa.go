package models

import (
	"time"

	"gorm.io/gorm"
)

// MutasiSiswa represents the Mutasi Siswa model
type MutasiSiswa struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	NomorPendaftaran   string          `gorm:"size:50;not null" json:"nomor_pendaftaran"`
	TahunPelajaranID   int             `gorm:"not null" json:"tahun_pelajaran_id"`
	TahunPelajaran     *TahunPelajaran `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	Semester           int             `gorm:"not null" json:"semester"`
	NamaLengkap        string         `gorm:"size:255;not null" json:"nama_lengkap"`
	NamaPanggilan      *string        `gorm:"size:100" json:"nama_panggilan"`
	NISN               *string        `gorm:"size:25" json:"nisn"`
	TempatLahir        string         `gorm:"size:100;not null" json:"tempat_lahir"`
	TanggalLahir       time.Time      `gorm:"type:date;not null" json:"tanggal_lahir"`
	JenisKelamin       string         `gorm:"size:20;not null" json:"jenis_kelamin"`
	Agama              string         `gorm:"size:50;not null" json:"agama"`
	GolonganDarah      *string        `gorm:"size:5" json:"golongan_darah"`
	AnakKe             *int           `json:"anak_ke"`
	JumlahSaudara      *int           `json:"jumlah_saudara"`
	StatusAnak         *string        `gorm:"size:50" json:"status_anak"`
	Alamat             string         `gorm:"type:text;not null" json:"alamat"`
	RT                 *string        `gorm:"size:10" json:"rt"`
	RW                 *string        `gorm:"size:10" json:"rw"`
	Kelurahan          *string        `gorm:"size:100" json:"kelurahan"`
	Kecamatan          *string        `gorm:"size:100" json:"kecamatan"`
	Kota               *string        `gorm:"size:100" json:"kota"`
	Provinsi           *string        `gorm:"size:100" json:"provinsi"`
	NamaAyah           *string        `gorm:"size:255" json:"nama_ayah"`
	NamaIbu            *string        `gorm:"size:255" json:"nama_ibu"`
	PendidikanAyah     *string        `gorm:"size:100" json:"pendidikan_ayah"`
	PendidikanIbu      *string        `gorm:"size:100" json:"pendidikan_ibu"`
	PekerjaanAyah      *string        `gorm:"size:255" json:"pekerjaan_ayah"`
	PekerjaanIbu       *string        `gorm:"size:255" json:"pekerjaan_ibu"`
	PenghasilanAyah    *float64       `gorm:"type:decimal(15,2)" json:"penghasilan_ayah"`
	PenghasilanIbu     *float64       `gorm:"type:decimal(15,2)" json:"penghasilan_ibu"`
	NomorHPOrtu        *string        `gorm:"size:20" json:"nomor_hp_ortu"`
	NamaWali           *string        `gorm:"size:255" json:"nama_wali"`
	PendidikanWali     *string        `gorm:"size:100" json:"pendidikan_wali"`
	HubunganWali       *string        `gorm:"size:100" json:"hubungan_wali"`
	PekerjaanWali      *string        `gorm:"size:255" json:"pekerjaan_wali"`
	NomorHPWali        *string        `gorm:"size:20" json:"nomor_hp_wali"`
	PindahanKelas      *int           `json:"pindahan_kelas"`
	AsalSekolah        *string        `gorm:"size:50" json:"asal_sekolah"`
	NamaAsalSekolah    *string        `gorm:"size:255" json:"nama_asal_sekolah"`
	Rapor              *string        `gorm:"size:255" json:"rapor"`
	AkteKelahiran      *string        `gorm:"size:255" json:"akte_kelahiran"`
	KartuKeluarga      *string        `gorm:"size:255" json:"kartu_keluarga"`
	SPTJM              *string        `gorm:"size:255" json:"sptjm"`
	CreatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for MutasiSiswa
func (m *MutasiSiswa) TableName() string {
	return "mutasi_siswa"
}
