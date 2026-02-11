package models

import (
	"time"
)

// KutipanKepsek represents the Kutipan Kepsek model
type KutipanKepsek struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaKepsek    string    `gorm:"not null" json:"nama_kepsek"`
	FotoKepsek    string    `gorm:"not null" json:"foto_kepsek"`
	KutipanKepsek string    `gorm:"type:text;not null" json:"kutipan_kepsek"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedByID   *uint     `json:"created_by_id"`
	UpdatedByID   *uint     `json:"updated_by_id"`
}

// TableName specifies the table name for KutipanKepsek
func (m *KutipanKepsek) TableName() string {
	return "kutipan_kepsek"
}
