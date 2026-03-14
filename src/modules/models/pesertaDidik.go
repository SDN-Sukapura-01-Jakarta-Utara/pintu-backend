package models

import (
	"time"

	"gorm.io/gorm"
)

// PesertaDidik represents the PesertaDidik model
type PesertaDidik struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Nama               string         `gorm:"not null" json:"nama"`
	NIS                string         `gorm:"not null" json:"nis"`
	JenisKelamin       string         `gorm:"not null" json:"jenis_kelamin"`
	NISN               string         `gorm:"not null" json:"nisn"`
	TempatLahir        string         `json:"tempat_lahir"`
	TanggalLahir       *time.Time     `json:"tanggal_lahir"`
	NIK                string         `json:"nik"`
	Agama              string         `json:"agama"`
	Alamat             string         `json:"alamat"`
	RT                 string         `json:"rt"`
	RW                 string         `json:"rw"`
	Kelurahan          string         `json:"kelurahan"`
	Kecamatan          string         `json:"kecamatan"`
	KodePos            string         `json:"kode_pos"`
	NamaAyah           string         `json:"nama_ayah"`
	NamaIbu            string         `json:"nama_ibu"`
	RombelID           *uint          `gorm:"column:rombel_id" json:"rombel_id"`
	TahunPelajaranID   *uint          `gorm:"column:tahun_pelajaran_id" json:"tahun_pelajaran_id"`
	Status             string         `gorm:"default:active" json:"status"`
	Username           string         `json:"username"`
	Password           string         `gorm:"not null" json:"password,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedByID        *uint          `json:"created_by_id"`
	UpdatedByID        *uint          `json:"updated_by_id"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// Foreign key relationships
	Rombel             *Rombel        `gorm:"foreignKey:RombelID" json:"rombel,omitempty"`
	TahunPelajaran     *TahunPelajaran `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	// Many-to-many relationship
	Roles              []Role         `gorm:"many2many:peserta_didik_roles" json:"roles,omitempty"`
}

// TableName specifies the table name for PesertaDidik
func (m *PesertaDidik) TableName() string {
	return "peserta_didik"
}
