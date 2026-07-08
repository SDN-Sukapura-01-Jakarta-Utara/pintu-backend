package models

import (
	"time"

	"gorm.io/gorm"
)

// RekapitulasiAbsensi represents the Rekapitulasi Absensi model for teacher attendance recording
type RekapitulasiAbsensi struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	PesertaDidikRombelID uint           `gorm:"column:peserta_didik_rombel_id;not null" json:"peserta_didik_rombel_id"`
	RombelID             *uint          `gorm:"column:rombel_id" json:"rombel_id"`
	TahunPelajaranID     uint           `gorm:"column:tahun_pelajaran_id;not null" json:"tahun_pelajaran_id"`
	Semester             int            `gorm:"column:semester;not null" json:"semester"`
	Tanggal              time.Time      `gorm:"column:tanggal;type:date;not null" json:"tanggal"`
	BidangStudiID        *uint          `gorm:"column:bidang_studi_id" json:"bidang_studi_id"`
	PertemuanKe          *int           `gorm:"column:pertemuan_ke" json:"pertemuan_ke"`
	Status               string         `gorm:"column:status;size:20;not null" json:"status"`
	WaktuAbsen           *time.Time     `gorm:"column:waktu_absen" json:"waktu_absen"`
	MetodeInput          string         `gorm:"column:metode_input;size:20;not null" json:"metode_input"`
	Keterangan           string         `gorm:"column:keterangan;type:text" json:"keterangan"`
	FileSurat            string         `gorm:"column:file_surat;type:text" json:"file_surat"`
	DicatatOlehID        *uint          `gorm:"column:dicatat_oleh_id" json:"dicatat_oleh_id"`
	CreatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relationships
	PesertaDidikRombel   *PesertaDidikRombel `gorm:"foreignKey:PesertaDidikRombelID" json:"peserta_didik_rombel,omitempty"`
	Rombel               *Rombel             `gorm:"foreignKey:RombelID" json:"rombel,omitempty"`
	TahunPelajaran       *TahunPelajaran     `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	BidangStudi          *BidangStudi        `gorm:"foreignKey:BidangStudiID" json:"bidang_studi,omitempty"`
	DicatatOleh          *User               `gorm:"foreignKey:DicatatOlehID" json:"dicatat_oleh,omitempty"`
}

// TableName specifies the table name for RekapitulasiAbsensi
func (m *RekapitulasiAbsensi) TableName() string {
	return "rekapitulasi_absensi"
}

