package models

import (
	"time"

	"gorm.io/gorm"
)

// KegiatanEkskul represents ekstrakurikuler activity session
type KegiatanEkskul struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	EkstrakurikulerID  uint            `gorm:"not null" json:"ekstrakurikuler_id"`
	TahunPelajaranID   uint            `gorm:"not null" json:"tahun_pelajaran_id"`
	TanggalKegiatan    time.Time       `gorm:"type:date;not null" json:"tanggal_kegiatan"`
	WaktuMulai         *string         `gorm:"type:time" json:"waktu_mulai"`
	WaktuSelesai       *string         `gorm:"type:time" json:"waktu_selesai"`
	MateriKegiatan     string          `gorm:"type:text;not null" json:"materi_kegiatan"`
	FotoKegiatan       *string         `gorm:"type:json" json:"foto_kegiatan"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	CreatedByID        *uint           `json:"created_by_id,omitempty"`
	UpdatedByID        *uint           `json:"updated_by_id,omitempty"`

	// Relations
	Ekstrakurikuler    *Ekstrakurikuler     `gorm:"foreignKey:EkstrakurikulerID" json:"ekstrakurikuler,omitempty"`
	TahunPelajaran     *TahunPelajaran      `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	AbsensiEkskul      []AbsensiEkskul      `gorm:"foreignKey:KegiatanEkskulID" json:"absensi_ekskul,omitempty"`
	AbsensiPelatih     []AbsensiPelatihEkskul `gorm:"foreignKey:KegiatanEkskulID" json:"absensi_pelatih,omitempty"`
}

// TableName specifies the table name for KegiatanEkskul model
func (KegiatanEkskul) TableName() string {
	return "kegiatan_ekskul"
}
