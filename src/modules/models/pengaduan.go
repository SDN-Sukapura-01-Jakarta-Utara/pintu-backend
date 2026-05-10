package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Pengaduan represents the Pengaduan model
type Pengaduan struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	IDTiket            string         `gorm:"unique;not null" json:"id_tiket"`
	TanggalPengajuan   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_pengajuan"`
	TipePelapor        string         `gorm:"size:50;default:'anonim'" json:"tipe_pelapor"`
	Nama               *string        `gorm:"size:255" json:"nama"`
	Email              *string        `gorm:"size:255" json:"email"`
	Telepon            *string        `gorm:"size:20" json:"telepon"`
	Kategori           string         `gorm:"size:100;not null" json:"kategori"`
	Prioritas          string         `gorm:"size:50;default:'Sedang'" json:"prioritas"`
	Judul              string         `gorm:"size:255;not null" json:"judul"`
	Deskripsi          string         `gorm:"type:text;not null" json:"deskripsi"`
	FilePengaduan     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_pengaduan"`
	JudulJawaban       *string        `gorm:"size:255" json:"judul_jawaban"`
	DeskripsiJawaban   *string        `gorm:"type:text" json:"deskripsi_jawaban"`
	FileJawaban        datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_jawaban"`
	TindakLanjut       *string        `gorm:"type:text" json:"tindak_lanjut"`
	FileTindakLanjut   datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_tindak_lanjut"`
	TanggalProses      *time.Time     `json:"tanggal_proses"`
	EmailTerkirim      bool           `gorm:"default:false" json:"email_terkirim"`
	TanggalSelesai     *time.Time     `json:"tanggal_selesai"`
	Status             string         `gorm:"size:50;default:'pending'" json:"status"`
	RepliedBy          *uint          `json:"replied_by"`
	CreatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DeletedByID        *uint          `json:"deleted_by_id"`
}

// TableName specifies the table name for Pengaduan
func (m *Pengaduan) TableName() string {
	return "pengaduan"
}