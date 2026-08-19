package models

import (
	"time"

	"gorm.io/gorm"
)

// AbsensiPelatihEkskul represents coach attendance for ekstrakurikuler activity
// Record exists = coach is present (hadir), no record = coach is absent (tidak hadir)
type AbsensiPelatihEkskul struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	KegiatanEkskulID uint            `gorm:"not null" json:"kegiatan_ekskul_id"`
	PelatihID        uint            `gorm:"not null" json:"pelatih_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	KegiatanEkskul   *KegiatanEkskul `gorm:"foreignKey:KegiatanEkskulID" json:"kegiatan_ekskul,omitempty"`
	Pelatih          *Pelatih        `gorm:"foreignKey:PelatihID" json:"pelatih,omitempty"`
}

// TableName specifies the table name for AbsensiPelatihEkskul model
func (AbsensiPelatihEkskul) TableName() string {
	return "absensi_pelatih_ekskul"
}
