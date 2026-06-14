package models

import (
	"time"

	"gorm.io/gorm"
)

// SettingLayananSPMB represents the Setting Layanan SPMB model
type SettingLayananSPMB struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	NamaKepalaSekolah   *string        `gorm:"column:nama_kepala_sekolah;size:255" json:"nama_kepala_sekolah"`
	NIPKepalaSekolah    *string        `gorm:"column:nip_kepala_sekolah;size:50" json:"nip_kepala_sekolah"`
	NamaKetuaPanitia    *string        `gorm:"column:nama_ketua_panitia;size:255" json:"nama_ketua_panitia"`
	NIPKetuaPanitia     *string        `gorm:"column:nip_ketua_panitia;size:50" json:"nip_ketua_panitia"`
	GrupWA              *string        `gorm:"column:grup_wa;type:text" json:"grup_wa"`
	CreatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedByID         *uint          `json:"created_by_id"`
	UpdatedByID         *uint          `json:"updated_by_id"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for SettingLayananSPMB
func (m *SettingLayananSPMB) TableName() string {
	return "setting_layanan_spmb"
}
