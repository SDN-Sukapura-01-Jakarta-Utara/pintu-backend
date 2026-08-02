package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FormulirResponseJawaban represents an answer to a form question
type FormulirResponseJawaban struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ResponseID   uint           `gorm:"column:response_id;not null" json:"response_id"`
	PertanyaanID uint           `gorm:"column:question_id;not null" json:"pertanyaan_id"`
	JawabanText  *string        `gorm:"type:text" json:"jawaban_text"`
	JawabanJSON  datatypes.JSON `gorm:"type:jsonb" json:"jawaban_json"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Response   *FormulirResponse   `gorm:"foreignKey:ResponseID" json:"response,omitempty"`
	Pertanyaan *FormulirPertanyaan `gorm:"foreignKey:PertanyaanID" json:"pertanyaan,omitempty"`
}

// TableName specifies the table name
func (m *FormulirResponseJawaban) TableName() string {
	return "form_response_answers"
}
