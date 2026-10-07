package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type FilePendukung []string

func (f *FilePendukung) Scan(value interface{}) error {
	if value == nil {
		*f = []string{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, f)
}

func (f FilePendukung) Value() (driver.Value, error) {
	if len(f) == 0 {
		return "[]", nil
	}
	return json.Marshal(f)
}

type MonitoringInklusiBulanan struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	PPIInklusiID          uint           `gorm:"not null" json:"ppi_inklusi_id"`
	Bulan                 int16          `gorm:"not null" json:"bulan"`
	Tahun                 int16          `gorm:"not null" json:"tahun"`
	DeskripsiPerkembangan string         `gorm:"type:text;not null" json:"deskripsi_perkembangan"`
	KendalaDitemui        string         `gorm:"type:text" json:"kendala_ditemui"`
	TindakLanjut          string         `gorm:"type:text" json:"tindak_lanjut"`
	FilePendukung         FilePendukung  `gorm:"type:jsonb;default:'[]'" json:"file_pendukung"`
	GuruPengisiID         *uint          `json:"guru_pengisi_id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedByID           *uint          `json:"created_by_id"`
	UpdatedByID           *uint          `json:"updated_by_id"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	PPIInklusi  *PPIInklusi  `gorm:"foreignKey:PPIInklusiID" json:"ppi_inklusi,omitempty"`
	GuruPengisi *Kepegawaian `gorm:"foreignKey:GuruPengisiID" json:"guru_pengisi,omitempty"`
}

func (MonitoringInklusiBulanan) TableName() string {
	return "monitoring_inklusi_bulanan"
}
