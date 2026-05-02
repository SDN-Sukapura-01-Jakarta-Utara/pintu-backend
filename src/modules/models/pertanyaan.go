package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Pertanyaan represents the Pertanyaan model
type Pertanyaan struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	IDTiket            string         `gorm:"unique;not null" json:"id_tiket"`
	TanggalPengajuan   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_pengajuan"`
	Nama               string         `gorm:"size:255;not null" json:"nama"`
	Email              string         `gorm:"size:255;not null" json:"email"`
	Telepon            string         `gorm:"size:20" json:"telepon"`
	Kategori           string         `gorm:"size:100;not null" json:"kategori"`
	Prioritas          string         `gorm:"size:50;default:'Sedang'" json:"prioritas"`
	Judul              string         `gorm:"size:255;not null" json:"judul"`
	Deskripsi          string         `gorm:"type:text;not null" json:"deskripsi"`
	FilePertanyaan     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_pertanyaan"`
	JudulJawaban       *string        `gorm:"size:255" json:"judul_jawaban"`
	DeskripsiJawaban   *string        `gorm:"type:text" json:"deskripsi_jawaban"`
	FileJawaban        datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_jawaban"`
	TanggalProses      *time.Time     `json:"tanggal_proses"`
	EmailTerkirim      bool           `gorm:"default:false" json:"email_terkirim"`
	TanggalSelesai     *time.Time     `json:"tanggal_selesai"`
	Status             string         `gorm:"size:50;default:'pending'" json:"status"`
	RepliedBy          *uint          `json:"replied_by"`
	CreatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DeletedByID        *uint          `json:"deleted_by_id"`
}

// TableName specifies the table name for Pertanyaan
func (m *Pertanyaan) TableName() string {
	return "pertanyaan"
}
