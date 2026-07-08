package models

import (
	"time"

	"gorm.io/gorm"
)

// PesertaDidikRombel represents the mapping between PesertaDidik and Rombel
type PesertaDidikRombel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	PesertaDidikID   uint           `gorm:"not null" json:"peserta_didik_id"`
	RombelID         uint           `gorm:"not null" json:"rombel_id"`
	TahunPelajaranID uint           `gorm:"not null" json:"tahun_pelajaran_id"`
	Status           string         `gorm:"default:active" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	CreatedByID      *uint          `json:"created_by_id"`
	UpdatedByID      *uint          `json:"updated_by_id"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// Foreign key relationships
	PesertaDidik     *PesertaDidik     `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
	Rombel           *Rombel           `gorm:"foreignKey:RombelID" json:"rombel,omitempty"`
	TahunPelajaran   *TahunPelajaran   `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
}

// TableName specifies the table name for PesertaDidikRombel
func (m *PesertaDidikRombel) TableName() string {
	return "peserta_didik_rombel"
}
