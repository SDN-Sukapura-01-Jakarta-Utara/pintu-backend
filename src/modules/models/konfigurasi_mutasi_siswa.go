package models

import (
	"time"

	"gorm.io/gorm"
)

// KonfigurasiMutasiSiswa represents the Konfigurasi Mutasi Siswa model
type KonfigurasiMutasiSiswa struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	TanggalBukaPendaftaran  time.Time      `gorm:"type:date;not null" json:"tanggal_buka_pendaftaran"`
	TanggalTutupPendaftaran time.Time      `gorm:"type:date;not null" json:"tanggal_tutup_pendaftaran"`
	NamaKepalaSekolah       string         `gorm:"size:255;not null" json:"nama_kepala_sekolah"`
	NIPKepalaSekolah        string         `gorm:"column:nip_kepala_sekolah;size:50;not null" json:"nip_kepala_sekolah"`
	NamaKetuaPanitia        string         `gorm:"size:255;not null" json:"nama_ketua_panitia"`
	NIPKetuaPanitia         string         `gorm:"column:nip_ketua_panitia;size:50;not null" json:"nip_ketua_panitia"`
	TemplateSPTJM           *string        `gorm:"size:255" json:"template_sptjm"`
	GrupWA                  *string        `gorm:"type:text" json:"grup_wa"`
	CreatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedByID             *uint          `json:"created_by_id"`
	UpdatedByID             *uint          `json:"updated_by_id"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for KonfigurasiMutasiSiswa
func (m *KonfigurasiMutasiSiswa) TableName() string {
	return "konfigurasi_mutasi_siswa"
}
