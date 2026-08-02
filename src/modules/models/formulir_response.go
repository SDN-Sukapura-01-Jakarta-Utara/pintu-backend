package models

import (
	"time"

	"gorm.io/gorm"
)

// FormulirResponse represents a form submission/response
type FormulirResponse struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	FormulirID         uint           `gorm:"column:form_id;not null" json:"formulir_id"`
	SubmittedByUserID  *uint          `gorm:"column:submitted_by_user_id" json:"submitted_by_user_id"`
	SubmittedAsRole    *string        `gorm:"column:submitted_as_role;type:varchar(50)" json:"submitted_as_role"`
	IPAddress          *string        `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent          *string        `gorm:"type:text" json:"user_agent"`
	SubmittedAt        time.Time      `json:"submitted_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Formulir      *Formulir                 `gorm:"foreignKey:FormulirID" json:"formulir,omitempty"`
	SubmittedBy   *User                     `gorm:"foreignKey:SubmittedByUserID" json:"submitted_by,omitempty"`
	Jawaban       []FormulirResponseJawaban `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE" json:"jawaban,omitempty"`
}

// TableName specifies the table name
func (m *FormulirResponse) TableName() string {
	return "form_responses"
}
