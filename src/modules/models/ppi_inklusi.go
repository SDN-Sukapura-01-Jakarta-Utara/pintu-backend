package models

import (
	"time"

	"gorm.io/gorm"
)

type PPIInklusi struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	AnakInklusiRombelID   uint           `gorm:"not null" json:"anak_inklusi_rombel_id"`
	TahunPelajaranID      uint           `gorm:"not null" json:"tahun_pelajaran_id"`
	AspekPembelajaran     string         `gorm:"type:varchar(100);not null" json:"aspek_pembelajaran"`
	NamaProgram           string         `gorm:"type:varchar(255);not null" json:"nama_program"`
	TujuanPembelajaran    string         `gorm:"type:text;not null" json:"tujuan_pembelajaran"`
	StrategiPembelajaran  string         `gorm:"type:text;not null" json:"strategi_pembelajaran"`
	TargetWaktu           string         `gorm:"type:varchar(100)" json:"target_waktu"`
	CapaianSaatIni        string         `gorm:"type:text" json:"capaian_saat_ini"`
	StatusCapaian         string         `gorm:"type:varchar(100);default:'Belum Dimulai'" json:"status_capaian"`
	GuruPembuatID         *uint          `json:"guru_pembuat_id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedByID           *uint          `json:"created_by_id"`
	UpdatedByID           *uint          `json:"updated_by_id"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	AnakInklusiRombel *AnakInklusiRombel `gorm:"foreignKey:AnakInklusiRombelID" json:"anak_inklusi_rombel,omitempty"`
	TahunPelajaran    *TahunPelajaran    `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	GuruPembuat       *Kepegawaian       `gorm:"foreignKey:GuruPembuatID" json:"guru_pembuat,omitempty"`
}

func (PPIInklusi) TableName() string {
	return "ppi_inklusi"
}
