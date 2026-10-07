package models

import (
	"time"

	"gorm.io/gorm"
)

type AnakInklusiRombel struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	AnakInklusiID           uint           `gorm:"not null" json:"anak_inklusi_id"`
	PesertaDidikRombelID    uint           `gorm:"not null" json:"peserta_didik_rombel_id"`
	GuruPendampingKhususID  *uint          `json:"guru_pendamping_khusus_id"`
	GuruKelasID             *uint          `json:"guru_kelas_id"`
	Catatan                 string         `gorm:"type:text" json:"catatan"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	CreatedByID             *uint          `json:"created_by_id"`
	UpdatedByID             *uint          `json:"updated_by_id"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	AnakInklusi          *AnakInklusi       `gorm:"foreignKey:AnakInklusiID" json:"anak_inklusi,omitempty"`
	PesertaDidikRombel   *PesertaDidikRombel `gorm:"foreignKey:PesertaDidikRombelID" json:"peserta_didik_rombel,omitempty"`
	GuruPendampingKhusus *Kepegawaian       `gorm:"foreignKey:GuruPendampingKhususID" json:"guru_pendamping_khusus,omitempty"`
	GuruKelas            *Kepegawaian       `gorm:"foreignKey:GuruKelasID" json:"guru_kelas,omitempty"`
}

func (AnakInklusiRombel) TableName() string {
	return "anak_inklusi_rombel"
}
