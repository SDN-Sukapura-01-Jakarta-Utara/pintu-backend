package models

import (
	"time"
)

// Absensi represents the Absensi model for student attendance tracking
type Absensi struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	PesertaDidikID   uint      `gorm:"column:peserta_didik_id;not null;index:idx_absensi_peserta_tanggal" json:"peserta_didik_id"`
	Tanggal          time.Time `gorm:"column:tanggal;type:date;not null;index:idx_absensi_peserta_tanggal" json:"tanggal"`
	JamDatang        *string   `gorm:"column:jam_datang;type:time" json:"jam_datang"`
	JamPulang        *string   `gorm:"column:jam_pulang;type:time" json:"jam_pulang"`
	Status           *string   `gorm:"column:status;size:20" json:"status"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	
	// Relationships
	PesertaDidik     *PesertaDidik `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
}

// TableName specifies the table name for Absensi
func (m *Absensi) TableName() string {
	return "absensi"
}
