package models

import (
	"time"

	"gorm.io/gorm"
)

// PesertaDidik represents the PesertaDidik model
type PesertaDidik struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Nama               string         `gorm:"not null" json:"nama"`
	NIS                string         `gorm:"not null;unique" json:"nis"`
	JenisKelamin       string         `gorm:"not null" json:"jenis_kelamin"`
	NISN               string         `gorm:"not null;unique" json:"nisn"`
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
	Status             string         `gorm:"default:active" json:"status"`
	Username           string         `json:"username"`
	Password           string         `json:"password,omitempty"`
	Photo              string         `json:"photo,omitempty"`
	Barcode            string         `json:"barcode,omitempty"`
	BarcodeGeneratedAt *time.Time     `json:"barcode_generated_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedByID        *uint          `json:"created_by_id"`
	UpdatedByID        *uint          `json:"updated_by_id"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// Many-to-many relationship
	Roles              []Role         `gorm:"many2many:peserta_didik_roles" json:"roles,omitempty"`
}

// TableName specifies the table name for PesertaDidik
func (m *PesertaDidik) TableName() string {
	return "peserta_didik"
}
