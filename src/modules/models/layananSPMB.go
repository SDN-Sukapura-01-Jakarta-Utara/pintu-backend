package models

import (
	"time"

	"gorm.io/gorm"
)

// LayananSPMB represents the Layanan SPMB model
type LayananSPMB struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	NamaOrangTua     string         `gorm:"size:255;not null" json:"nama_orang_tua"`
	NomorTelepon     string         `gorm:"size:20;not null" json:"nomor_telepon"`
	Alamat           string         `gorm:"type:text;not null" json:"alamat"`
	NamaLengkapMurid string         `gorm:"size:255;not null" json:"nama_lengkap_murid"`
	Keperluan        string         `gorm:"type:text;not null" json:"keperluan"`
	TanggalLaporan   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_laporan"`
	Status           string         `gorm:"size:50;default:'pending'" json:"status"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for LayananSPMB
func (m *LayananSPMB) TableName() string {
	return "layanan_spmb"
}
