package models

import (
	"time"

	"gorm.io/gorm"
)

// AbsensiEkskul represents student attendance for ekstrakurikuler activity
type AbsensiEkskul struct {
	ID                   uint                 `gorm:"primaryKey" json:"id"`
	KegiatanEkskulID     uint                 `gorm:"not null" json:"kegiatan_ekskul_id"`
	PesertaDidikRombelID uint                 `gorm:"not null" json:"peserta_didik_rombel_id"`
	Status               string               `gorm:"type:varchar(20);not null" json:"status"` // hadir, sakit, izin, alpha
	Keterangan           *string              `gorm:"type:text" json:"keterangan"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	DeletedAt            gorm.DeletedAt       `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	KegiatanEkskul       *KegiatanEkskul      `gorm:"foreignKey:KegiatanEkskulID" json:"kegiatan_ekskul,omitempty"`
	PesertaDidikRombel   *PesertaDidikRombel  `gorm:"foreignKey:PesertaDidikRombelID" json:"peserta_didik_rombel,omitempty"`
}

// TableName specifies the table name for AbsensiEkskul model
func (AbsensiEkskul) TableName() string {
	return "absensi_ekskul"
}
