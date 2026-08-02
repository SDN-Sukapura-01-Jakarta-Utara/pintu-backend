package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FormulirPertanyaan represents the FormulirPertanyaan model
type FormulirPertanyaan struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	FormulirID      uint           `gorm:"column:form_id;not null" json:"formulir_id"`
	Urutan          int            `gorm:"not null" json:"urutan"`
	Label           string         `gorm:"type:text;not null" json:"label"`
	Placeholder     *string        `gorm:"type:varchar(255)" json:"placeholder"`
	Tipe            string         `gorm:"type:form_question_type;not null" json:"tipe"`
	IsRequired      bool           `gorm:"default:false" json:"is_required"`
	Options         datatypes.JSON `gorm:"type:jsonb" json:"options"`
	ValidationRules datatypes.JSON `gorm:"type:jsonb" json:"validation_rules"`
	FileConfig      datatypes.JSON `gorm:"type:jsonb" json:"file_config"`
	Dokumen         *string        `gorm:"type:varchar(255)" json:"dokumen"`
	Link            *string        `gorm:"type:text" json:"link"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for FormulirPertanyaan
func (m *FormulirPertanyaan) TableName() string {
	return "form_questions"
}
