package models

import (
	"time"

	"gorm.io/gorm"
)

type AnakInklusi struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	PesertaDidikID         uint           `gorm:"not null" json:"peserta_didik_id"`
	JenisHambatan          string         `gorm:"type:varchar(100);not null" json:"jenis_hambatan"`
	TanggalIdentifikasi    *time.Time     `json:"tanggal_identifikasi"`
	TanggalDiagnosa        *time.Time     `json:"tanggal_diagnosa"`
	TanggalKadaluarsaSurat *time.Time     `json:"tanggal_kadaluarsa_surat"`
	FileSuratDokter        string         `gorm:"type:varchar(255)" json:"file_surat_dokter"`
	Status                 string         `gorm:"type:varchar(50);not null;default:'identified'" json:"status"`
	Catatan                string         `gorm:"type:text" json:"catatan"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	CreatedByID            *uint          `json:"created_by_id"`
	UpdatedByID            *uint          `json:"updated_by_id"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	PesertaDidik *PesertaDidik `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
}

func (AnakInklusi) TableName() string {
	return "anak_inklusi"
}
