package models

import (
	"time"
)

// KonfigurasiAbsensi represents the Konfigurasi Absensi model
type KonfigurasiAbsensi struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	JamDatangMulai    string    `gorm:"column:jam_datang_mulai;type:time;not null" json:"jam_datang_mulai"`
	JamMaxDatang      string    `gorm:"column:jam_max_datang;type:time;not null" json:"jam_max_datang"`
	JamDatangSelesai  string    `gorm:"column:jam_datang_selesai;type:time;not null" json:"jam_datang_selesai"`
	JamPulangMulai    string    `gorm:"column:jam_pulang_mulai;type:time;not null" json:"jam_pulang_mulai"`
	JamPulangSelesai  string    `gorm:"column:jam_pulang_selesai;type:time;not null" json:"jam_pulang_selesai"`
	NamaKepsek        *string   `gorm:"column:nama_kepsek;size:200" json:"nama_kepsek"`
	NIPKepsek         *string   `gorm:"column:nip_kepsek;size:50" json:"nip_kepsek"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for KonfigurasiAbsensi
func (m *KonfigurasiAbsensi) TableName() string {
	return "konfigurasi_absensi"
}
